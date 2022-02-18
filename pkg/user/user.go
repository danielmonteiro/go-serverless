package user

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/danielmonteiro/golang-serverless/pkg/validators"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	tableName              = "golang-serverless-user"
	ErrorGetUser           = Error("error get user")
	ErrorCreateUser        = Error("error create user")
	ErrorUpdateUser        = Error("error update user")
	ErrorDeleteUser        = Error("error delete user")
	ErrorUnmarshalItem     = Error("error unmarshal item")
	ErrorMarshalItem       = Error("error marshal item")
	ErrorInvalidUserData   = Error("invalid user data")
	ErrorInvalidEmail      = Error("invalid email")
	ErrorUserAlreadyExists = Error("user already exists")
	ErrorUserDoesNotExist  = Error("user does not exist")
)

func GetUser(email string, dynamoClient dynamodbiface.DynamoDBAPI) (*User, error) {
	query := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	}

	result, err := dynamoClient.GetItem(query)
	if err != nil {
		fmt.Println("error when getting user", err)
		return nil, ErrorGetUser
	}

	user := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, user)
	if err != nil {
		fmt.Println("error when UnmarshalMap user", err)
		return nil, ErrorUnmarshalItem
	}

	if user != nil && len(user.Email) == 0 {
		return nil, ErrorUserDoesNotExist
	}

	return user, nil
}

func CreateUser(req events.APIGatewayProxyRequest, dynamoClient dynamodbiface.DynamoDBAPI) (*User, error) {
	var u User
	err := json.Unmarshal([]byte(req.Body), &u)
	if err != nil {
		fmt.Println("error when Unmarshal user", err)
		return nil, ErrorInvalidUserData
	}
	if !validators.IsEmailValid(u.Email) {
		return nil, ErrorInvalidEmail
	}

	existingUser, errGet := GetUser(u.Email, dynamoClient)
	if errGet != nil && errGet != ErrorUserDoesNotExist {
		return nil, ErrorGetUser
	}
	if existingUser != nil {
		return nil, ErrorUserAlreadyExists
	}

	user, errMarshal := dynamodbattribute.MarshalMap(u)
	if errMarshal != nil {
		fmt.Println("error when MarshalMap user", err)
		return nil, ErrorMarshalItem
	}

	input := &dynamodb.PutItemInput{
		Item:      user,
		TableName: aws.String(tableName),
	}

	_, errPut := dynamoClient.PutItem(input)
	if errPut != nil {
		fmt.Println("error when creating user", err)
		return nil, ErrorCreateUser
	}
	return &u, nil
}

func UpdateUser(req events.APIGatewayProxyRequest, dynamoClient dynamodbiface.DynamoDBAPI) (*User, error) {
	var u User
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, ErrorInvalidEmail
	}

	_, err := GetUser(u.Email, dynamoClient)
	if err != nil {
		return nil, err
	}

	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, ErrorUnmarshalItem
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynamoClient.PutItem(input)
	if err != nil {
		return nil, ErrorUpdateUser
	}
	return &u, nil
}

func DeleteUser(email string, dynamoClient dynamodbiface.DynamoDBAPI) error {
	_, errGet := GetUser(email, dynamoClient)
	if errGet != nil {
		return errGet
	}

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := dynamoClient.DeleteItem(input)
	if err != nil {
		return ErrorDeleteUser
	}

	return nil
}
