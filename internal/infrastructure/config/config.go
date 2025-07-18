package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// MongoDB settings
	MongoURI         string
	MongoDatabase    string
	MongoTimeout     time.Duration
	MongoMaxPoolSize uint64
	MongoMinPoolSize uint64

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

	mongoTimeout, _ := time.ParseDuration(getEnv("MONGO_TIMEOUT", "10s"))
	mongoMaxPoolSize, _ := strconv.ParseUint(getEnv("MONGO_MAX_POOL_SIZE", "100"), 10, 64)
	mongoMinPoolSize, _ := strconv.ParseUint(getEnv("MONGO_MIN_POOL_SIZE", "5"), 10, 64)

	// Environment
	environment := getEnv("ENVIRONMENT", "development")

	jwtExpirationStr := getEnv("JWT_EXPIRATION", "24h")
	jwtExpiration, err := time.ParseDuration(jwtExpirationStr)
	if err != nil {
		log.Printf("Warning: invalid JWT_EXPIRATION value %q: %v. Using default value 24h.", jwtExpirationStr, err)
		jwtExpiration = 24 * time.Hour
	}

	return &Config{
		// MongoDB settings
		MongoURI:         getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase:    getEnv("MONGO_DATABASE", "fastfood_10soat_g22_tc4"),
		MongoTimeout:     mongoTimeout,
		MongoMaxPoolSize: mongoMaxPoolSize,
		MongoMinPoolSize: mongoMinPoolSize,

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
