package utils

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/config"
	"math/rand"
	"time"
)

func String(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randBytes := make([]byte, length)
	for i := range randBytes {
		randBytes[i] = config.CookieCharSet[r.Intn(len(config.CookieCharSet))]
	}
	return string(randBytes)
}
