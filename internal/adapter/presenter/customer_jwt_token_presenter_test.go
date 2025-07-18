package presenter

import (
	"encoding/json"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	mock_port "github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/port/mocks"

	"github.com/stretchr/testify/require"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/dto"
)

func TestCustomerJwtTokenPresenter_Present_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jwtService := mock_port.NewMockIAuthenticationService(ctrl)
	presenter := NewCustomerJwtTokenPresenter(jwtService)

	customer := &entity.Customer{ID: "999"}
	jwtService.EXPECT().
		GenerateToken("999").
		Return("atoken", "Bearer", int64(12345), nil)

	input := dto.PresenterInput{Result: customer}
	data, err := presenter.Present(input)
	require.NoError(t, err)

	// Convert data to expected output struct
	var resp JWTResponse
	err = json.Unmarshal(data, &resp)
	require.NoError(t, err)
	require.Equal(t, "atoken", resp.AccessToken)
	require.Equal(t, "Bearer", resp.TokenType)
	require.Equal(t, int64(12345), resp.ExpiresIn)
}

func TestCustomerJwtTokenPresenter_Present_GenerateTokenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jwtService := mock_port.NewMockIAuthenticationService(ctrl)
	presenter := NewCustomerJwtTokenPresenter(jwtService)

	customer := &entity.Customer{ID: "2"}
	jwtService.EXPECT().
		GenerateToken("2").
		Return("", "", int64(0), errors.New("fail"))

	input := dto.PresenterInput{Result: customer}
	data, err := presenter.Present(input)
	require.Nil(t, data)
	require.Error(t, err)
}

func TestCustomerJwtTokenPresenter_Present_InvalidType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jwtService := mock_port.NewMockIAuthenticationService(ctrl)
	presenter := NewCustomerJwtTokenPresenter(jwtService)

	input := dto.PresenterInput{Result: 42} // not a *entity.Customer
	data, err := presenter.Present(input)
	require.Nil(t, data)
	require.Error(t, err)
}
