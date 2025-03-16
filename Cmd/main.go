package main

import (
	"UbicaBus/UbicaBusBackend/infrastructure/delivery"
	"UbicaBus/UbicaBusBackend/infrastructure/persistence"
	"log"
)

func main() {
	// Inicializa la conexión a la base de datos (patrón Singleton)
	if err := persistence.InitDB(); err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	go delivery.StartMQTT()

	delivery.StartServer()

}
