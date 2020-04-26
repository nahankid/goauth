package main

import (
	"auth/lib"
	"auth/types"
	"context"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Help function to generate an IAM policy
func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}

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

	// Optional output with custom properties of the String, Number or Boolean type.
	authResponse.Context = map[string]interface{}{
		"sub": principalID,
	}
	return authResponse
}

func handleRequest(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := event.AuthorizationToken
	parse, e := lib.ValidateToken(token)

	if e != nil || !parse.Valid {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}

	claims := parse.Claims.(*types.CustomClaims)
	log.Println("Claims", claims)

	return generatePolicy(claims.Subject, "Allow", event.MethodArn), nil
}

func main() {
	lambda.Start(handleRequest)
}
