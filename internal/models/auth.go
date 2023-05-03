package models

import (
	"html"
)

type AuthResponse struct {
	Email     string `json:"email" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type EditPasswordRequest struct {
	PasswordOld string `json:"passwordOld" validate:"required"`
	Password    string `json:"password" validate:"required"`
	RepeatPw    string `json:"repeatPw" validate:"required"`
}

func (form *EditPasswordRequest) Sanitize() {
	form.Password = html.EscapeString(form.Password)
}
