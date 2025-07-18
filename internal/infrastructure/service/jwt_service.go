package service

import (
	"time"

	"github.com/google/uuid"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/config"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	secretKey  []byte
	issuer     string
	audience   []string
	expiration time.Duration
}

func NewJWTService(cfg *config.Config) *JwtService {
	return &JwtService{
		secretKey:  []byte(cfg.JWTSecret),
		expiration: cfg.JWTExpiration,
		issuer:     cfg.JWTIssuer,
		audience:   []string{cfg.JWTAudience},
	}
}

func (s *JwtService) GenerateToken(userIdentifier string) (string, string, int64, error) {
	expiresAt := time.Now().Add(s.expiration)
	jwtTokenId := uuid.New().String()

	tokenClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    s.issuer,
		Subject:   userIdentifier,
		ID:        jwtTokenId,
		Audience:  s.audience,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	signedToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", "", 0, err
	}

	return "Bearer", signedToken, expiresAt.UnixMilli(), nil
}
