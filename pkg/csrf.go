package pkg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgErrors "github.com/pkg/errors"
	"os"
	"strconv"
	"strings"
	"time"
)

var secret []byte

const DefaultCSRFSecret = "sadsa2sadl149891mppinadpon"

func init() {
	csrfSecret := os.Getenv("CSRF_SECRET")
	if csrfSecret == "" {
		secret = []byte(DefaultCSRFSecret)
		return
	}
	secret = []byte(csrfSecret)
}

func CreateMAC(data []byte) ([]byte, error) {
	h := hmac.New(sha256.New, secret)
	if _, err := h.Write(data); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func CreateCSRF(cookie string) (string, error) {
	csrfExpire := time.Now().Add(time.Minute * 10).Unix()
	tokenData := []byte(fmt.Sprintf("%s.%d", cookie, csrfExpire))
	mac, err := CreateMAC(tokenData)
	if err != nil {
		return "", err
	}
	token := hex.EncodeToString(mac) + "." + strconv.FormatInt(csrfExpire, 10)
	fmt.Println(token)
	return token, nil
}

func CheckCSRF(cookie string, csrfToken string) error {
	data := strings.Split(csrfToken, ".")
	expire, err := strconv.ParseInt(data[1], 10, 64)
	if err != nil || expire < time.Now().Unix() {
		return pkgErrors.WithMessage(errors.ErrWrongCSRF, "bad token time")
	}

	tokenData := []byte(fmt.Sprintf("%s.%d", cookie, expire))
	tokenMAC, err := hex.DecodeString(data[0])
	if err != nil {
		return err
	}

	mac, err := CreateMAC(tokenData)
	if err != nil {
		return err
	}
	if hmac.Equal(mac, tokenMAC) {
		return nil
	}
	return pkgErrors.WithMessage(errors.ErrWrongCSRF, "not equal with expected")
}
