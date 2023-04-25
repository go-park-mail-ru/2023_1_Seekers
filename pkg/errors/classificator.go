package errors

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	ErrInvalidForm          = errors.New("invalid form")
	ErrPwDontMatch          = errors.New("passwords dont match")
	ErrInvalidLogin         = errors.New("invalid login")
	ErrWrongPw              = errors.New("wrong password")
	ErrUserExists           = errors.New("user already exists")
	ErrFailedGetSession     = errors.New("failed get session")
	ErrFailedDeleteSession  = errors.New("failed delete session")
	ErrFailedGetUser        = errors.New("failed to get user")
	ErrInternal             = errors.New("internal server error")
	ErrTooShortPw           = errors.New("too short password")
	ErrInvalidEmail         = errors.New("invalid email address")
	ErrUserNotFound         = errors.New("user not found")
	ErrGetFile              = errors.New("failed get file")
	ErrNoKey                = errors.New("no key")
	ErrNoBucket             = errors.New("no bucket")
	ErrInvalidURL           = errors.New("invalid url address")
	ErrFolderNotFound       = errors.New("folder not found")
	ErrMessageNotFound      = errors.New("message not found")
	ErrNoValidEmails        = errors.New("no valid emails")
	ErrWrongContentType     = errors.New("unsupported content type")
	ErrFailedAuth           = errors.New("failed auth")
	ErrWrongCSRF            = errors.New("wrong csrf token")
	ErrFolderAlreadyExists  = errors.New("folder already exists")
	ErrDeleteDefaultFolder  = errors.New("can't delete default folder")
	ErrEditDefaultFolder    = errors.New("can't edit default folder")
	ErrInvalidFolderName    = errors.New("invalid folder name")
	ErrMoveToSameFolder     = errors.New("can't move message to same folder")
	ErrMoveToDraftFolder    = errors.New("can't move message to draft folder")
	ErrMoveFromDraftFolder  = errors.New("can't move message from draft folder")
	ErrSomeEmailsAreInvalid = errors.New("some emails are invalid")
)

var Codes = map[error]int{
	ErrInvalidForm:          http.StatusForbidden,
	ErrPwDontMatch:          http.StatusUnauthorized,
	ErrInvalidLogin:         http.StatusUnauthorized,
	ErrWrongPw:              http.StatusUnauthorized,
	ErrUserExists:           http.StatusConflict,
	ErrFailedGetSession:     http.StatusUnauthorized,
	ErrFailedDeleteSession:  http.StatusUnauthorized,
	ErrInternal:             http.StatusInternalServerError,
	ErrTooShortPw:           http.StatusForbidden,
	ErrInvalidEmail:         http.StatusUnauthorized,
	ErrUserNotFound:         http.StatusNotFound,
	ErrFailedGetUser:        http.StatusUnauthorized,
	ErrGetFile:              http.StatusBadRequest,
	ErrNoKey:                http.StatusBadRequest,
	ErrNoBucket:             http.StatusBadRequest,
	ErrInvalidURL:           http.StatusBadRequest,
	ErrFolderNotFound:       http.StatusNotFound,
	ErrMessageNotFound:      http.StatusNotFound,
	ErrNoValidEmails:        http.StatusBadRequest,
	ErrWrongContentType:     http.StatusBadRequest,
	ErrFailedAuth:           http.StatusUnauthorized,
	ErrWrongCSRF:            http.StatusBadRequest,
	ErrFolderAlreadyExists:  http.StatusBadRequest,
	ErrDeleteDefaultFolder:  http.StatusBadRequest,
	ErrEditDefaultFolder:    http.StatusBadRequest,
	ErrInvalidFolderName:    http.StatusBadRequest,
	ErrMoveToSameFolder:     http.StatusBadRequest,
	ErrMoveToDraftFolder:    http.StatusBadRequest,
	ErrMoveFromDraftFolder:  http.StatusBadRequest,
	ErrSomeEmailsAreInvalid: http.StatusBadRequest,
}

var LogLevels = map[error]logrus.Level{
	ErrInvalidForm:          logrus.WarnLevel,
	ErrPwDontMatch:          logrus.WarnLevel,
	ErrInvalidLogin:         logrus.WarnLevel,
	ErrWrongPw:              logrus.WarnLevel,
	ErrUserExists:           logrus.WarnLevel,
	ErrFailedGetSession:     logrus.WarnLevel,
	ErrFailedDeleteSession:  logrus.WarnLevel,
	ErrInternal:             logrus.ErrorLevel,
	ErrTooShortPw:           logrus.WarnLevel,
	ErrInvalidEmail:         logrus.WarnLevel,
	ErrUserNotFound:         logrus.WarnLevel,
	ErrFailedGetUser:        logrus.WarnLevel,
	ErrGetFile:              logrus.WarnLevel,
	ErrNoKey:                logrus.WarnLevel,
	ErrNoBucket:             logrus.ErrorLevel,
	ErrInvalidURL:           logrus.WarnLevel,
	ErrFolderNotFound:       logrus.WarnLevel,
	ErrMessageNotFound:      logrus.WarnLevel,
	ErrNoValidEmails:        logrus.WarnLevel,
	ErrWrongContentType:     logrus.WarnLevel,
	ErrFailedAuth:           logrus.WarnLevel,
	ErrWrongCSRF:            logrus.WarnLevel,
	ErrFolderAlreadyExists:  logrus.WarnLevel,
	ErrDeleteDefaultFolder:  logrus.WarnLevel,
	ErrEditDefaultFolder:    logrus.WarnLevel,
	ErrInvalidFolderName:    logrus.WarnLevel,
	ErrMoveToSameFolder:     logrus.WarnLevel,
	ErrMoveToDraftFolder:    logrus.WarnLevel,
	ErrMoveFromDraftFolder:  logrus.WarnLevel,
	ErrSomeEmailsAreInvalid: logrus.WarnLevel,
}

func Code(err error) int {
	code, ok := Codes[err]
	if !ok {
		return http.StatusInternalServerError
	}

	return code
}

func LogLevel(err error) logrus.Level {
	level, ok := LogLevels[err]
	if !ok {
		return logrus.ErrorLevel
	}

	return level
}
