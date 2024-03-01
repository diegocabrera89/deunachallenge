package response

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/diegocabrera89/ms-payment-core/models"
)

// SuccessResponse returns a success HTTP response.
func SuccessResponse(statusCode int, data []byte, message string) (events.APIGatewayProxyResponse, error) {
	response := models.APIResponse{
		Status:  statusCode,
		Message: message,
		Data:    json.RawMessage(data),
	}

	return buildResponse(response)
}

// ErrorResponse returns an error HTTP response.
func ErrorResponse(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	response := models.APIResponse{
		Status:  statusCode,
		Message: message,
	}

	return buildResponse(response)
}

// buildResponse constructs an HTTP response from the APIResponse structure.
func buildResponse(response models.APIResponse) (events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: response.Status,
		Body:       string(body),
	}, nil
}
