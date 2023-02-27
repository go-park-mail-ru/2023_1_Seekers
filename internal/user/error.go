package user

import "errors"

var (
	ErrUserExists   = errors.New("such user exists")
	ErrUserNotFound = errors.New("user not found")
)
