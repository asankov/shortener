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
)

const (
	idField      = "id"
	urlField     = "url"
	metricsField = "metrics"
	clicksField  = "clicks"

	region = "eu-west-1"

	maxAllowedConflicts = 50
)

// Database represents a DynamoDB database.
type Database struct {
	client *dynamodb.Client
	random *random.Random

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

// New creates a new database will config loaded from the environment.
//
// It returns an error if not possible to do so.
func New() (*Database, error) {
	client, err := buildDynamoDBClient()
	if err != nil {
		return nil, err
	}
	return &Database{
		client: client,
		random: random.New(),
	}, nil
}

// SetLogger sets the logger used in the Database.
func (d *Database) SetLogger(l *slog.Logger) {
	d.logger = l
}

// GetByID looks up a link by ID and returns it.
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

// GetAll returns all links.
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

// Create creates a new link with the provided ID and URL.
func (d *Database) Create(id string, url string) error {
	return d.saveLink(id, url, &links.Metrics{Clicks: 0}, &saveOptions{conditionalExpression: aws.String("attribute_not_exists(id)")})
}

// Delete deletes the link with the given ID.
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

// IncrementClicks increments the clicks for the link with the given ID.
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

// GenerateID generates a new ID and ensures that it is not already in use.
func (d *Database) GenerateID() (string, error) {
	var (
		conflictCount        int = 0
		allowedConflictCount int = 4
		idLength             int = 3
	)
	for {
		if conflictCount > allowedConflictCount {
			idLength++
			allowedConflictCount *= 2
		}

		if conflictCount > maxAllowedConflicts {
			return "", links.ErrIDNotGenerated
		}

		id := d.random.ID(idLength)
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
