package http

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgErr "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	logger, ok := r.Context().Value(pkg.ContextHandlerLog).(*pkg.Logger)
	if !ok {
		log.Error("failed to get logger for handler", r.URL.Path)
		log.Error(err)
	} else {
		logger.Error(err)
	}

	causeErr := pkgErr.Cause(err)
	code := errors.Code(causeErr)
	customErr := errors.New(code, causeErr)
	SendJSON(w, r, code, customErr)
}
