package http

import (
	"net/http"
)

// SendError пока возможно и не так надо
func SendError(w http.ResponseWriter, status int, err error) {
	SendJson(w, status, err.Error())
}
