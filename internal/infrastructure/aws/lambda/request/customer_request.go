package request

import "github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/dto"

type CustomerRequest struct {
	CPF string `json:"cpf"`
}

func (c CustomerRequest) ToGetCustomerInput() dto.GetCustomerInput {
	return dto.GetCustomerInput{
		CPF: c.CPF,
	}
}
