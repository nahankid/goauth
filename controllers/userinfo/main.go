package main

import (
	"auth/db"
	"auth/lib"
	"auth/models"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	token := request.Headers["Authorization"]
	fmt.Println("token in userinfo handler", token)

	// userId := common.GetStore().Get(token, false)

	userID := []byte(token)

	fmt.Println("Userer", userID)

	// Open Database Connection
	pgConn := db.PGConn{}
	db2, err := pgConn.GetConnection()
	defer db2.Close()
	if err != nil {
		return lib.APIResponse(http.StatusInternalServerError, err.Error())
	}

	var userFilter models.User
	u, _ := binary.Uvarint(userID)
	userFilter.ID = uint(u)

	var user models.User
	db2.Where(userFilter).Find(&user)

	res, err := json.Marshal(user)
	if err != nil {
		return lib.APIResponse(http.StatusInternalServerError, err.Error())
	}

	return lib.APIResponse(http.StatusOK, string(res))
}

func main() {
	lambda.Start(handler)
}
