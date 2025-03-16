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
	clientInstance    *mongo.Client
	clientInstanceErr error
	mongoOnce         sync.Once
)

// InitDB inicializa la conexión a MongoDB usando el patrón Singleton.
func InitDB() error {
	mongoOnce.Do(func() {
		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
		client, err := mongo.NewClient(clientOptions)
		if err != nil {
			clientInstanceErr = err
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = client.Connect(ctx)
		if err != nil {
			clientInstanceErr = err
			return
		}

		if err = client.Ping(ctx, nil); err != nil {
			clientInstanceErr = err
			return
		}

		clientInstance = client
		log.Println("Conexión a MongoDB establecida exitosamente.")
	})

	return clientInstanceErr
}
