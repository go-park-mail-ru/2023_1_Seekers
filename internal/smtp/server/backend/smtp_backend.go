package backend

import (
	"bytes"
	"fmt"
	"github.com/emersion/go-smtp"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth"
	_mail "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
	pkgSmtp "github.com/go-park-mail-ru/2023_1_Seekers/pkg/smtp"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/validation"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
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
	if password != s.cfg.SmtpServer.SecretPassword {
		_, _, err := s.authClient.SignIn(&models.FormLogin{
			Login:    username,
			Password: password,
		})
		if err != nil {
			return errors.Wrap(err, "failed smtp auth")
		}
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
	if err := validation.ValidateEmail(to); err == nil {
		s.to = append(s.to, to)
	}

	return nil
}

func (s *Session) Data(r io.Reader) error {
	bytesMail, err := io.ReadAll(r)
	fmt.Println(string(bytesMail))
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

	var signedMail []byte
	if domainFrom == s.cfg.Mail.PostDomain {
		//2. sign DKIM
		signedMail, err = pkgSmtp.SignDKIM(bytesMail, s.cfg.Mail.PostDomain, s.cfg.SmtpServer.DkimPrivateKeyFile)
		if err != nil {
			return errors.Wrap(err, "failed to sign")
		}
	} else {
		// Verify other service
		if err := pkgSmtp.VerifyDKIM(bytes.NewReader(bytesMail), domainFrom); err != nil {
			return errors.Wrap(err, "smtp - DATA")
		}
	}

	invalidLocalRecipients := make([]string, 0)
	// validate recipients (anti spam), remove invalid domains from recipients
	for _, to := range s.to {
		domainTo, err := pkgSmtp.ParseDomain(to)
		if err != nil {
			return errors.Wrap(err, "send - failed get domain")
		}

		if domainTo == s.cfg.Mail.PostDomain {
			if _, err := s.userClient.GetInfoByEmail(to); err != nil {
				invalidLocalRecipients = append(invalidLocalRecipients, to)
			}
		}
	}

	// Delete invalid recipients
	for _, v := range invalidLocalRecipients {
		s.to, _ = common.RemoveFromSlice(s.to, v)
	}

	var fromUser *models.User
	var msgData *pkgSmtp.Message

	if domainFrom != s.cfg.Mail.PostDomain {
		msgData, err = pkgSmtp.ParseMail(bytesMail)
		if err != nil {
			return errors.Wrap(err, "failed parse mail")
		}

		fromUser, err = s.userClient.GetByEmail(s.from)
		if err != nil {
			// Create External user
			firstName, lastName := pkgSmtp.GetFirstLastNameFromAddr(msgData.FromName)

			fromUser, err = s.userClient.Create(&models.User{
				Email:      s.from,
				Password:   s.cfg.UserService.ExternalUserPassword,
				FirstName:  firstName,
				LastName:   lastName,
				IsExternal: true,
			})
			if err != nil {
				return errors.Wrap(err, "smtp send message : create external user")
			}

			_, err = s.mailClient.CreateDefaultFolders(fromUser.UserID)
			if err != nil {
				return errors.Wrap(err, "smtp send message : create default folders for external user")
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
			err = s.DialAndSend(signedMail, to)
			if err != nil {
				return errors.Wrap(err, "send - dial and send")
			}
		} else {
			batchRecipients = append(batchRecipients, to)
		}
	}

	if batchRecipients != nil {
		message := models.FormMessage{
			FromUser:         fromUser.Email,
			Recipients:       batchRecipients,
			Title:            msgData.Subject,
			Text:             msgData.HTMLBody,
			ReplyToMessageID: nil,
			Attachments:      msgData.Attaches,
		}

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
	err := pkgSmtp.MxRecordSendMostPriority(to, func(sHostName string) (isSend bool, err error) {
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
