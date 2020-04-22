package types

import "time"

// TokenSet struct
type TokenSet struct {
	AccessToken     string
	ExpireAt        time.Time
	RefreshToken    string
	RefreshExpireAt time.Time
}
