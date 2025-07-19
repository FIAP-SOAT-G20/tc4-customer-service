package response

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func NewAPIGatewayProxyResponse(data []byte) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		Body:            string(data),
		Headers:         map[string]string{"Content-Type": "application/json"},
		IsBase64Encoded: false,
	}
}
