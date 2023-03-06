package pkg

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func HandleError(w http.ResponseWriter, status int, err error) {
	customErr := errors.New(status, err)
	log.Error(customErr)
	SendError(w, customErr)
}
