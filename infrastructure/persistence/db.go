package persistence

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	initOnce   sync.Once // Para garantizar que InitDB solo se ejecuta una vez
	clientOnce sync.Once // Para cerrar la conexión una sola vez
)

// InitDB inicializa la conexión a MongoDB de forma segura
func InitDB() (*mongo.Client, error) {
	var err error

	initOnce.Do(func() { // Garantiza que solo se ejecuta una vez
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		clientOptions := options.Client().ApplyURI("mongodb+srv://root:Elizabeth3004@cluster0.9rjse.mongodb.net/users")
		client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Println("Error conectando a la base de datos:", err)
			return
		}
		log.Println("Conexión a MongoDB exitosa.")
	})

	return client, err
}

// GetDB retorna la base de datos 'Development'
func GetDB() *mongo.Database {
	if client == nil {
		log.Fatal("La conexión a la base de datos no ha sido inicializada")
	}
	return client.Database("Development")
}

// CloseDB cierra la conexión con MongoDB cuando la aplicación termina
func CloseDB() {
	clientOnce.Do(func() {
		if client != nil {
			if err := client.Disconnect(context.Background()); err != nil {
				log.Println("Error cerrando la conexión a MongoDB:", err)
			} else {
				log.Println("Conexión a MongoDB cerrada correctamente.")
			}
		}
	})
}
