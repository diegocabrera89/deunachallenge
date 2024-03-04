package helpers

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/diegocabrera89/ms-payment-core/logs"
)

// MarshallItem converts an object to DynamoDB format.
func MarshallItem(item interface{}) (map[string]types.AttributeValue, error) {
	// Converts any type of object to a DynamoDB attribute map.
	attributeMap, err := attributevalue.MarshalMap(item)
	if err != nil {
		return nil, err
	}
	return attributeMap, nil
}

// UnmarshalMapToType is a generic function that deserializes a map to a given type.
func UnmarshalMapToType(inputMap map[string]types.AttributeValue, outputType interface{}) error {
	return attributevalue.UnmarshalMap(inputMap, outputType)
}

// UnmarshalListOfMaps is a generic function that deserializes a map to a given type.
func UnmarshalListOfMaps(responseValidateMerchant *dynamodb.QueryOutput, outputType dynamodb.QueryOutput) error {
	//Create a map slice to store the results of the query.
	items := make([]map[string]types.AttributeValue, len(responseValidateMerchant.Items))
	for i, item := range responseValidateMerchant.Items {
		items[i] = item
	}
	return attributevalue.UnmarshalListOfMaps(items, outputType)
}

// GetPrimaryKey get primary key value from DynamoDB.
func GetPrimaryKey(namePrimaryKey, valuePrimaryKey string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		namePrimaryKey: &types.AttributeValueMemberS{Value: valuePrimaryKey},
	}
}

// BuildUpdateExpression build expression for update item from DynamoDB.
func BuildUpdateExpression(updateValues map[string]interface{}, skipFields []string, ctx context.Context, request events.APIGatewayProxyRequest) (expression.UpdateBuilder, error) {
	logs.LogTrackingInfo("UpdateCustomerRepository BuildUpdateExpression", ctx, request)
	updateBuilder := expression.UpdateBuilder{}

	for fieldName, value := range updateValues {
		if !SkipUpdatingFields(fieldName, skipFields) {
			updateBuilder = updateBuilder.Set(expression.Name(ToLowerCase(fieldName)), expression.Value(value))
		}
	}
	return updateBuilder, nil
}
