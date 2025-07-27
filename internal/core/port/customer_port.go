package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/dto"
)

type CustomerController interface {
	List(ctx context.Context, presenter Presenter, input dto.ListCustomersInput) ([]byte, error)
	Create(ctx context.Context, presenter Presenter, input dto.CreateCustomerInput) ([]byte, error)
	Get(ctx context.Context, presenter Presenter, input dto.GetCustomerInput) ([]byte, error)
	GetByCPF(ctx context.Context, presenter Presenter, input dto.GetCustomerByCPFInput) ([]byte, error)
	Update(ctx context.Context, presenter Presenter, input dto.UpdateCustomerInput) ([]byte, error)
	Delete(ctx context.Context, presenter Presenter, input dto.DeleteCustomerInput) ([]byte, error)
}

type CustomerUseCase interface {
	List(ctx context.Context, input dto.ListCustomersInput) ([]*entity.Customer, int64, error)
	Create(ctx context.Context, input dto.CreateCustomerInput) (*entity.Customer, error)
	Get(ctx context.Context, input dto.GetCustomerInput) (*entity.Customer, error)
	GetByCPF(ctx context.Context, i dto.GetCustomerByCPFInput) (*entity.Customer, error)
	Update(ctx context.Context, input dto.UpdateCustomerInput) (*entity.Customer, error)
	Delete(ctx context.Context, input dto.DeleteCustomerInput) (*entity.Customer, error)
}

type CustomerGateway interface {
	FindByID(ctx context.Context, id int) (*entity.Customer, error)
	FindByCPF(ctx context.Context, cpf string) (*entity.Customer, error)
	FindAll(ctx context.Context, name string, page, limit int) ([]*entity.Customer, int64, error)
	Create(ctx context.Context, customer *entity.Customer) error
	Update(ctx context.Context, customer *entity.Customer) error
	Delete(ctx context.Context, id int) error
}

type CustomerDataSource interface {
	FindByID(ctx context.Context, id int) (*entity.Customer, error)
	FindByCPF(ctx context.Context, cpf string) (*entity.Customer, error)
	FindAll(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entity.Customer, int64, error)
	Create(ctx context.Context, product *entity.Customer) error
	Update(ctx context.Context, product *entity.Customer) error
	Delete(ctx context.Context, id int) error
}
