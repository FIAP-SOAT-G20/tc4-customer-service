package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/dto"
	mockport "github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/port/mocks"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/usecase"
)

// Helper function to create mock customers for tests
func createMockCustomers() []*entity.Customer {
	currentTime := time.Now()
	return []*entity.Customer{
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

func TestCustomersUseCase_List(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGateway := mockport.NewMockCustomerGateway(ctrl)
	useCase := usecase.NewCustomerUseCase(mockGateway)
	ctx := context.Background()
	mockCustomers := createMockCustomers()

	tests := []struct {
		name        string
		input       dto.ListCustomersInput
		setupMocks  func()
		checkResult func(*testing.T, []*entity.Customer, int64, error)
	}{
		{
			name: "should list products successfully",
			input: dto.ListCustomersInput{
				Page:  1,
				Limit: 10,
			},
			setupMocks: func() {
				mockGateway.EXPECT().
					FindAll(ctx, "", 1, 10).
					Return(mockCustomers, int64(2), nil)
			},
			checkResult: func(t *testing.T, customers []*entity.Customer, total int64, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, customers)
				assert.Equal(t, len(mockCustomers), len(customers))
				assert.Equal(t, int64(2), total)
			},
		},
		{
			name: "should return error when repository fails",
			input: dto.ListCustomersInput{
				Page:  1,
				Limit: 10,
			},
			setupMocks: func() {
				mockGateway.EXPECT().
					FindAll(ctx, "", 1, 10).
					Return(nil, int64(0), assert.AnError)
			},
			checkResult: func(t *testing.T, customers []*entity.Customer, total int64, err error) {
				assert.Error(t, err)
				assert.Nil(t, customers)
				assert.Equal(t, int64(0), total)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
		{
			name: "should filter by name",
			input: dto.ListCustomersInput{
				Name:  "Test",
				Page:  1,
				Limit: 10,
			},
			setupMocks: func() {
				mockGateway.EXPECT().
					FindAll(ctx, "Test", 1, 10).
					Return(mockCustomers, int64(2), nil)
			},
			checkResult: func(t *testing.T, customers []*entity.Customer, total int64, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, customers)
				assert.Equal(t, len(mockCustomers), len(customers))
				assert.Equal(t, int64(2), total)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			customers, total, err := useCase.List(ctx, tt.input)

			// Assert
			tt.checkResult(t, customers, total, err)
		})
	}
}

func TestCustomerUseCase_Create(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGateway := mockport.NewMockCustomerGateway(ctrl)
	useCase := usecase.NewCustomerUseCase(mockGateway)
	ctx := context.Background()
	mockCustomers := createMockCustomers()

	tests := []struct {
		name        string
		input       dto.CreateCustomerInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Customer, error)
	}{
		{
			name: "should create customer successfully",
			input: dto.CreateCustomerInput{
				Name:  mockCustomers[0].Name,
				Email: mockCustomers[0].Email,
				CPF:   mockCustomers[0].CPF,
			},
			setupMocks: func() {
				mockGateway.EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil)
			},
			checkResult: func(t *testing.T, customer *entity.Customer, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, customer)
				assert.Equal(t, mockCustomers[0].Name, customer.Name)
				assert.Equal(t, mockCustomers[0].Email, customer.Email)
				assert.Equal(t, mockCustomers[0].CPF, customer.CPF)
			},
		},
		{
			name: "should return error when gateway fails",
			input: dto.CreateCustomerInput{
				Name:  mockCustomers[0].Name,
				Email: mockCustomers[0].Email,
				CPF:   mockCustomers[0].CPF,
			},
			setupMocks: func() {
				mockGateway.EXPECT().
					Create(ctx, gomock.Any()).
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, customer *entity.Customer, err error) {
				assert.Error(t, err)
				assert.Nil(t, customer)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			customer, err := useCase.Create(ctx, tt.input)

			// Assert
			tt.checkResult(t, customer, err)
		})
	}
}

func TestCustomerUseCase_Get(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGateway := mockport.NewMockCustomerGateway(ctrl)
	useCase := usecase.NewCustomerUseCase(mockGateway)
	ctx := context.Background()
	mockCustomers := createMockCustomers()

	tests := []struct {
		name        string
		input       dto.GetCustomerInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Customer, error)
	}{
		{
			name:  "should get customer successfully",
			input: dto.GetCustomerInput{ID: "123"},
			setupMocks: func() {
				mockGateway.EXPECT().
					FindByID(ctx, "123").
					Return(mockCustomers[0], nil)
			},
			checkResult: func(t *testing.T, customer *entity.Customer, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, customer)
				assert.Equal(t, mockCustomers[0].ID, customer.ID)
				assert.Equal(t, mockCustomers[0].Name, customer.Name)
				assert.Equal(t, mockCustomers[0].Email, customer.Email)
				assert.Equal(t, mockCustomers[0].CreatedAt, customer.CreatedAt)
				assert.Equal(t, mockCustomers[0].UpdatedAt, customer.UpdatedAt)
			},
		},
		{
			name:  "should return not found error when customer doesn't exist",
			input: dto.GetCustomerInput{ID: "123"},
			setupMocks: func() {
				mockGateway.EXPECT().
					FindByID(ctx, "123").
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, customer *entity.Customer, err error) {
				assert.Error(t, err)
				assert.Nil(t, customer)
				assert.IsType(t, &domain.NotFoundError{}, err)
			},
		},
		{
			name:  "should return internal error when gateway fails",
			input: dto.GetCustomerInput{ID: "123"},
			setupMocks: func() {
				mockGateway.EXPECT().
					FindByID(ctx, "123").
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, customer *entity.Customer, err error) {
				assert.Error(t, err)
				assert.Nil(t, customer)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			customer, err := useCase.Get(ctx, tt.input)

			// Assert
			tt.checkResult(t, customer, err)
		})
	}
}

func TestCustomerUseCase_Update(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGateway := mockport.NewMockCustomerGateway(ctrl)
	useCase := usecase.NewCustomerUseCase(mockGateway)
	ctx := context.Background()
	mockCustomers := createMockCustomers()

	tests := []struct {
		name        string
		input       dto.UpdateCustomerInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Customer, error)
	}{
		{
			name: "should update customer successfully",
			input: dto.UpdateCustomerInput{
				ID:    "123",
				Name:  "New Name",
				Email: "new.name@email.com",
			},
			setupMocks: func() {
				mockGateway.EXPECT().
					FindByID(ctx, "123").
					Return(mockCustomers[0], nil)

				mockGateway.EXPECT().
					Update(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, p *entity.Customer) error {
						assert.Equal(t, "New Name", p.Name)
						assert.Equal(t, "new.name@email.com", p.Email)
						return nil
					})
			},
			checkResult: func(t *testing.T, customer *entity.Customer, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, customer)
				assert.Equal(t, "New Name", customer.Name)
				assert.Equal(t, "new.name@email.com", customer.Email)
				assert.Equal(t, mockCustomers[0].CreatedAt, customer.CreatedAt)
			},
		},
		{
			name: "should return error when customer not found",
			input: dto.UpdateCustomerInput{
				ID:    "123",
				Name:  "New Name",
				Email: "new.name@email.com",
			},
			setupMocks: func() {
				mockGateway.EXPECT().
					FindByID(ctx, "123").
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, customer *entity.Customer, err error) {
				assert.Error(t, err)
				assert.Nil(t, customer)
				assert.IsType(t, &domain.NotFoundError{}, err)
			},
		},
		{
			name: "should return error when gateway find fails",
			input: dto.UpdateCustomerInput{
				ID:    "123",
				Name:  "New Name",
				Email: "new.name@email.com",
			},
			setupMocks: func() {
				mockGateway.EXPECT().
					FindByID(ctx, "123").
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, customer *entity.Customer, err error) {
				assert.Error(t, err)
				assert.Nil(t, customer)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
		{
			name: "should return error when gateway update fails",
			input: dto.UpdateCustomerInput{
				ID:    "123",
				Name:  "New Name",
				Email: "new.name@email.com",
			},
			setupMocks: func() {
				mockGateway.EXPECT().
					FindByID(ctx, "123").
					Return(mockCustomers[0], nil)

				mockGateway.EXPECT().
					Update(ctx, gomock.Any()).
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, customer *entity.Customer, err error) {
				assert.Error(t, err)
				assert.Nil(t, customer)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			customer, err := useCase.Update(ctx, tt.input)

			// Assert
			tt.checkResult(t, customer, err)
		})
	}
}

func TestCustomerUseCase_Delete(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGateway := mockport.NewMockCustomerGateway(ctrl)
	useCase := usecase.NewCustomerUseCase(mockGateway)
	ctx := context.Background()

	tests := []struct {
		name        string
		input       dto.DeleteCustomerInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Customer, error)
	}{
		{
			name:  "should delete customer successfully",
			input: dto.DeleteCustomerInput{ID: "123"},
			setupMocks: func() {
				mockGateway.EXPECT().
					FindByID(ctx, "123").
					Return(&entity.Customer{ID: "123"}, nil)

				mockGateway.EXPECT().
					Delete(ctx, "123").
					Return(nil)
			},
			checkResult: func(t *testing.T, customer *entity.Customer, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, customer)
				assert.Equal(t, "123", customer.ID)
			},
		},
		{
			name:  "should return not found error when customer doesn't exist",
			input: dto.DeleteCustomerInput{ID: "123"},
			setupMocks: func() {
				mockGateway.EXPECT().
					FindByID(ctx, "123").
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, customer *entity.Customer, err error) {
				assert.Error(t, err)
				assert.Nil(t, customer)
				assert.IsType(t, &domain.NotFoundError{}, err)
			},
		},
		{
			name:  "should return error when gateway fails on find",
			input: dto.DeleteCustomerInput{ID: "123"},
			setupMocks: func() {
				mockGateway.EXPECT().
					FindByID(ctx, "123").
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, customer *entity.Customer, err error) {
				assert.Error(t, err)
				assert.Nil(t, customer)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
		{
			name:  "should return error when gateway fails on delete",
			input: dto.DeleteCustomerInput{ID: "123"},
			setupMocks: func() {
				mockGateway.EXPECT().
					FindByID(ctx, "123").
					Return(&entity.Customer{}, nil)

				mockGateway.EXPECT().
					Delete(ctx, "123").
					Return(assert.AnError)
			},
			checkResult: func(t *testing.T, customer *entity.Customer, err error) {
				assert.Error(t, err)
				assert.Nil(t, customer)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			customer, err := useCase.Delete(ctx, tt.input)

			// Assert
			tt.checkResult(t, customer, err)
		})
	}
}

func TestCustomerUseCase_GetByCPF(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGateway := mockport.NewMockCustomerGateway(ctrl)
	useCase := usecase.NewCustomerUseCase(mockGateway)
	ctx := context.Background()
	mockCustomers := createMockCustomers()

	tests := []struct {
		name        string
		input       dto.GetCustomerByCPFInput
		setupMocks  func()
		checkResult func(*testing.T, *entity.Customer, error)
	}{
		{
			name:  "should get customer by CPF successfully",
			input: dto.GetCustomerByCPFInput{CPF: "12345678901"},
			setupMocks: func() {
				mockGateway.EXPECT().
					FindByCPF(ctx, "12345678901").
					Return(mockCustomers[0], nil)
			},
			checkResult: func(t *testing.T, customer *entity.Customer, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, customer)
				assert.Equal(t, mockCustomers[0].ID, customer.ID)
				assert.Equal(t, mockCustomers[0].Name, customer.Name)
				assert.Equal(t, mockCustomers[0].Email, customer.Email)
				assert.Equal(t, mockCustomers[0].CPF, customer.CPF)
				assert.Equal(t, mockCustomers[0].CreatedAt, customer.CreatedAt)
				assert.Equal(t, mockCustomers[0].UpdatedAt, customer.UpdatedAt)
			},
		},
		{
			name:  "should return not found error when customer doesn't exist",
			input: dto.GetCustomerByCPFInput{CPF: "99999999999"},
			setupMocks: func() {
				mockGateway.EXPECT().
					FindByCPF(ctx, "99999999999").
					Return(nil, nil)
			},
			checkResult: func(t *testing.T, customer *entity.Customer, err error) {
				assert.Error(t, err)
				assert.Nil(t, customer)
				assert.IsType(t, &domain.NotFoundError{}, err)
			},
		},
		{
			name:  "should return internal error when gateway fails",
			input: dto.GetCustomerByCPFInput{CPF: "12345678901"},
			setupMocks: func() {
				mockGateway.EXPECT().
					FindByCPF(ctx, "12345678901").
					Return(nil, assert.AnError)
			},
			checkResult: func(t *testing.T, customer *entity.Customer, err error) {
				assert.Error(t, err)
				assert.Nil(t, customer)
				assert.IsType(t, &domain.InternalError{}, err)
			},
		},
		{
			name:  "should use default CPF when CPF is empty",
			input: dto.GetCustomerByCPFInput{CPF: ""},
			setupMocks: func() {
				mockGateway.EXPECT().
					FindByCPF(ctx, "000.000.000-00").
					Return(mockCustomers[0], nil)
			},
			checkResult: func(t *testing.T, customer *entity.Customer, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, customer)
				assert.Equal(t, mockCustomers[0].ID, customer.ID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMocks()

			// Act
			customer, err := useCase.GetByCPF(ctx, tt.input)

			// Assert
			tt.checkResult(t, customer, err)
		})
	}
}
