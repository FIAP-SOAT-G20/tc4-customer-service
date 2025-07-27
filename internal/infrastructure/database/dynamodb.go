package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDatabase struct {
	Client    *dynamodb.Client
	TableName string
	logger    *logger.Logger
}

func NewDynamoConnection(cfg *config.Config, l *logger.Logger) (*DynamoDatabase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Load AWS configuration
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(cfg.DynamoRegion))
	if err != nil {
		return nil, err
	}

	// Create DynamoDB client
	client := dynamodb.NewFromConfig(awsCfg)

	l.InfoContext(ctx, "Successfully connected to DynamoDB",
		"table", cfg.DynamoTableName,
		"region", cfg.DynamoRegion)

	return &DynamoDatabase{
		Client:    client,
		TableName: cfg.DynamoTableName,
		logger:    l,
	}, nil
}

// NewDynamoTestConnection creates a DynamoDB connection for testing with local endpoint
func NewDynamoTestConnection(cfg *config.Config, l *logger.Logger, endpoint string) (*DynamoDatabase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use static credentials for local DynamoDB
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx,
		awsConfig.WithRegion(cfg.DynamoRegion),
		awsConfig.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     "dummy",
				SecretAccessKey: "dummy",
				SessionToken:    "",
			},
		}),
	)
	if err != nil {
		return nil, err
	}

	// Create DynamoDB client with custom endpoint for local testing
	client := dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	l.InfoContext(ctx, "Successfully connected to DynamoDB Local",
		"table", cfg.DynamoTableName,
		"region", cfg.DynamoRegion,
		"endpoint", endpoint)

	return &DynamoDatabase{
		Client:    client,
		TableName: cfg.DynamoTableName,
		logger:    l,
	}, nil
}

func (d *DynamoDatabase) Close(ctx context.Context) error {
	d.logger.InfoContext(ctx, "Closing DynamoDB connection")
	return nil
}

func (d *DynamoDatabase) Ping(ctx context.Context) error {
	// Simple health check by describing the table
	_, err := d.Client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(d.TableName),
	})
	return err
}

// LogOperation logs DynamoDB operations for debugging
func (d *DynamoDatabase) LogOperation(ctx context.Context, operation, table string, duration time.Duration, err error) {
	if err != nil {
		d.logger.ErrorContext(ctx, "DynamoDB operation failed",
			slog.String("operation", operation),
			slog.String("table", table),
			slog.Duration("duration", duration),
			slog.String("error", err.Error()),
		)
		return
	}

	d.logger.DebugContext(ctx, "DynamoDB operation completed",
		slog.String("operation", operation),
		slog.String("table", table),
		slog.Duration("duration", duration),
	)
}
