package models

import "github.com/jinzhu/gorm"

// User model
type User struct {
	gorm.Model
	Email         string
	Password      string `json:"-"`
	EmailVerified bool
	LoginTries    int
}
