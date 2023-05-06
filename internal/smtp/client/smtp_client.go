package client

import (
	"bytes"
	"fmt"
	"github.com/emersion/go-message/mail"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	pkgSmtp "github.com/go-park-mail-ru/2023_1_Seekers/pkg/smtp"
	"github.com/pkg/errors"
	"io"
	"strings"
	"time"
)

func Address2Slice(addrs []*mail.Address) []string {
	result := make([]string, len(addrs))
	for _, v := range addrs {
		result = append(result, v.Address)
	}
	return result
}

func SendMail(from *models.User, to, subject, message, smtpDomain, secret string) error {
	var b bytes.Buffer

	login, err := pkgSmtp.ParseLogin(from.Email)
	if err != nil {
		return errors.Wrap(err, "failed get login from email")
	}

	auth := sasl.NewPlainClient("", login, secret)
	fmt.Println(login, ":", secret)

	addrFrom := []*mail.Address{{from.FirstName + " " + from.LastName, from.Email}}
	addrTo := []*mail.Address{{Address: to}}

	// Create our mail header
	var h mail.Header
	h.SetDate(time.Now())
	h.SetAddressList("From", addrFrom)
	h.SetAddressList("To", addrTo)
	h.SetSubject(subject)

	mailWriter, err := mail.CreateWriter(&b, h)
	if err != nil {
		return errors.Wrap(err, "failed create mail writer")
	}

	// Create a text part
	textWriter, err := mailWriter.CreateInline()
	if err != nil {
		return errors.Wrap(err, "failed create text writer")
	}

	//var textHeader mail.InlineHeader
	//textHeader.SetContentType("text/plain", nil)
	//w, err := textWriter.CreatePart(textHeader)
	//if err != nil {
	//	return errors.Wrap(err, "failed create part of message")
	//}
	//io.WriteString(w, "Message text")
	//
	//w.Close()

	var textHeaderHtml mail.InlineHeader
	textHeaderHtml.SetContentType("text/html", nil)
	partWriter, err := textWriter.CreatePart(textHeaderHtml)
	if err != nil {
		return errors.Wrap(err, "failed create part of message")
	}
	//пока без html поэтому обрамим сообщение в заголовок
	io.WriteString(partWriter, "<h1>"+message+"</h1>")

	partWriter.Close()
	textWriter.Close()

	// Create an attachment
	var attachHeader mail.AttachmentHeader
	attachHeader.Set("Content-Type", "plain/text; charset=utf-8")
	attachHeader.SetFilename("mailbox.txt")
	attachWriter, err := mailWriter.CreateAttachment(attachHeader)
	if err != nil {
		return errors.Wrap(err, "failed create attach")
	}

	attachWriter.Write([]byte("Здравствуйте, наш почтовый сервис работает в тестовом режиме, данное сообщение пришло с целью тестирования вложений.\nХорошего дня!"))

	attachWriter.Close()
	mailWriter.Close()

	err = smtp.SendMail(smtpDomain+":25", auth, from.Email, Address2Slice(addrTo), strings.NewReader(b.String()))
	if err != nil {
		errors.Wrap(err, "failed to send mail")
	}

	return nil
}
