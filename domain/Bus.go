package domain

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Bus representa la entidad de un autobús en la base de datos.
type Bus struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Placa       string             `bson:"placa"`
	ConductorID primitive.ObjectID `bson:"conductor"`
	RutaID      primitive.ObjectID `bson:"ruta"`
	FechaInicio time.Time          `bson:"fecha_inicio"`
	FechaFin    time.Time          `bson:"fecha_fin"`
}

// CrearBus inserta un nuevo documento en la colección "buses".
func CrearBus(ctx context.Context, db *mongo.Database, bus *Bus) error {
	bus.ID = primitive.NewObjectID()

	collection := db.Collection("buses")
	_, err := collection.InsertOne(ctx, bus)
	if err != nil {
		log.Println("Error al insertar bus:", err)
		return err
	}
	return nil
}

// EditarBus actualiza los campos no vacíos de un Bus existente y retorna el documento actualizado.
func EditarBus(ctx context.Context, db *mongo.Database, bus *Bus) (*Bus, error) {
	collection := db.Collection("buses")

	updateFields := bson.M{}

	if bus.Placa != "" {
		updateFields["placa"] = bus.Placa
	}
	if !bus.ConductorID.IsZero() {
		updateFields["conductor"] = bus.ConductorID
	}
	if !bus.RutaID.IsZero() {
		updateFields["ruta"] = bus.RutaID
	}
	if !bus.FechaInicio.IsZero() {
		updateFields["fecha_inicio"] = bus.FechaInicio
	}
	if !bus.FechaFin.IsZero() {
		updateFields["fecha_fin"] = bus.FechaFin
	}

	// Si no hay nada que actualizar, devolvemos el documento existente
	if len(updateFields) == 0 {
		var existing Bus
		if err := collection.FindOne(ctx, bson.M{"_id": bus.ID}).Decode(&existing); err != nil {
			log.Println("Bus no encontrado:", err)
			return nil, err
		}
		return &existing, nil
	}

	filter := bson.M{"_id": bus.ID}
	update := bson.M{"$set": updateFields}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated Bus
	if err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("Bus no encontrado")
		} else {
			log.Println("Error al editar bus:", err)
		}
		return nil, err
	}

	return &updated, nil
}
