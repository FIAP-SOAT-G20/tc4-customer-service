package presenter

import (
	"encoding/json"
	"errors"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/port"
)

type customerJsonPresenter struct{}

// CustomerJsonResponse represents the response of a customer
func NewCustomerJsonPresenter() port.Presenter {
	return &customerJsonPresenter{}
}

// ToCustomerJsonResponse convert entity.Customer to CustomerJsonResponse
func ToCustomerJsonResponse(customer *entity.Customer) CustomerJsonResponse {
	return CustomerJsonResponse{
		ID:        customer.ID,
		Name:      customer.Name,
		Email:     customer.Email,
		CPF:       customer.CPF,
		CreatedAt: customer.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: customer.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Present write the response to the client
func (p *customerJsonPresenter) Present(pp dto.PresenterInput) ([]byte, error) {
	switch v := pp.Result.(type) {
	case *entity.Customer:
		output := ToCustomerJsonResponse(v)
		return json.Marshal(output)
	case []*entity.Customer:
		customerOutputs := make([]CustomerJsonResponse, len(v))
		for i, customer := range v {
			customerOutputs[i] = ToCustomerJsonResponse(customer)
		}

		output := &CustomerJsonPaginatedResponse{
			JsonPagination: JsonPagination{
				Total: pp.Total,
				Page:  pp.Page,
				Limit: pp.Limit,
			},
			Customers: customerOutputs,
		}

		return json.Marshal(output)
	default:
		return nil, domain.NewInternalError(errors.New(domain.ErrInternalError))
	}
}
