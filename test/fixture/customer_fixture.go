package fixture

import (
	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/domain/entity"
)

var (
	SampleCPF      = "12345678900"
	SampleCustomer = &entity.Customer{
		CPF:  SampleCPF,
		Name: "Test User",
	}
)
