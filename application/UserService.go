package application

import (
	"UbicaBus/UbicaBusBackend/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserService define la lógica de negocio relacionada a los usuarios.
type UserService struct {
	// Aquí se podría inyectar un repositorio para acceder a la colección 'users'
}

// NewUserService es el constructor de UserService.
func NewUserService() *UserService {
	return &UserService{}
}

// GetAllUsers simula la obtención de todos los usuarios.
// En un caso real, aquí se consultaría la colección 'users' de MongoDB.
func (s *UserService) GetAllUsers() []domain.User {
	return []domain.User{
		{
			ID:       primitive.NewObjectID(),
			Name:     "Ned Stark",
			Email:    "sean_bean@gameofthron.es",
			Password: "$2b$12$UREFwsRUoyF0CRqGNK0LzO0HM/jLhgUCNNIJ9RJAqMUQ74crlJ1Vu",
		},
		{
			ID:       primitive.NewObjectID(),
			Name:     "Ned Starksdasd",
			Email:    "sean_bean@gameofthron.espn",
			Password: "$2b$12$UREFwsRUoyF0CRqGNK0LzO0HM/jLhgUCNNIJ9RJAqMUQ74crlJ1Vu",
		},
	}
}
