package backend

import (
	"bytes"
	"fmt"
	"github.com/emersion/go-message"
	"github.com/emersion/go-message/mail"
	"github.com/emersion/go-msgauth/dkim"
	"github.com/emersion/go-smtp"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth"
	_mail "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	pkgSmtp "github.com/go-park-mail-ru/2023_1_Seekers/pkg/smtp"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"mime"
	"strings"
)

type SmtpBackend struct {
	cfg        *config.Config
	mailClient _mail.UseCaseI
	userClient user.UseCaseI
	authClient auth.UseCaseI
}

func NewSmtpBackend(c *config.Config, mailC _mail.UseCaseI, userC user.UseCaseI, authC auth.UseCaseI) *SmtpBackend {
	return &SmtpBackend{cfg: c, mailClient: mailC, userClient: userC, authClient: authC}
}

type Session struct {
	cfg        *config.Config
	mailClient _mail.UseCaseI
	userClient user.UseCaseI
	authClient auth.UseCaseI
	port       string
	username   string
	password   string
	isAuth     bool
	heloDomain string

	from string
	to   []string
}

func (bkd *SmtpBackend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{cfg: bkd.cfg, mailClient: bkd.mailClient, userClient: bkd.userClient, authClient: bkd.authClient}, nil
}

func (s *Session) AuthPlain(username, password string) error {
	_, _, err := s.authClient.SignIn(&models.FormLogin{
		Login:    username,
		Password: password,
	})
	if err != nil {
		return errors.Wrap(err, "failed smtp auth")
	}
	s.username = username
	s.password = password
	s.isAuth = true
	return nil
}

func (s *Session) Mail(from string, _ *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string) error {
	s.to = append(s.to, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	bytesMail, err := io.ReadAll(r)
	if err != nil {
		return errors.Wrap(err, "failed read message")
	}

	domainFrom, err := pkgSmtp.ParseDomain(s.from)
	if err != nil {
		return errors.New("failed to parse domain from:" + err.Error())
	}

	if !s.isAuth && domainFrom == s.cfg.Mail.PostDomain {
		return errors.New("auth required")
	}

	fmt.Println("here")

	var signedMail []byte
	if domainFrom == s.cfg.Mail.PostDomain {
		//2. sign DKIM
		signedMail, err = pkgSmtp.SignDKIM(bytesMail, s.cfg.Mail.PostDomain, s.cfg.SmtpServer.DkimPrivateKeyFile)
		if err != nil {
			return errors.Wrap(err, "failed to sign")
		}
	} else {
		// Verify other service
		verifications, err := dkim.Verify(bytes.NewReader(bytesMail))
		if err != nil {
			return errors.Wrap(err, "failed to verify dkim")
		}

		var isValidSignature = false
		for _, v := range verifications {
			if v.Err == nil {
				if v.Domain == domainFrom {
					isValidSignature = true
				}
			} else {
				if v.Domain == domainFrom {
					isValidSignature = false
				}
			}
		}
		if !isValidSignature {
			return errors.New("failed sign dkim")
		}
	}

	// validate recipients (anti spam)
	for _, to := range s.to {
		domainTo, err := pkgSmtp.ParseDomain(to)
		if err != nil {
			return errors.Wrap(err, "send - failed get domain")
		}

		if domainTo == s.cfg.Mail.PostDomain {
			if _, err := s.userClient.GetInfoByEmail(to); err != nil {
				return errors.Wrap(err, "failed to send mail")
			}
		}
	}

	var fromUser *models.User
	var subject, messageBody string
	if domainFrom != s.cfg.Mail.PostDomain {
		entity, err := message.Read(bytes.NewReader(bytesMail))
		if err != nil {
			return errors.Wrap(err, "failed read message")
		}
		subject = entity.Header.Get("Subject")

		decoder := &mime.WordDecoder{CharsetReader: message.CharsetReader}
		decodedSubject, err := decoder.Decode(subject)
		if err == nil {
			//if subject was utf-8
			subject = decodedSubject
		}

		messageBody, err = pkgSmtp.GetMessageBody(bytesMail)
		if err != nil {
			log.Debug("message not multipart")
			bytesBody, err := io.ReadAll(entity.Body)
			if err != nil {
				return errors.Wrap(err, "failed read body")
			}
			messageBody = string(bytesBody)
		}

		fromUser, err = s.userClient.GetByEmail(s.from)
		fmt.Println("DEBUG FROM USER")
		fmt.Println(fromUser, err)
		if err != nil {
			// trying to get info about sender, sometimes its not defined and errors can be ignored
			addr, _ := mail.ParseAddress(entity.Header.Get("From"))
			personalInfo := strings.Split(addr.Name, " ")

			var firstName, lastName string
			if len(personalInfo) >= 2 {
				firstName = personalInfo[0]
				lastName = personalInfo[0]
			}

			fromUser, err = s.userClient.Create(&models.User{
				Email:      s.from,
				Password:   "______",
				FirstName:  firstName,
				LastName:   lastName,
				IsExternal: true,
			})

			fmt.Println("DEBUG")
			fmt.Println("CREATE USER")
			fmt.Println(fromUser, err)

			if err != nil {
				return errors.Wrap(err, "failed create external user")
			}
		}
	}

	var batchRecipients []string

	for _, to := range s.to {
		// 3. dial and send
		domainTo, err := pkgSmtp.ParseDomain(to)
		if err != nil {
			return errors.Wrap(err, "send - failed get domain")
		}

		if domainTo != s.cfg.Mail.PostDomain {
			log.Debug("sending to other service ....")
			err = s.DialAndSend(signedMail, to)
			if err != nil {
				return err
			}
		} else {
			batchRecipients = append(batchRecipients, to)
		}
	}

	if batchRecipients != nil {
		log.Debug("store this letter to mailbx service ...")

		message := models.FormMessage{
			Recipients:       batchRecipients,
			Title:            subject,
			Text:             messageBody,
			ReplyToMessageID: nil,
		}
		fmt.Println("DEBUG")
		fmt.Println("SENDING LETTER")
		fmt.Println(fromUser.UserID, message)
		_, err = s.mailClient.SendMessage(fromUser.UserID, message)
		if err != nil {
			return errors.Wrap(err, "failed send message to mailbx service")
		}
	}
	return nil
}

func (s *Session) Reset() {
}

func (s *Session) Logout() error {
	s.username = ""
	s.password = ""
	return nil
}

func (s *Session) DialAndSend(email []byte, to string) error {
	var err error

	err = pkgSmtp.MxRecordSendMostPriority(to, func(sHostName string) (isSend bool, err error) {
		err = pkgSmtp.SendMailRaw(sHostName, s.cfg.SmtpServer.Port, s.heloDomain, nil, s.from, to, &email)
		if err != nil {
			return false, errors.Wrap(err, "failed to send raw mail")
		}
		return true, nil
	})

	if err != nil {
		return errors.Wrap(err, "failed dial and send")
	}

	log.Debug("success send email")
	return nil
}
