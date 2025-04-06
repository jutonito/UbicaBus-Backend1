package delivery

import (
	"fmt"
	"net/http"

	"UbicaBus/UbicaBusBackend/application"

	"github.com/gin-gonic/gin"
)

// UserHandler maneja las peticiones relacionadas con usuarios
type UserHandler struct {
	UserService *application.UserService
}

// NewUserHandler crea un nuevo manejador de usuarios
func NewUserHandler(userService *application.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// RegisterUserHandler maneja el registro de un usuario
func (h *UserHandler) RegisterUserHandler(c *gin.Context) {
	var request struct {
		Nombre     string `json:"nombre"`
		Password   string `json:"password"`
		RolID      string `json:"rol_id"`
		CompaniaID string `json:"compania_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos"})
		return
	}

	userID, err := h.UserService.RegisterUser(request.Nombre, request.Password, request.RolID, request.CompaniaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar usuario: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuario registrado correctamente",
		"user_id": userID.Hex(),
	})
}

// StartServer inicia el servidor HTTP y registra rutas con Gin
func StartServer(userService *application.UserService) {
	r := gin.Default()

	// Crear el manejador de usuarios
	userHandler := NewUserHandler(userService)

	// Registrar rutas
	r.POST("/register", userHandler.RegisterUserHandler)

	// Iniciar servidor con Gin
	fmt.Println("Iniciando servidor en el puerto 8080...")
	if err := r.Run(":8080"); err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
