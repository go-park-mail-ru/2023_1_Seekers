package validation

import (
	pkgErrors "github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/pkg/errors"
	"strings"
)

func Login(login, postAtDomain string) (string, error) {
	email := login
	if len(login) > 30 || len(login) < 3 {
		return "", errors.WithMessage(pkgErrors.ErrInvalidLogin, "Invalid login len")
	}

	if !strings.Contains(login, postAtDomain) {
		email += postAtDomain
	} else {
		idx := strings.Index(login, postAtDomain)
		if idx+len(postAtDomain) < len(login) ||
			strings.Index(login, "@") < idx {
			return "", pkgErrors.ErrInvalidLogin
		}
	}

	if err := ValidateEmail(email); err != nil {
		return "", err
	}
	return email, nil
}
