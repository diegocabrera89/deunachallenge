package dynamodbcore

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/diegocabrera89/ms-payment-core/helpers"
	"github.com/diegocabrera89/ms-payment-core/logs"
	"log"
)

// CoreRepository defines the interface for repository operations.
type CoreRepository interface {
	PutItemCore(ctx context.Context, request events.APIGatewayProxyRequest, item map[string]types.AttributeValue) error
	GetItemCore(ctx context.Context, request events.APIGatewayProxyRequest, fieldNameFilterByID string, fieldValueFilterByID string) (*dynamodb.GetItemOutput, error)
	DeleteItemCore(ctx context.Context, request events.APIGatewayProxyRequest, fieldNameFilterByID string, fieldValueFilterByID string) error
	UpdateItemCore(ctx context.Context, request events.APIGatewayProxyRequest, itemObject interface{}, fieldNameFilterByID string, fieldValueFilterByID string, skipFields []string) error
	GetItemByFieldCore(ctx context.Context, request events.APIGatewayProxyRequest, fieldNameFilterByID string, fieldValueFilterByID string, globalSecondaryIndex string, fieldNameFilterStatus string, fieldValueFilterStatus string) (*dynamodb.QueryOutput, error)
}

// DynamoDBRepository implements DynamoDBRepository for DynamoDB.
type DynamoDBRepository struct {
	client DynamoDBClientInterface
	table  string
}

// NewDynamoDBRepository createHandler a new DynamoDBRepository instance.
func NewDynamoDBRepository(tableName string, region string) (*DynamoDBRepository, error) {
	defaultConfig, err := config.LoadDefaultConfig(context.TODO(), func(opts *config.LoadOptions) error {
		opts.Region = region
		return nil
	})
	if err != nil {
		log.Fatal("Error when load default config", err)
		return nil, err
	}

	client := &DynamoDBClient{
		client: dynamodb.NewFromConfig(defaultConfig),
	}

	return &DynamoDBRepository{
		client: client,
		table:  tableName,
	}, nil
}

// PutItemCore put item in DynamoDB.
func (d DynamoDBRepository) PutItemCore(ctx context.Context, request events.APIGatewayProxyRequest, item map[string]types.AttributeValue) error {
	logs.LogTrackingInfo("PutItemCore", ctx, request)
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: &d.table,
	}
	logs.LogTrackingInfoData("PutItemCore input", input, ctx, request)
	_, err := d.client.PutItem(ctx, input)
	if err != nil {
		logs.LogTrackingError("CreateItemRepository", "PutItem", ctx, request, err)
		return err
	}
	return nil
}

// GetItemCore get item from DynamoDB.
func (d DynamoDBRepository) GetItemCore(ctx context.Context, request events.APIGatewayProxyRequest, fieldNameFilterByID string, fieldValueFilterByID string) (*dynamodb.GetItemOutput, error) {
	logs.LogTrackingInfo("GetItemCore", ctx, request)
	input := &dynamodb.GetItemInput{
		Key:       helpers.GetPrimaryKey(fieldNameFilterByID, fieldValueFilterByID),
		TableName: aws.String(d.table),
	}
	response, err := d.client.GetItem(context.TODO(), input)
	if err != nil {
		logs.LogTrackingError("GetItemCore", "GetItem", ctx, request, err)
		return &dynamodb.GetItemOutput{}, nil
	}
	return response, nil
}

// DeleteItemCore item from DynamoDB.
func (d DynamoDBRepository) DeleteItemCore(ctx context.Context, request events.APIGatewayProxyRequest, fieldNameFilterByID string, fieldValueFilterByID string) error {
	logs.LogTrackingInfo("DeleteItemCore", ctx, request)
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(d.table),
		Key:       helpers.GetPrimaryKey(fieldNameFilterByID, fieldValueFilterByID),
	}
	_, err := d.client.DeleteItem(context.TODO(), input)
	if err != nil {
		logs.LogTrackingError("DeleteItemCore", "DeleteItem", ctx, request, err)
		return err
	}

	return err
}

// UpdateItemCore item from DynamoDB.
func (d DynamoDBRepository) UpdateItemCore(ctx context.Context, request events.APIGatewayProxyRequest, itemObject interface{}, fieldNameFilterByID string, fieldValueFilterByID string, skipFields []string) error {
	logs.LogTrackingInfo("UpdateItemCore", ctx, request)
	updateValues := helpers.BuildUpdateValues(itemObject, ctx, request)
	updateExpression, errorBuildUpdateExpression := helpers.BuildUpdateExpression(updateValues, skipFields, ctx, request)
	if errorBuildUpdateExpression != nil {
		logs.LogTrackingError("UpdateItemCore", "BuildUpdateExpression", ctx, request, errorBuildUpdateExpression)
	}

	cond := expression.Equal(
		expression.Name(fieldNameFilterByID),
		expression.Value(fieldValueFilterByID))

	expr, errorExpression := expression.NewBuilder().WithUpdate(updateExpression).WithCondition(cond).Build()
	if errorExpression != nil {
		logs.LogTrackingError("UpdateItemCore", "expression.NewBuilder", ctx, request, errorExpression)
	}
	logs.LogTrackingInfoData("UpdateItemCore", expr, ctx, request)
	updateItemInput := &dynamodb.UpdateItemInput{
		Key:                       helpers.GetPrimaryKey(fieldNameFilterByID, fieldValueFilterByID),
		TableName:                 aws.String(d.table),
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}
	_, errorUpdateItem := d.client.UpdateItem(context.TODO(), updateItemInput)

	return errorUpdateItem
}

// GetItemByFieldCore get item from DynamoDB.
func (d DynamoDBRepository) GetItemByFieldCore(ctx context.Context, request events.APIGatewayProxyRequest, fieldNameFilterByID string, fieldValueFilterByID string, globalSecondaryIndex string, fieldNameFilterStatus string, fieldValueFilterStatus string) (*dynamodb.QueryOutput, error) {
	logs.LogTrackingInfo("GetItemByFieldCore", ctx, request)

	input := &dynamodb.QueryInput{
		IndexName:              aws.String(globalSecondaryIndex),
		TableName:              aws.String(d.table),
		KeyConditionExpression: aws.String(fieldNameFilterByID + " = :" + fieldNameFilterByID + " and " + fieldNameFilterStatus + " = :" + fieldNameFilterStatus),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":" + fieldNameFilterByID:   &types.AttributeValueMemberS{Value: fieldValueFilterByID},
			":" + fieldNameFilterStatus: &types.AttributeValueMemberS{Value: fieldValueFilterStatus},
		},
	}
	logs.LogTrackingInfoData("GetItemByFieldCore input", input, ctx, request)
	response, err := d.client.GetItemByField(context.TODO(), input)
	logs.LogTrackingInfoData("GetItemByFieldCore response", response, ctx, request) //TODO

	if err != nil {
		logs.LogTrackingError("GetItemByFieldCore", "GetItemByField", ctx, request, err)
		return &dynamodb.QueryOutput{}, nil
	}
	return response, nil
}
