package utils

import (
	"crypto/rand"

	"github.com/hegonal/hegonal-backend/pkg/configs"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateSessionString() (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!@#$%^&*()_+-=<>?[]{}|~"
	lettersLen := len(letters)
	bytes, err := GenerateRandomBytes(configs.SessionLength)
	if err != nil {
		return "", err
	}

	result := make([]byte, configs.SessionLength)
	for i, b := range bytes {
		result[i] = letters[int(b)%lettersLen]
	}
	return string(result), nil
}