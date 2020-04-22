package main

// LoginRequest holds login request params
type LoginRequest struct {
	Email           string `validate:"required,email"`
	Password        string `validate:"required"`
	CaptchaID       string
	CaptchaResponse string
}
