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
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

    // Getting id from path parameters
    pathParamId := request.PathParameters["id"]

    fmt.Println("Derived pathParamId from path params:", pathParamId)

    // GetItem request
    result, err := svc.GetItem(ctx, &dynamodb.GetItemInput{
        TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: pathParamId},
        },
    })

    // Checking for errors, return error
    if err != nil {
        fmt.Println(err.Error())
        return events.APIGatewayProxyResponse{StatusCode: 500}, nil
    }

    // Checking if item exists
    if result.Item == nil || len(result.Item) == 0 {
        return events.APIGatewayProxyResponse{StatusCode: 404}, nil
    }

    // Created item of type Item
    item := Item{}

    // Unmarshal result.Item into item
    err = attributevalue.UnmarshalMap(result.Item, &item)

    if err != nil {
        fmt.Printf("Failed to UnmarshalMap result.Item: %v\n", err)
        return events.APIGatewayProxyResponse{StatusCode: 500}, nil
    }

    // Marshal to JSON
    marshalledItem, err := json.Marshal(item)
    if err != nil {
        return events.APIGatewayProxyResponse{StatusCode: 500}, nil
    }

    // Return marshalled item
    return events.APIGatewayProxyResponse{Body: string(marshalledItem), StatusCode: 200}, nil
}

func main() {
    lambda.Start(Handler)
}
