package pkg

import (
	"math/rand"
	"time"
)

func String(length int) string {
	charSet := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randBytes := make([]byte, length)
	for i := range randBytes {
		randBytes[i] = charSet[r.Intn(len(charSet))]
	}
	return string(randBytes)
}
