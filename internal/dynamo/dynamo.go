package dynamo

import (
	"context"

	"github.com/asankov/shortener/pkg/links"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Database struct {
	client *dynamodb.Client
}

func buildDynamoDBClient() (*dynamodb.Client, error) {
	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	dynamodbClient := dynamodb.NewFromConfig(awsConfig, func(opt *dynamodb.Options) {
		// opt.Region = awsConfig.Region
		opt.Region = "eu-west-1"
	})

	return dynamodbClient, nil
}

func New() (*Database, error) {
	client, err := buildDynamoDBClient()
	if err != nil {
		return nil, err
	}
	return &Database{
		client: client,
	}, nil

}

func (d *Database) GetByID(id string) (*links.Link, error) {
	out, err := d.client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String("links"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}

	return &links.Link{
		ID:  out.Item["id"].(*types.AttributeValueMemberS).Value,
		URL: out.Item["url"].(*types.AttributeValueMemberS).Value,
	}, nil
}
func (d *Database) GetAll() ([]*links.Link, error)     { return nil, nil }
func (d *Database) Create(id string, url string) error { return nil }
func (d *Database) Delete(id string) error             { return nil }
