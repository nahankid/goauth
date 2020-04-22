package main

import (
	"auth/db"
	"auth/lib"
	"auth/models"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/go-playground/validator.v9"
)

type signupRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate signup request
	var sr signupRequest
	err := json.Unmarshal([]byte(request.Body), &sr)
	if err != nil {
		return lib.APIResponse(http.StatusBadRequest, lib.StatusText(lib.JSONParseError))
	}

	v := validator.New()
	err = v.Struct(sr)

	if err != nil {
		return lib.APIResponse(http.StatusBadRequest, lib.StatusText(lib.JSONParseError))
	}

	// Open Database Connection
	pgConn := db.PGConn{}
	db2, err := pgConn.GetConnection()
	defer db2.Close()
	if err != nil {
		return lib.APIResponse(http.StatusInternalServerError, err.Error())
	}

	db2.AutoMigrate(&models.User{})

	// Check whether the user already exists
	var users []models.User
	filter := &models.User{}
	filter.Email = sr.Email
	db2.Where(filter).Find(&users)
	if users != nil && len(users) > 0 {
		return lib.APIResponse(http.StatusBadRequest, lib.StatusText(lib.UserAlreadyExist))
	}

	// Create the user
	newUser := models.User{Email: sr.Email, Password: lib.GenerateHash(sr.Password)}
	db2.Create(&newUser)

	// Send token via email
	// expireAt := time.Now().Add(1 * time.Hour)
	// token, err := lib.CreateToken(newUser, "Mail", expireAt)

	// var mailRequest model.SendVerificationMailRequest
	// mailRequest.UserId = newUser.ID
	// mailRequest.Token = token
	// mailRequest.Email = newUser.Email
	// emailData, _ := json.Marshal(mailRequest)
	// s := string(emailData)
	// u := string(os.Getenv("email_queue_url"))

	// sess, err := session.NewSession(&aws.Config{
	// 	Region: aws.String("eu-west-1")},
	// )
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// sqsClient := sqs.New(sess)
	// sqsClient.ServiceName = os.Getenv("email_queue")
	// input := sqs.SendMessageInput{
	// 	MessageBody: &s,
	// 	QueueUrl:    &u,
	// }
	// _, err = sqsClient.SendMessage(&input)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	res, err := json.Marshal(newUser)
	if err != nil {
		return lib.APIResponse(http.StatusInternalServerError, err.Error())
	}

	return lib.APIResponse(http.StatusCreated, string(res))
}

func main() {
	lambda.Start(handler)
}
