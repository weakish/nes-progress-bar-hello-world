package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            "{\"value\": 42, \"max\": 100}",
		IsBase64Encoded: false,
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call
	lambda.Start(handler)
}
