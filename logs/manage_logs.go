package logs

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/diegocabrera89/ms-payment-core/constantscore"
	"github.com/diegocabrera89/ms-payment-core/tracking"
	"log"
	"reflect"
)

// LogTrackingInfo print log info with format.
func LogTrackingInfo(nameFunction string, ctx context.Context, request events.APIGatewayProxyRequest) {
	awsRequestID := tracking.GetRequestId(ctx, request)
	log.Println(awsRequestID + constantscore.LogInfo + nameFunction)
}

// LogTrackingInfoData print log info with format.
func LogTrackingInfoData(nameFunction string, object interface{}, ctx context.Context, request events.APIGatewayProxyRequest) {
	awsRequestID := tracking.GetRequestId(ctx, request)
	isNilEmpty, objectFormat := IsObjectNilEmpty(object, ctx, request)
	if isNilEmpty {
		log.Println(awsRequestID + constantscore.LogInfo + nameFunction)
	} else {
		log.Println(awsRequestID+constantscore.LogInfo+nameFunction, objectFormat)
	}
}

// LogTrackingError print log error with format.
func LogTrackingError(nameFunction string, causeMsg string, ctx context.Context, request events.APIGatewayProxyRequest, err error) {
	awsRequestID := tracking.GetRequestId(ctx, request)
	if len(causeMsg) != 0 {
		log.Println(awsRequestID+constantscore.LogError+nameFunction+" "+causeMsg, err)
	} else {
		log.Println(awsRequestID+constantscore.LogError+nameFunction, err)
	}
}

// IsObjectNilEmpty validate if objet is nil or empty.
func IsObjectNilEmpty(objet interface{}, ctx context.Context, request events.APIGatewayProxyRequest) (bool, string) {
	if objet == nil {
		return true, ""
	}

	valor := reflect.ValueOf(objet)
	switch valor.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		if valor.Len() == 0 {
			return true, ""
		}
	case reflect.Ptr:
		if valor.IsNil() {
			return true, ""
		}
	}

	// Convert to JSON
	jsonData, err := json.Marshal(objet)
	if err != nil {
		LogTrackingError("IsObjectNilEmpty", "JSON Marshal", ctx, request, err)
		return false, ""
	}

	return false, string(jsonData)
}
