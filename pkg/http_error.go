package pkg

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func SendError(w http.ResponseWriter, r *http.Request, error *errors.JSONError) {
	SendJSON(w, r, error.Code, error)
}

func HandleError(w http.ResponseWriter, r *http.Request, status int, logErr error, UnwrappedErr error) {
	customErr := errors.New(status, UnwrappedErr)
	logger, ok := r.Context().Value(ContextHandlerLog).(*Logger)
	if !ok {
		log.Error("failed to get logger for handler", r.URL.Path)
		log.Error(logErr)
	} else {
		logger.Error(logErr)
	}
	SendError(w, r, customErr)
}
