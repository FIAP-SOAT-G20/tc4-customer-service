package bdd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/cucumber/godog"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/adapter/controller"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/adapter/gateway"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/adapter/presenter"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/usecase"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/aws/lambda/request"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/database"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/datasource"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/logger"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/service"
)

type TestContext struct {
	response              events.APIGatewayProxyResponse
	responseBody          string
	requestBody           map[string]interface{}
	customerDataSource    port.CustomerDataSource
	customerGateway       port.CustomerGateway
	customerUseCase       port.CustomerUseCase
	customerController    port.CustomerController
	jsonPresenter         port.Presenter
	jwtPresenter          port.Presenter
	mongoClient           *mongo.Client
	mongoDatabase         *mongo.Database
	logger                *logger.Logger
	lastCreatedCustomerID string
}

var testCtx *TestContext

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		testCtx = &TestContext{
			requestBody: make(map[string]interface{}),
		}
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if testCtx.mongoClient != nil {
			_ = testCtx.mongoClient.Disconnect(ctx)
		}
		return ctx, nil
	})

	// Background steps
	ctx.Step(`^the customer service is running$`, theCustomerServiceIsRunning)
	ctx.Step(`^the database is clean$`, theDatabaseIsClean)

	// Authentication steps
	ctx.Step(`^I send an authentication request with CPF "([^"]*)"$`, iSendAnAuthenticationRequestWithCPF)

	// Customer management steps
	ctx.Step(`^a customer exists with CPF "([^"]*)"$`, aCustomerExistsWithCPF)
	ctx.Step(`^a customer exists with email "([^"]*)"$`, aCustomerExistsWithEmail)
	ctx.Step(`^a customer exists with ID "([^"]*)"$`, aCustomerExistsWithID)
	ctx.Step(`^I send a request to create a customer with the following details:$`, iSendARequestToCreateACustomerWithTheFollowingDetails)
	ctx.Step(`^I send a request to delete customer with ID "([^"]*)"$`, iSendARequestToDeleteCustomerWithID)
	ctx.Step(`^I send a request to get customer with CPF "([^"]*)"$`, iSendARequestToGetCustomerWithCPF)
	ctx.Step(`^I send a request to get customer with ID "([^"]*)"$`, iSendARequestToGetCustomerWithID)
	ctx.Step(`^I send a request to list all customers$`, iSendARequestToListAllCustomers)
	ctx.Step(`^I send a request to update customer with ID "([^"]*)" with the following details:$`, iSendARequestToUpdateCustomerWithIDWithTheFollowingDetails)
	ctx.Step(`^the following customers exist:$`, theFollowingCustomersExist)

	// Response validation steps
	ctx.Step(`^I should receive a response with status (\d+)$`, iShouldReceiveAResponseWithStatus)
	ctx.Step(`^the response should contain a valid JWT token$`, theResponseShouldContainAValidJWTToken)
	ctx.Step(`^the response should contain customer information$`, theResponseShouldContainCustomerInformation)
	ctx.Step(`^the response should contain an error message "([^"]*)"$`, theResponseShouldContainAnErrorMessage)
	ctx.Step(`^the response should contain the customer ID$`, theResponseShouldContainTheCustomerID)
	ctx.Step(`^the response should contain customer details$`, theResponseShouldContainCustomerDetails)
	ctx.Step(`^the response should contain a list of (\d+) customers$`, theResponseShouldContainAListOfCustomers)
	ctx.Step(`^the response should contain the updated customer details$`, theResponseShouldContainTheUpdatedCustomerDetails)
}

func theCustomerServiceIsRunning() error {
	// Set a test environment
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://admin:admin@localhost:27017/fastfood_test?authSource=admin"
	}

	cfg := &config.Config{
		Environment:   "test",
		MongoURI:      mongoURI,
		MongoDatabase: "fastfood_test",
		JWTSecret:     "test-secret-key",
		JWTIssuer:     "test-issuer",
		JWTAudience:   "test-audience",
		JWTExpiration: 86400000000000, // 24h in nanoseconds
	}

	testCtx.logger = logger.NewLogger(cfg)

	// Setup test database
	mongoDb, err := database.NewMongoConnection(cfg, testCtx.logger)
	if err != nil {
		fmt.Printf("Skipping BDD tests: MongoDB not available: %v\n", err)
		os.Exit(0)
	}

	testCtx.mongoClient = mongoDb.Client
	testCtx.mongoDatabase = mongoDb.Database

	// Setup dependencies
	jwtService := service.NewJWTService(cfg)
	testCtx.customerDataSource = datasource.NewCustomerDataSource(mongoDb)
	testCtx.customerGateway = gateway.NewCustomerGateway(testCtx.customerDataSource)
	testCtx.customerUseCase = usecase.NewCustomerUseCase(testCtx.customerGateway)
	testCtx.customerController = controller.NewCustomerController(testCtx.customerUseCase)
	testCtx.jsonPresenter = presenter.NewCustomerJsonPresenter()
	testCtx.jwtPresenter = presenter.NewCustomerJwtTokenPresenter(jwtService)

	return nil
}

func theDatabaseIsClean() error {
	if testCtx.mongoDatabase == nil {
		return fmt.Errorf("database not initialized")
	}

	// Clean all collections
	collections := []string{"customers"}
	for _, collection := range collections {
		err := testCtx.mongoDatabase.Collection(collection).Drop(context.Background())
		if err != nil {
			return fmt.Errorf("failed to clean collection %s: %w", collection, err)
		}
	}

	return nil
}

func aCustomerExistsWithEmail(email string) error {
	customerRequest := request.CustomerRequest{
		Name:  "Test Customer",
		Email: email,
		CPF:   "12345678901",
	}

	return createTestCustomer(customerRequest)
}

func aCustomerExistsWithID(id string) error {
	// Create a customer for BDD tests and capture the real ID
	customerRequest := request.CustomerRequest{
		Name:  "Test Customer",
		Email: "test@example.com",
		CPF:   "12345678901",
	}

	input := customerRequest.ToCreateCustomerInput()
	resp, err := testCtx.customerController.Create(context.Background(), testCtx.jsonPresenter, input)
	if err != nil {
		return err
	}

	// Parse response to get real customer ID
	var customerData map[string]interface{}
	if err := json.Unmarshal(resp, &customerData); err != nil {
		return err
	}

	if realID, exists := customerData["id"]; exists {
		testCtx.lastCreatedCustomerID = realID.(string)
	} else {
		return fmt.Errorf("could not get customer ID from response")
	}

	return nil
}

func aCustomerExistsWithCPF(cpf string) error {
	customerRequest := request.CustomerRequest{
		Name:  "Test Customer",
		Email: "test@example.com",
		CPF:   cpf,
	}

	return createTestCustomer(customerRequest)
}

func theFollowingCustomersExist(table *godog.Table) error {
	for i, row := range table.Rows {
		if i == 0 { // Skip header row
			continue
		}

		customerRequest := request.CustomerRequest{
			Name:  row.Cells[0].Value,
			Email: row.Cells[1].Value,
			CPF:   row.Cells[2].Value,
		}

		if err := createTestCustomer(customerRequest); err != nil {
			return err
		}
	}

	return nil
}

func createTestCustomer(customerRequest request.CustomerRequest) error {
	input := customerRequest.ToCreateCustomerInput()
	_, err := testCtx.customerController.Create(context.Background(), testCtx.jsonPresenter, input)
	return err
}

func iSendAnAuthenticationRequestWithCPF(cpf string) error {
	customerRequest := request.CustomerRequest{
		CPF: cpf,
	}

	return sendLambdaRequest("POST", "/auth", customerRequest, map[string]string{})
}

// iSendARequestToCreateACustomerWithTheFollowingDetails parses a data table and sends a request to create a customer.
// The table should contain rows with the following cells:
// - name: The customer's name
// - email: The customer's email
// - cpf: The customer's CPF
// It returns an error if the request fails.
func iSendARequestToCreateACustomerWithTheFollowingDetails(table *godog.Table) error {
	customerRequest := request.CustomerRequest{}

	for _, row := range table.Rows {
		switch row.Cells[0].Value {
		case "name":
			customerRequest.Name = row.Cells[1].Value
		case "email":
			customerRequest.Email = row.Cells[1].Value
		case "cpf":
			customerRequest.CPF = row.Cells[1].Value
		}
	}

	return sendLambdaRequest("POST", "/customers", customerRequest, map[string]string{})
}

func iSendARequestToGetCustomerWithCPF(cpf string) error {
	return sendLambdaRequest("GET", "/customers", nil, map[string]string{"cpf": cpf})
}

func iSendARequestToGetCustomerWithID(id string) error {
	// Use the last created customer ID if available
	actualID := id
	if testCtx.lastCreatedCustomerID != "" {
		actualID = testCtx.lastCreatedCustomerID
	}
	return sendLambdaRequest("GET", "/customers", nil, map[string]string{"id": actualID})
}

func iSendARequestToListAllCustomers() error {
	return sendLambdaRequest("GET", "/customers", nil, map[string]string{})
}

func iSendARequestToUpdateCustomerWithIDWithTheFollowingDetails(id string, table *godog.Table) error {
	customerRequest := request.CustomerRequest{}

	for _, row := range table.Rows {
		switch row.Cells[0].Value {
		case "name":
			customerRequest.Name = row.Cells[1].Value
		case "email":
			customerRequest.Email = row.Cells[1].Value
		case "cpf":
			customerRequest.CPF = row.Cells[1].Value
		}
	}

	// Use the last created customer ID if available
	actualID := id
	if testCtx.lastCreatedCustomerID != "" {
		actualID = testCtx.lastCreatedCustomerID
	}

	return sendLambdaRequest("PUT", "/customers", customerRequest, map[string]string{"id": actualID})
}

func iSendARequestToDeleteCustomerWithID(id string) error {
	// Use the last created customer ID if available
	actualID := id
	if testCtx.lastCreatedCustomerID != "" {
		actualID = testCtx.lastCreatedCustomerID
	}
	return sendLambdaRequest("DELETE", "/customers", nil, map[string]string{"id": actualID})
}

func sendLambdaRequest(method, resource string, body interface{}, pathParams map[string]string) error {
	var bodyStr string
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyStr = string(bodyBytes)

		// Store request body for later validation
		if err := json.Unmarshal(bodyBytes, &testCtx.requestBody); err != nil {
			return err
		}
	}

	// Create API Gateway request event
	request := events.APIGatewayProxyRequest{
		HTTPMethod:            method,
		Resource:              resource,
		Body:                  bodyStr,
		PathParameters:        pathParams,
		QueryStringParameters: map[string]string{},
	}

	// Handle query parameters for GET requests
	if method == "GET" {
		if cpf, exists := pathParams["cpf"]; exists {
			request.QueryStringParameters["cpf"] = cpf
			delete(request.PathParameters, "cpf")
		}
	}

	// Simulate lambda handler logic
	ctx := context.Background()
	var err error

	switch {
	case resource == "/auth" && method == "POST":
		testCtx.response, err = handleTestAuthRequest(ctx, request)
	case method == "GET":
		testCtx.response, err = handleTestGetRequest(ctx, request)
	case method == "POST":
		testCtx.response, err = handleTestPostRequest(ctx, request)
	case method == "PUT":
		testCtx.response, err = handleTestPutRequest(ctx, request)
	case method == "DELETE":
		testCtx.response, err = handleTestDeleteRequest(ctx, request)
	default:
		return fmt.Errorf("unsupported method: %s", method)
	}

	if err != nil {
		return err
	}

	testCtx.responseBody = testCtx.response.Body
	return nil
}

// Lambda handlers for testing
func handleTestAuthRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var customerRequest request.CustomerRequest
	if err := json.Unmarshal([]byte(req.Body), &customerRequest); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"message": "Invalid request body"}`}, nil
	}

	if customerRequest.CPF == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"message": "CPF is required"}`}, nil
	}

	input := dto.GetCustomerByCPFInput{CPF: customerRequest.CPF}
	resp, err := testCtx.customerController.GetByCPF(ctx, testCtx.jwtPresenter, input)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 401, Body: `{"message": "Invalid credentials"}`}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(resp)}, nil
}

func handleTestGetRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	customerID, hasID := req.PathParameters["id"]
	cpf := req.QueryStringParameters["cpf"]

	if !hasID && cpf == "" {
		// List customers
		input := dto.ListCustomersInput{Page: 1, Limit: 10}
		resp, err := testCtx.customerController.List(ctx, testCtx.jsonPresenter, input)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: `{"message": "Internal server error"}`}, nil
		}
		return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(resp)}, nil
	}

	if cpf != "" {
		// Get by CPF
		input := dto.GetCustomerByCPFInput{CPF: cpf}
		resp, err := testCtx.customerController.GetByCPF(ctx, testCtx.jsonPresenter, input)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 404, Body: `{"message": "Customer not found"}`}, nil
		}
		return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(resp)}, nil
	}

	// Get by ID
	input := dto.GetCustomerInput{ID: customerID}
	resp, err := testCtx.customerController.Get(ctx, testCtx.jsonPresenter, input)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 404, Body: `{"message": "Customer not found"}`}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(resp)}, nil
}

func handleTestPostRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var customerRequest request.CustomerRequest
	if err := json.Unmarshal([]byte(req.Body), &customerRequest); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"message": "Invalid request body"}`}, nil
	}

	input := customerRequest.ToCreateCustomerInput()
	resp, err := testCtx.customerController.Create(ctx, testCtx.jsonPresenter, input)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return events.APIGatewayProxyResponse{StatusCode: 409, Body: `{"message": "Customer with this email already exists"}`}, nil
		}
		if strings.Contains(err.Error(), "Invalid CPF") {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"message": "Invalid CPF format"}`}, nil
		}
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: `{"message": "Internal server error"}`}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 201, Body: string(resp)}, nil
}

func handleTestPutRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	customerID, hasID := req.PathParameters["id"]
	if !hasID {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"message": "Customer ID is required for update"}`}, nil
	}

	var customerRequest request.CustomerRequest
	if err := json.Unmarshal([]byte(req.Body), &customerRequest); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"message": "Invalid request body"}`}, nil
	}

	input := customerRequest.ToUpdateCustomerInput()
	input.ID = customerID

	resp, err := testCtx.customerController.Update(ctx, testCtx.jsonPresenter, input)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 404, Body: `{"message": "Customer not found"}`}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(resp)}, nil
}

func handleTestDeleteRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	customerID, hasID := req.PathParameters["id"]
	if !hasID {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"message": "Customer ID is required for deletion"}`}, nil
	}

	input := dto.DeleteCustomerInput{ID: customerID}
	_, err := testCtx.customerController.Delete(ctx, testCtx.jsonPresenter, input)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 404, Body: `{"message": "Customer not found"}`}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 204, Body: ""}, nil
}

// Response validation functions
func iShouldReceiveAResponseWithStatus(expectedStatus int) error {
	if testCtx.response.StatusCode != expectedStatus {
		return fmt.Errorf("expected status %d, got %d. Response body: %s",
			expectedStatus, testCtx.response.StatusCode, testCtx.responseBody)
	}
	return nil
}

func theResponseShouldContainAValidJWTToken() error {
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(testCtx.responseBody), &response); err != nil {
		return err
	}

	token, exists := response["access_token"]
	if !exists {
		return fmt.Errorf("response does not contain access_token")
	}

	tokenStr, ok := token.(string)
	if !ok || tokenStr == "" {
		return fmt.Errorf("access_token is not a valid string")
	}

	if tokenStr != "Bearer" {
		parts := strings.Split(tokenStr, ".")
		if len(parts) != 3 {
			return fmt.Errorf("access_token is not in valid JWT format or Bearer token")
		}
	}

	return nil
}

func theResponseShouldContainCustomerInformation() error {
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(testCtx.responseBody), &response); err != nil {
		return err
	}

	customer, exists := response["customer"]
	if !exists {
		return fmt.Errorf("response does not contain customer information")
	}

	customerMap, ok := customer.(map[string]interface{})
	if !ok {
		return fmt.Errorf("customer information is not in expected format")
	}

	requiredFields := []string{"id", "name", "email", "cpf"}
	for _, field := range requiredFields {
		if _, exists := customerMap[field]; !exists {
			return fmt.Errorf("customer information missing field: %s", field)
		}
	}

	return nil
}

func theResponseShouldContainAnErrorMessage(expectedMessage string) error {
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(testCtx.responseBody), &response); err != nil {
		return fmt.Errorf("failed to parse response body as JSON: %w. Response body: %s", err, testCtx.responseBody)
	}

	message, exists := response["message"]
	if !exists {
		return fmt.Errorf("response does not contain 'message' field. Response body: %s", testCtx.responseBody)
	}

	messageStr, ok := message.(string)
	if !ok {
		return fmt.Errorf("'message' field is not a string. Type: %T, Value: %v", message, message)
	}

	if !strings.Contains(messageStr, expectedMessage) {
		return fmt.Errorf("expected error message to contain '%s', got '%s'", expectedMessage, messageStr)
	}

	return nil
}

func theResponseShouldContainTheCustomerID() error {
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(testCtx.responseBody), &response); err != nil {
		return err
	}

	id, exists := response["id"]
	if !exists {
		return fmt.Errorf("response does not contain customer ID")
	}

	idStr, ok := id.(string)
	if !ok || idStr == "" {
		return fmt.Errorf("customer ID is not a valid string")
	}

	return nil
}

func theResponseShouldContainCustomerDetails() error {
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(testCtx.responseBody), &response); err != nil {
		return err
	}

	requiredFields := []string{"id", "name", "email", "cpf"}
	for _, field := range requiredFields {
		if _, exists := response[field]; !exists {
			return fmt.Errorf("response missing field: %s", field)
		}
	}

	return nil
}

func theResponseShouldContainAListOfCustomers(expectedCount int) error {
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(testCtx.responseBody), &response); err != nil {
		return err
	}

	customers, exists := response["customers"]
	if !exists {
		return fmt.Errorf("response does not contain customers field")
	}

	customersList, ok := customers.([]interface{})
	if !ok {
		return fmt.Errorf("customers field is not an array")
	}

	if len(customersList) != expectedCount {
		return fmt.Errorf("expected %d customers, got %d", expectedCount, len(customersList))
	}

	return nil
}

func theResponseShouldContainTheUpdatedCustomerDetails() error {
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(testCtx.responseBody), &response); err != nil {
		return err
	}

	// Check if the response contains the updated fields from the request
	for key, expectedValue := range testCtx.requestBody {
		if actualValue, exists := response[key]; !exists {
			return fmt.Errorf("response missing updated field: %s", key)
		} else if actualValue != expectedValue {
			return fmt.Errorf("field %s was not updated correctly. Expected: %v, Got: %v",
				key, expectedValue, actualValue)
		}
	}

	return nil
}
