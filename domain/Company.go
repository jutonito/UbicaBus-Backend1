package domain

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Company representa la entidad de compañía en la base de datos.
type Company struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Nombre      string             `bson:"nombre"`
	Descripcion string             `bson:"descripcion"`
}

// CrearCompania inserta una nueva compañía en la colección "Companias".
func CrearCompania(ctx context.Context, db *mongo.Database, comp *Company) error {
	comp.ID = primitive.NewObjectID()

	collection := db.Collection("Companias")
	_, err := collection.InsertOne(ctx, comp)
	if err != nil {
		log.Println("Error al insertar compañía:", err)
		return err
	}
	return nil
}

// EditarCompania actualiza los campos no vacíos de una Company existente y devuelve el documento actualizado.
func EditarCompania(ctx context.Context, db *mongo.Database, comp *Company) (*Company, error) {
	collection := db.Collection("Companias")

	updateFields := bson.M{}
	if comp.Nombre != "" {
		updateFields["nombre"] = comp.Nombre
	}
	if comp.Descripcion != "" {
		updateFields["descripcion"] = comp.Descripcion
	}

	// Si no hay campos para actualizar, devolvemos el documento existente
	if len(updateFields) == 0 {
		var existing Company
		if err := collection.FindOne(ctx, bson.M{"_id": comp.ID}).Decode(&existing); err != nil {
			log.Println("Compañía no encontrada:", err)
			return nil, err
		}
		return &existing, nil
	}

	filter := bson.M{"_id": comp.ID}
	update := bson.M{"$set": updateFields}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated Company
	if err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("Compañía no encontrada")
		} else {
			log.Println("Error al editar compañía:", err)
		}
		return nil, err
	}

	return &updated, nil
}
