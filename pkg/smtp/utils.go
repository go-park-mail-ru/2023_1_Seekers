package smtp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/emersion/go-message"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	pkgJson "github.com/go-park-mail-ru/2023_1_Seekers/pkg/json"
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

func GetMessageBody(mailBody []byte) (string, []models.Attachment, error) {
	fmt.Println(string(mailBody))
	m, err := message.Read(bytes.NewReader(mailBody))
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to read mail body")
	}
	var messageBody string
	var htmlBody string
	var attaches []models.Attachment
	if mr := m.MultipartReader(); mr != nil {
		// This is a multipart message
		for {
			p, err := mr.NextPart()
			if err != nil {
				break
			}

			t, _, _ := p.Header.ContentType()
			fmt.Println(t)
			fmt.Println(p.Header)

			disp, headers, err := p.Header.ContentDisposition()

			fmt.Println(disp)
			fmt.Println(headers)

			if err == nil && disp == "attachment" {
				fmt.Println("GOT ATTACH")
				bytesBody, err := io.ReadAll(p.Body)
				if err != nil {
					return "", nil, errors.Wrap(err, "failed read attach content")
				}
				data := base64.StdEncoding.EncodeToString(bytesBody)
				filename := headers["filename"]
				attaches = append(attaches, models.Attachment{
					FileName: filename,
					FileData: data,
				})
				fmt.Println("GOT ATTACH", attaches)
			}

			//if p.Header.ContentDisposition() == "attachment" {
			//	fmt.Println("attachment")
			//	fmt.Println(p.Header)
			//	fmt.Println(p.Body)
			//}

			if t == "text/html" {
				bytesBody, err := io.ReadAll(p.Body)
				if err != nil {
					return "", nil, errors.Wrap(err, "failed read text/html content")
				}
				htmlBody = string(bytesBody)
			} else if t == "text/plain" {
				bytesBody, err := io.ReadAll(p.Body)
				if err != nil {
					return "", nil, errors.Wrap(err, "failed read text/plain content")
				}
				messageBody = string(bytesBody)
			}
		}
	} else {
		t, _, err := m.Header.ContentType()
		if err != nil {
			return "", nil, errors.Wrap(err, "failed get content type of non multipart message")
		}
		if t == "text/plain" || t == "text/html" {
			bytesBody, err := io.ReadAll(m.Body)
			if err != nil {
				return "", nil, errors.Wrap(err, "failed read non multipart message body")
			}

			messageBody = string(bytesBody)
		}
	}

	if len(htmlBody) > 0 {
		messageBody = htmlBody
	}

	messageBody = strings.ReplaceAll(messageBody, "<HTML><BODY>", "")
	messageBody = strings.ReplaceAll(messageBody, "</BODY></HTML>", "")
	body, err := pkgJson.Escape(messageBody)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed escape body")
	}

	return body, attaches, nil
}
