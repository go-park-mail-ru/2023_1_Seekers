package crypto

import (
	"bytes"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/rand"
	pkgErr "github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

func GetSalt(saltLen int) ([]byte, error) {
	salt, err := rand.String(saltLen)
	if err != nil {
		return nil, err
	}

	return []byte(salt), nil
}

func Hash(salt []byte, str string) []byte {
	hash := argon2.IDKey([]byte(str), salt, 1, 1024, 4, 32)
	return append(salt, hash...)
}

func HashPw(password string, saultLen int) (string, error) {
	salt, err := GetSalt(saultLen)
	if err != nil {
		return "", pkgErr.Wrap(err, "failed get salt")
	}

	hashPw := string(Hash(salt, password))
	return hashPw, nil
}

func ComparePw2Hash(password, hash string, saultLen int) bool {
	if len(hash) < saultLen {
		return false
	}

	salt := hash[0:saultLen]
	newHash := Hash([]byte(salt), password)
	return bytes.Equal(newHash, []byte(hash))
}
