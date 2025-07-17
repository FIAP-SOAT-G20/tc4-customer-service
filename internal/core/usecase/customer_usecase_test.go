package usecase

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/test/fixture"

	"github.com/stretchr/testify/assert"

	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/domain"
	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/dto"
	mockport "github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/port/mocks"
)

func TestCustomerUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	tests := []struct {
		name         string
		input        dto.GetCustomerInput
		mockSetup    func(g *mockport.MockCustomerGateway)
		wantCustomer *entity.Customer
		wantErrType  interface{}
	}{
		{
			name:  "success",
			input: dto.GetCustomerInput{CPF: fixture.SampleCPF},
			mockSetup: func(g *mockport.MockCustomerGateway) {
				g.EXPECT().
					FindOne(ctx, fixture.SampleCPF).
					Return(fixture.SampleCustomer, nil)
			},
			wantCustomer: fixture.SampleCustomer,
			wantErrType:  nil,
		},
		{
			name:  "empty CPF uses default",
			input: dto.GetCustomerInput{CPF: ""},
			mockSetup: func(g *mockport.MockCustomerGateway) {
				g.EXPECT().
					FindOne(ctx, "000.000.000-00").
					Return(nil, nil)
			},
			wantCustomer: nil,
			wantErrType:  &domain.InvalidInputError{},
		},
		{
			name:  "not found",
			input: dto.GetCustomerInput{CPF: fixture.SampleCPF},
			mockSetup: func(g *mockport.MockCustomerGateway) {
				g.EXPECT().
					FindOne(ctx, fixture.SampleCPF).
					Return(nil, nil)
			},
			wantCustomer: nil,
			wantErrType:  &domain.NotFoundError{},
		},
		{
			name:  "internal error",
			input: dto.GetCustomerInput{CPF: fixture.SampleCPF},
			mockSetup: func(g *mockport.MockCustomerGateway) {
				g.EXPECT().
					FindOne(ctx, fixture.SampleCPF).
					Return(nil, errors.New("db error"))
			},
			wantCustomer: nil,
			wantErrType:  &domain.InternalError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGateway := mockport.NewMockCustomerGateway(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(mockGateway)
			}

			uc := NewCustomerUseCase(mockGateway)
			customer, err := uc.Get(ctx, tt.input)

			assert.Equal(t, tt.wantCustomer, customer)

			if tt.wantErrType != nil {
				assert.ErrorAs(t, err, &tt.wantErrType)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
