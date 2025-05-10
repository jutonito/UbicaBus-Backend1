package domain

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Location representa un punto geográfico.
type Location struct {
	Lat float64 `bson:"lat"`
	Lng float64 `bson:"lng"`
}

// Waypoint representa una parada intermedia en una ruta.
type Waypoint struct {
	Lat         float64 `bson:"lat"`
	Lng         float64 `bson:"lng"`
	Descripcion string  `bson:"descripcion"`
}

// Route representa la entidad de ruta en la base de datos.
type Route struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Nombre         string             `bson:"nombre"`
	Descripcion    string             `bson:"descripcion"`
	Origen         Location           `bson:"origen"`
	Destino        Location           `bson:"destino"`
	ModoTransporte string             `bson:"modo_transporte"`
	Waypoints      []Waypoint         `bson:"waypoints"`
}

// CrearRoute inserta una nueva ruta en la colección "ruta".
func CrearRoute(ctx context.Context, db *mongo.Database, r *Route) error {
	r.ID = primitive.NewObjectID()

	collection := db.Collection("ruta")
	_, err := collection.InsertOne(ctx, r)
	if err != nil {
		log.Println("Error al insertar ruta:", err)
		return err
	}
	return nil
}

// EditarRoute actualiza los campos no vacíos de una Route existente
// y devuelve el documento actualizado.
func EditarRoute(ctx context.Context, db *mongo.Database, r *Route) (*Route, error) {
	collection := db.Collection("ruta")

	updateFields := bson.M{}
	if r.Nombre != "" {
		updateFields["nombre"] = r.Nombre
	}
	if r.Descripcion != "" {
		updateFields["descripcion"] = r.Descripcion
	}
	if r.ModoTransporte != "" {
		updateFields["modo_transporte"] = r.ModoTransporte
	}
	// Para campos complejos (origen, destino, waypoints), verificamos si no están vacíos:
	if r.Origen != (Location{}) {
		updateFields["origen"] = r.Origen
	}
	if r.Destino != (Location{}) {
		updateFields["destino"] = r.Destino
	}
	if len(r.Waypoints) > 0 {
		updateFields["waypoints"] = r.Waypoints
	}

	// Si no hay campos para actualizar, devolvemos el documento tal cual está
	if len(updateFields) == 0 {
		var existing Route
		if err := collection.FindOne(ctx, bson.M{"_id": r.ID}).Decode(&existing); err != nil {
			log.Println("Ruta no encontrada:", err)
			return nil, err
		}
		return &existing, nil
	}

	filter := bson.M{"_id": r.ID}
	update := bson.M{"$set": updateFields}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated Route
	if err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("Ruta no encontrada")
		} else {
			log.Println("Error al editar ruta:", err)
		}
		return nil, err
	}

	return &updated, nil
}
