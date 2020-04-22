package main

import (
	"auth/db"
	"auth/lib"
	"auth/models"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Validate Login Request
	var user models.User
	var loginRequest LoginRequest
	err := json.Unmarshal([]byte(request.Body), &loginRequest)
	if err != nil {
		return lib.APIResponse(http.StatusBadRequest, err.Error())
	}

	v := validator.New()
	validateErr := v.Struct(loginRequest)

	if validateErr != nil {
		return lib.APIResponse(http.StatusBadRequest, err.Error())
	}

	// Open Database Connection
	pgConn := db.PGConn{}
	db2, err := pgConn.GetConnection()
	defer db2.Close()
	if err != nil {
		return lib.APIResponse(http.StatusInternalServerError, err.Error())
	}

	filter := &models.User{Email: loginRequest.Email}
	db2.Where(filter).First(&user)

	if user.ID <= 0 {
		return lib.APIResponse(http.StatusUnauthorized, lib.StatusText(lib.UserNameOrPasswordWrong))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))

	if err != nil {
		res, _ := lib.APIResponse(http.StatusUnauthorized, lib.StatusText(lib.UserNameOrPasswordWrong))
		return loginFailResponse(res, user, db2)
	}

	tokenSet, err := lib.CreateTokens(user)

	if err != nil {
		res, _ := lib.APIResponse(http.StatusInternalServerError, lib.StatusText(lib.DatabaseError))
		return loginFailResponse(res, user, db2)
	}

	user.LoginTries = 0
	db2.Save(&user)

	loginResponse := LoginResponse{AccessToken: tokenSet.AccessToken, RefreshToken: tokenSet.RefreshToken}

	res, err := json.Marshal(loginResponse)
	if err != nil {
		return lib.APIResponse(http.StatusInternalServerError, err.Error())
	}

	return lib.APIResponse(http.StatusOK, string(res))
}

func loginFailResponse(response events.APIGatewayProxyResponse, user models.User, dbConn *gorm.DB) (events.APIGatewayProxyResponse, error) {
	if user.ID > 0 {
		user.LoginTries++
		dbConn.Save(user)
		if user.LoginTries >= 5 {
			return lib.APIResponse(http.StatusUnauthorized, lib.StatusText(lib.CaptchaNeeded))
		}
	}
	return response, nil
}

func main() {
	lambda.Start(handler)
}
