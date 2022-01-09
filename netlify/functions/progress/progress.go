package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Payload struct {
	Url          string `json:"url"`
	ValueKeyName string `json:"value_key_name"`
	MaxKeyName   string `json:"max_key_name"`
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	body := request.Body
	params := Payload{}
	err := json.Unmarshal([]byte(body), &params)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       "{\"error\": \"Invalid request body\"}",
		}, nil
	}
	url := params.Url
	valueKeyName := params.ValueKeyName
	maxKeyName := params.MaxKeyName

	resp, err := http.Get(url)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 502,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       fmt.Sprintf("{\"error\": \"Could not get response from upstream url: %s\"}", url),
		}, nil
	}
	defer resp.Body.Close()

	upstreamStatusCode := resp.StatusCode
	upstreamBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 502,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       "{\"error\": \"Could not read upstream response body.\"}",
		}, nil
	}
	if upstreamStatusCode != 200 {
		return &events.APIGatewayProxyResponse{
			StatusCode: 502,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       fmt.Sprintf("{\"code\": %d, \"error\": %q}", upstreamStatusCode, upstreamBody[:]),
		}, nil
	}

	var upstreamJson interface{}
	err = json.Unmarshal(upstreamBody, &upstreamJson)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 502,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       fmt.Sprintf("{\"error\": \"Unable to decoe upstream response body: %s\"}", upstreamBody[:]),
		}, nil
	}
	value, ok := upstreamJson.(map[string]interface{})[valueKeyName].(float64)
	if !ok {
		return &events.APIGatewayProxyResponse{
			StatusCode: 502,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       "{\"error\": \"Failed to parse upstream response json.\"}",
		}, nil
	}
	max, ok := upstreamJson.(map[string]interface{})[maxKeyName].(float64)
	if !ok {
		return &events.APIGatewayProxyResponse{
			StatusCode: 502,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       "{\"error\": \"Failed to parse upstream response json.\"}",
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       fmt.Sprintf("{\"value\": %f, \"max\": %f}", value, max),
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call
	lambda.Start(handler)
}
