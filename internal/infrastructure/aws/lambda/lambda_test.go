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

	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/domain"
	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/dto"
	mockport "github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/port/mocks"
)

//go:embed golden/success_response.golden
var expectedGolden []byte

func TestHandleRequest_Success(t *testing.T) {
	fmt.Println("Starting TestHandleRequest_Success")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockController := mockport.NewMockCustomerController(ctrl)
	mockPresenter := mockport.NewMockPresenter(ctrl)

	customerController = mockController
	pr = mockPresenter

	// Prepare input and expected request/response
	reqInput := dto.GetCustomerInput{CPF: "12345678900"}
	customerReq := struct {
		Cpf string `json:"cpf"`
	}{Cpf: reqInput.CPF}

	body, err := json.Marshal(customerReq)
	assert.NoError(t, err)

	lambdaReq := events.APIGatewayProxyRequest{Body: string(body)}

	expectedResp := []byte(`{"name":"Test User"}`)

	mockController.
		EXPECT().
		Get(gomock.Any(), pr, reqInput).
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
	req := events.APIGatewayProxyRequest{Body: invalidBody}

	resp, _ := handleRequest(context.Background(), req)
	assert.Equal(t, 400, resp.StatusCode)
	assert.Contains(t, resp.Body, "invalid character")
}

func TestHandleRequest_ControllerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockController := mockport.NewMockCustomerController(ctrl)
	customerController = mockController

	reqInput := dto.GetCustomerInput{CPF: "wrong"}
	customerReq := struct {
		Cpf string `json:"cpf"`
	}{Cpf: reqInput.CPF}
	body, _ := json.Marshal(customerReq)

	lambdaReq := events.APIGatewayProxyRequest{Body: string(body)}

	mockController.
		EXPECT().
		Get(gomock.Any(), gomock.Any(), reqInput).
		Return(nil, &domain.NotFoundError{Message: "not found"}).
		Times(1)

	resp, _ := handleRequest(context.Background(), lambdaReq)
	assert.Equal(t, 404, resp.StatusCode)
	assert.Contains(t, resp.Body, "not found")
}
