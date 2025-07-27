package lambda

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/dto"
	mockport "github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/port/mocks"
)

//go:embed golden/success_response.golden
var expectedGolden []byte

func TestHandleRequest_GetByID_Success(t *testing.T) {
	fmt.Println("Starting TestHandleRequest_GetByID_Success")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockController := mockport.NewMockCustomerController(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)

	customerController = mockController
	jsonPresenter = mockPresenter

	// Prepare input and expected request/response for GET by ID
	customerID := "123"
	reqInput := dto.GetCustomerInput{ID: 123}

	lambdaReq := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		PathParameters: map[string]string{
			"id": customerID,
		},
	}

	expectedResp := []byte(`{"name":"Test User"}`)

	mockController.
		EXPECT().
		Get(gomock.Any(), jsonPresenter, reqInput).
		Return(expectedResp, nil).
		Times(1)

	// Call handleRequest
	got, err := handleRequest(context.Background(), lambdaReq)
	assert.NoError(t, err)
	assert.Equal(t, string(expectedResp), got.Body)

	assert.JSONEq(t, string(expectedGolden), got.Body)
}

func TestHandleRequest_InvalidJSON(t *testing.T) {
	invalidBody := "{ invalid json }"
	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Body:       invalidBody,
	}

	resp, _ := handleRequest(context.Background(), req)
	assert.Equal(t, 400, resp.StatusCode)
	assert.Contains(t, resp.Body, "invalid character")
}

func TestHandleRequest_ControllerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockController := mockport.NewMockCustomerController(ctrl)
	customerController = mockController

	customerID := "wrong"

	lambdaReq := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		PathParameters: map[string]string{
			"id": customerID,
		},
	}

	// Test expects an invalid ID to return 500 error (not 404) since conversion fails
	resp, _ := handleRequest(context.Background(), lambdaReq)
	assert.Equal(t, 500, resp.StatusCode)
	assert.Contains(t, resp.Body, "invalid syntax")
}

func TestHandleRequest_CreateCustomer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockController := mockport.NewMockCustomerController(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)

	customerController = mockController
	jsonPresenter = mockPresenter

	customerReq := struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		CPF   string `json:"cpf"`
	}{
		Name:  "John Doe",
		Email: "john@example.com",
		CPF:   "12345678900",
	}

	body, _ := json.Marshal(customerReq)
	lambdaReq := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Body:       string(body),
	}

	expectedResp := []byte(`{"id":"123","name":"John Doe"}`)

	mockController.
		EXPECT().
		Create(gomock.Any(), jsonPresenter, gomock.Any()).
		Return(expectedResp, nil).
		Times(1)

	resp, err := handleRequest(context.Background(), lambdaReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, string(expectedResp), resp.Body)
}

func TestHandleRequest_UnsupportedMethod(t *testing.T) {
	lambdaReq := events.APIGatewayProxyRequest{
		HTTPMethod: "PATCH",
	}

	resp, _ := handleRequest(context.Background(), lambdaReq)
	assert.Equal(t, 400, resp.StatusCode)
	assert.Contains(t, resp.Body, "HTTP method PATCH not supported")
}

func TestHandleRequest_Auth_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockController := mockport.NewMockCustomerController(ctrl)
	mockJwtPresenter := mockport.NewMockPresenter(ctrl)

	customerController = mockController
	jwtPresenter = mockJwtPresenter

	customerReq := struct {
		CPF string `json:"cpf"`
	}{
		CPF: "12345678900",
	}

	body, _ := json.Marshal(customerReq)
	lambdaReq := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Resource:   "/auth",
		Body:       string(body),
	}

	expectedResp := []byte(`{"token":"jwt.token.here"}`)

	mockController.
		EXPECT().
		GetByCPF(gomock.Any(), jwtPresenter, gomock.Any()).
		Return(expectedResp, nil).
		Times(1)

	resp, err := handleRequest(context.Background(), lambdaReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, string(expectedResp), resp.Body)
}

func TestHandleRequest_Auth_MissingCPF(t *testing.T) {
	customerReq := struct {
		Name string `json:"name"`
	}{
		Name: "John Doe",
	}

	body, _ := json.Marshal(customerReq)
	lambdaReq := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Resource:   "/auth",
		Body:       string(body),
	}

	resp, _ := handleRequest(context.Background(), lambdaReq)
	assert.Equal(t, 400, resp.StatusCode)
	assert.Contains(t, resp.Body, "CPF is required for authentication")
}
