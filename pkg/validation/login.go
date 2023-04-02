package validation

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"strings"
)

func Login(login string) (string, error) {
	email := login
	if !strings.Contains(login, config.PostAtDomain) {
		if strings.Contains(login, "@") || strings.Contains(login, ".") {
			return "", errors.ErrInvalidLogin
		} else {
			email += config.PostAtDomain
		}
	} else {
		idx := strings.Index(login, config.PostAtDomain)
		if idx+len(config.PostAtDomain) < len(login) ||
			strings.Index(login, "@") < idx || strings.Index(login, ".") < idx {
			return "", errors.ErrInvalidLogin
		}
	}
	return email, nil
}
