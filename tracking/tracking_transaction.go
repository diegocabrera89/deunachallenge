package tracking

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
)

// GetRequestId get request id reference.
func GetRequestId(ctx context.Context, request events.APIGatewayProxyRequest) string {
	return request.RequestContext.RequestID
}
