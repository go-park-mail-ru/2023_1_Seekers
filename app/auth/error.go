package auth

import "errors"

var (
	ErrInvalidMethodPost = errors.New("invalid method, not post")
	ErrInvalidForm       = errors.New("invalid form, cant decode")
	ErrInvalidPw         = errors.New("invalid password")
	ErrPwDontMatch       = errors.New("passwords dont match")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserExists        = errors.New("user already exists")
)
