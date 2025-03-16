package application

import "UbicaBus/UbicaBusBackend/domain"

// UserService define la lógica de negocio relacionada a los usuarios.
type UserService struct {
	// Aquí se podrían inyectar repositorios o dependencias necesarias.
}

// NewUserService es el constructor del UserService.
func NewUserService() *UserService {
	return &UserService{}
}

// GetAllUsers simula la obtención de todos los usuarios.
func (s *UserService) GetAllUsers() []domain.User {
	// En un caso real, aquí se interactuaría con la capa de persistencia.
	return []domain.User{
		{Id: "1", Name: "Juan"},
		{Id: "2", Name: "Juan Nada"},
	}
}
