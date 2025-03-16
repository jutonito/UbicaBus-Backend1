package delivery

import (
	"UbicaBus/UbicaBusBackend/application"
	"net/http"
	"os"

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

	// Obtén el puerto de la variable de entorno PORT o usa 8080 por defecto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Ejecuta el servidor en el puerto especificado
	router.Run(":" + port)
}
