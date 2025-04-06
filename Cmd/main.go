package main

import (
	"log"

	"UbicaBus/UbicaBusBackend/application"
	"UbicaBus/UbicaBusBackend/infrastructure/delivery"
	"UbicaBus/UbicaBusBackend/infrastructure/persistence"
)

func main() {
	log.Print("El servidor está corriendo!")

	// Inicializa la conexión a la base de datos
	client, err := persistence.InitDB()
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	// Selecciona la base de datos 'Development'
	db := client.Database("Development")

	// Crear el servicio de usuario
	userService := application.NewUserService(db)

	// Iniciar servidor con el servicio de usuario
	delivery.StartServer(userService)
}
