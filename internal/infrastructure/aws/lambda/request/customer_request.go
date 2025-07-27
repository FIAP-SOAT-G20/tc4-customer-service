package request

import (
	"strconv"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/dto"
)

type CustomerRequest struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	CPF   string `json:"cpf,omitempty"`
}

func (c CustomerRequest) ToGetCustomerInput() dto.GetCustomerInput {
	id, _ := strconv.Atoi(c.ID)
	return dto.GetCustomerInput{
		ID: id,
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
	id, _ := strconv.Atoi(c.ID)
	return dto.UpdateCustomerInput{
		ID:    id,
		Name:  c.Name,
		Email: c.Email,
	}
}
