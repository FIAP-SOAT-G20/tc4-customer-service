package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/dto"
)

type CustomerController interface {
	Get(ctx context.Context, presenter Presenter, input dto.GetCustomerInput) ([]byte, error)
}

type CustomerUseCase interface {
	Get(ctx context.Context, input dto.GetCustomerInput) (*entity.Customer, error)
}

type CustomerGateway interface {
	FindOne(ctx context.Context, cpf string) (*entity.Customer, error)
}

type CustomerDataSource interface {
	FindOne(ctx context.Context, filters dto.CustomerDatasourceFilter) (*entity.Customer, error)
}
