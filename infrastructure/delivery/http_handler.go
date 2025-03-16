package delivery

import (
	"UbicaBus/UbicaBusBackend/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StartServer inicia el servidor HTTP y configura los endpoints.
func StartServer() {
	router := gin.Default()

	// Inicializa el servicio de usuario
	userService := application.NewUserService()

	// Endpoint para obtener todos los usuarios
	router.GET("/users", func(c *gin.Context) {
		users := userService.GetAllUsers()
		c.JSON(http.StatusOK, gin.H{
			"users": users,
		})
	})

	// Endpoint de ejemplo "Hola Mundo"
	router.GET("/hola", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"mensaje": "Hola Mundo"})
	})

	// Endpoint para WebSockets
	// Nota: en este endpoint se establece la conexión WebSocket
	router.GET("/ws", WebsocketHandler)

	// Ejecuta el servidor en el puerto 8080
	router.Run(":8080")
}
