package domain

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User representa la entidad de usuario en la base de datos.
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Nombre    string             `bson:"nombre"`
	Password  string             `bson:"password"`
	RolID     primitive.ObjectID `bson:"rol"`
	Compania  primitive.ObjectID `bson:"compania"`
	CreatedAt time.Time          `bson:"created_at"`
}

// CrearUsuario inserta un nuevo usuario en la colección "usuarios"
func CrearUsuario(ctx context.Context, db *mongo.Database, user *User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()

	collection := db.Collection("usuarios")

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Println("Error al insertar usuario:", err)
		return err
	}
	return nil
}

func EditarUsuario(ctx context.Context, db *mongo.Database, user *User) (*User, error) {
	collection := db.Collection("usuarios")

	// Esto es un set { "Key": value }
	updateFields := bson.M{}

	if user.Nombre != "" {
		updateFields["nombre"] = user.Nombre
	}
	if user.Password != "" {
		updateFields["password"] = user.Password
	}
	if !user.RolID.IsZero() {
		updateFields["rol"] = user.RolID
	}
	if !user.Compania.IsZero() {
		updateFields["compania"] = user.Compania
	}

	if len(updateFields) == 0 {
		var existing User
		if err := collection.FindOne(ctx, bson.M{"_id": user.ID}).Decode(&existing); err != nil {
			log.Println("Not Found User", err)
			return nil, err
		}
	}

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": updateFields}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated User
	if err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("User not found")
		} else {
			log.Println("Error al editar usuario:", err)
		}
		return nil, err
	}

	return &updated, nil
}

// HashPassword encripta la contraseña usando SHA-256
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
