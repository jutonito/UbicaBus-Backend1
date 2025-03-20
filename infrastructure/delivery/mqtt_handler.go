package delivery

import (
	"log"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
)

// StartMQTT inicia el cliente MQTT, se conecta al broker y se suscribe a un tópico.
func StartMQTT() {
	// Crea un nuevo servidor MQTT.
	server := mqtt.New(nil)

	// Agrega un listener TCP en el puerto 1883.
	tcp := listeners.NewTCP(listeners.Config{
		ID:      "tcp1",
		Address: ":1883",
	})
	if err := server.AddListener(tcp); err != nil {
		log.Fatalf("Error al agregar listener: %v", err)
	}

	// Inicia el servidor en una goroutine.
	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("Error al servir: %v", err)
		}
	}()

	log.Println("Broker MQTT corriendo en el puerto 1883...")
	// Mantiene el programa en ejecución.
	select {}
}
