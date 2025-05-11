package application

import (
	"context"
	"errors"
	"time"

	"UbicaBus/UbicaBusBackend/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// BusService maneja la lógica de negocio relacionada con los buses.
type BusService struct {
	DB *mongo.Database
}

// NewBusService crea una nueva instancia de BusService.
func NewBusService(db *mongo.Database) *BusService {
	return &BusService{DB: db}
}

// GetAllBuses obtiene todos los buses.
func (s *BusService) GetAllBuses() ([]domain.Bus, error) {
	return domain.GetAllBuses(context.TODO(), s.DB)
}

// GetBusByID retorna un bus por su ID.
func (s *BusService) GetBusByID(idHex string) (*domain.Bus, error) {
	if idHex == "" {
		return nil, errors.New("ID de bus es obligatorio")
	}
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, errors.New("ID de bus inválido")
	}
	return domain.GetBusByID(context.TODO(), s.DB, id)
}

// SearchBusesByPlaca busca buses por placa exacta.
func (s *BusService) SearchBusesByPlaca(placa string) ([]domain.Bus, error) {
	if placa == "" {
		return nil, errors.New("la placa es obligatoria")
	}
	return domain.GetBusesByPlaca(context.TODO(), s.DB, placa)
}

// RegisterBus crea un nuevo bus.
func (s *BusService) RegisterBus(
	placa, conductorIDHex, rutaIDHex string,
	fechaInicio, fechaFin time.Time,
) (primitive.ObjectID, error) {
	if placa == "" || conductorIDHex == "" || rutaIDHex == "" {
		return primitive.NilObjectID, errors.New("placa, conductor y ruta son obligatorios")
	}
	condID, err := primitive.ObjectIDFromHex(conductorIDHex)
	if err != nil {
		return primitive.NilObjectID, errors.New("ID de conductor inválido")
	}
	rutaID, err := primitive.ObjectIDFromHex(rutaIDHex)
	if err != nil {
		return primitive.NilObjectID, errors.New("ID de ruta inválido")
	}
	bus := domain.Bus{
		Placa:       placa,
		ConductorID: condID,
		RutaID:      rutaID,
		FechaInicio: fechaInicio,
		FechaFin:    fechaFin,
	}
	if err := domain.CrearBus(context.TODO(), s.DB, &bus); err != nil {
		return primitive.NilObjectID, err
	}
	return bus.ID, nil
}

// EditBus actualiza un bus existente.
func (s *BusService) EditBus(
	idHex, placa, conductorIDHex, rutaIDHex string,
	fechaInicio, fechaFin *time.Time,
) (*domain.Bus, error) {
	if idHex == "" {
		return nil, errors.New("ID de bus es obligatorio")
	}
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, errors.New("ID de bus inválido")
	}
	b := &domain.Bus{ID: id}
	if placa != "" {
		b.Placa = placa
	}
	if conductorIDHex != "" {
		if cid, err := primitive.ObjectIDFromHex(conductorIDHex); err == nil {
			b.ConductorID = cid
		} else {
			return nil, errors.New("ID de conductor inválido")
		}
	}
	if rutaIDHex != "" {
		if rid, err := primitive.ObjectIDFromHex(rutaIDHex); err == nil {
			b.RutaID = rid
		} else {
			return nil, errors.New("ID de ruta inválido")
		}
	}
	if fechaInicio != nil {
		b.FechaInicio = *fechaInicio
	}
	if fechaFin != nil {
		b.FechaFin = *fechaFin
	}
	return domain.EditarBus(context.TODO(), s.DB, b)
}

// DeleteBus elimina un bus por su ID.
func (s *BusService) DeleteBus(idHex string) error {
	if idHex == "" {
		return errors.New("ID de bus es obligatorio")
	}
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return errors.New("ID de bus inválido")
	}
	return domain.DeleteBus(context.TODO(), s.DB, id)
}
