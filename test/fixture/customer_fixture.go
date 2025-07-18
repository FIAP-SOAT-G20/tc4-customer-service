package fixture

import (
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain/entity"
)

var (
	SampleCPF      = "12345678900"
	SampleCustomer = &entity.Customer{
		CPF:  SampleCPF,
		Name: "Test User",
	}
)
