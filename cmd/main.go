package main

import (
	"fmt"
	"os"

	"github.com/danielmonteiro/golang-serverless/pkg/handlers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	dynamoClient dynamodbiface.DynamoDBAPI
)

func main() {
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		fmt.Println("Error when initializing new session", err)
		return
	}

	dynamoClient = dynamodb.New(awsSession)
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return handlers.GetUser(req, dynamoClient)
	case "POST":
		return handlers.CreateUser(req, dynamoClient)
	case "PUT":
		return handlers.UpdatetUser(req, dynamoClient)
	case "DELETE":
		return handlers.DeleteUser(req, dynamoClient)
	default:
		return handlers.UnhandledMethod()
	}
}
