package mail

import (
	"errors"
	"net/http"
)

var (
	ErrHttpGetMethod           = errors.New("a get request was expected")
	ErrInvalidURL              = errors.New("invalid url address")
	ErrFailedGetInboxMessages  = errors.New("failed to get inbox messages")
	ErrFailedGetOutboxMessages = errors.New("failed to get outbox messages")
	ErrFailedGetFolderMessages = errors.New("failed to get folder messages")
	ErrFailedGetUser           = errors.New("failed to get user")
)

var MailErrors = map[error]int{
	ErrHttpGetMethod:           http.StatusBadRequest,
	ErrInvalidURL:              http.StatusNotFound,
	ErrFailedGetInboxMessages:  http.StatusBadRequest,
	ErrFailedGetOutboxMessages: http.StatusBadRequest,
	ErrFailedGetFolderMessages: http.StatusBadRequest,
	ErrFailedGetUser:           http.StatusBadRequest,
}
