package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain/entity"
)

// CustomerModel represents the customer model in MongoDB
type CustomerModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	CPF       string             `bson:"cpf"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

// ToEntity converts MongoDB model to domain entity
func (cm *CustomerModel) ToEntity() *entity.Customer {
	return &entity.Customer{
		ID:        cm.ID.Hex(),
		Name:      cm.Name,
		Email:     cm.Email,
		CPF:       cm.CPF,
		CreatedAt: cm.CreatedAt,
		UpdatedAt: cm.UpdatedAt,
	}
}

// FromEntity converts domain entity to MongoDB model
func FromEntity(customer *entity.Customer) *CustomerModel {
	cm := &CustomerModel{
		Name:      customer.Name,
		Email:     customer.Email,
		CPF:       customer.CPF,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}

	if customer.ID != "" {
		if objectID, err := primitive.ObjectIDFromHex(customer.ID); err == nil {
			cm.ID = objectID
		}
	}

	return cm
}
