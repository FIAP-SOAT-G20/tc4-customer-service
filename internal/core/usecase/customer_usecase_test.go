package usecase

import (
	"context"
	"errors"
	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/test/fixture"
	"go.uber.org/mock/gomock"
	"testing"

	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/domain"
	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/dto"
	mock_port "github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/port/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCustomerUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	tests := []struct {
		name         string
		input        dto.GetCustomerInput
		mockSetup    func(g *mock_port.MockCustomerGateway)
		wantCustomer *entity.Customer
		wantErrType  interface{}
	}{
		{
			name:  "success",
			input: dto.GetCustomerInput{CPF: fixture.SampleCPF},
			mockSetup: func(g *mock_port.MockCustomerGateway) {
				g.EXPECT().
					FindOne(ctx, fixture.SampleCPF).
					Return(fixture.SampleCustomer, nil)
			},
			wantCustomer: fixture.SampleCustomer,
			wantErrType:  nil,
		},
		{
			name:  "invalid input (empty CPF)",
			input: dto.GetCustomerInput{CPF: ""},
			mockSetup: func(g *mock_port.MockCustomerGateway) {
				// no call expected
			},
			wantCustomer: nil,
			wantErrType:  &domain.InvalidInputError{},
		},
		{
			name:  "not found",
			input: dto.GetCustomerInput{CPF: fixture.SampleCPF},
			mockSetup: func(g *mock_port.MockCustomerGateway) {
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
			mockSetup: func(g *mock_port.MockCustomerGateway) {
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
			mockGateway := mock_port.NewMockCustomerGateway(ctrl)
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
