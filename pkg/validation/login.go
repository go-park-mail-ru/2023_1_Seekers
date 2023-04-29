package validation

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"strings"
)

func Login(login, postAtDomain string) (string, error) {
	email := login
	if !strings.Contains(login, postAtDomain) {
		if strings.Contains(login, "@") || strings.Contains(login, ".") {
			return "", errors.ErrInvalidLogin
		} else {
			email += postAtDomain
		}
	} else {
		idx := strings.Index(login, postAtDomain)
		if idx+len(postAtDomain) < len(login) ||
			strings.Index(login, "@") < idx || strings.Index(login, ".") < idx {
			return "", errors.ErrInvalidLogin
		}
	}
	return email, nil
}
