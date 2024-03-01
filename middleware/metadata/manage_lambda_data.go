package metadata

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/diegocabrera89/ms-payment-core/logs"
	"github.com/diegocabrera89/ms-payment-core/tracking"
	"log"
	"os"
	"strings"
)

// InputData show input data lambda.
func InputData(ctx context.Context, request events.APIGatewayProxyRequest) {
	awsRequestID := tracking.GetRequestId(ctx, request)
	errorGetEnvironmentVariables := GetEnvironmentVariables(os.Environ(), ctx, request)
	if errorGetEnvironmentVariables != nil {
		logs.LogTrackingError("InputData", "GetEnvironmentVariables", ctx, request, errorGetEnvironmentVariables)
	}
	requestJSON, errorMarshal := json.Marshal(request)
	if errorMarshal != nil {
		logs.LogTrackingError("InputData", "errMarshal", ctx, request, errorMarshal)
	}
	log.Println(awsRequestID + " [INPUT-DATA] " + string(requestJSON))
}

// GetEnvironmentVariables get environment variables from AWS.
func GetEnvironmentVariables(variables []string, ctx context.Context, request events.APIGatewayProxyRequest) error {
	result := make(map[string]interface{})
	splitSize := 2
	for _, pair := range variables {
		parts := strings.SplitN(pair, "=", splitSize)
		if len(parts) == splitSize {
			result[parts[0]] = parts[1]
		}
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		logs.LogTrackingError("InputData", "JSON Marshal", ctx, request, err)
		return err
	}
	awsRequestID := tracking.GetRequestId(ctx, request)
	log.Println(awsRequestID + " [ENV-DATA] " + string(jsonData))
	return err
}

// MiddlewareMetadata to print the lambda's metadata before each call to the handler.
func MiddlewareMetadata(HandlerMiddleware func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		InputData(ctx, request)
		// Call to actual handling function.
		return HandlerMiddleware(ctx, request)
	}
}
