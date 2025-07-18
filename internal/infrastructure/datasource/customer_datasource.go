package datasource

import (
	"context"
	"errors"
	"time"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/database"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/datasource/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (ds *customerDataSource) FindByID(ctx context.Context, id string) (*entity.Customer, error) {
	startTime := time.Now()

	mongoFilter := bson.M{"_id": id}

	var customerModel model.CustomerModel
	err := ds.collection.FindOne(ctx, mongoFilter).Decode(&customerModel)

	duration := time.Since(startTime)
	ds.db.LogOperation(ctx, "FindByID", customersCollection, duration, err)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return customerModel.ToEntity(), nil
}

func (ds *customerDataSource) FindByCPF(ctx context.Context, cpf string) (*entity.Customer, error) {
	startTime := time.Now()

	mongoFilter := bson.M{"cpf": cpf}

	var customerModel model.CustomerModel
	err := ds.collection.FindOne(ctx, mongoFilter).Decode(&customerModel)

	duration := time.Since(startTime)
	ds.db.LogOperation(ctx, "FindByCPF", customersCollection, duration, err)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return customerModel.ToEntity(), nil
}

func (ds *customerDataSource) FindAll(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entity.Customer, int64, error) {
	startTime := time.Now()

	// Build MongoDB filter
	mongoFilter := bson.M{}
	for key, value := range filters {
		mongoFilter[key] = value
	}

	// Count total documents
	total, err := ds.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, err
	}

	// Calculate skip for pagination
	skip := (page - 1) * limit

	// Find documents with pagination
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	cursor, err := ds.collection.Find(ctx, mongoFilter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var customerModels []model.CustomerModel
	if err = cursor.All(ctx, &customerModels); err != nil {
		return nil, 0, err
	}

	duration := time.Since(startTime)
	ds.db.LogOperation(ctx, "FindAll", customersCollection, duration, err)

	// Convert to entities
	customers := make([]*entity.Customer, len(customerModels))
	for i, customerModel := range customerModels {
		customers[i] = customerModel.ToEntity()
	}

	return customers, total, nil
}

func (ds *customerDataSource) Create(ctx context.Context, customer *entity.Customer) error {
	startTime := time.Now()

	customerModel := model.FromEntity(customer)
	_, err := ds.collection.InsertOne(ctx, customerModel)

	duration := time.Since(startTime)
	ds.db.LogOperation(ctx, "Create", customersCollection, duration, err)

	return err
}

func (ds *customerDataSource) Update(ctx context.Context, customer *entity.Customer) error {
	startTime := time.Now()

	customerModel := model.FromEntity(customer)
	filter := bson.M{"_id": customer.ID}
	update := bson.M{"$set": customerModel}

	_, err := ds.collection.UpdateOne(ctx, filter, update)

	duration := time.Since(startTime)
	ds.db.LogOperation(ctx, "Update", customersCollection, duration, err)

	return err
}

func (ds *customerDataSource) Delete(ctx context.Context, id string) error {
	startTime := time.Now()

	filter := bson.M{"_id": id}
	_, err := ds.collection.DeleteOne(ctx, filter)

	duration := time.Since(startTime)
	ds.db.LogOperation(ctx, "Delete", customersCollection, duration, err)

	return err
}
