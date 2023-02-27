package http

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2023_1_Seekers/config"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func SendJson(w http.ResponseWriter, status int, dataStruct any) {
	dataJson, err := json.Marshal(dataStruct)
	if err != nil {
		log.Error("failed to marshal", err)
		return
	}

	w.Header().Set("Content-Type", config.ContentTypeJSON)
	w.WriteHeader(status)

	_, err = w.Write(dataJson)
	if err != nil {
		log.Error("failed to send", err)
	}
	return
}
