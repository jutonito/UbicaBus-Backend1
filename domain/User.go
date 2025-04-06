package domain

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

// HashPassword encripta la contraseña usando SHA-256
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
