package utils

import (
	"net/http"
)

// SendError пока возможно и не так надо
func SendError(w http.ResponseWriter, status int, err error) {
	//w.Header().Set("Content-Type", config.NetTypeJSON)
	//w.WriteHeader(status)
	//
	//_, err = w.Write(dataJson)
	//if err != nil {
	//	return err
	//}
	SendJson(w, status, err)
}
