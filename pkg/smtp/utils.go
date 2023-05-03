package smtp

import (
	"bytes"
	"fmt"
	"github.com/emersion/go-message"
	"github.com/pkg/errors"
	"io"
	"strings"
)

func ParseDomain(emailAddr string) (string, error) {
	result := strings.Split(emailAddr, "@")
	if len(result) != 2 {
		return "", fmt.Errorf("invalid address | %s", emailAddr)
	}
	domain := result[1]
	return domain, nil
}

func ParseLogin(emailAddr string) (string, error) {
	result := strings.Split(emailAddr, "@")
	if len(result) != 2 {
		return "", fmt.Errorf("invalid address | %s", emailAddr)
	}
	login := result[0]
	return login, nil
}

func GetMessageBody(mailBody []byte) (string, error) {
	m, err := message.Read(bytes.NewReader(mailBody))
	if err != nil {
		return "", errors.Wrap(err, "failed to read mail body")
	}
	var messageBody string
	var htmlBody string
	if mr := m.MultipartReader(); mr != nil {
		// This is a multipart message
		for {
			p, err := mr.NextPart()
			if err != nil {
				break
			}

			t, _, _ := p.Header.ContentType()
			if t == "text/html" {
				bytesBody, err := io.ReadAll(p.Body)
				if err != nil {
					return "", errors.Wrap(err, "failed read text/html content")
				}
				htmlBody = string(bytesBody)
			} else if t == "text/plain" {
				bytesBody, err := io.ReadAll(p.Body)
				if err != nil {
					return "", errors.Wrap(err, "failed read text/plain content")
				}
				messageBody = string(bytesBody)
			}
		}
	}

	if len(messageBody) == 0 {
		messageBody = htmlBody
	}

	return messageBody, nil
}
