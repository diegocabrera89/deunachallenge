package helpers

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/diegocabrera89/ms-payment-core/logs"
	"reflect"
	"unicode"
)

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
		logs.LogTrackingError("IsObjectNilEmpty", "JSON Marshal", ctx, request, err)
		return false, ""
	}

	return false, string(jsonData)
}

// SkipUpdatingFields skip the fields to update.
func SkipUpdatingFields(currentField string, skipFields []string) bool {
	for _, item := range skipFields {
		if item == ToLowerCase(currentField) {
			return true
		}
	}
	return false
}

// ToLowerCase covert string to lowerCae.
func ToLowerCase(input string) string {
	if len(input) == 0 {
		return input
	}
	firstChar := string(unicode.ToLower(rune(input[0])))
	return firstChar + input[1:]
}

// BuildUpdateValues build field and values for update item from DynamoDB.
func BuildUpdateValues(object interface{}, ctx context.Context, request events.APIGatewayProxyRequest) map[string]interface{} {
	// Create map to store updated values.
	updateValues := make(map[string]interface{})

	// Obtain reflection of the Customer structure.
	objectReflect := reflect.ValueOf(object)

	// Iterate over object fields and add them to updateValues
	for i := 0; i < objectReflect.NumField(); i++ {
		field := objectReflect.Type().Field(i)
		fieldValue := objectReflect.Field(i).Interface()
		updateValues[field.Name] = fieldValue
	}

	return updateValues
}
