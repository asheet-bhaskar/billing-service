package utils

import (
	"time"

	"math/rand"

	"github.com/google/uuid"
)

func GetNewUUID() string {
	return uuid.New().String()
}

func IsValidUUID(id string) bool {
	err := uuid.Validate(id)
	return err == nil
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
