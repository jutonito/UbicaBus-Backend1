package domain


import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// BusLocation representa la entidad de localización de un bus.
type BusLocation struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	BusID        primitive.ObjectID `bson:"bus_id"`
	Localizacion Location           `bson:"localizacion"`
	CreatedAt    time.Time          `bson:"created_at"`
}

// CrearBusLocation inserta una nueva localización de bus.
// Valida que el bus exista en la colección "buses".
func CrearBusLocation(ctx context.Context, db *mongo.Database, bl *BusLocation) error {
	// Verificar existencia del bus
	if err := db.Collection("buses").FindOne(ctx, bson.M{"_id": bl.BusID}).Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("bus_id no existe")
		}
		return err
	}

	// Asignar metadatos
	bl.ID = primitive.NewObjectID()
	bl.CreatedAt = time.Now()

	// Insertar documento
	if _, err := db.Collection("BusLocations").InsertOne(ctx, bl); err != nil {
		log.Println("Error al insertar bus location:", err)
		return err
	}
	return nil
}

// GetBusLocationsByBusID retorna todas las localizaciones de un bus.
func GetBusLocationsByBusID(ctx context.Context, db *mongo.Database, busID primitive.ObjectID) ([]BusLocation, error) {
	cursor, err := db.Collection("BusLocations").Find(ctx, bson.M{"bus_id": busID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var out []BusLocation
	for cursor.Next(ctx) {
		var bl BusLocation
		if err := cursor.Decode(&bl); err != nil {
			log.Println("Decode error:", err)
			continue
		}
		out = append(out, bl)
	}
	return out, cursor.Err()
}

// GetAllBusLocations retorna todas las localizaciones.
func GetAllBusLocations(ctx context.Context, db *mongo.Database) ([]BusLocation, error) {
	cursor, err := db.Collection("BusLocations").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var out []BusLocation
	for cursor.Next(ctx) {
		var bl BusLocation
		if err := cursor.Decode(&bl); err != nil {
			log.Println("Decode error:", err)
			continue
		}
		out = append(out, bl)
	}
	return out, cursor.Err()
}

// DeleteBusLocation elimina una localización por su ID.
func DeleteBusLocation(ctx context.Context, db *mongo.Database, id primitive.ObjectID) error {
	if _, err := db.Collection("BusLocations").DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		log.Println("Error al eliminar bus location:", err)
		return err
	}
	return nil
}
