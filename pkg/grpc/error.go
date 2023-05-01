package grpc

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/logger"
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

	logLevel := errors.LogLevel(causeErr)

	globalLogger, ok := ctx.Value(common.ContextHandlerLog).(*logger.Logger)
	if !ok {
		log.Error("failed to get logger from ctx")
		log.Error(err)
	} else {
		globalLogger.Log(logLevel, err)
	}

	return status.Error(errors.GRPCCode(causeErr), causeErr.Error())
}
