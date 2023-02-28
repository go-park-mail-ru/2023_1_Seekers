package auth

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidMethodPost    = errors.New("invalid method, not post")
	ErrInvalidForm          = errors.New("invalid form, cant decode")
	ErrPwDontMatch          = errors.New("passwords dont match")
	ErrUserNotFound         = errors.New("user not found")
	ErrUserExists           = errors.New("user already exists")
	ErrSessionNotFound      = errors.New("session not found")
	ErrSessionExists        = errors.New("session exists")
	ErrFailedSignUp         = errors.New("failed to sign up")
	ErrFailedCreateProfile  = errors.New("failed to create profile")
	ErrFailedCreateSession  = errors.New("failed to create session")
	ErrFailedSignIn         = errors.New("failed to sign in")
	ErrFailedLogout         = errors.New("failed logout")
	ErrFailedLogoutNoCookie = errors.New("failed logout")
	ErrFailedLogoutSession  = errors.New("failed logout")
	ErrFailedAuth           = errors.New("failed auth")
	ErrFailedGetSession     = errors.New("failed get session")
)

var AuthErrors = map[error]int{
	ErrInvalidMethodPost:    http.StatusMethodNotAllowed,
	ErrInvalidForm:          http.StatusForbidden,
	ErrPwDontMatch:          http.StatusUnauthorized,
	ErrUserNotFound:         http.StatusUnauthorized,
	ErrUserExists:           http.StatusConflict,
	ErrSessionNotFound:      http.StatusUnauthorized,
	ErrSessionExists:        http.StatusConflict,
	ErrFailedSignUp:         http.StatusConflict,
	ErrFailedCreateProfile:  http.StatusBadRequest,
	ErrFailedCreateSession:  http.StatusUnauthorized,
	ErrFailedSignIn:         http.StatusUnauthorized,
	ErrFailedLogout:         http.StatusBadRequest,
	ErrFailedLogoutNoCookie: http.StatusUnauthorized,
	ErrFailedLogoutSession:  http.StatusUnauthorized,
	ErrFailedAuth:           http.StatusUnauthorized,
	ErrFailedGetSession:     http.StatusUnauthorized,
}
