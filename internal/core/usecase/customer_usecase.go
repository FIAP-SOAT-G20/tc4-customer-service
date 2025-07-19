package usecase

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/port"
)

type customerUseCase struct {
	gateway port.CustomerGateway
}

// NewCustomerUseCase creates a new CreateCustomerUseCase
func NewCustomerUseCase(gateway port.CustomerGateway) port.CustomerUseCase {
	return &customerUseCase{gateway}
}

// Get return a customer by his CPF
func (uc *customerUseCase) Get(ctx context.Context, i dto.GetCustomerInput) (*entity.Customer, error) {
	var cpf = i.CPF
	if cpf == "" {
		cpf = "000.000.000-00"
	}

	customers, err := uc.gateway.FindOne(ctx, cpf)
	if err != nil {
		return nil, domain.NewInternalError(err)
	}

	if customers == nil {
		return nil, domain.NewNotFoundError("customer not found")
	}

	return customers, nil
}
