package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// Permite cualquier origen; en producción, restringe el origen según sea necesario.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WebsocketHandler gestiona la conexión WebSocket.
func WebsocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al establecer conexión WebSocket: %v", err)
		return
	}
	defer conn.Close()

	for {
		// Lee un mensaje del cliente
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Envía el mismo mensaje de vuelta (echo)
		if err := conn.WriteMessage(msgType, msg); err != nil {
			break
		}
	}
}
