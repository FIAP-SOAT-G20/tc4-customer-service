package presenter

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/port"
)

type customerJwtTokenPresenter struct {
	jwtService port.IAuthenticationService
}

// CustomerJsonResponse represents the response of a customer
func NewCustomerJwtTokenPresenter(jwtService port.IAuthenticationService) port.Presenter {
	return &customerJwtTokenPresenter{jwtService: jwtService}
}

// ToCustomerJsonResponse convert entity.Customer to CustomerJsonResponse
func toCustomerJsonResponse(accessToken, tokenType string, expiresIn int64) JWTResponse {
	return JWTResponse{
		AccessToken: accessToken,
		TokenType:   tokenType,
		ExpiresIn:   expiresIn,
	}
}

// Present write the response to the client
func (p *customerJwtTokenPresenter) Present(pp dto.PresenterInput) ([]byte, error) {
	switch v := pp.Result.(type) {
	case *entity.Customer:
		accessToken, tokenType, expiresIn, err := p.jwtService.GenerateToken(strconv.Itoa(v.ID))
		if err != nil {
			return nil, domain.NewInternalError(errors.New(domain.ErrInternalError))
		}
		output := toCustomerJsonResponse(accessToken, tokenType, expiresIn)
		return json.Marshal(output)
	default:
		return nil, domain.NewInternalError(errors.New(domain.ErrInternalError))
	}
}
