package datasource_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/database"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/datasource"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type CustomerDynamoDataSourceIntegrationTestSuite struct {
	suite.Suite
	ctx        context.Context
	db         *database.DynamoDatabase
	dataSource port.CustomerDataSource
}

func (suite *CustomerDynamoDataSourceIntegrationTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	cfg := &config.Config{
		DynamoTableName: getTestDynamoTableName(),
		DynamoRegion:    getTestDynamoRegion(),
		Environment:     "test",
	}

	l := logger.NewLogger(cfg)

	var err error
	testEndpoint := getTestDynamoEndpoint()
	if testEndpoint == "" {
		testEndpoint = "http://localhost:8000" // default for local testing
	}
	suite.db, err = database.NewDynamoTestConnection(cfg, l, testEndpoint)
	require.NoError(suite.T(), err, "Failed to connect to test DynamoDB")

	suite.dataSource = datasource.NewCustomerDynamoDataSource(suite.db)

	// Create test table if it doesn't exist
	suite.createTestTableIfNotExists()
}

func (suite *CustomerDynamoDataSourceIntegrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		// Clean up test data
		suite.clearTestTable()
	}
}

func (suite *CustomerDynamoDataSourceIntegrationTestSuite) SetupTest() {
	// Clear table before each test
	suite.clearTestTable()
}

func (suite *CustomerDynamoDataSourceIntegrationTestSuite) TearDownTest() {
	// Clear table after each test
	suite.clearTestTable()
}

func (suite *CustomerDynamoDataSourceIntegrationTestSuite) createTestTableIfNotExists() {
	// Check if table exists
	_, err := suite.db.Client.DescribeTable(suite.ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(suite.db.TableName),
	})

	if err != nil {
		// Table doesn't exist, create it
		_, err = suite.db.Client.CreateTable(suite.ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(suite.db.TableName),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("id"),
					KeyType:       types.KeyTypeHash,
				},
			},
			AttributeDefinitions: []types.AttributeDefinition{
				{
					AttributeName: aws.String("id"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("cpf"),
					AttributeType: types.ScalarAttributeTypeS,
				},
			},
			GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
				{
					IndexName: aws.String("cpf-index"),
					KeySchema: []types.KeySchemaElement{
						{
							AttributeName: aws.String("cpf"),
							KeyType:       types.KeyTypeHash,
						},
					},
					Projection: &types.Projection{
						ProjectionType: types.ProjectionTypeAll,
					},
					ProvisionedThroughput: &types.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(5),
						WriteCapacityUnits: aws.Int64(5),
					},
				},
			},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
		})

		if err != nil {
			suite.T().Logf("Warning: Could not create test table: %v", err)
		} else {
			// Wait for the table to be active
			waiter := dynamodb.NewTableExistsWaiter(suite.db.Client)
			err = waiter.Wait(suite.ctx, &dynamodb.DescribeTableInput{
				TableName: aws.String(suite.db.TableName),
			}, 30*time.Second, func(o *dynamodb.TableExistsWaiterOptions) {
				o.MaxDelay = 30 * time.Second
				o.MinDelay = 2 * time.Second
			})
			if err != nil {
				suite.T().Logf("Warning: Table creation wait failed: %v", err)
			}
		}
	}
}

func (suite *CustomerDynamoDataSourceIntegrationTestSuite) clearTestTable() {
	// Scan and delete all items
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(suite.db.TableName),
	}

	result, err := suite.db.Client.Scan(suite.ctx, scanInput)
	if err != nil {
		suite.T().Logf("Warning: Could not scan table for cleanup: %v", err)
		return
	}

	for _, item := range result.Items {
		if id, exists := item["id"]; exists {
			_, err := suite.db.Client.DeleteItem(suite.ctx, &dynamodb.DeleteItemInput{
				TableName: aws.String(suite.db.TableName),
				Key: map[string]types.AttributeValue{
					"id": id,
				},
			})
			if err != nil {
				suite.T().Logf("Warning: Could not delete item during cleanup: %v", err)
			}
		}
	}
}

func TestCustomerDynamoDataSourceIntegrationTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	// Skip integration tests in coverage mode or when explicitly requested
	if os.Getenv("COVERAGE_MODE") == "true" || os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests in coverage mode")
	}

	// Skip if no DynamoDB endpoint is configured
	if getTestDynamoEndpoint() == "" {
		t.Skip("Skipping DynamoDB integration tests - no endpoint configured")
	}

	suite.Run(t, new(CustomerDynamoDataSourceIntegrationTestSuite))
}

func getTestDynamoTableName() string {
	if tableName := os.Getenv("DYNAMODB_TABLE_NAME"); tableName != "" {
		return tableName
	}
	if tableName := os.Getenv("TEST_DYNAMODB_TABLE_NAME"); tableName != "" {
		return tableName
	}
	return "tc4-customer-service-test-customers"
}

func getTestDynamoRegion() string {
	if region := os.Getenv("DYNAMODB_REGION"); region != "" {
		return region
	}
	if region := os.Getenv("TEST_DYNAMODB_REGION"); region != "" {
		return region
	}
	return "us-east-1"
}

func getTestDynamoEndpoint() string {
	return os.Getenv("TEST_DYNAMODB_ENDPOINT")
}
