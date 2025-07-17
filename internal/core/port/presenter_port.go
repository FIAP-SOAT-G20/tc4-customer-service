package port

import "github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/dto"

type Presenter interface {
	Present(dto.PresenterInput) ([]byte, error)
}
