package validation

import (
	"net/mail"
)

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}
