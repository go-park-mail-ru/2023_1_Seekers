package mail

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidURL         = errors.New("invalid url address")
	ErrFailedGetUser      = errors.New("failed to get user")
	ErrInvalidMessageForm = errors.New("invalid message form")
	ErrFolderNotFound     = errors.New("folder not found")
	ErrMessageNotFound    = errors.New("message not found")
	ErrNoValidEmails      = errors.New("no valid emails")
)

var mailErrors = map[error]int{
	ErrInvalidURL:         http.StatusBadRequest,
	ErrFailedGetUser:      http.StatusBadRequest,
	ErrInvalidMessageForm: http.StatusBadRequest,
	ErrFolderNotFound:     http.StatusNotFound,
	ErrMessageNotFound:    http.StatusNotFound,
	ErrNoValidEmails:      http.StatusBadRequest,
}

func GetStatusForError(err error) int {
	status := mailErrors[err]
	if status == 0 {
		return http.StatusInternalServerError
	}

	return status
}
