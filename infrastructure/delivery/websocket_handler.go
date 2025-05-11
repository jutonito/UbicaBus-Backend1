// delivery/websocket.go
package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

// WebsocketHandler registra la conexión en el hub.
func WebsocketHandler(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.String(http.StatusInternalServerError, "WebSocket error: %v", err)
			return
		}
		hub.register <- conn
		defer func() { hub.unregister <- conn }()

		// Simplemente mantenemos viva la conexión:
		for {
			if _, _, err := conn.NextReader(); err != nil {
				break
			}
		}
	}
}
