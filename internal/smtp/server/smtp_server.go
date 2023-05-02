package server

import (
	"crypto/tls"
	"github.com/emersion/go-smtp"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/smtp/server/backend"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func RunSmtpServer(cfg *config.Config, mailClient mail.UseCaseI, userClient user.UseCaseI, authClient auth.UseCaseI) error {
	smtpBackend := backend.NewSmtpBackend(cfg, mailClient, userClient, authClient)

	s := smtp.NewServer(smtpBackend)

	s.Addr = ":" + cfg.SmtpServer.Port
	s.Domain = cfg.SmtpServer.Domain
	s.ReadTimeout = cfg.SmtpServer.ReadTimeout
	s.WriteTimeout = cfg.SmtpServer.WriteTimeout
	s.MaxMessageBytes = cfg.SmtpServer.MaxMessageBytes
	s.MaxRecipients = cfg.SmtpServer.MaxRecipients
	s.AllowInsecureAuth = *cfg.SmtpServer.AllowInsecureAuth

	cert, err := tls.LoadX509KeyPair(cfg.SmtpServer.CertFile, cfg.SmtpServer.KeyFile)
	if err != nil {
		return errors.Wrap(err, "failed load tls keys")
	}
	s.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}

	log.Println("Starting server at", s.Addr)
	if err = s.ListenAndServe(); err != nil {
		errors.Wrap(err, "failed serve")
	}
	return errors.New("server stopped")
}
