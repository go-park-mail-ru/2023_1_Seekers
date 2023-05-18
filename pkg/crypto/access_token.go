package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgErrors "github.com/pkg/errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var acTokSecret []byte
var acTokTTL = time.Second * 60 * 3
var acTokSep = ":"

const defaultAccessTokenSecret = "1234567890abcdef"

func init() {
	envSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	if envSecret == "" {
		acTokSecret = []byte(defaultAccessTokenSecret)
		return
	}
	acTokSecret = []byte(envSecret)
}

func encrypt(plaintext string) (string, error) {
	byteMsg := []byte(plaintext)
	block, err := aes.NewCipher(acTokSecret)
	if err != nil {
		return "", pkgErrors.Wrap(err, "could not create new cipher")
	}

	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", pkgErrors.Wrap(err, "could not encrypt")
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func decrypt(message string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", pkgErrors.Wrap(err, "could not base64 decode")
	}

	block, err := aes.NewCipher(acTokSecret)
	if err != nil {
		return "", pkgErrors.Wrap(err, "could not create new cipher")
	}

	if len(cipherText) < aes.BlockSize {
		return "", pkgErrors.Wrap(err, "invalid ciphertext block size")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

func EncryptAccessToken(userID uint64) (string, error) {
	tokenExpire := time.Now().Add(acTokTTL).Unix()
	tokenData := fmt.Sprintf("%d%s%d", userID, acTokSep, tokenExpire)
	return encrypt(tokenData)
}

func DecryptAccessToken(cipherToken string) (uint64, error) {
	accessToken, err := decrypt(cipherToken)
	if err != nil {
		return 0, pkgErrors.WithMessage(errors.ErrWrongAccessToken, "failed decrypt token: "+err.Error())
	}

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
