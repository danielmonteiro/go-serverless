package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/danielmonteiro/golang-serverless/pkg/user"
)

type ErrorBody struct {
	ErrorMessage *string `json:"error,omitempty"`
}

func GetUser(req events.APIGatewayProxyRequest, dynamoClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	email := req.QueryStringParameters["email"]
	u, err := user.GetUser(email, dynamoClient)
	if err != nil {
		switch err {
		case user.ErrorUserDoesNotExist:
			return apiResponse(http.StatusNotFound, ErrorBody{aws.String(err.Error())})
		default:
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}
	}
	return apiResponse(http.StatusOK, u)
}

func CreateUser(req events.APIGatewayProxyRequest, dynamoClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	user, err := user.CreateUser(req, dynamoClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, user)
}

func UpdatetUser(req events.APIGatewayProxyRequest, dynamoClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	result, err := user.UpdateUser(req, dynamoClient)
	if err != nil {
		switch err {
		case user.ErrorUserDoesNotExist:
			return apiResponse(http.StatusNotFound, ErrorBody{aws.String(err.Error())})
		default:
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}
	}
	return apiResponse(http.StatusOK, result)
}

func DeleteUser(req events.APIGatewayProxyRequest, dynamoClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	email := req.QueryStringParameters["email"]
	err := user.DeleteUser(email, dynamoClient)

	if err != nil {
		switch err {
		case user.ErrorUserDoesNotExist:
			return apiResponse(http.StatusNotFound, ErrorBody{aws.String(err.Error())})
		default:
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}
	}
	return apiResponse(http.StatusOK, nil)
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusBadRequest, nil)
}
