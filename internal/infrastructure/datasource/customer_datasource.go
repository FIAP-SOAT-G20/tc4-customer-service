package datasource

import (
	"context"
	"errors"
	"time"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/database"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/datasource/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type customerDataSource struct {
	db         *database.MongoDatabase
	collection *mongo.Collection
}

const customersCollection = "customers"

func NewCustomerDataSource(db *database.MongoDatabase) port.CustomerDataSource {
	return &customerDataSource{
		db:         db,
		collection: db.Collection(customersCollection),
	}
}

func (ds *customerDataSource) FindOne(ctx context.Context, filter dto.CustomerDatasourceFilter) (*entity.Customer, error) {
	startTime := time.Now()

	// Build MongoDB filter
	mongoFilter := bson.M{}
	if filter.CPF != nil {
		mongoFilter["cpf"] = *filter.CPF
	}

	var customerModel model.CustomerModel
	err := ds.collection.FindOne(ctx, mongoFilter).Decode(&customerModel)

	duration := time.Since(startTime)
	ds.db.LogOperation(ctx, "FindOne", customersCollection, duration, err)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Not found
		}
		return nil, err
	}

	return customerModel.ToEntity(), nil
}
