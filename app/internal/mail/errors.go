package mail

import (
	"errors"
	"net/http"
)

var (
	HttpGetMethodError         = errors.New("a get request was expected")
	InvalidURL                 = errors.New("invalid url address")
	ErrFailedGetInboxMessages  = errors.New("failed to get inbox messages")
	ErrFailedGetOutboxMessages = errors.New("failed to get outbox messages")
	ErrFailedGetFolderMessages = errors.New("failed to get folder messages")
)

var MailErrors = map[error]int{
	HttpGetMethodError:         http.StatusBadRequest,
	InvalidURL:                 http.StatusNotFound,
	ErrFailedGetInboxMessages:  http.StatusBadRequest,
	ErrFailedGetOutboxMessages: http.StatusBadRequest,
	ErrFailedGetFolderMessages: http.StatusBadRequest,
}
