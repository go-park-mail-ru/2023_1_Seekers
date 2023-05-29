package rand

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"path/filepath"
)

func String(length int, withUpper bool) (string, error) {
	charSet := "0123456789abcdefghijklmnopqrstuvwxyz"
	if withUpper {
		charSet += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

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

func FileName(prefix, suffix string) string {
	randBytes := make([]byte, 20)
	rand.Read(randBytes)
	return filepath.Join(prefix + hex.EncodeToString(randBytes) + suffix)
}
