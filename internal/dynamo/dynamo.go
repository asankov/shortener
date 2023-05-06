package dynamo

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/asankov/shortener/pkg/links"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sirupsen/logrus"
)

var (
	tableName = aws.String("links")

	idField  = "id"
	urlField = "url"

	region = "eu-west-1"
)

type Database struct {
	client *dynamodb.Client

	logger *logrus.Logger
	rand   *rand.Rand
}

func buildDynamoDBClient() (*dynamodb.Client, error) {
	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	dynamodbClient := dynamodb.NewFromConfig(awsConfig, func(opt *dynamodb.Options) {
		// opt.Region = awsConfig.Region
		opt.Region = region
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
		rand:   rand.New(rand.NewSource(time.Now().Unix())),
	}, nil
}

func (d *Database) SetLogger(l *logrus.Logger) {
	d.logger = l
}

func (d *Database) GetByID(id string) (*links.Link, error) {
	out, err := d.client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: tableName,
		Key: map[string]types.AttributeValue{
			idField: &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(out.Item) == 0 {
		return nil, links.ErrLinkNotFound
	}

	return &links.Link{
		ID:  out.Item[idField].(*types.AttributeValueMemberS).Value,
		URL: out.Item[urlField].(*types.AttributeValueMemberS).Value,
	}, nil
}

func (d *Database) GetAll() ([]*links.Link, error) {
	scanOutput, err := d.client.Scan(context.Background(), &dynamodb.ScanInput{
		TableName: tableName,
	})
	if err != nil {
		return nil, err
	}
	if len(scanOutput.LastEvaluatedKey) > 0 {
		// TODO: implement pagination
		d.logger.Warnln("there are more results")
	}
	result := make([]*links.Link, 0, len(scanOutput.Items))
	for _, item := range scanOutput.Items {
		result = append(result, &links.Link{
			ID:  item[idField].(*types.AttributeValueMemberS).Value,
			URL: item[urlField].(*types.AttributeValueMemberS).Value,
		})
	}
	return result, nil
}

func (d *Database) Create(id string, url string) error {
	idValue, err := attributevalue.Marshal(id)
	if err != nil {
		return err
	}
	urlValue, err := attributevalue.Marshal(url)
	if err != nil {
		return err
	}

	if _, err = d.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: tableName,
		Item: map[string]types.AttributeValue{
			idField:  idValue,
			urlField: urlValue,
		},
		// TODO: does that work?
		ConditionExpression: aws.String("attribute_not_exists(id)"),
	}); err != nil {
		return err
	}

	return nil
}

func (d *Database) Delete(id string) error {
	idValue, err := attributevalue.Marshal(id)
	if err != nil {
		return err
	}
	_, err = d.client.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
		TableName: tableName,
		Key: map[string]types.AttributeValue{
			idField: idValue,
		},
	})
	return err
}

func (d *Database) GenerateID() (string, error) {
	var (
		conflictCount        int
		allowedConflictCount int = 4
		idLength             int = 3

		maxAllowedConflicts = 50
	)
	for {
		if conflictCount > allowedConflictCount {
			idLength++
			allowedConflictCount *= 2
		}

		if conflictCount > maxAllowedConflicts {
			return "", links.ErrIDNotGenerated
		}

		id := d.randomID(idLength)
		// This is not optimal as DynamoDB reads are not free.
		// Probably best to substitute it with some sort of cache at some point.
		_, err := d.GetByID(id)

		// An item with this ID is not found, so we can safely use it.
		if err != nil && errors.Is(err, links.ErrLinkNotFound) {
			return id, nil
		}
		conflictCount++
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (d *Database) randomID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[d.rand.Intn(len(letterBytes))]
	}
	return string(b)
}
