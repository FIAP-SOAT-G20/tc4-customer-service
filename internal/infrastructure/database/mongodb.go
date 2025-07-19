package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDatabase struct {
	Client   *mongo.Client
	Database *mongo.Database
	mongoURI string
	logger   *logger.Logger
}

func NewMongoConnection(cfg *config.Config, l *logger.Logger) (*MongoDatabase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.MongoTimeout)
	defer cancel()

	// Configure MongoDB client options
	clientOptions := options.Client().
		ApplyURI(cfg.MongoURI).
		SetMaxPoolSize(cfg.MongoMaxPoolSize).
		SetMinPoolSize(cfg.MongoMinPoolSize).
		SetMaxConnIdleTime(5 * time.Minute).
		SetServerSelectionTimeout(cfg.MongoTimeout)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test the connection
	ctx, cancel = context.WithTimeout(context.Background(), cfg.MongoTimeout)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database := client.Database(cfg.MongoDatabase)

	l.InfoContext(ctx, "Successfully connected to MongoDB",
		"database", cfg.MongoDatabase,
		"uri", cfg.MongoURI)

	return &MongoDatabase{
		Client:   client,
		Database: database,
		mongoURI: cfg.MongoURI,
		logger:   l,
	}, nil
}

func (m *MongoDatabase) Close(ctx context.Context) error {
	if m.Client != nil {
		m.logger.InfoContext(ctx, "Closing MongoDB connection")
		return m.Client.Disconnect(ctx)
	}
	return nil
}

func (m *MongoDatabase) Collection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}

func (m *MongoDatabase) Ping(ctx context.Context) error {
	return m.Client.Ping(ctx, readpref.Primary())
}

// LogOperation logs MongoDB operations for debugging
func (m *MongoDatabase) LogOperation(ctx context.Context, operation, collection string, duration time.Duration, err error) {
	if err != nil {
		m.logger.ErrorContext(ctx, "MongoDB operation failed",
			slog.String("operation", operation),
			slog.String("collection", collection),
			slog.Duration("duration", duration),
			slog.String("error", err.Error()),
		)
		return
	}

	m.logger.DebugContext(ctx, "MongoDB operation completed",
		slog.String("operation", operation),
		slog.String("collection", collection),
		slog.Duration("duration", duration),
	)
}
