package smtp

import (
	"bytes"
	"github.com/emersion/go-message"
	"github.com/emersion/go-message/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/pkg/errors"
	"io"
	"mime"
	"strings"
)

type Message struct {
	FromName string
	Subject  string
	Body     string
	Attaches []models.Attachment
}

func ParseMail(bytesMail []byte) (*Message, error) {
	var resultMessage Message

	entity, err := message.Read(bytes.NewReader(bytesMail))
	if err != nil {
		return nil, errors.Wrap(err, "failed read message")
	}

	resultMessage.Subject = entity.Header.Get("Subject")

	//if subject was utf-8 encoded
	decoder := &mime.WordDecoder{CharsetReader: message.CharsetReader}
	decodedSubject, err := decoder.Decode(resultMessage.Subject)
	if err == nil {
		resultMessage.Subject = decodedSubject
	}

	messageBody, attaches, err := GetMessageBody(bytesMail)
	if err != nil {
		bytesBody, err := io.ReadAll(entity.Body)
		if err != nil {
			return nil, errors.Wrap(err, "failed read body")
		}

		messageBody = string(bytesBody)
	}

	resultMessage.Body = messageBody
	resultMessage.Attaches = attaches

	fromString := entity.Header.Get("From")
	addr, err := mail.ParseAddress(fromString)

	if err == nil {
		resultMessage.FromName = addr.Name
	}

	return &resultMessage, nil
}

func GetFirstLastNameFromAddr(addr string) (string, string) {
	personalInfo := strings.Split(addr, " ")

	var firstName, lastName string
	if len(personalInfo) >= 1 {
		firstName = personalInfo[0]
		if len(personalInfo) >= 2 {
			lastName = personalInfo[1]
		}
	}

	return firstName, lastName
}
