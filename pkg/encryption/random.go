package encryption

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomString generates a random string using crypto/rand
func GenerateRandomString(length int) (string, error) {
	data := make([]byte, length)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}
