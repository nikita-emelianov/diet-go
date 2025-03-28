package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	// Get the token from the Authorization header
	tokenString := event.AuthorizationToken

	// Check if token exists and has correct format
	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
		fmt.Println("Token missing or invalid format")
		return generatePolicy("user", "Deny", event.MethodArn), nil
	}

	// Extract token value
	token := strings.TrimPrefix(tokenString, "Bearer ")

	// Compare with expected token from environment variable
	expectedToken := os.Getenv("AUTH_TOKEN")
	if expectedToken == "" {
		fmt.Println("Warning: AUTH_TOKEN environment variable not set")
		return generatePolicy("user", "Deny", event.MethodArn), nil
	}

	// Validate token
	if token == expectedToken {
		return generatePolicy("user", "Allow", event.MethodArn), nil
	}

	return generatePolicy("user", "Deny", event.MethodArn), nil
}

func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	// Generate IAM policy
	authResponse := events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: principalID,
	}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}

	return authResponse
}

func main() {
	lambda.Start(Handler)
}
