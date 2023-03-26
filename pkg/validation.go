package pkg

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"strings"
)

func ValidateLogin(login string) (string, error) {
	email := login
	if !strings.Contains(login, config.PostAtDomain) {
		if strings.Contains(login, "@") || strings.Contains(login, ".") {
			return "", auth.ErrInvalidLogin
		} else {
			email += config.PostAtDomain
		}
	} else {
		idx := strings.Index(login, config.PostAtDomain)
		if idx+len(config.PostAtDomain) < len(login) ||
			strings.Index(login, "@") < idx || strings.Index(login, ".") < idx {
			return "", auth.ErrInvalidLogin
		}
	}
	return email, nil
}

//func ValidateMessage(message *models.FormMessage) error {
//	if len(message.Recipients) == 0 {
//		return mail.ErrInvalidRecipient
//	}
//
//	for _, login := range message.Recipients {
//		_, err := ValidateLogin(login)
//
//		if err != nil {
//			return mail.ErrInvalidRecipient
//		}
//	}
//
//	return nil
//}

func CheckImageContentType(contentType string) bool {
	if contentType == ContentTypePNG || contentType == ContentTypeSVG ||
		contentType == ContentTypeWEBP || contentType == ContentTypeJPEG {
		return true
	}
	return false
}
