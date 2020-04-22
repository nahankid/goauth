package types

import "github.com/dgrijalva/jwt-go"

// CustomClaims struct
type CustomClaims struct {
	jwt.StandardClaims
	Type string
}
