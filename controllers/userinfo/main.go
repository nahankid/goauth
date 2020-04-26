package main

import (
	"auth/db"
	"auth/lib"
	"auth/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sub := request.RequestContext.Authorizer["sub"]

	var str string
	var ok bool

	if str, ok = sub.(string); !ok || str == "" {
		log.Println("Subject in Authorizer Response is either empty or not a string")
		return lib.APIResponse(http.StatusBadRequest, "Invalid token")
	}

	userID, err := strconv.Atoi(str)
	if err != nil || userID <= 0 {
		return lib.APIResponse(http.StatusBadRequest, err.Error())
	}

	// Open Database Connection
	pgConn := db.PGConn{}
	db2, err := pgConn.GetConnection()
	defer db2.Close()
	if err != nil {
		return lib.APIResponse(http.StatusInternalServerError, err.Error())
	}

	var filter models.User
	filter.ID = uint(userID)

	var user models.User
	db2.Where(filter).Find(&user)

	res, err := json.Marshal(user)
	if err != nil {
		return lib.APIResponse(http.StatusInternalServerError, err.Error())
	}

	return lib.APIResponse(http.StatusOK, string(res))
}

func main() {
	lambda.Start(handler)
}
