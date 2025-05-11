package delivery

import (
	"encoding/json"
	"log"

	"UbicaBus/UbicaBusBackend/application" // Asegúrate que esta ruta sea correcta

	mqtt "github.com/mochi-mqtt/server/v2" // Asegúrate de que go.mod apunte a v2
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets" // Asegúrate de importar packets
)

// Asumo que la estructura Hub está definida en otro archivo en este paquete 'delivery'
// y tiene un canal 'broadcast chan []byte' correctamente inicializado y gestionado (ej. con un método Run).

// MessageHook maneja los mensajes MQTT entrantes (PUBLISH).
// Interactúa con la lógica de negocio (DB) y el Hub de WebSockets.
type MessageHook struct {
	mqtt.HookBase                                 // Requerido para ser un hook
	blService     *application.BusLocationService // Servicio para lógica de negocio (guardar en DB)
	hub           *Hub                            // Hub para reenvío a WebSockets
}

// ID identifica este hook.
func (h *MessageHook) ID() string { return "message-hook" }

// Provides indica los eventos que este hook maneja.
// Basado en el código fuente que proporcionaste, el evento para un mensaje PUBLISH entrante es mqtt.OnPublish.
// Este hook implementa el método OnPublish.
func (h *MessageHook) Provides(b byte) bool {
	// *** CORRECTO SEGUN TU CODIGO FUENTE ***
	return b&mqtt.OnPublish != 0
}

// OnPublish es llamado por el broker cuando recibe un paquete PUBLISH de un cliente.
// Es el punto de hook para interceptar mensajes entrantes.
func (h *MessageHook) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	// Log inicial del paquete PUBLISH recibido en el hook.
	// *** CORRECCION FINAL: Acceso directo a QoS, Retain, Dup, PacketID desde pk.Packet ***
	log.Printf("MQTT [Client %s]: << RECEIVED PUBLISH (HOOK) Topic='%s' (QoS %d, Retain %t, DUP %t), PacketID=%d, Payload size=%d bytes",
		cl.ID, pk.TopicName, pk.FixedHeader.Qos, pk.FixedHeader.Retain, pk.FixedHeader.Dup, pk.PacketID, len(pk.Payload)) // Acceso directo

	// 1) Parsear el payload JSON esperado (BusID, Lat, Lng).
	var msg struct {
		BusID string  `json:"bus_id"`
		Lat   float64 `json:"lat"`
		Lng   float64 `json:"lng"`
	}

	// Validar si hay payload.
	if len(pk.Payload) == 0 {
		log.Printf("MQTT [Client %s]: PUBLISH con payload vacío en tópico '%s'. Hook ignora procesamiento JSON.", cl.ID, pk.TopicName)
		return pk, nil // Devuelve el paquete original.
	}

	if err := json.Unmarshal(pk.Payload, &msg); err != nil {
		// Log del error de parsing.
		log.Printf("MQTT [Client %s]: !!! ERROR JSON Unmarshal para PUBLISH en tópico '%s': %v. Payload recibido (muestra): '%s'",
			cl.ID, pk.TopicName, err, string(pk.Payload)) // Truncar payload si es muy largo en un log real
		return pk, nil // Devuelve el paquete original para no detener el broker.
	}

	// Log confirmando el parsing exitoso.
	log.Printf("MQTT [Client %s]: >>> Parsed PUBLISH payload OK for BusID %s: Lat=%f, Lng=%f from topic '%s'", cl.ID, msg.BusID, msg.Lat, msg.Lng, pk.TopicName)

	// Opcional: Filtrar por tópico si necesario.
	// const expectedTopic = "gps/bus/location" // Define el tópico que esperas
	// if pk.TopicName != expectedTopic {
	//     log.Printf("MQTT [Client %s]: PUBLISH en tópico '%s' ignorado por hook. Esperado '%s'.", cl.ID, pk.TopicName, expectedTopic)
	//     return pk, nil // Pasa el paquete sin procesar por este hook.
	// }
	// log.Printf("MQTT [Client %s]: Processing PUBLISH from expected topic '%s'.", cl.ID, pk.TopicName)

	// 2) Guardar la ubicación en la base de datos.
	_, err := h.blService.RegisterBusLocation(msg.BusID, msg.Lat, msg.Lng)
	if err != nil {
		log.Printf("MQTT [Client %s]: !!! ERROR al guardar ubicación de BusID %s en DB: %v", cl.ID, msg.BusID, err)
	} else {
		log.Printf("MQTT [Client %s]: DB save successful for BusID %s.", cl.ID, msg.BusID)
	}

	// 3) Reenviar a WebSockets.
	broadcastPayload := pk.Payload // Asumiendo que el payload MQTT JSON es el mismo que necesitas para WS

	select {
	case h.hub.broadcast <- broadcastPayload:
		log.Printf("MQTT [Client %s]: >>> Sent location of BusID %s to WS Hub broadcast channel.", cl.ID, msg.BusID)
	default:
		log.Printf("MQTT [Client %s]: !!! WARNING: WS Hub broadcast channel FULL. Location for BusID %s NOT sent to WS.", cl.ID, msg.BusID)
	}

	// *** FUNDAMENTAL ***: Retornar el paquete original y nil error para que el broker lo siga procesando.
	return pk, nil
}

// StartMQTT configura y arranca el broker MQTT con el MessageHook.
// El parametro 'topic' se mantiene por compatibilidad con la firma anterior, pero no se usa en el hook por defecto.
func StartMQTT(brokerAddr string, topic string, blService *application.BusLocationService, hub *Hub) {
	log.Println("INFO: Initializing MQTT Broker...")
	server := mqtt.New(nil) // Configuración por defecto

	// Registrar el hook.
	log.Println("INFO: Registering MessageHook...")
	if err := server.AddHook(&MessageHook{blService: blService, hub: hub}, nil); err != nil {
		log.Fatalf("FATAL: Error al registrar el MessageHook: %v", err)
	} else {
		log.Println("INFO: MessageHook registered successfully. It will intercept PUBLISH events.")
	}

	// Añadir listener TCP.
	log.Printf("INFO: Adding TCP listener on %s...", brokerAddr)
	tcpListener := listeners.NewTCP(listeners.Config{
		ID:      "tcp1",
		Address: brokerAddr,
	})
	if err := server.AddListener(tcpListener); err != nil {
		log.Fatalf("FATAL: Error al agregar listener TCP MQTT en %s: %v", brokerAddr, err)
	} else {
		log.Printf("INFO: Listener TCP MQTT agregado exitosamente en %s.", brokerAddr)
	}

	// Iniciar el servidor en goroutine.
	log.Printf("INFO: Starting MQTT broker service...")
	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("FATAL: MQTT Serve error: %v", err)
		}
	}()

	log.Printf("INFO: MQTT Broker started successfully on %s.", brokerAddr)
}
