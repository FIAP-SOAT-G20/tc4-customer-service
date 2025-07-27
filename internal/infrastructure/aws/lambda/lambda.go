package lambda

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/adapter/controller"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/adapter/gateway"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/adapter/presenter"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/usecase"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/aws/lambda/request"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/aws/lambda/response"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/database"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/datasource"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/logger"
	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/infrastructure/service"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/port"
)

var customerDataSource port.CustomerDataSource
var customerGateway port.CustomerGateway
var customerUseCase port.CustomerUseCase
var customerController port.CustomerController
var jsonPresenter port.Presenter
var jwtPresenter port.Presenter
var l *logger.Logger

// init function is called in a lambda cold start. So, at this moment is initialized
// all structures and also the database connection
func init() {
	fmt.Println("🟠 Initializing lambda presenter")
	cfg := config.LoadConfig()
	l = logger.NewLogger(cfg)

	if cfg.Environment == "test" {
		return
	}

	// Skip initialization during test execution
	if testing.Testing() {
		return
	}

	db, err := database.NewDynamoConnection(cfg, l)
	if err != nil {
		panic(err)
	}
	jwtService := service.NewJWTService(cfg)
	customerDataSource = datasource.NewCustomerDynamoDataSource(db)
	customerGateway = gateway.NewCustomerGateway(customerDataSource)
	customerUseCase = usecase.NewCustomerUseCase(customerGateway)
	customerController = controller.NewCustomerController(customerUseCase)
	jsonPresenter = presenter.NewCustomerJsonPresenter()
	jwtPresenter = presenter.NewCustomerJwtTokenPresenter(jwtService)
}

// StartLambda is the function that tells lambda which function should be call to start lambda.
func StartLambda() {
	fmt.Println("🟢 Lambda is ready to receive requests!")
	lambda.Start(handleRequest)
}

// handleRequest responsible to handle lambda events
func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	l.InfoContext(ctx, "Starting lambda handler",
		"httpMethod", req.HTTPMethod,
		"resource", req.Resource,
		"pathParameters", req.PathParameters,
		"queryStringParameters", req.QueryStringParameters,
		"isBase64Encoded", req.IsBase64Encoded,
		"body", req.Body)

	// Check if it's an authentication request
	if req.Resource == "/auth" && req.HTTPMethod == "POST" {
		return handleAuthRequest(ctx, req)
	}

	switch req.HTTPMethod {
	case "GET":
		return handleGetRequest(ctx, req)
	case "POST":
		return handlePostRequest(ctx, req)
	case "PUT":
		return handlePutRequest(ctx, req)
	case "DELETE":
		return handleDeleteRequest(ctx, req)
	default:
		return response.NewAPIGatewayProxyResponseError(&domain.InvalidInputError{
			Message: fmt.Sprintf("HTTP method %s not supported", req.HTTPMethod),
		}), nil
	}
}

// handleGetRequest handles GET requests for customers
func handleGetRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Check if it's a list request (no ID in path) or get by ID/CPF
	customerID, hasID := req.PathParameters["id"]
	cpf := req.QueryStringParameters["cpf"]

	// List customers with pagination
	if !hasID && cpf == "" {
		page := 1
		limit := 10
		if pageStr := req.QueryStringParameters["page"]; pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}
		if limitStr := req.QueryStringParameters["limit"]; limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
				limit = l
			}
		}

		input := dto.ListCustomersInput{
			Page:  page,
			Limit: limit,
		}

		resp, err := customerController.List(ctx, jsonPresenter, input)
		if err != nil {
			l.ErrorContext(ctx, "Failed to list customers", "error", err)
			return response.NewAPIGatewayProxyResponseError(err), nil
		}
		return response.NewAPIGatewayProxyResponse(resp), nil
	}

	// Get by CPF
	if cpf != "" {
		input := dto.GetCustomerByCPFInput{CPF: cpf}
		resp, err := customerController.GetByCPF(ctx, jsonPresenter, input)
		if err != nil {
			l.ErrorContext(ctx, "Failed to get customer by CPF", "cpf", cpf, "error", err)
			return response.NewAPIGatewayProxyResponseError(err), nil
		}
		return response.NewAPIGatewayProxyResponse(resp), nil
	}

	// Get by ID
	id, err := strconv.Atoi(customerID)
	if err != nil {
		l.ErrorContext(ctx, "Invalid customer ID", "id", customerID, "error", err)
		return response.NewAPIGatewayProxyResponseError(err), nil
	}
	input := dto.GetCustomerInput{ID: id}
	resp, err := customerController.Get(ctx, jsonPresenter, input)
	if err != nil {
		l.ErrorContext(ctx, "Failed to get customer", "id", customerID, "error", err)
		return response.NewAPIGatewayProxyResponseError(err), nil
	}
	return response.NewAPIGatewayProxyResponse(resp), nil
}

// handlePostRequest handles POST requests to create customers
func handlePostRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var customerRequest request.CustomerRequest
	var body = []byte(req.Body)
	var err error

	if req.IsBase64Encoded {
		body, err = base64.StdEncoding.DecodeString(req.Body)
		if err != nil {
			return response.NewAPIGatewayProxyResponseError(err), nil
		}
	}

	err = json.Unmarshal(body, &customerRequest)
	if err != nil {
		return response.NewAPIGatewayProxyResponseError(&domain.InvalidInputError{Message: err.Error()}), nil
	}

	input := customerRequest.ToCreateCustomerInput()
	resp, err := customerController.Create(ctx, jsonPresenter, input)
	if err != nil {
		l.ErrorContext(ctx, "Failed to create customer", "error", err)
		return response.NewAPIGatewayProxyResponseError(err), nil
	}

	return response.NewAPIGatewayProxyResponse(resp), nil
}

// handlePutRequest handles PUT requests to update customers
func handlePutRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	customerID, hasID := req.PathParameters["id"]
	if !hasID {
		return response.NewAPIGatewayProxyResponseError(&domain.InvalidInputError{
			Message: "Customer ID is required for update",
		}), nil
	}

	var customerRequest request.CustomerRequest
	var body = []byte(req.Body)
	var err error

	if req.IsBase64Encoded {
		body, err = base64.StdEncoding.DecodeString(req.Body)
		if err != nil {
			return response.NewAPIGatewayProxyResponseError(err), nil
		}
	}

	err = json.Unmarshal(body, &customerRequest)
	if err != nil {
		return response.NewAPIGatewayProxyResponseError(&domain.InvalidInputError{Message: err.Error()}), nil
	}

	id, err := strconv.Atoi(customerID)
	if err != nil {
		l.ErrorContext(ctx, "Invalid customer ID", "id", customerID, "error", err)
		return response.NewAPIGatewayProxyResponseError(err), nil
	}

	input := customerRequest.ToUpdateCustomerInput()
	input.ID = id

	resp, err := customerController.Update(ctx, jsonPresenter, input)
	if err != nil {
		l.ErrorContext(ctx, "Failed to update customer", "id", customerID, "error", err)
		return response.NewAPIGatewayProxyResponseError(err), nil
	}

	return response.NewAPIGatewayProxyResponse(resp), nil
}

// handleDeleteRequest handles DELETE requests to delete customers
func handleDeleteRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	customerID, hasID := req.PathParameters["id"]
	if !hasID {
		return response.NewAPIGatewayProxyResponseError(&domain.InvalidInputError{
			Message: "Customer ID is required for deletion",
		}), nil
	}

	id, err := strconv.Atoi(customerID)
	if err != nil {
		l.ErrorContext(ctx, "Invalid customer ID", "id", customerID, "error", err)
		return response.NewAPIGatewayProxyResponseError(err), nil
	}

	input := dto.DeleteCustomerInput{ID: id}
	resp, err := customerController.Delete(ctx, jsonPresenter, input)
	if err != nil {
		l.ErrorContext(ctx, "Failed to delete customer", "id", customerID, "error", err)
		return response.NewAPIGatewayProxyResponseError(err), nil
	}

	return response.NewAPIGatewayProxyResponse(resp), nil
}

// handleAuthRequest handles authentication requests that return JWT tokens
func handleAuthRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var customerRequest request.CustomerRequest
	var body = []byte(req.Body)
	var err error

	if req.IsBase64Encoded {
		body, err = base64.StdEncoding.DecodeString(req.Body)
		if err != nil {
			return response.NewAPIGatewayProxyResponseError(err), nil
		}
	}

	err = json.Unmarshal(body, &customerRequest)
	if err != nil {
		return response.NewAPIGatewayProxyResponseError(&domain.InvalidInputError{Message: err.Error()}), nil
	}

	// Authentication typically uses CPF lookup
	if customerRequest.CPF == "" {
		return response.NewAPIGatewayProxyResponseError(&domain.InvalidInputError{
			Message: "CPF is required for authentication",
		}), nil
	}

	input := dto.GetCustomerByCPFInput{CPF: customerRequest.CPF}
	resp, err := customerController.GetByCPF(ctx, jwtPresenter, input)
	if err != nil {
		l.ErrorContext(ctx, "Failed to authenticate customer", "cpf", customerRequest.CPF, "error", err)
		return response.NewAPIGatewayProxyResponseError(err), nil
	}

	return response.NewAPIGatewayProxyResponse(resp), nil
}
