package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

type Item struct {
    Id      string `json:"id,omitempty"`
    Title   string `json:"title"`
    Details string `json:"details"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // Create context
    ctx := context.TODO()

    // Load AWS configuration
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        fmt.Println("Error loading AWS config:", err.Error())
        return events.APIGatewayProxyResponse{StatusCode: 500}, nil
    }

    // Create DynamoDB client
    svc := dynamodb.NewFromConfig(cfg)

    // Build the scan input parameters
    params := &dynamodb.ScanInput{
        TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
    }

    // Scan table
    result, err := svc.Scan(ctx, params)

    // Checking for errors, return error
    if err != nil {
        fmt.Println("Scan API call failed:", err.Error())
        return events.APIGatewayProxyResponse{StatusCode: 500}, nil
    }

    var itemArray []Item

    // Unmarshal the entire result Items at once
    err = attributevalue.UnmarshalListOfMaps(result.Items, &itemArray)
    if err != nil {
        fmt.Println("Got error unmarshalling:", err.Error())
        return events.APIGatewayProxyResponse{StatusCode: 500}, nil
    }

    fmt.Println("itemArray:", itemArray)

    itemArrayString, err := json.Marshal(itemArray)
    if err != nil {
        fmt.Println("Got error marshalling result:", err.Error())
        return events.APIGatewayProxyResponse{StatusCode: 500}, nil
    }

    return events.APIGatewayProxyResponse{Body: string(itemArrayString), StatusCode: 200}, nil
}

func main() {
    lambda.Start(Handler)
}
