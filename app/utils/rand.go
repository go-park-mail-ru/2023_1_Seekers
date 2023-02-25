package utils

import (
	"math/rand"
	"time"
)

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func String(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randBytes := make([]byte, length)
	for i := range randBytes {
		randBytes[i] = charset[r.Intn(len(charset))]
	}
	return string(randBytes)
}
