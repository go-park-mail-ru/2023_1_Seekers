package pkg

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func SendError(w http.ResponseWriter, error *errors.JSONError) {
	SendJSON(w, error.Code, error)
}

func HandleError(w http.ResponseWriter, r *http.Request, status int, err error) {
	customErr := errors.New(status, err)
	logger, ok := r.Context().Value(ContextHandlerLog).(*Logger)
	if !ok {
		log.Error("failed to get logger for handler", r.URL.Path)
	} else {
		logger.Error(err)
	}
	SendError(w, customErr)
}
