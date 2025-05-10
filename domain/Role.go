package domain

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Role representa la entidad de rol en la base de datos.
type Role struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Nombre      string             `bson:"nombre"`
	Descripcion string             `bson:"descripcion"`
}

// CrearRole inserta un nuevo rol en la colección "roles".
func CrearRole(ctx context.Context, db *mongo.Database, role *Role) error {
	role.ID = primitive.NewObjectID()

	collection := db.Collection("roles")
	_, err := collection.InsertOne(ctx, role)
	if err != nil {
		log.Println("Error al insertar rol:", err)
		return err
	}
	return nil
}

// EditarRole actualiza los campos no vacíos de un Role existente y devuelve el documento actualizado.
func EditarRole(ctx context.Context, db *mongo.Database, role *Role) (*Role, error) {
	collection := db.Collection("roles")

	updateFields := bson.M{}
	if role.Nombre != "" {
		updateFields["nombre"] = role.Nombre
	}
	if role.Descripcion != "" {
		updateFields["descripcion"] = role.Descripcion
	}

	// Si no hay campos para actualizar, devolvemos el documento existente
	if len(updateFields) == 0 {
		var existing Role
		if err := collection.FindOne(ctx, bson.M{"_id": role.ID}).Decode(&existing); err != nil {
			log.Println("Rol no encontrado:", err)
			return nil, err
		}
		return &existing, nil
	}

	filter := bson.M{"_id": role.ID}
	update := bson.M{"$set": updateFields}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated Role
	if err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("Rol no encontrado")
		} else {
			log.Println("Error al editar rol:", err)
		}
		return nil, err
	}

	return &updated, nil
}
