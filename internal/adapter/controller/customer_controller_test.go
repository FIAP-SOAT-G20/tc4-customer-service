package controller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/adapter/controller"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/dto"
	mockport "github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/port/mocks"
)

func TestCustomerController_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockport.NewMockCustomerUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	customerController := controller.NewCustomerController(mockUseCase)

	ctx := context.Background()
	input := dto.ListCustomersInput{
		Page:  1,
		Limit: 10,
	}

	mockCustomers := []*entity.Customer{
		{ID: 1, Name: "Customer 1", Email: "customer1@test.com", CPF: "123.456.789-01"},
		{ID: 2, Name: "Customer 2", Email: "customer2@test.com", CPF: "123.456.789-02"},
	}

	tests := []struct {
		name        string
		setupMocks  func()
		checkResult func(*testing.T, []byte, error)
	}{
		{
			name: "should list customers successfully",
			setupMocks: func() {
				mockUseCase.EXPECT().
					List(ctx, input).
					Return(mockCustomers, int64(2), nil)

				mockPresenter.EXPECT().
					Present(dto.PresenterInput{
						Total:  int64(2),
						Page:   1,
						Limit:  10,
						Result: mockCustomers,
					}).
					Return([]byte(`{"customers":[{"id":"1","name":"Customer 1"}]}`), nil)
			},
			checkResult: func(t *testing.T, result []byte, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Contains(t, string(result), "customers")
			},
		},
		{
			name: "should return error when use case fails",
			setupMocks: func() {
				mockUseCase.EXPECT().
					List(ctx, input).
					Return(nil, int64(0), errors.New("use case error"))
			},
			checkResult: func(t *testing.T, result []byte, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, "use case error", err.Error())
			},
		},
		{
			name: "should return error when presenter fails",
			setupMocks: func() {
				mockUseCase.EXPECT().
					List(ctx, input).
					Return(mockCustomers, int64(2), nil)

				mockPresenter.EXPECT().
					Present(dto.PresenterInput{
						Total:  int64(2),
						Page:   1,
						Limit:  10,
						Result: mockCustomers,
					}).
					Return(nil, errors.New("presenter error"))
			},
			checkResult: func(t *testing.T, result []byte, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, "presenter error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			result, err := customerController.List(ctx, mockPresenter, input)

			tt.checkResult(t, result, err)
		})
	}
}

func TestCustomerController_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockport.NewMockCustomerUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	customerController := controller.NewCustomerController(mockUseCase)

	ctx := context.Background()
	input := dto.CreateCustomerInput{
		Name:  "New Customer",
		Email: "new@test.com",
		CPF:   "123.456.789-00",
	}

	mockCustomer := &entity.Customer{
		ID:    1,
		Name:  "New Customer",
		Email: "new@test.com",
		CPF:   "123.456.789-00",
	}

	tests := []struct {
		name        string
		setupMocks  func()
		checkResult func(*testing.T, []byte, error)
	}{
		{
			name: "should create customer successfully",
			setupMocks: func() {
				mockUseCase.EXPECT().
					Create(ctx, input).
					Return(mockCustomer, nil)

				mockPresenter.EXPECT().
					Present(dto.PresenterInput{
						Result: mockCustomer,
					}).
					Return([]byte(`{"id":"1","name":"New Customer"}`), nil)
			},
			checkResult: func(t *testing.T, result []byte, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Contains(t, string(result), "New Customer")
			},
		},
		{
			name: "should return error when use case fails",
			setupMocks: func() {
				mockUseCase.EXPECT().
					Create(ctx, input).
					Return(nil, errors.New("use case error"))
			},
			checkResult: func(t *testing.T, result []byte, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, "use case error", err.Error())
			},
		},
		{
			name: "should return error when presenter fails",
			setupMocks: func() {
				mockUseCase.EXPECT().
					Create(ctx, input).
					Return(mockCustomer, nil)

				mockPresenter.EXPECT().
					Present(dto.PresenterInput{
						Result: mockCustomer,
					}).
					Return(nil, errors.New("presenter error"))
			},
			checkResult: func(t *testing.T, result []byte, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, "presenter error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			result, err := customerController.Create(ctx, mockPresenter, input)

			tt.checkResult(t, result, err)
		})
	}
}

func TestCustomerController_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockport.NewMockCustomerUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	customerController := controller.NewCustomerController(mockUseCase)

	ctx := context.Background()
	input := dto.GetCustomerInput{ID: 123}

	mockCustomer := &entity.Customer{
		ID:    123,
		Name:  "Test Customer",
		Email: "test@test.com",
		CPF:   "123.456.789-01",
	}

	tests := []struct {
		name        string
		setupMocks  func()
		checkResult func(*testing.T, []byte, error)
	}{
		{
			name: "should get customer successfully",
			setupMocks: func() {
				mockUseCase.EXPECT().
					Get(ctx, input).
					Return(mockCustomer, nil)

				mockPresenter.EXPECT().
					Present(dto.PresenterInput{
						Result: mockCustomer,
					}).
					Return([]byte(`{"id":"123","name":"Test Customer"}`), nil)
			},
			checkResult: func(t *testing.T, result []byte, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Contains(t, string(result), "Test Customer")
			},
		},
		{
			name: "should return error when use case fails",
			setupMocks: func() {
				mockUseCase.EXPECT().
					Get(ctx, input).
					Return(nil, errors.New("customer not found"))
			},
			checkResult: func(t *testing.T, result []byte, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, "customer not found", err.Error())
			},
		},
		{
			name: "should return error when presenter fails",
			setupMocks: func() {
				mockUseCase.EXPECT().
					Get(ctx, input).
					Return(mockCustomer, nil)

				mockPresenter.EXPECT().
					Present(dto.PresenterInput{
						Result: mockCustomer,
					}).
					Return(nil, errors.New("presenter error"))
			},
			checkResult: func(t *testing.T, result []byte, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, "presenter error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			result, err := customerController.Get(ctx, mockPresenter, input)

			tt.checkResult(t, result, err)
		})
	}
}

func TestCustomerController_GetByCPF(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockport.NewMockCustomerUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	customerController := controller.NewCustomerController(mockUseCase)

	ctx := context.Background()
	input := dto.GetCustomerByCPFInput{CPF: "123.456.789-01"}

	mockCustomer := &entity.Customer{
		ID:    123,
		Name:  "Test Customer",
		Email: "test@test.com",
		CPF:   "123.456.789-01",
	}

	tests := []struct {
		name        string
		setupMocks  func()
		checkResult func(*testing.T, []byte, error)
	}{
		{
			name: "should get customer by CPF successfully",
			setupMocks: func() {
				mockUseCase.EXPECT().
					GetByCPF(ctx, input).
					Return(mockCustomer, nil)

				mockPresenter.EXPECT().
					Present(dto.PresenterInput{
						Result: mockCustomer,
					}).
					Return([]byte(`{"id":"123","cpf":"123.456.789-01"}`), nil)
			},
			checkResult: func(t *testing.T, result []byte, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Contains(t, string(result), "123.456.789-01")
			},
		},
		{
			name: "should return error when use case fails",
			setupMocks: func() {
				mockUseCase.EXPECT().
					GetByCPF(ctx, input).
					Return(nil, errors.New("customer not found"))
			},
			checkResult: func(t *testing.T, result []byte, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, "customer not found", err.Error())
			},
		},
		{
			name: "should return error when presenter fails",
			setupMocks: func() {
				mockUseCase.EXPECT().
					GetByCPF(ctx, input).
					Return(mockCustomer, nil)

				mockPresenter.EXPECT().
					Present(dto.PresenterInput{
						Result: mockCustomer,
					}).
					Return(nil, errors.New("presenter error"))
			},
			checkResult: func(t *testing.T, result []byte, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, "presenter error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			result, err := customerController.GetByCPF(ctx, mockPresenter, input)

			tt.checkResult(t, result, err)
		})
	}
}

func TestCustomerController_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockport.NewMockCustomerUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	customerController := controller.NewCustomerController(mockUseCase)

	ctx := context.Background()
	input := dto.UpdateCustomerInput{
		ID:    123,
		Name:  "Updated Customer",
		Email: "updated@test.com",
	}

	mockCustomer := &entity.Customer{
		ID:    123,
		Name:  "Updated Customer",
		Email: "updated@test.com",
		CPF:   "123.456.789-01",
	}

	tests := []struct {
		name        string
		setupMocks  func()
		checkResult func(*testing.T, []byte, error)
	}{
		{
			name: "should update customer successfully",
			setupMocks: func() {
				mockUseCase.EXPECT().
					Update(ctx, input).
					Return(mockCustomer, nil)

				mockPresenter.EXPECT().
					Present(dto.PresenterInput{
						Result: mockCustomer,
					}).
					Return([]byte(`{"id":"123","name":"Updated Customer"}`), nil)
			},
			checkResult: func(t *testing.T, result []byte, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Contains(t, string(result), "Updated Customer")
			},
		},
		{
			name: "should return error when use case fails",
			setupMocks: func() {
				mockUseCase.EXPECT().
					Update(ctx, input).
					Return(nil, errors.New("customer not found"))
			},
			checkResult: func(t *testing.T, result []byte, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, "customer not found", err.Error())
			},
		},
		{
			name: "should return error when presenter fails",
			setupMocks: func() {
				mockUseCase.EXPECT().
					Update(ctx, input).
					Return(mockCustomer, nil)

				mockPresenter.EXPECT().
					Present(dto.PresenterInput{
						Result: mockCustomer,
					}).
					Return(nil, errors.New("presenter error"))
			},
			checkResult: func(t *testing.T, result []byte, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, "presenter error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			result, err := customerController.Update(ctx, mockPresenter, input)

			tt.checkResult(t, result, err)
		})
	}
}

func TestCustomerController_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockport.NewMockCustomerUseCase(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)
	customerController := controller.NewCustomerController(mockUseCase)

	ctx := context.Background()
	input := dto.DeleteCustomerInput{ID: 123}

	mockCustomer := &entity.Customer{
		ID:    123,
		Name:  "Deleted Customer",
		Email: "deleted@test.com",
		CPF:   "123.456.789-01",
	}

	tests := []struct {
		name        string
		setupMocks  func()
		checkResult func(*testing.T, []byte, error)
	}{
		{
			name: "should delete customer successfully",
			setupMocks: func() {
				mockUseCase.EXPECT().
					Delete(ctx, input).
					Return(mockCustomer, nil)

				mockPresenter.EXPECT().
					Present(dto.PresenterInput{
						Result: mockCustomer,
					}).
					Return([]byte(`{"id":"123","name":"Deleted Customer"}`), nil)
			},
			checkResult: func(t *testing.T, result []byte, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Contains(t, string(result), "Deleted Customer")
			},
		},
		{
			name: "should return error when use case fails",
			setupMocks: func() {
				mockUseCase.EXPECT().
					Delete(ctx, input).
					Return(nil, errors.New("customer not found"))
			},
			checkResult: func(t *testing.T, result []byte, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, "customer not found", err.Error())
			},
		},
		{
			name: "should return error when presenter fails",
			setupMocks: func() {
				mockUseCase.EXPECT().
					Delete(ctx, input).
					Return(mockCustomer, nil)

				mockPresenter.EXPECT().
					Present(dto.PresenterInput{
						Result: mockCustomer,
					}).
					Return(nil, errors.New("presenter error"))
			},
			checkResult: func(t *testing.T, result []byte, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, "presenter error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			result, err := customerController.Delete(ctx, mockPresenter, input)

			tt.checkResult(t, result, err)
		})
	}
}
