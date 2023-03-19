package user

import (
	"errors"
	"net/http"
)

var (
	ErrUserExists       = errors.New("such user exists")
	ErrUserNotFound     = errors.New("user not found")
	ErrTooShortPw       = errors.New("password too short")
	ErrInvalidEmail     = errors.New("invalid email address")
	ErrInvalidForm      = errors.New("invalid form")
	ErrFailedGetUser    = errors.New("failed to get user")
	ErrInternal         = errors.New("internal server error")
	ErrEmptyContentType = errors.New("content type not presented")
	ErrWrongContentType = errors.New("unsupported content type")
	ErrFailedDelete     = errors.New("failed delete user")
	ErrFailedEditInfo   = errors.New("failed edit info")
	ErrFailedEditPw     = errors.New("failed edit password")
)
var Errors = map[error]int{
	ErrInvalidForm:      http.StatusForbidden,
	ErrTooShortPw:       http.StatusForbidden,
	ErrInvalidEmail:     http.StatusUnauthorized,
	ErrUserNotFound:     http.StatusUnauthorized,
	ErrUserExists:       http.StatusConflict,
	ErrFailedGetUser:    http.StatusUnauthorized,
	ErrInternal:         http.StatusInternalServerError,
	ErrEmptyContentType: http.StatusForbidden,
	ErrWrongContentType: http.StatusForbidden,
	ErrFailedDelete:     http.StatusInternalServerError,
	ErrFailedEditInfo:   http.StatusBadRequest,
	ErrFailedEditPw:     http.StatusBadRequest,
}
