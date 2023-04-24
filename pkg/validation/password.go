package validation

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgErr "github.com/pkg/errors"
)

func Password(password string, minLen int) error {
	if len(password) < minLen {
		return pkgErr.WithMessage(errors.ErrTooShortPw, "failed validate")
	}
	return nil
}
