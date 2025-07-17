package gateway

import (
	"context"

	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/dto"
	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/port"
)

type customerGateway struct {
	dataSource port.CustomerDataSource
}

func NewCustomerGateway(dataSource port.CustomerDataSource) port.CustomerGateway {
	return &customerGateway{dataSource}
}

func (g *customerGateway) FindOne(ctx context.Context, cpf string) (*entity.Customer, error) {
	return g.dataSource.FindOne(ctx, dto.CustomerDatasourceFilter{CPF: &cpf})
}
