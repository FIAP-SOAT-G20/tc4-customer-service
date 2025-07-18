package controller

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/port"
)

type customerController struct {
	useCase port.CustomerUseCase
}

func NewCustomerController(useCase port.CustomerUseCase) port.CustomerController {
	return &customerController{useCase}
}

func (c *customerController) Get(ctx context.Context, p port.Presenter, i dto.GetCustomerInput) ([]byte, error) {
	customer, err := c.useCase.Get(ctx, i)
	if err != nil {
		return nil, err
	}

	return p.Present(dto.PresenterInput{
		Result: customer,
	})
}
