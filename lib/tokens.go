package lib

import (
	"auth/models"
	"auth/types"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// CreateTokens returns tokenset for user
func CreateTokens(user models.User) (types.TokenSet, error) {
	accessTokenExpireAt := time.Now().Add(1 * time.Hour)
	tokenStr, signErr := CreateToken(user, "Access", accessTokenExpireAt)

	if signErr != nil {
		return types.TokenSet{}, signErr
	}

	refreshTokenExpireAt := time.Now().Add(24 * time.Hour)
	refreshTokenStr, signErr := CreateToken(user, "Refresh", refreshTokenExpireAt)

	if signErr != nil {
		return types.TokenSet{}, signErr
	}
	return types.TokenSet{AccessToken: tokenStr, ExpireAt: accessTokenExpireAt, RefreshToken: refreshTokenStr, RefreshExpireAt: refreshTokenExpireAt}, nil
}

// ValidateToken validates token
func ValidateToken(token string) (*jwt.Token, error) {

	return jwt.ParseWithClaims(token, &types.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("jwt_key")), nil
	})
}

// CreateToken creates token
func CreateToken(user models.User, tokenType string, expireTime time.Time) (string, error) {
	var claim types.CustomClaims
	claim.Id = string(user.ID)
	claim.Type = tokenType
	expiresAt := expireTime
	claim.ExpiresAt = expiresAt.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	jwtKey := os.Getenv("jwt_key")
	tokenStr, signErr := token.SignedString([]byte(jwtKey))
	return tokenStr, signErr
}
