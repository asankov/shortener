package dynamo

import (
	"context"
	"errors"
	"strconv"

	"github.com/asankov/shortener/internal/links"
	"github.com/asankov/shortener/internal/random"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"golang.org/x/exp/slog"
)

var (
	tableName = aws.String("links")

	idField      = "id"
	urlField     = "url"
	metricsField = "metrics"
	clicksField  = "clicks"

	region = "eu-west-1"
)

type Database struct {
	client *dynamodb.Client

	logger *slog.Logger
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
	}, nil
}

func (d *Database) SetLogger(l *slog.Logger) {
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
		d.logger.Warn("there are more results")
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
	return d.saveLink(id, url, &links.Metrics{Clicks: 0}, &saveOptions{conditionalExpression: aws.String("attribute_not_exists(id)")})
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

func (d *Database) IncrementClicks(id string) error {
	link, err := d.GetByID(id)
	if err != nil {
		return err
	}

	return d.saveLink(link.ID, link.URL, &links.Metrics{Clicks: link.Metrics.Clicks + 1}, nil)
}

type saveOptions struct {
	conditionalExpression *string
}

func (d *Database) saveLink(id, url string, metrics *links.Metrics, opts *saveOptions) error {
	if opts == nil {
		opts = &saveOptions{}
	}

	idValue, err := attributevalue.Marshal(id)
	if err != nil {
		return err
	}
	urlValue, err := attributevalue.Marshal(url)
	if err != nil {
		return err
	}

	metricsValue := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			clicksField: &types.AttributeValueMemberN{Value: strconv.Itoa(metrics.Clicks)},
		},
	}

	putItemInput := &dynamodb.PutItemInput{
		TableName: tableName,
		Item: map[string]types.AttributeValue{
			idField:      idValue,
			urlField:     urlValue,
			metricsField: metricsValue,
		},
	}
	if opts.conditionalExpression != nil {
		putItemInput.ConditionExpression = opts.conditionalExpression
	}

	if _, err = d.client.PutItem(context.Background(), putItemInput); err != nil {
		if opts.conditionalExpression != nil {
			var ccfe *types.ConditionalCheckFailedException
			if errors.As(err, &ccfe) {
				return links.ErrLinkAlreadyExists
			}
		}
		return err
	}

	return nil
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

		id := random.ID(idLength)
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
