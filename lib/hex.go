package lib

import (
	"crypto/rand"
	"encoding/hex"
)

// RandomHex returns a n*2 digit random hex
func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
