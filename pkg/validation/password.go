package validation

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgErr "github.com/pkg/errors"
)

func Password(password string) error {
	if len(password) < config.PasswordMinLen {
		return pkgErr.WithMessage(errors.ErrTooShortPw, "failed validate")
	}
	return nil
}
