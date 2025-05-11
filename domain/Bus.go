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

func GetAllBuses(ctx context.Context, db *mongo.Database) ([]Bus, error) {
	coll := db.Collection("buses")
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error al obtener buses:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var out []Bus
	for cursor.Next(ctx) {
		var b Bus
		if err := cursor.Decode(&b); err != nil {
			log.Println("Error al decodificar bus:", err)
			continue
		}
		out = append(out, b)
	}
	if err := cursor.Err(); err != nil {
		log.Println("Cursor error en buses:", err)
		return nil, err
	}
	return out, nil
}

// GetBusByID busca un bus por su ObjectID.
func GetBusByID(ctx context.Context, db *mongo.Database, id primitive.ObjectID) (*Bus, error) {
	coll := db.Collection("buses")
	var b Bus
	if err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&b); err != nil {
		log.Println("Bus no encontrado:", err)
		return nil, err
	}
	return &b, nil
}

// GetBusesByPlaca retorna todos los buses cuya placa coincide exactamente.
func GetBusesByPlaca(ctx context.Context, db *mongo.Database, placa string) ([]Bus, error) {
	coll := db.Collection("buses")
	cursor, err := coll.Find(ctx, bson.M{"placa": placa})
	if err != nil {
		log.Println("Error al buscar buses por placa:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var out []Bus
	for cursor.Next(ctx) {
		var b Bus
		if err := cursor.Decode(&b); err != nil {
			log.Println("Error al decodificar bus:", err)
			continue
		}
		out = append(out, b)
	}
	if err := cursor.Err(); err != nil {
		log.Println("Cursor error en búsqueda de buses:", err)
		return nil, err
	}
	return out, nil
}

// DeleteBus elimina un bus por su ObjectID.
func DeleteBus(ctx context.Context, db *mongo.Database, id primitive.ObjectID) error {
	coll := db.Collection("buses")
	if _, err := coll.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		log.Println("Error al eliminar bus:", err)
		return err
	}
	return nil
}
