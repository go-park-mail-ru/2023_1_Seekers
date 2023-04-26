package grpc

import (
	"context"
	pkgErr "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

func CauseError(err error) error {
	s, _ := status.FromError(pkgErr.Cause(err))

	return pkgErr.New(s.Message())
}

func HandleError(ctx context.Context, err error) error {
	causeErr := pkgErr.Cause(err)

	log.Error(err)
	// TODO прокидывать логгер в микросервисы и правильно логировать
	//logLevel := errors.LogLevel(causeErr)
	//globalLogger, ok := r.Context().Value(common.ContextHandlerLog).(*logger.Logger)
	//if !ok {
	//	log.Error("failed to get logger for handler", r.URL.Path)
	//	log.Error(err)
	//} else {
	//	globalLogger.Log(logLevel, err)
	//}

	return causeErr
}
