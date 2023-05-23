package client

import (
	"bytes"
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
	for i, v := range addrs {
		result[i] = v.Address
	}
	return result
}

func SendMail(from *models.User, to string, message *models.MessageInfo, smtpDomain, secret string) error {
	var b bytes.Buffer

	login, err := pkgSmtp.ParseLogin(from.Email)
	if err != nil {
		return errors.Wrap(err, "failed get login from email")
	}

	auth := sasl.NewPlainClient("", login, secret)

	addrFrom := []*mail.Address{{from.FirstName + " " + from.LastName, from.Email}}
	addrTo := []*mail.Address{{Address: to}}

	// Create our mail header
	var h mail.Header
	h.SetDate(time.Now())
	h.SetAddressList("From", addrFrom)
	h.SetAddressList("To", addrTo)
	h.SetSubject(message.Title)

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

	io.WriteString(partWriter, message.Text)

	partWriter.Close()
	textWriter.Close()

	for _, attach := range message.Attachments {
		// Create an attachment
		var attachHeader mail.AttachmentHeader

		attachHeader.Set("Content-Type", attach.Type)
		//attachHeader.Set("Content-Transfer-Encoding", "base64")
		attachHeader.SetFilename(attach.FileName)
		attachWriter, err := mailWriter.CreateAttachment(attachHeader)
		if err != nil {
			return errors.Wrap(err, "failed create attach")
		}

		attachWriter.Write(attach.FileData)
		attachWriter.Close()
	}

	mailWriter.Close()

	err = smtp.SendMail(smtpDomain+"localhost:25", auth, from.Email, Address2Slice(addrTo), strings.NewReader(b.String()))
	if err != nil {
		errors.Wrap(err, "failed to send mail")
	}

	return nil
}
