package errors

import (
	"github.com/pkg/errors"
	"net/http"
)

var (
	ErrInvalidForm         = errors.New("invalid form")
	ErrPwDontMatch         = errors.New("passwords dont match")
	ErrInvalidLogin        = errors.New("invalid login")
	ErrWrongPw             = errors.New("wrong password")
	ErrUserExists          = errors.New("user already exists")
	ErrFailedGetSession    = errors.New("failed get session")
	ErrFailedDeleteSession = errors.New("failed delete session")
	ErrFailedGetUser       = errors.New("failed to get user")
	ErrInternal            = errors.New("internal server error")
	ErrTooShortPw          = errors.New("too short password")
	ErrInvalidEmail        = errors.New("invalid email address")
	ErrUserNotFound        = errors.New("user not found")
	ErrGetFile             = errors.New("failed get file")
	ErrNoKey               = errors.New("no key")
	ErrNoBucket            = errors.New("no bucket")
	ErrInvalidURL          = errors.New("invalid url address")
	ErrInvalidMessageForm  = errors.New("invalid message form")
	ErrFolderNotFound      = errors.New("folder not found")
	ErrMessageNotFound     = errors.New("message not found")
	ErrNoValidEmails       = errors.New("no valid emails")
	ErrWrongContentType    = errors.New("unsupported content type")
	ErrFailedAuth          = errors.New("failed auth")
)

var Codes = map[error]int{
	ErrInvalidForm:         http.StatusForbidden,
	ErrPwDontMatch:         http.StatusUnauthorized,
	ErrInvalidLogin:        http.StatusUnauthorized,
	ErrWrongPw:             http.StatusUnauthorized,
	ErrUserExists:          http.StatusConflict,
	ErrFailedGetSession:    http.StatusUnauthorized,
	ErrFailedDeleteSession: http.StatusUnauthorized,
	ErrInternal:            http.StatusInternalServerError,
	ErrTooShortPw:          http.StatusForbidden,
	ErrInvalidEmail:        http.StatusUnauthorized,
	ErrUserNotFound:        http.StatusNotFound,
	ErrFailedGetUser:       http.StatusUnauthorized,
	ErrGetFile:             http.StatusBadRequest,
	ErrNoKey:               http.StatusBadRequest,
	ErrNoBucket:            http.StatusBadRequest,
	ErrInvalidURL:          http.StatusBadRequest,
	ErrFailedGetUser:       http.StatusBadRequest,
	ErrInvalidMessageForm:  http.StatusBadRequest,
	ErrFolderNotFound:      http.StatusNotFound,
	ErrMessageNotFound:     http.StatusNotFound,
	ErrNoValidEmails:       http.StatusBadRequest,
	ErrWrongContentType:    http.StatusBadRequest,
	ErrFailedAuth:          http.StatusUnauthorized,
}

func Code(err error) int {
	code, ok := Codes[err]
	if !ok {
		return http.StatusInternalServerError
	}

	return code
}
