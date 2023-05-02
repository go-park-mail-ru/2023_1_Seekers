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
	ErrCantEditSentMessage  = errors.New("can't edit sent message")
)

var Codes = map[string]int{
	ErrInvalidForm.Error():          http.StatusForbidden,
	ErrPwDontMatch.Error():          http.StatusUnauthorized,
	ErrInvalidLogin.Error():         http.StatusUnauthorized,
	ErrWrongPw.Error():              http.StatusUnauthorized,
	ErrUserExists.Error():           http.StatusConflict,
	ErrFailedGetSession.Error():     http.StatusUnauthorized,
	ErrFailedDeleteSession.Error():  http.StatusUnauthorized,
	ErrInternal.Error():             http.StatusInternalServerError,
	ErrTooShortPw.Error():           http.StatusForbidden,
	ErrInvalidEmail.Error():         http.StatusUnauthorized,
	ErrUserNotFound.Error():         http.StatusNotFound,
	ErrFailedGetUser.Error():        http.StatusUnauthorized,
	ErrGetFile.Error():              http.StatusBadRequest,
	ErrNoKey.Error():                http.StatusBadRequest,
	ErrNoBucket.Error():             http.StatusBadRequest,
	ErrInvalidURL.Error():           http.StatusBadRequest,
	ErrFolderNotFound.Error():       http.StatusNotFound,
	ErrMessageNotFound.Error():      http.StatusNotFound,
	ErrNoValidEmails.Error():        http.StatusBadRequest,
	ErrWrongContentType.Error():     http.StatusBadRequest,
	ErrFailedAuth.Error():           http.StatusUnauthorized,
	ErrWrongCSRF.Error():            http.StatusBadRequest,
	ErrFolderAlreadyExists.Error():  http.StatusBadRequest,
	ErrDeleteDefaultFolder.Error():  http.StatusBadRequest,
	ErrEditDefaultFolder.Error():    http.StatusBadRequest,
	ErrInvalidFolderName.Error():    http.StatusBadRequest,
	ErrMoveToSameFolder.Error():     http.StatusBadRequest,
	ErrMoveToDraftFolder.Error():    http.StatusBadRequest,
	ErrMoveFromDraftFolder.Error():  http.StatusBadRequest,
	ErrSomeEmailsAreInvalid.Error(): http.StatusBadRequest,
}

var LogLevels = map[string]logrus.Level{
	ErrInvalidForm.Error():          logrus.WarnLevel,
	ErrPwDontMatch.Error():          logrus.WarnLevel,
	ErrInvalidLogin.Error():         logrus.WarnLevel,
	ErrWrongPw.Error():              logrus.WarnLevel,
	ErrUserExists.Error():           logrus.WarnLevel,
	ErrFailedGetSession.Error():     logrus.WarnLevel,
	ErrFailedDeleteSession.Error():  logrus.WarnLevel,
	ErrInternal.Error():             logrus.ErrorLevel,
	ErrTooShortPw.Error():           logrus.WarnLevel,
	ErrInvalidEmail.Error():         logrus.WarnLevel,
	ErrUserNotFound.Error():         logrus.WarnLevel,
	ErrFailedGetUser.Error():        logrus.WarnLevel,
	ErrGetFile.Error():              logrus.WarnLevel,
	ErrNoKey.Error():                logrus.WarnLevel,
	ErrNoBucket.Error():             logrus.ErrorLevel,
	ErrInvalidURL.Error():           logrus.WarnLevel,
	ErrFolderNotFound.Error():       logrus.WarnLevel,
	ErrMessageNotFound.Error():      logrus.WarnLevel,
	ErrNoValidEmails.Error():        logrus.WarnLevel,
	ErrWrongContentType.Error():     logrus.WarnLevel,
	ErrFailedAuth.Error():           logrus.WarnLevel,
	ErrWrongCSRF.Error():            logrus.WarnLevel,
	ErrFolderAlreadyExists.Error():  logrus.WarnLevel,
	ErrDeleteDefaultFolder.Error():  logrus.WarnLevel,
	ErrEditDefaultFolder.Error():    logrus.WarnLevel,
	ErrInvalidFolderName.Error():    logrus.WarnLevel,
	ErrMoveToSameFolder.Error():     logrus.WarnLevel,
	ErrMoveToDraftFolder.Error():    logrus.WarnLevel,
	ErrMoveFromDraftFolder.Error():  logrus.WarnLevel,
	ErrSomeEmailsAreInvalid.Error(): logrus.WarnLevel,
}

func Code(err error) int {
	code, ok := Codes[err.Error()]
	if !ok {
		return http.StatusInternalServerError
	}

	return code
}

func LogLevel(err error) logrus.Level {
	level, ok := LogLevels[err.Error()]
	if !ok {
		return logrus.ErrorLevel
	}

	return level
}
