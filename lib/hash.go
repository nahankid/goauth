package lib

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// GenerateHash for password
func GenerateHash(pwd string) (hashed string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}
