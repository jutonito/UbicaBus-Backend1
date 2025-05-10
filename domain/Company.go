package domain

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Company representa la entidad de compañía en la base de datos.
type Company struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Nombre      string             `bson:"nombre"`
	Descripcion string             `bson:"descripcion"`
}

// CrearCompania inserta una nueva compañía en la colección "Companias".
func CrearCompania(ctx context.Context, db *mongo.Database, comp *Company) error {
	comp.ID = primitive.NewObjectID()

	collection := db.Collection("Companias")
	_, err := collection.InsertOne(ctx, comp)
	if err != nil {
		log.Println("Error al insertar compañía:", err)
		return err
	}
	return nil
}

// EditarCompania actualiza los campos no vacíos de una Company existente y devuelve el documento actualizado.
func EditarCompania(ctx context.Context, db *mongo.Database, comp *Company) (*Company, error) {
	collection := db.Collection("Companias")

	updateFields := bson.M{}
	if comp.Nombre != "" {
		updateFields["nombre"] = comp.Nombre
	}
	if comp.Descripcion != "" {
		updateFields["descripcion"] = comp.Descripcion
	}

	// Si no hay campos para actualizar, devolvemos el documento existente
	if len(updateFields) == 0 {
		var existing Company
		if err := collection.FindOne(ctx, bson.M{"_id": comp.ID}).Decode(&existing); err != nil {
			log.Println("Compañía no encontrada:", err)
			return nil, err
		}
		return &existing, nil
	}

	filter := bson.M{"_id": comp.ID}
	update := bson.M{"$set": updateFields}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated Company
	if err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("Compañía no encontrada")
		} else {
			log.Println("Error al editar compañía:", err)
		}
		return nil, err
	}

	return &updated, nil
}

// GetAllCompanies retorna todas las compañías en la colección "Companias".
func GetAllCompanies(ctx context.Context, db *mongo.Database) ([]Company, error) {
	collection := db.Collection("Companias")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error al obtener compañías:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	companies := make([]Company, 0)
	for cursor.Next(ctx) {
		var comp Company
		if err := cursor.Decode(&comp); err != nil {
			log.Println("Error al decodificar compañía:", err)
			continue
		}
		companies = append(companies, comp)
	}
	if err := cursor.Err(); err != nil {
		log.Println("Cursor error en compañías:", err)
		return nil, err
	}

	return companies, nil
}

// GetCompanyByID busca una compañía por su ObjectID.
func GetCompanyByID(ctx context.Context, db *mongo.Database, id primitive.ObjectID) (*Company, error) {
	collection := db.Collection("Companias")

	var comp Company
	if err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&comp); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("Compañía no encontrada:", err)
		}
		return nil, err
	}
	return &comp, nil
}

// GetCompaniesByName retorna todas las compañías cuyo nombre coincide exactamente.
func GetCompaniesByName(ctx context.Context, db *mongo.Database, nombre string) ([]Company, error) {
	collection := db.Collection("Companias")

	filter := bson.M{"nombre": nombre}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println("Error al buscar compañías por nombre:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var companies []Company
	for cursor.Next(ctx) {
		var comp Company
		if err := cursor.Decode(&comp); err != nil {
			log.Println("Error al decodificar compañía:", err)
			continue
		}
		companies = append(companies, comp)
	}
	if err := cursor.Err(); err != nil {
		log.Println("Cursor error en búsqueda de compañías:", err)
		return nil, err
	}

	return companies, nil
}

// DeleteCompany elimina una compañía por su ObjectID.
func DeleteCompany(ctx context.Context, db *mongo.Database, id primitive.ObjectID) error {
	collection := db.Collection("Companias")

	if _, err := collection.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		log.Println("Error al eliminar compañía:", err)
		return err
	}
	return nil
}
