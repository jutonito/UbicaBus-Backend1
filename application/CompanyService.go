package application

import (
    "context"
    "errors"

    "UbicaBus/UbicaBusBackend/domain"

    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

// CompanyService maneja la lógica de negocio relacionada con compañias.
type CompanyService struct {
    DB *mongo.Database
}

// NewCompanyService crea una nueva instancia de CompanyService.
func NewCompanyService(db *mongo.Database) *CompanyService {
    return &CompanyService{DB: db}
}

// GetAllCompanies obtiene todas las compañias.
func (s *CompanyService) GetAllCompanies() ([]domain.Company, error) {
    return domain.GetAllCompanies(context.TODO(), s.DB)
}

// GetCompanyByID busca una compañía por su ID.
func (s *CompanyService) GetCompanyByID(idHex string) (*domain.Company, error) {
    if idHex == "" {
        return nil, errors.New("ID de compañía es obligatorio")
    }
    id, err := primitive.ObjectIDFromHex(idHex)
    if err != nil {
        return nil, errors.New("ID de compañía inválido")
    }
    return domain.GetCompanyByID(context.TODO(), s.DB, id)
}

// SearchCompaniesByName busca compañías por nombre exacto.
func (s *CompanyService) SearchCompaniesByName(nombre string) ([]domain.Company, error) {
    if nombre == "" {
        return nil, errors.New("el nombre de la compañía es obligatorio")
    }
    return domain.GetCompaniesByName(context.TODO(), s.DB, nombre)
}

// RegisterCompany crea una nueva compañía.
func (s *CompanyService) RegisterCompany(nombre, descripcion string) (primitive.ObjectID, error) {
    if nombre == "" {
        return primitive.NilObjectID, errors.New("el nombre es obligatorio")
    }
    comp := domain.Company{
        Nombre:      nombre,
        Descripcion: descripcion,
    }
    if err := domain.CrearCompania(context.TODO(), s.DB, &comp); err != nil {
        return primitive.NilObjectID, err
    }
    return comp.ID, nil
}

// EditCompany actualiza una compañía existente.
func (s *CompanyService) EditCompany(idHex, nombre, descripcion string) (*domain.Company, error) {
    if idHex == "" {
        return nil, errors.New("ID de compañía es obligatorio")
    }
    id, err := primitive.ObjectIDFromHex(idHex)
    if err != nil {
        return nil, errors.New("ID de compañía inválido")
    }
    comp := &domain.Company{ID: id}
    if nombre != "" {
        comp.Nombre = nombre
    }
    if descripcion != "" {
        comp.Descripcion = descripcion
    }
    return domain.EditarCompania(context.TODO(), s.DB, comp)
}

// DeleteCompany elimina una compañía por su ID.
func (s *CompanyService) DeleteCompany(idHex string) error {
    if idHex == "" {
        return errors.New("ID de compañía es obligatorio")
    }
    id, err := primitive.ObjectIDFromHex(idHex)
    if err != nil {
        return errors.New("ID de compañía inválido")
    }
    return domain.DeleteCompany(context.TODO(), s.DB, id)
}
