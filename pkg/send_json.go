package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/config"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func SendJson(w http.ResponseWriter, status int, dataStruct any) error {
	dataJson, err := json.Marshal(dataStruct)
	if err != nil {
		log.Error("failed to marshal", err)
		return fmt.Errorf("failed to marshal %w", err)
	}

	w.Header().Set("Content-Type", config.ContentTypeJSON)
	w.WriteHeader(status)

	_, err = w.Write(dataJson)
	if err != nil {
		log.Error("failed to send", err)
		return fmt.Errorf("failed to send json %w", err)
	}
	return nil
}
