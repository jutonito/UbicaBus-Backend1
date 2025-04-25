package application

import (
	"UbicaBus/UbicaBusBackend/domain"
	"UbicaBus/UbicaBusBackend/infrastructure/persistence" // Import para crear el repo (si es necesario aquí)
	"go.mongodb.org/mongo-driver/mongo"             // Import para crear el repo (si es necesario aquí)
)

// Definición de la estructura RouteService
type RouteService struct {
	repo domain.RouteRepository
}

// Constructor para RouteService que recibe la interfaz RouteRepository
func NewRouteService(repo domain.RouteRepository) *RouteService {
	return &RouteService{repo: repo}
}

// Otra forma del constructor si la capa de aplicación necesita crear el repositorio
// En este caso, aún dependemos de la implementación concreta aquí.
// Una mejor práctica sería inyectar la dependencia del repositorio desde una capa superior.
// func NewRouteService(db *mongo.Database) *RouteService {
// 	repo := persistence.NewRouteRepository(db)
// 	return &RouteService{repo: repo}
// }

func (s *RouteService) GetAllRoutes() ([]domain.Route, error) {
	return s.repo.GetAllRoutes() // Corregido el nombre del método para que coincida
}

// Aquí puedes agregar otros métodos para la lógica de negocio relacionada con las rutas
// Por ejemplo:
// func (s *RouteService) GetRouteByID(id string) (*domain.Route, error) {
// 	// ... lógica para obtener una ruta por ID
// }
// func (s *RouteService) CreateRoute(route *domain.Route) error {
// 	// ... lógica para crear una nueva ruta
// }