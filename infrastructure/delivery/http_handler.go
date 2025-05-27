package delivery

import (
	"fmt"
	"net/http"
	"time"

	"UbicaBus/UbicaBusBackend/application"
	"UbicaBus/UbicaBusBackend/domain"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type BusLocationHandler struct {
	BLService *application.BusLocationService
}

type CreateBusLocationReq struct {
	BusID string  `json:"bus_id" binding:"required"`
	Lat   float64 `json:"lat" binding:"required"`
	Lng   float64 `json:"lng" binding:"required"`
}

// NewBusLocationHandler crea nuevo BusLocationHandler
func NewBusLocationHandler(bls *application.BusLocationService) *BusLocationHandler {
	return &BusLocationHandler{BLService: bls}
}

type BusHandler struct {
	BusService *application.BusService
}

type CreateBusReq struct {
	Placa       string    `json:"placa" binding:"required"`
	ConductorID string    `json:"conductor_id" binding:"required"`
	RutaID      string    `json:"ruta_id" binding:"required"`
	FechaInicio time.Time `json:"fecha_inicio" binding:"required"`
	FechaFin    time.Time `json:"fecha_fin" binding:"required"`
}

func NewBusHandler(bs *application.BusService) *BusHandler {
	return &BusHandler{BusService: bs}
}

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

func (h *BusHandler) GetAllBusesHandler(c *gin.Context) {
	buses, err := h.BusService.GetAllBuses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, buses)
}

// GetBusByIDHandler retorna un bus por su ID.
func (h *BusHandler) GetBusByIDHandler(c *gin.Context) {
	id := c.Param("id")
	bus, err := h.BusService.GetBusByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bus)
}

// SearchBusesByPlacaHandler busca buses por placa (?placa=...).
func (h *BusHandler) SearchBusesByPlacaHandler(c *gin.Context) {
	placa := c.Query("placa")
	if placa == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'placa' is required"})
		return
	}
	buses, err := h.BusService.SearchBusesByPlaca(placa)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, buses)
}

// RegisterBusHandler crea un nuevo bus.
func (h *BusHandler) RegisterBusHandler(c *gin.Context) {
	var req CreateBusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.BusService.RegisterBus(
		req.Placa,
		req.ConductorID,
		req.RutaID,
		req.FechaInicio,
		req.FechaFin,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Bus creado correctamente",
		"bus_id":  id.Hex(),
	})
}

// EditBusHandler actualiza un bus existente.
func (h *BusHandler) EditBusHandler(c *gin.Context) {
	id := c.Param("id")
	var req CreateBusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := h.BusService.EditBus(
		id,
		req.Placa,
		req.ConductorID,
		req.RutaID,
		&req.FechaInicio,
		&req.FechaFin,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteBusHandler elimina un bus por su ID.
func (h *BusHandler) DeleteBusHandler(c *gin.Context) {
	id := c.Param("id")
	if err := h.BusService.DeleteBus(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Bus %s eliminado", id)})
}

// GetAllBusLocationsHandler devuelve todas las localizaciones de buses
func (h *BusLocationHandler) GetAllBusLocationsHandler(c *gin.Context) {
	locations, err := h.BLService.GetAllBusLocations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, locations)
}

// GetBusLocationsByBusIDHandler devuelve localizaciones para un bus específico
func (h *BusLocationHandler) GetBusLocationsByBusIDHandler(c *gin.Context) {
	busID := c.Param("bus_id")
	if busID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bus_id requerido"})
		return
	}
	locations, err := h.BLService.GetBusLocationsByBusID(busID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, locations)
}

// RegisterBusLocationHandler registra una nueva localización
func (h *BusLocationHandler) RegisterBusLocationHandler(c *gin.Context) {
	var req CreateBusLocationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.BLService.RegisterBusLocation(req.BusID, req.Lat, req.Lng)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Localización registrada", "id": id.Hex()})
}

// DeleteBusLocationHandler elimina una localización por su ID
func (h *BusLocationHandler) DeleteBusLocationHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id requerido"})
		return
	}
	err := h.BLService.DeleteBusLocation(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Localización %s eliminada", id)})
}

// StartServer inicia el servidor HTTP y registra rutas con Gin
func StartServer(userService *application.UserService, routeService *application.RouteService, companyService *application.CompanyService, roleService *application.RoleService, busService *application.BusService, busLocService *application.BusLocationService) {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Cambia según tu frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Crear el manejador de usuarios
	userHandler := NewUserHandler(userService)
	routeHandler := NewRouteHandler(routeService)
	companyHandler := NewCompanyHandler(companyService)
	roleHandler := NewRoleHandler(roleService)
	busHandler := NewBusHandler(busService)
	busLocHandler := NewBusLocationHandler(busLocService)

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
	r.GET("/buses", busHandler.GetAllBusesHandler)
	r.GET("/buses/search", busHandler.SearchBusesByPlacaHandler) // ?placa=...
	r.GET("/buses/:id", busHandler.GetBusByIDHandler)
	r.POST("/buses", busHandler.RegisterBusHandler)
	r.PUT("/buses/:id", busHandler.EditBusHandler)
	r.DELETE("/buses/:id", busHandler.DeleteBusHandler)
	r.GET("/buslocations", busLocHandler.GetAllBusLocationsHandler)
	r.GET("/buslocations/:bus_id", busLocHandler.GetBusLocationsByBusIDHandler)
	r.POST("/buslocations", busLocHandler.RegisterBusLocationHandler)
	// Para eliminar por id de la localización, no por bus_id:
	r.DELETE("/buslocations/:id", busLocHandler.DeleteBusLocationHandler)

	// Iniciar servidor con Gin
	fmt.Println("Iniciando servidor en el puerto 8080...")
	if err := r.Run(":8080"); err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
