package dto

import "github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain/entity"

type CreateCustomerInput struct {
	Name  string
	Email string
	CPF   string
}

func (i CreateCustomerInput) ToEntity() *entity.Customer {
	return &entity.Customer{
		Name:  i.Name,
		Email: i.Email,
		CPF:   i.CPF,
	}
}

type UpdateCustomerInput struct {
	ID    string
	Name  string
	Email string
}

type GetCustomerInput struct {
	ID string
}

type GetCustomerByCPFInput struct {
	CPF string
}

type DeleteCustomerInput struct {
	ID string
}

type ListCustomersInput struct {
	Name  string
	Page  int
	Limit int
}

type FindCustomerByCPFInput struct {
	CPF string
}
