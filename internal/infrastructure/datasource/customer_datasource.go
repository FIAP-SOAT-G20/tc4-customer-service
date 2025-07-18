package datasource

import (
	"context"
	"time"

	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/dto"
	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/port"
	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/infrastructure/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CustomerModel represents the customer model in MongoDB
type CustomerModel struct {
	ID        uint64    `bson:"_id,omitempty"`
	Name      string    `bson:"name"`
	Email     string    `bson:"email"`
	CPF       string    `bson:"cpf"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

// ToEntity converts MongoDB model to domain entity
func (cm *CustomerModel) ToEntity() *entity.Customer {
	return &entity.Customer{
		ID:        cm.ID,
		Name:      cm.Name,
		Email:     cm.Email,
		CPF:       cm.CPF,
		CreatedAt: cm.CreatedAt,
		UpdatedAt: cm.UpdatedAt,
	}
}

// FromEntity converts domain entity to MongoDB model
func (cm *CustomerModel) FromEntity(customer *entity.Customer) {
	cm.ID = customer.ID
	cm.Name = customer.Name
	cm.Email = customer.Email
	cm.CPF = customer.CPF
	cm.CreatedAt = customer.CreatedAt
	cm.UpdatedAt = customer.UpdatedAt
}

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

	var customerModel CustomerModel
	err := ds.collection.FindOne(ctx, mongoFilter).Decode(&customerModel)

	duration := time.Since(startTime)
	ds.db.LogOperation(ctx, "FindOne", customersCollection, duration, err)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Not found
		}
		return nil, err
	}

	return customerModel.ToEntity(), nil
}
