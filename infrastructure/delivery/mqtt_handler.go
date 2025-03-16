package delivery

import (
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// StartMQTT inicia el cliente MQTT, se conecta al broker y se suscribe a un tópico.
func StartMQTT() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883") // Ajusta la URL del broker si es necesario.
	opts.SetClientID("mi-backend-mqtt-client")
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Mensaje recibido en el tópico '%s': %s", msg.Topic(), msg.Payload())
	})
	opts.OnConnect = func(client mqtt.Client) {
		log.Println("Conectado al broker MQTT")
		// Suscribirse al tópico "mi/topico"
		if token := client.Subscribe("mi/topico", 1, nil); token.Wait() && token.Error() != nil {
			log.Fatalf("Error en la suscripción al tópico MQTT: %v", token.Error())
		} else {
			log.Println("Suscrito al tópico: mi/topico")
		}
	}

	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()
	if token.Error() != nil {
		log.Fatalf("Error al conectar al broker MQTT: %v", token.Error())
	}

	// Mantiene la conexión viva
	for {
		time.Sleep(1 * time.Second)
	}
}
