package application

import (
	"context"
	"errors"

	"UbicaBus/UbicaBusBackend/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RouteService maneja la lógica de negocio relacionada con las rutas.
// Ahora trabaja directamente con las funciones de dominio para crear y editar rutas.
type RouteService struct {
	DB *mongo.Database
}

// NewRouteService crea una nueva instancia de RouteService
func NewRouteService(db *mongo.Database) *RouteService {
	return &RouteService{DB: db}
}

// GetAllRoutes obtiene todas las rutas de la base de datos.
func (s *RouteService) GetAllRoutes() ([]domain.Route, error) {
	// Asumimos que existe un método en dominio para listar rutas.
	// Si no, se podría implementar usando s.DB.Collection("ruta").Find(...)
	return domain.GetAllRoutes(context.TODO(), s.DB)
}

// RegisterRoute crea una nueva ruta con los datos proporcionados.
func (s *RouteService) RegisterRoute(
	nombre, descripcion, modoTransporte string,
	origenLat, origenLng, destinoLat, destinoLng float64,
	waypoints []domain.Waypoint,
) (primitive.ObjectID, error) {
	// Validaciones básicas
	if nombre == "" || modoTransporte == "" {
		return primitive.NilObjectID, errors.New("nombre y modo de transporte son obligatorios")
	}

	// Construir la entidad de dominio
	route := domain.Route{
		Nombre:         nombre,
		Descripcion:    descripcion,
		ModoTransporte: modoTransporte,
		Origen: domain.Location{
			Lat: origenLat,
			Lng: origenLng,
		},
		Destino: domain.Location{
			Lat: destinoLat,
			Lng: destinoLng,
		},
		Waypoints: waypoints,
	}

	// Insertar en la base de datos
	if err := domain.CrearRoute(context.TODO(), s.DB, &route); err != nil {
		return primitive.NilObjectID, err
	}

	return route.ID, nil
}

// EditRoute actualiza una ruta existente con los campos proporcionados.
func (s *RouteService) EditRoute(
	idHex, nombre, descripcion, modoTransporte string,
	origen *domain.Location, destino *domain.Location,
	waypoints []domain.Waypoint,
) (*domain.Route, error) {
	// Validar ID
	if idHex == "" {
		return nil, errors.New("ID de ruta es obligatorio")
	}
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, errors.New("ID de ruta inválido")
	}

	// Preparar entidad con solo los campos a actualizar
	r := &domain.Route{ID: id}
	if nombre != "" {
		r.Nombre = nombre
	}
	if descripcion != "" {
		r.Descripcion = descripcion
	}
	if modoTransporte != "" {
		r.ModoTransporte = modoTransporte
	}
	if origen != nil {
		r.Origen = *origen
	}
	if destino != nil {
		r.Destino = *destino
	}
	if len(waypoints) > 0 {
		r.Waypoints = waypoints
	}

	// Llamar a la función de dominio para actualizar
	updated, err := domain.EditarRoute(context.TODO(), s.DB, r)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *RouteService) GetRoutesByName(nombre string) ([]domain.Route, error) {
	if nombre == "" {
		return nil, errors.New("el nombre de la ruta es obligatorio")
	}

	return domain.GetRoutesByName(context.TODO(), s.DB, nombre)
}