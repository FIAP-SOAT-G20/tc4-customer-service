package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// DynamoDB settings
	DynamoTableName string
	DynamoRegion    string

	// Environment
	Environment string

	// JWT Settings
	JWTSecret     string
	JWTIssuer     string
	JWTAudience   string
	JWTExpiration time.Duration
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: .env file not found or error loading it: %v", err)
	}

	// Environment
	environment := getEnv("ENVIRONMENT", "development")

	jwtExpirationStr := getEnv("JWT_EXPIRATION", "24h")
	jwtExpiration, err := time.ParseDuration(jwtExpirationStr)
	if err != nil {
		log.Printf("Warning: invalid JWT_EXPIRATION value %q: %v. Using default value 24h.", jwtExpirationStr, err)
		jwtExpiration = 24 * time.Hour
	}

	return &Config{
		// DynamoDB settings
		DynamoTableName: getEnv("DYNAMODB_TABLE_NAME", "tc4-customer-service-dev-customers"),
		DynamoRegion:    getEnv("DYNAMODB_REGION", "us-east-1"),

		// Environment
		Environment: environment,

		// JWT Settings
		JWTSecret:     getEnv("JWT_SECRET", "SUPER_SECRET_KEY_DONT_TELL_ANYONE"),
		JWTIssuer:     getEnv("JWT_ISSUER", "https://fast-food-auth-abc12345.execute-api.us-east-1.amazonaws.com/prod"),
		JWTAudience:   getEnv("JWT_AUDIENCE", "https://fast-food-api-def67890.execute-api.us-east-1.amazonaws.com/prod"),
		JWTExpiration: jwtExpiration,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
