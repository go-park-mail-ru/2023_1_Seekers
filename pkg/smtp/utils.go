package smtp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/DusanKasan/parsemail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	pkgJson "github.com/go-park-mail-ru/2023_1_Seekers/pkg/json"
	"github.com/pkg/errors"
	"io"
	"strings"
)

type Message struct {
	FromName  string
	FromEmail string
	Subject   string
	HTMLBody  string
	PlainBody string
	Attaches  []models.Attachment
}

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

func GetMessageData(mailBody []byte) (*Message, error) {
	res := &Message{}
	r := bytes.NewReader(mailBody)
	email, err := parsemail.Parse(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed parse email message")
	}

	decodedSubject, err := base64.StdEncoding.DecodeString(email.Subject)
	if err == nil {
		res.FromName = string(decodedSubject)
	} else {
		res.Subject = email.Subject
	}

	res.FromEmail = email.From[0].Address

	decodedName, err := base64.StdEncoding.DecodeString(email.From[0].Name)
	if err == nil {
		res.FromName = string(decodedName)
	} else {
		res.FromName = email.From[0].Name
	}

	htmlBody := email.HTMLBody

	decodedHtml, err := base64.StdEncoding.DecodeString(htmlBody)
	if err == nil {
		htmlBody = string(decodedHtml)
	}
	res.HTMLBody = htmlBody

	plainBody := email.TextBody
	decodedPlain, err := base64.StdEncoding.DecodeString(plainBody)
	if err != nil {
		plainBody = string(decodedPlain)
	}
	res.PlainBody = plainBody

	var attaches []models.Attachment
	for _, a := range email.Attachments {
		bytesBody, err := io.ReadAll(a.Data)
		if err != nil {
			return nil, errors.Wrap(err, "failed read attach data")
		}

		attachData := string(bytesBody)
		_, err = base64.StdEncoding.DecodeString(attachData)
		if err != nil {
			attachData = base64.StdEncoding.EncodeToString(bytesBody)
		}

		attaches = append(attaches, models.Attachment{
			FileName: a.Filename,
			FileData: attachData,
		})
	}
	res.Attaches = attaches
	res.HTMLBody = strings.ReplaceAll(res.HTMLBody, "<HTML><BODY>", "")
	res.HTMLBody = strings.ReplaceAll(res.HTMLBody, "</BODY></HTML>", "")
	res.HTMLBody, err = pkgJson.Escape(res.HTMLBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed escape html body")
	}

	res.PlainBody, err = pkgJson.Escape(res.PlainBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed escape plain text body")
	}

	return res, nil
}

//func GetMessageBody(mailBody []byte) (string, []models.Attachment, error) {
//	r := bytes.NewReader(mailBody)
//	fmt.Println("parsing mail...")
//	email, err := parsemail.Parse(r)
//	fmt.Println(email.Subject)
//	fmt.Println(email.From)
//	fmt.Println(email.To)
//	fmt.Println(email.HTMLBody)
//	fmt.Println(email.ReplyTo)
//	for _, a := range email.Attachments {
//		fmt.Println(a.Filename)
//		fmt.Println(a.ContentType)
//		a.Data
//		//and read a.Data
//	}
//	m, err := message.Read(bytes.NewReader(mailBody))
//	if err != nil {
//		return "", nil, errors.Wrap(err, "failed to read mail body")
//	}
//	var messageBody string
//	var htmlBody string
//	var attaches []models.Attachment
//	if mr := m.MultipartReader(); mr != nil {
//		// This is a multipart message
//		for {
//			p, err := mr.NextPart()
//			if err != nil {
//				break
//			}
//
//			t, _, _ := p.Header.ContentType()
//			fmt.Println(t)
//			fmt.Println(p.Header)
//
//			disp, headers, err := p.Header.ContentDisposition()
//
//			fmt.Println(disp)
//			fmt.Println(headers)
//
//			if err == nil && disp == "attachment" {
//				fmt.Println("GOT ATTACH")
//				bytesBody, err := io.ReadAll(p.Body)
//				if err != nil {
//					return "", nil, errors.Wrap(err, "failed read attach content")
//				}
//				data := base64.StdEncoding.EncodeToString(bytesBody)
//				filename := headers["filename"]
//				attaches = append(attaches, models.Attachment{
//					FileName: filename,
//					FileData: data,
//				})
//				fmt.Println("GOT ATTACH", attaches)
//			}
//
//			//if p.Header.ContentDisposition() == "attachment" {
//			//	fmt.Println("attachment")
//			//	fmt.Println(p.Header)
//			//	fmt.Println(p.Body)
//			//}
//
//			if t == "text/html" {
//				bytesBody, err := io.ReadAll(p.Body)
//				if err != nil {
//					return "", nil, errors.Wrap(err, "failed read text/html content")
//				}
//				htmlBody = string(bytesBody)
//			} else if t == "text/plain" {
//				bytesBody, err := io.ReadAll(p.Body)
//				if err != nil {
//					return "", nil, errors.Wrap(err, "failed read text/plain content")
//				}
//				messageBody = string(bytesBody)
//			}
//		}
//	} else {
//		t, _, err := m.Header.ContentType()
//		if err != nil {
//			return "", nil, errors.Wrap(err, "failed get content type of non multipart message")
//		}
//		if t == "text/plain" || t == "text/html" {
//			bytesBody, err := io.ReadAll(m.Body)
//			if err != nil {
//				return "", nil, errors.Wrap(err, "failed read non multipart message body")
//			}
//
//			messageBody = string(bytesBody)
//		}
//	}
//
//	if len(htmlBody) > 0 {
//		messageBody = htmlBody
//	}
//
//	messageBody = strings.ReplaceAll(messageBody, "<HTML><BODY>", "")
//	messageBody = strings.ReplaceAll(messageBody, "</BODY></HTML>", "")
//	body, err := pkgJson.Escape(messageBody)
//	if err != nil {
//		return "", nil, errors.Wrap(err, "failed escape body")
//	}
//
//	return body, attaches, nil
//}
