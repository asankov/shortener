package dynamo

import (
	"context"
	"fmt"

	"github.com/asankov/shortener/internal/users"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"golang.org/x/crypto/bcrypt"
)

var (
	usersTableName = aws.String("users")
)

const (
	emailField    = "email"
	passwordField = "password"
	rolesField    = "roles"
)

func (d *Database) CreateUser(email, password string, roles []users.Role) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	fmt.Println(string(hashedPassword))

	emailValue, err := attributevalue.Marshal(email)
	if err != nil {
		return err
	}
	hashedPasswordValue, err := attributevalue.Marshal(hashedPassword)
	if err != nil {
		return err
	}
	rolesValue, err := attributevalue.Marshal(roles)
	if err != nil {
		return err
	}

	if _, err = d.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: usersTableName,
		Item: map[string]types.AttributeValue{
			emailField:    emailValue,
			passwordField: hashedPasswordValue,
			rolesField:    rolesValue,
		},
		// TODO: does that work?
		ConditionExpression: aws.String("attribute_not_exists(email)"),
	}); err != nil {
		return err
	}
	return nil
}

func (d *Database) GetUser(email, password string) (*users.User, error) {
	out, err := d.client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: usersTableName,
		Key: map[string]types.AttributeValue{
			emailField: &types.AttributeValueMemberS{Value: email},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(out.Item) == 0 {
		return nil, users.ErrUserNotFound
	}

	hashedPassword := out.Item[passwordField].(*types.AttributeValueMemberS).Value
	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return nil, err
	}

	return &users.User{
		Email: out.Item[emailField].(*types.AttributeValueMemberS).Value,
		// TODO
		Roles: []users.Role{},
	}, nil
}

func (d *Database) ShouldCreateInitialUser() (bool, error) {
	scanOutput, err := d.client.Scan(context.Background(), &dynamodb.ScanInput{
		TableName: usersTableName,
		Limit:     aws.Int32(1),
	})
	if err != nil {
		return false, err
	}
	if scanOutput.Count == 0 {
		return true, nil
	}
	return false, nil
}
