package pkg

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/rand"
	pkgErr "github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

func GetSalt() ([]byte, error) {
	salt, err := rand.String(config.PasswordSaltLen)
	if err != nil {
		return nil, err
	}

	return []byte(salt), nil
}

func Hash(salt []byte, str string) []byte {
	hash := argon2.IDKey([]byte(str), salt, 1, 1024, 4, 32)
	return append(salt, hash...)
}

func HashPw(password string) (string, error) {
	salt, err := GetSalt()
	if err != nil {
		return "", pkgErr.Wrap(err, "failed get salt")
	}

	hashPw := string(Hash(salt, password))
	return hashPw, nil
}

func ComparePw2Hash(password, hash string) bool {
	salt := hash[0:config.PasswordSaltLen]
	newHash := Hash([]byte(salt), password)
	return password == string(newHash)
}
