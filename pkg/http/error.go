package http

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/logger"
	pkgErr "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"syscall"
)

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	causeErr := pkgErr.Cause(err)
	code := errors.HttpCode(causeErr)
	customErr := errors.New(code, causeErr)
	logLevel := errors.LogLevel(causeErr)

	globalLogger, ok := r.Context().Value(common.ContextHandlerLog).(*logger.Logger)
	if !ok {
		log.Error("failed to get logger for handler", r.URL.Path)
		log.Error(err)
	} else {
		globalLogger.Log(logLevel, err)
	}

	if pkgErr.Is(causeErr, syscall.EPIPE) || pkgErr.Is(causeErr, http.ErrHandlerTimeout) {
		return
	}

	SendJSON(w, r, code, customErr)
}
