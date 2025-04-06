package application

import (
	"context"
	"errors"

	"UbicaBus/UbicaBusBackend/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserService maneja la lógica de negocio relacionada con los usuarios.
type UserService struct {
	DB *mongo.Database
}

// NewUserService crea una nueva instancia de UserService
func NewUserService(db *mongo.Database) *UserService {
	return &UserService{DB: db}
}

// RegisterUser registra un nuevo usuario en la base de datos.
func (s *UserService) RegisterUser(nombre, password, rolID, companiaID string) (primitive.ObjectID, error) {
	if nombre == "" || password == "" || rolID == "" || companiaID == "" {
		return primitive.NilObjectID, errors.New("todos los campos son obligatorios")
	}

	// Convertir a ObjectID
	rolObjID, err := primitive.ObjectIDFromHex(rolID)
	if err != nil {
		return primitive.NilObjectID, errors.New("ID de rol inválido")
	}

	companiaObjID, err := primitive.ObjectIDFromHex(companiaID)
	if err != nil {
		return primitive.NilObjectID, errors.New("ID de compañía inválido")
	}

	// Encriptar contraseña
	hashedPassword := domain.HashPassword(password)

	// Crear usuario
	user := domain.User{
		Nombre:   nombre,
		Password: hashedPassword,
		RolID:    rolObjID,
		Compania: companiaObjID,
	}

	// Insertar en la BD
	err = domain.CrearUsuario(context.TODO(), s.DB, &user)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return user.ID, nil
}
