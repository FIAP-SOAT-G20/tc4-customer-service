package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database settings
	DBHost         string
	DBPort         int
	DBUser         string
	DBPassword     string
	DBName         string
	DBMaxOpenConns int
	DBMaxIdleConns int
	DBMaxLifetime  time.Duration

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

	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	dbMaxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	dbMaxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "25"))
	dbMaxLifetime, _ := time.ParseDuration(getEnv("DB_CONN_MAX_LIFETIME", "5m"))

	// Environment
	environment := getEnv("ENVIRONMENT", "development")

	jwtExpirationStr := getEnv("JWT_EXPIRATION", "24h")
	jwtExpiration, err := time.ParseDuration(jwtExpirationStr)
	if err != nil {
		log.Printf("Warning: invalid JWT_EXPIRATION value %q: %v. Using default value 24h.", jwtExpirationStr, err)
		jwtExpiration = 24 * time.Hour
	}

	return &Config{
		// Database settings
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         dbPort,
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "fastfood_10soat_g18_tc2"),
		DBMaxOpenConns: dbMaxOpenConns,
		DBMaxIdleConns: dbMaxIdleConns,
		DBMaxLifetime:  dbMaxLifetime,

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
