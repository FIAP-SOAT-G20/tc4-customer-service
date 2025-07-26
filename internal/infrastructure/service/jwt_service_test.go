package service

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/config"
)

func TestNewJWTService(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:     "test-secret",
		JWTExpiration: time.Hour,
		JWTIssuer:     "test-issuer",
		JWTAudience:   "test-audience",
	}

	service := NewJWTService(cfg)

	assert.NotNil(t, service)
	assert.Equal(t, []byte(cfg.JWTSecret), service.secretKey)
	assert.Equal(t, cfg.JWTExpiration, service.expiration)
	assert.Equal(t, cfg.JWTIssuer, service.issuer)
	assert.Equal(t, []string{cfg.JWTAudience}, service.audience)
}

func TestJwtService_GenerateToken(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:     "test-secret-key",
		JWTExpiration: time.Hour,
		JWTIssuer:     "test-issuer",
		JWTAudience:   "test-audience",
	}

	service := NewJWTService(cfg)

	tests := []struct {
		name           string
		userIdentifier string
		wantErr        bool
	}{
		{
			name:           "should generate token successfully",
			userIdentifier: "user123",
			wantErr:        false,
		},
		{
			name:           "should generate token with empty identifier",
			userIdentifier: "",
			wantErr:        false,
		},
		{
			name:           "should generate token with special characters",
			userIdentifier: "user@domain.com",
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenType, token, expiresAt, err := service.GenerateToken(tt.userIdentifier)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
				assert.Empty(t, tokenType)
				assert.Equal(t, int64(0), expiresAt)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "Bearer", tokenType)
				assert.NotEmpty(t, token)
				assert.Greater(t, expiresAt, time.Now().UnixMilli())

				// Verify the token can be parsed
				parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					return service.secretKey, nil
				})
				assert.NoError(t, err)
				assert.True(t, parsedToken.Valid)

				// Verify claims
				claims, ok := parsedToken.Claims.(jwt.MapClaims)
				assert.True(t, ok)
				if tt.userIdentifier == "" {
					// For empty identifier, sub should be nil or empty
					sub := claims["sub"]
					assert.True(t, sub == nil || sub == "")
				} else {
					assert.Equal(t, tt.userIdentifier, claims["sub"])
				}
				assert.Equal(t, service.issuer, claims["iss"])
				assert.Contains(t, claims["aud"], service.audience[0])
				assert.NotEmpty(t, claims["jti"]) // JWT ID should be set
			}
		})
	}
}

func TestJwtService_GenerateToken_InvalidSecretKey(t *testing.T) {
	service := &JwtService{
		secretKey:  []byte("test-key"),
		expiration: time.Hour,
		issuer:     "test-issuer",
		audience:   []string{"test-audience"},
	}

	tokenType, token, expiresAt, err := service.GenerateToken("test-user")
	assert.NoError(t, err)
	assert.Equal(t, "Bearer", tokenType)
	assert.NotEmpty(t, token)
	assert.Greater(t, expiresAt, time.Now().UnixMilli())
}
