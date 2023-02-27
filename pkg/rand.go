package pkg

import (
	"crypto/rand"
	"math/big"
)

func String(length int) (string, error) {
	charSet := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randBytes := make([]byte, length)
	for i := range randBytes {
		res, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
		if err != nil {
			return "", err
		}
		randBytes[i] = charSet[res.Int64()]
	}
	return string(randBytes), nil
}
