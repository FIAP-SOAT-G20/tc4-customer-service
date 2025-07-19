package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain"
)

type errorResponse struct {
	Title   string `json:"title"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewErrorResponse(title, status, message string) errorResponse {
	return errorResponse{
		Title:   title,
		Status:  status,
		Message: message,
	}
}

func NewAPIGatewayProxyResponseError(err error) events.APIGatewayProxyResponse {
	var title string
	var status int
	var internal *domain.InternalError
	var validation *domain.ValidationError
	var notfound *domain.NotFoundError
	var invalidInput *domain.InvalidInputError
	switch {
	case errors.As(err, &internal):
		title = internal.Message
		status = http.StatusInternalServerError
	case errors.As(err, &validation):
		title = validation.Message
		status = http.StatusPreconditionFailed
	case errors.As(err, &notfound):
		title = notfound.Message
		status = http.StatusNotFound
	case errors.As(err, &invalidInput):
		title = invalidInput.Message
		status = http.StatusBadRequest
	default:
		title = "Unknown error"
		status = http.StatusInternalServerError

	}
	errorResponse := NewErrorResponse(title, http.StatusText(status), err.Error())
	jsn, _ := json.Marshal(errorResponse)
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(jsn),
	}
}
