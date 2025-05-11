package application

import (
	"context"
	"errors"

	"UbicaBus/UbicaBusBackend/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// BusLocationService maneja la lógica de negocio relacionada con las localizaciones de buses.
// Trabaja directamente con las funciones de dominio para crear, obtener y eliminar localizaciones.
type BusLocationService struct {
	DB *mongo.Database
}

// NewBusLocationService crea una nueva instancia de BusLocationService
func NewBusLocationService(db *mongo.Database) *BusLocationService {
	return &BusLocationService{DB: db}
}

// GetAllBusLocations obtiene todas las localizaciones de buses de la base de datos.
func (s *BusLocationService) GetAllBusLocations() ([]domain.BusLocation, error) {
	return domain.GetAllBusLocations(context.TODO(), s.DB)
}

// GetBusLocationsByBusID obtiene todas las localizaciones asociadas a un bus.
func (s *BusLocationService) GetBusLocationsByBusID(busIDHex string) ([]domain.BusLocation, error) {
	if busIDHex == "" {
		return nil, errors.New("busID es obligatorio")
	}
	busID, err := primitive.ObjectIDFromHex(busIDHex)
	if err != nil {
		return nil, errors.New("busID inválido")
	}
	return domain.GetBusLocationsByBusID(context.TODO(), s.DB, busID)
}

// RegisterBusLocation crea una nueva localización para un bus.
func (s *BusLocationService) RegisterBusLocation(
	busIDHex string,
	lat, lng float64,
) (primitive.ObjectID, error) {
	// Validaciones básicas
	if busIDHex == "" {
		return primitive.NilObjectID, errors.New("busID es obligatorio")
	}
	busID, err := primitive.ObjectIDFromHex(busIDHex)
	if err != nil {
		return primitive.NilObjectID, errors.New("busID inválido")
	}

	// Construir la entidad de dominio
	bl := &domain.BusLocation{
		BusID: busID,
		Localizacion: domain.Location{
			Lat: lat,
			Lng: lng,
		},
	}

	// Llamar a la función de dominio para insertar
	if err := domain.CrearBusLocation(context.TODO(), s.DB, bl); err != nil {
		return primitive.NilObjectID, err
	}

	return bl.ID, nil
}

// DeleteBusLocation elimina una localización por su ID.
func (s *BusLocationService) DeleteBusLocation(idHex string) error {
	if idHex == "" {
		return errors.New("ID de localización es obligatorio")
	}
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return errors.New("ID de localización inválido")
	}
	return domain.DeleteBusLocation(context.TODO(), s.DB, id)
}
