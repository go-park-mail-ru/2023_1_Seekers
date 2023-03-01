package pkg

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/pkg/errors"

	//log "github.com/sirupsen/logrus"
	"net/http"
)

func SendJSON(w http.ResponseWriter, status int, dataStruct any) {
	dataJSON, err := json.Marshal(dataStruct)

	if err != nil {
		//log.Error("failed to marshal", err)
		SendError(w, errors.NewWrappedErr(http.StatusInternalServerError, "failed to marshal", err))
		return
	}

	w.Header().Set("Content-Type", config.ContentTypeJSON)
	w.WriteHeader(status)

	_, err = w.Write(dataJSON)

	if err != nil {
		//log.Error("failed to send", err)
		SendError(w, errors.NewWrappedErr(http.StatusInternalServerError, "failed to send", err))
		return
	}
}
