package smtp

import (
	"bytes"
	"github.com/emersion/go-message"
	"github.com/pkg/errors"
	"mime"
	"net/mail"
	"strings"
)

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

	fromString := entity.Header.Get("From")
	addr, err := mail.ParseAddress(fromString)

	if err == nil {
		resultMessage.FromName = addr.Name
	}
	resultMessage.FromEmail = addr.Address

	msgData, err := GetMessageData(bytesMail)
	if err != nil {
		return nil, errors.Wrap(err, "failed read body")
	}

	if msgData.Subject != resultMessage.Subject {
		resultMessage.Subject = msgData.Subject
	}

	if msgData.HTMLBody == "" {
		resultMessage.HTMLBody = msgData.PlainBody
	} else {
		resultMessage.HTMLBody = msgData.HTMLBody
	}

	resultMessage.Attaches = msgData.Attaches
	resultMessage.FromName = msgData.FromName

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
