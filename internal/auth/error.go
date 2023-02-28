package auth

import (
	"errors"
	"net/http"
)

var (
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
	ErrFailedDeleteSession  = errors.New("failed delete session")
)

var AuthErrors = map[error]int{
	ErrInvalidForm:          http.StatusForbidden,
	ErrPwDontMatch:          http.StatusUnauthorized,
	ErrUserNotFound:         http.StatusUnauthorized,
	ErrUserExists:           http.StatusConflict,
	ErrSessionNotFound:      http.StatusUnauthorized,
	ErrSessionExists:        http.StatusConflict,
	ErrFailedSignUp:         http.StatusConflict,
	ErrFailedCreateProfile:  http.StatusInternalServerError,
	ErrFailedCreateSession:  http.StatusUnauthorized,
	ErrFailedSignIn:         http.StatusUnauthorized,
	ErrFailedLogout:         http.StatusBadRequest,
	ErrFailedLogoutNoCookie: http.StatusUnauthorized,
	ErrFailedLogoutSession:  http.StatusUnauthorized,
	ErrFailedAuth:           http.StatusUnauthorized,
	ErrFailedGetSession:     http.StatusUnauthorized,
	ErrFailedDeleteSession:  http.StatusUnauthorized,
}
