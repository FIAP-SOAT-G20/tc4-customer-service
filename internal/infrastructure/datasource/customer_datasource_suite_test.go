package datasource_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/database"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/datasource"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/logger"
)

type CustomerDataSourceIntegrationTestSuite struct {
	suite.Suite
	ctx        context.Context
	db         *database.MongoDatabase
	dataSource port.CustomerDataSource
	collection *mongo.Collection
}

func (suite *CustomerDataSourceIntegrationTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	cfg := &config.Config{
		MongoURI:         getTestMongoURI(),
		MongoDatabase:    "fastfood_test",
		MongoTimeout:     30 * time.Second,
		MongoMaxPoolSize: 10,
		MongoMinPoolSize: 1,
	}

	l := logger.NewLogger(cfg)

	var err error
	suite.db, err = database.NewMongoConnection(cfg, l)
	require.NoError(suite.T(), err, "Failed to connect to test MongoDB")

	suite.dataSource = datasource.NewCustomerDataSource(suite.db)
	suite.collection = suite.db.Collection("customers")
}

func (suite *CustomerDataSourceIntegrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		if err := suite.db.Close(suite.ctx); err != nil {
			suite.T().Logf("Error closing database connection: %v", err)
		}
	}
}

func (suite *CustomerDataSourceIntegrationTestSuite) SetupTest() {
	_ = suite.collection.Drop(suite.ctx)
}

func (suite *CustomerDataSourceIntegrationTestSuite) TearDownTest() {
	_ = suite.collection.Drop(suite.ctx)
}

func TestCustomerDataSourceIntegrationTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	// Skip integration tests in coverage mode or when explicitly requested
	if os.Getenv("COVERAGE_MODE") == "true" || os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests in coverage mode")
	}

	suite.Run(t, new(CustomerDataSourceIntegrationTestSuite))
}

func getTestMongoURI() string {
	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		return uri
	}
	if uri := os.Getenv("TEST_MONGODB_URI"); uri != "" {
		return uri
	}
	return "mongodb://admin:admin@localhost:27017/fastfood_test?authSource=admin"
}
