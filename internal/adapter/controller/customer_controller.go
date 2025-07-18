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

func (c *customerController) List(ctx context.Context, presenter port.Presenter, input dto.ListCustomersInput) ([]byte, error) {
	customers, total, err := c.useCase.List(ctx, input)
	if err != nil {
		return nil, err
	}

	return presenter.Present(dto.PresenterInput{
		Total:  total,
		Page:   input.Page,
		Limit:  input.Limit,
		Result: customers,
	})
}

func (c *customerController) Create(ctx context.Context, presenter port.Presenter, input dto.CreateCustomerInput) ([]byte, error) {
	customer, err := c.useCase.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return presenter.Present(dto.PresenterInput{
		Result: customer,
	})
}

func (c *customerController) Get(ctx context.Context, presenter port.Presenter, input dto.GetCustomerInput) ([]byte, error) {
	customer, err := c.useCase.Get(ctx, input)
	if err != nil {
		return nil, err
	}

	return presenter.Present(dto.PresenterInput{
		Result: customer,
	})
}

func (c *customerController) Update(ctx context.Context, presenter port.Presenter, input dto.UpdateCustomerInput) ([]byte, error) {
	customer, err := c.useCase.Update(ctx, input)
	if err != nil {
		return nil, err
	}

	return presenter.Present(dto.PresenterInput{
		Result: customer,
	})
}

func (c *customerController) GetByCPF(ctx context.Context, presenter port.Presenter, input dto.GetCustomerByCPFInput) ([]byte, error) {
	customer, err := c.useCase.GetByCPF(ctx, input)
	if err != nil {
		return nil, err
	}

	return presenter.Present(dto.PresenterInput{
		Result: customer,
	})
}

func (c *customerController) Delete(ctx context.Context, presenter port.Presenter, input dto.DeleteCustomerInput) ([]byte, error) {
	customer, err := c.useCase.Delete(ctx, input)
	if err != nil {
		return nil, err
	}

	return presenter.Present(dto.PresenterInput{
		Result: customer,
	})
}
