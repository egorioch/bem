package jwt_auth

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomSecret() (string, error) {
	key := make([]byte, 32) // 32 байта = 256 бит
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}
