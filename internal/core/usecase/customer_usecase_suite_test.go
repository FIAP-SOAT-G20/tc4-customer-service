package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/port"
	mockport "github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/port/mocks"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/usecase"
)

type CustomerUsecaseSuiteTest struct {
	suite.Suite
	mockCustomers []*entity.Customer
	mockGateway   *mockport.MockCustomerGateway
	useCase       port.CustomerUseCase
	ctx           context.Context
}

func (s *CustomerUsecaseSuiteTest) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	s.mockGateway = mockport.NewMockCustomerGateway(ctrl)
	s.useCase = usecase.NewCustomerUseCase(s.mockGateway)
	s.ctx = context.Background()
	currentTime := time.Now()
	s.mockCustomers = []*entity.Customer{
		{
			ID:        "123",
			Name:      "Test Customer 1",
			Email:     "test.customer.1@email.com",
			CPF:       "12345678901",
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		},
		{
			ID:        "321",
			Name:      "Test Customer 2",
			Email:     "test.customer.2@email.com",
			CPF:       "12345678902",
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		},
	}
}

func TestCustomerUsecaseSuiteTest(t *testing.T) {
	suite.Run(t, new(CustomerUsecaseSuiteTest))
}
