package delivery

import (
	"fmt"
	"net/http"

	"UbicaBus/UbicaBusBackend/application"
	"UbicaBus/UbicaBusBackend/domain"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	RoleService *application.RoleService
}

type CreateRoleReq struct {
	Nombre      string `json:"nombre" binding:"required"`
	Descripcion string `json:"descripcion"`
}

func NewRoleHandler(rs *application.RoleService) *RoleHandler {
	return &RoleHandler{RoleService: rs}
}

type CompanyHandler struct {
	CompanyService *application.CompanyService
}

type CreateCompanyReq struct {
	Nombre      string `json:"nombre" binding:"required"`
	Descripcion string `json:"descripcion"`
}

func NewCompanyHandler(cs *application.CompanyService) *CompanyHandler {
	return &CompanyHandler{CompanyService: cs}
}

// UserHandler maneja las peticiones relacionadas con usuarios
type UserHandler struct {
	UserService *application.UserService
}

type EditUserReq struct {
	Nombre     string `json:"nombre"`
	Password   string `json:"password"`
	RolID      string `json:"rol_id"`
	CompaniaID string `json:"compania_id"`
}

type RouteHandler struct {
	RouteService *application.RouteService
}

type CreateRouteReq struct {
	Nombre         string            `json:"nombre" binding:"required"`
	Descripcion    string            `json:"descripcion"`
	ModoTransporte string            `json:"modo_transporte" binding:"required"`
	OrigenLat      float64           `json:"origen_lat" binding:"required"`
	OrigenLng      float64           `json:"origen_lng" binding:"required"`
	DestinoLat     float64           `json:"destino_lat" binding:"required"`
	DestinoLng     float64           `json:"destino_lng" binding:"required"`
	Waypoints      []domain.Waypoint `json:"waypoints"`
}

// NewRouteHandler crea un nuevo manejador de rutas.
func NewRouteHandler(routeService *application.RouteService) *RouteHandler {
	return &RouteHandler{RouteService: routeService}
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

func (h *UserHandler) EditUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario requerido"})
		return
	}

	var request EditUserReq
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos"})
		return
	}

	updated, err := h.UserService.EditUser(id, request.Nombre, request.Password, request.RolID, request.CompaniaID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Error al editar usuario: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *RouteHandler) GetAllRoutesHandler(c *gin.Context) {
	routes, err := h.RouteService.GetAllRoutes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, routes)
}

func (h *RouteHandler) GetRoutesByNameHandler(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'name' is required"})
		return
	}

	routes, err := h.RouteService.GetRoutesByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(routes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("no routes found with name '%s'", name)})
		return
	}
	c.JSON(http.StatusOK, routes)
}

func (h *RouteHandler) RegisterRouteHandler(c *gin.Context) {
	var req CreateRouteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.RouteService.RegisterRoute(
		req.Nombre,
		req.Descripcion,
		req.ModoTransporte,
		req.OrigenLat,
		req.OrigenLng,
		req.DestinoLat,
		req.DestinoLng,
		req.Waypoints,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Ruta creada correctamente",
		"route_id": id.Hex(),
	})
}

func (h *RouteHandler) EditRouteHandler(c *gin.Context) {
	routeID := c.Param("id")
	if routeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de ruta requerido"})
		return
	}

	var req CreateRouteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Optional: si quieres permitir no enviar campos obligatorios para ediciones parciales,
	// quita el binding:"required" de CreateRouteReq y sólo valídalos en el service.

	updated, err := h.RouteService.EditRoute(
		routeID,
		req.Nombre,
		req.Descripcion,
		req.ModoTransporte,
		&domain.Location{Lat: req.OrigenLat, Lng: req.OrigenLng},
		&domain.Location{Lat: req.DestinoLat, Lng: req.DestinoLng},
		req.Waypoints,
	)
	if err != nil {
		status := http.StatusInternalServerError
		// podrías distinguir ErrNotFound para 404, etc.
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *CompanyHandler) GetAllCompaniesHandler(c *gin.Context) {
	companies, err := h.CompanyService.GetAllCompanies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, companies)
}

// GetCompanyByIDHandler retorna una compañía por su ID.
func (h *CompanyHandler) GetCompanyByIDHandler(c *gin.Context) {
	idHex := c.Param("id")
	comp, err := h.CompanyService.GetCompanyByID(idHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comp)
}

// SearchCompaniesByNameHandler busca compañías por nombre exacto (?name=...).
func (h *CompanyHandler) SearchCompaniesByNameHandler(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'name' is required"})
		return
	}
	companies, err := h.CompanyService.SearchCompaniesByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, companies)
}

// RegisterCompanyHandler crea una nueva compañía.
func (h *CompanyHandler) RegisterCompanyHandler(c *gin.Context) {
	var req CreateCompanyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.CompanyService.RegisterCompany(req.Nombre, req.Descripcion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":    "Compañía creada correctamente",
		"company_id": id.Hex(),
	})
}

// EditCompanyHandler actualiza una compañía existente.
func (h *CompanyHandler) EditCompanyHandler(c *gin.Context) {
	idHex := c.Param("id")
	var req CreateCompanyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := h.CompanyService.EditCompany(idHex, req.Nombre, req.Descripcion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteCompanyHandler elimina una compañía por su ID.
func (h *CompanyHandler) DeleteCompanyHandler(c *gin.Context) {
	idHex := c.Param("id")
	if err := h.CompanyService.DeleteCompany(idHex); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Compañía %s eliminada", idHex)})
}

func (h *RoleHandler) GetAllRolesHandler(c *gin.Context) {
	roles, err := h.RoleService.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

// GetRoleByIDHandler retorna un rol por su ID.
func (h *RoleHandler) GetRoleByIDHandler(c *gin.Context) {
	idHex := c.Param("id")
	role, err := h.RoleService.GetRoleByID(idHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

// SearchRolesByNameHandler busca roles por nombre exacto (?name=...).
func (h *RoleHandler) SearchRolesByNameHandler(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'name' is required"})
		return
	}
	roles, err := h.RoleService.SearchRolesByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

// RegisterRoleHandler crea un nuevo rol.
func (h *RoleHandler) RegisterRoleHandler(c *gin.Context) {
	var req CreateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.RoleService.RegisterRole(req.Nombre, req.Descripcion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Rol creado correctamente",
		"role_id": id.Hex(),
	})
}

// EditRoleHandler actualiza un rol existente.
func (h *RoleHandler) EditRoleHandler(c *gin.Context) {
	idHex := c.Param("id")
	var req CreateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := h.RoleService.EditRole(idHex, req.Nombre, req.Descripcion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteRoleHandler elimina un rol por su ID.
func (h *RoleHandler) DeleteRoleHandler(c *gin.Context) {
	idHex := c.Param("id")
	if err := h.RoleService.DeleteRole(idHex); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Rol %s eliminado", idHex)})
}

// StartServer inicia el servidor HTTP y registra rutas con Gin
func StartServer(userService *application.UserService, routeService *application.RouteService, companyService *application.CompanyService, roleService *application.RoleService) {
	r := gin.Default()

	// Crear el manejador de usuarios
	userHandler := NewUserHandler(userService)
	routeHandler := NewRouteHandler(routeService)
	companyHandler := NewCompanyHandler(companyService)
	roleHandler := NewRoleHandler(roleService)

	// Registrar rutas
	r.POST("/register", userHandler.RegisterUserHandler)
	r.PUT("/user/:id", userHandler.EditUser)
	r.GET("/routes", routeHandler.GetAllRoutesHandler) // devuelve todas las rutas
	r.GET("/routes/search", routeHandler.GetRoutesByNameHandler)
	r.POST("/routes", routeHandler.RegisterRouteHandler) // Crear ruta
	r.PUT("/routes/:id", routeHandler.EditRouteHandler)
	r.GET("/companies", companyHandler.GetAllCompaniesHandler)
	r.GET("/companies/search", companyHandler.SearchCompaniesByNameHandler) // ?name=...
	r.GET("/companies/:id", companyHandler.GetCompanyByIDHandler)
	r.POST("/companies", companyHandler.RegisterCompanyHandler)
	r.PUT("/companies/:id", companyHandler.EditCompanyHandler)
	r.DELETE("/companies/:id", companyHandler.DeleteCompanyHandler)
	r.GET("/roles", roleHandler.GetAllRolesHandler)
	r.GET("/roles/search", roleHandler.SearchRolesByNameHandler)
	r.GET("/roles/:id", roleHandler.GetRoleByIDHandler)
	r.POST("/roles", roleHandler.RegisterRoleHandler)
	r.PUT("/roles/:id", roleHandler.EditRoleHandler)
	r.DELETE("/roles/:id", roleHandler.DeleteRoleHandler)

	// Iniciar servidor con Gin
	fmt.Println("Iniciando servidor en el puerto 8080...")
	if err := r.Run(":8080"); err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
