// application/role_service.go

package application

import (
	"context"
	"errors"

	"UbicaBus/UbicaBusBackend/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RoleService maneja la lógica de negocio relacionada con los roles.
type RoleService struct {
	DB *mongo.Database
}

// NewRoleService crea una nueva instancia de RoleService.
func NewRoleService(db *mongo.Database) *RoleService {
	return &RoleService{DB: db}
}

// GetAllRoles obtiene todos los roles.
func (s *RoleService) GetAllRoles() ([]domain.Role, error) {
	return domain.GetAllRoles(context.TODO(), s.DB)
}

// GetRoleByID busca un rol por su ID.
func (s *RoleService) GetRoleByID(idHex string) (*domain.Role, error) {
	if idHex == "" {
		return nil, errors.New("ID de rol es obligatorio")
	}
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, errors.New("ID de rol inválido")
	}
	return domain.GetRoleByID(context.TODO(), s.DB, id)
}

// SearchRolesByName busca roles por nombre exacto.
func (s *RoleService) SearchRolesByName(nombre string) ([]domain.Role, error) {
	if nombre == "" {
		return nil, errors.New("el nombre del rol es obligatorio")
	}
	return domain.GetRolesByName(context.TODO(), s.DB, nombre)
}

// RegisterRole crea un nuevo rol.
func (s *RoleService) RegisterRole(nombre, descripcion string) (primitive.ObjectID, error) {
	if nombre == "" {
		return primitive.NilObjectID, errors.New("el nombre es obligatorio")
	}
	r := domain.Role{
		Nombre:      nombre,
		Descripcion: descripcion,
	}
	if err := domain.CrearRole(context.TODO(), s.DB, &r); err != nil {
		return primitive.NilObjectID, err
	}
	return r.ID, nil
}

// EditRole actualiza un rol existente.
func (s *RoleService) EditRole(idHex, nombre, descripcion string) (*domain.Role, error) {
	if idHex == "" {
		return nil, errors.New("ID de rol es obligatorio")
	}
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, errors.New("ID de rol inválido")
	}
	r := &domain.Role{ID: id}
	if nombre != "" {
		r.Nombre = nombre
	}
	if descripcion != "" {
		r.Descripcion = descripcion
	}
	return domain.EditarRole(context.TODO(), s.DB, r)
}

// DeleteRole elimina un rol por su ID.
func (s *RoleService) DeleteRole(idHex string) error {
	if idHex == "" {
		return errors.New("ID de rol es obligatorio")
	}
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return errors.New("ID de rol inválido")
	}
	return domain.DeleteRole(context.TODO(), s.DB, id)
}
