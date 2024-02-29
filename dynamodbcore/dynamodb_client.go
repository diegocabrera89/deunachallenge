package dynamodbcore

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DynamoDBClientInterface defines an interface for DynamoDB operations used in the repository.
type DynamoDBClientInterface interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
}

// DynamoDBClient implements the DynamoDBClientInterface interface using the actual DynamoDB client.
type DynamoDBClient struct {
	client *dynamodb.Client
}

// PutItem implements DynamoDB's PutItem operation.
func (c *DynamoDBClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return c.client.PutItem(ctx, params, optFns...)
}

// GetItem implements DynamoDB's GetItem operation.
func (c *DynamoDBClient) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return c.client.GetItem(ctx, params, optFns...)
}

// DeleteItem implements DynamoDB's DeleteItem operation.
func (c *DynamoDBClient) DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	return c.client.DeleteItem(ctx, params, optFns...)
}

// UpdateItem implements DynamoDB's DeleteItem operation.
func (c *DynamoDBClient) UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	return c.client.UpdateItem(ctx, params, optFns...)
}
