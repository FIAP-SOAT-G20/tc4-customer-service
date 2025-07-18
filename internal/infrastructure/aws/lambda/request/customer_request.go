package request

import "github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/dto"

type CustomerRequest struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	CPF   string `json:"cpf,omitempty"`
}

func (c CustomerRequest) ToGetCustomerInput() dto.GetCustomerInput {
	return dto.GetCustomerInput{
		ID: c.ID,
	}
}

func (c CustomerRequest) ToCreateCustomerInput() dto.CreateCustomerInput {
	return dto.CreateCustomerInput{
		Name:  c.Name,
		Email: c.Email,
		CPF:   c.CPF,
	}
}

func (c CustomerRequest) ToUpdateCustomerInput() dto.UpdateCustomerInput {
	return dto.UpdateCustomerInput{
		ID:    c.ID,
		Name:  c.Name,
		Email: c.Email,
	}
}
