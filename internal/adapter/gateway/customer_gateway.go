package gateway

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/port"
)

type customerGateway struct {
	dataSource port.CustomerDataSource
}

func NewCustomerGateway(dataSource port.CustomerDataSource) port.CustomerGateway {
	return &customerGateway{dataSource}
}

func (g *customerGateway) FindByID(ctx context.Context, id string) (*entity.Customer, error) {
	return g.dataSource.FindByID(ctx, id)
}

func (g *customerGateway) FindByCPF(ctx context.Context, cpf string) (*entity.Customer, error) {
	return g.dataSource.FindByCPF(ctx, cpf)
}

func (g *customerGateway) FindAll(ctx context.Context, name string, page, limit int) ([]*entity.Customer, int64, error) {
	filters := make(map[string]interface{})
	if name != "" {
		filters["name"] = name
	}

	return g.dataSource.FindAll(ctx, filters, page, limit)
}

func (g *customerGateway) Create(ctx context.Context, customer *entity.Customer) error {
	return g.dataSource.Create(ctx, customer)
}

func (g *customerGateway) Update(ctx context.Context, customer *entity.Customer) error {
	return g.dataSource.Update(ctx, customer)
}

func (g *customerGateway) Delete(ctx context.Context, id string) error {
	return g.dataSource.Delete(ctx, id)
}
