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
    "github.com/google/uuid"
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
        return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error loading AWS config"}, nil
    }

    // Create DynamoDB client
    svc := dynamodb.NewFromConfig(cfg)

    // New uuid for item id
    itemUuid := uuid.New().String()
    fmt.Println("Generated new item uuid:", itemUuid)

    // Unmarshal to Item to access object properties
    itemString := request.Body
    itemStruct := Item{}
    err = json.Unmarshal([]byte(itemString), &itemStruct)
    if err != nil {
        fmt.Println("Error unmarshalling request:", err.Error())
        return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid JSON"}, nil
    }

    if itemStruct.Title == "" {
        return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Title is required"}, nil
    }

    // Create new item
    item := Item{
        Id:      itemUuid,
        Title:   itemStruct.Title,
        Details: itemStruct.Details,
    }

    // Marshal to dynamodb item
    av, err := attributevalue.MarshalMap(item)
    if err != nil {
        fmt.Println("Error marshalling item:", err.Error())
        return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error marshalling item"}, nil
    }

    tableName := os.Getenv("DYNAMODB_TABLE")
    if tableName == "" {
        fmt.Println("DYNAMODB_TABLE environment variable not set")
        return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Missing table configuration"}, nil
    }

    // Build put item input
    fmt.Printf("Putting item to table %s: %v\n", tableName, av)
    input := &dynamodb.PutItemInput{
        Item: map[string]types.AttributeValue{
            "id":      &types.AttributeValueMemberS{Value: itemUuid},
            "title":   &types.AttributeValueMemberS{Value: item.Title},
            "details": &types.AttributeValueMemberS{Value: item.Details},
        },
        TableName: aws.String(tableName),
    }

    // PutItem request
    _, err = svc.PutItem(ctx, input)

    // Checking for errors
    if err != nil {
        fmt.Printf("Got error calling PutItem: %v\n", err)
        return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error saving item to database: %v", err)}, nil
    }

    // Marshal item to return
    itemMarshalled, err := json.Marshal(item)
    if err != nil {
        fmt.Println("Error marshalling response:", err.Error())
        return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error creating response"}, nil
    }
    fmt.Println("Returning item:", string(itemMarshalled))

    // Returning response
    return events.APIGatewayProxyResponse{
        Body:       string(itemMarshalled),
        StatusCode: 200,
    }, nil
}

func main() {
    lambda.Start(Handler)
}
