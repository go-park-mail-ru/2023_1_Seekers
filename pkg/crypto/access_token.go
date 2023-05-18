package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgErrors "github.com/pkg/errors"
	"os"
	"strconv"
	"strings"
	"time"
)

var acTokSecret []byte
var acTokTTL = time.Second * 20
var acTokSep = ":"

const defaultAccessTokenSecret = "nnPCdw3M2B1TfJhoaY2mL736p2vCUc47"

func init() {
	envSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	if envSecret == "" {
		acTokSecret = []byte(defaultAccessTokenSecret)
		return
	}
	acTokSecret = []byte(envSecret)
}

func encrypt(plaintext string) string {
	newCipher, err := aes.NewCipher(acTokSecret)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(newCipher)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		panic(err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return string(ciphertext)
}

func decrypt(ciphertext string) string {
	newCipher, err := aes.NewCipher(acTokSecret)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(newCipher)
	if err != nil {
		panic(err)
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		panic(err)
	}

	return string(plaintext)
}

func EncryptAccessToken(userID uint64) string {
	tokenExpire := time.Now().Add(acTokTTL).Unix()
	tokenData := fmt.Sprintf("%d%s%d", userID, acTokSep, tokenExpire)
	return encrypt(tokenData)
}

func DecryptAccessToken(cipherToken string) (uint64, error) {
	accessToken := decrypt(cipherToken)

	data := strings.Split(accessToken, acTokSep)
	if len(data) < 2 {
		return 0, pkgErrors.WithMessage(errors.ErrWrongAccessToken, "wrong token len after split")
	}

	expire, err := strconv.ParseInt(data[1], 10, 64)
	if err != nil || expire < time.Now().Unix() {
		return 0, pkgErrors.WithMessage(errors.ErrWrongAccessToken, "bad token time")
	}

	userID, err := strconv.ParseUint(data[0], 10, 64)
	if err != nil {
		return 0, pkgErrors.WithMessage(errors.ErrWrongAccessToken, "failed parse userID")
	}

	return userID, nil
}
