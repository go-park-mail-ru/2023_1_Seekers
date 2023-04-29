package main

import (
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	_fStorageRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage/repository/S3"
	_fStorageServer "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage/server"
	_fStorageUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage/usecase"
	_middleware "github.com/go-park-mail-ru/2023_1_Seekers/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/connectors"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/logger"
	promMetrics "github.com/go-park-mail-ru/2023_1_Seekers/pkg/metrics/prometheus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	var configFile string

	flag.StringVar(&configFile, "config", "cmd/config/debug.yml", "-config=./cmd/configs/debug.yml")
	flag.Parse()

	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}

	globalLogger := logger.Init(log.InfoLevel, *cfg.Logger.LogsUseStdOut, cfg.Logger.LogsFileServiceFileName, cfg.Logger.LogsTimeFormat, cfg.Project.ProjectBaseDir, cfg.Logger.LogsDir)

	endpoint := aws.String(cfg.S3.S3Endpoint)
	region := aws.String(cfg.S3.S3Region)
	disableSSL := aws.Bool(true)
	s3ForcePathStyle := aws.Bool(true)
	creds := credentials.NewStaticCredentials(cfg.S3.S3AccessKey, cfg.S3.S3ASecretKey, "")
	s3Session, err := connectors.NewS3(endpoint, region, disableSSL, s3ForcePathStyle, creds)
	if err != nil {
		log.Fatalf("Failecd create S3 session : %v", err)
	}

	fStorageRepo := _fStorageRepo.New(s3Session)
	fStorageUC := _fStorageUCase.New(fStorageRepo)

	metrics, err := promMetrics.NewMetricsGRPCServer("file_service")
	if err != nil {
		log.Fatal("file service - failed create metrics server", err)
	}
	middleware := _middleware.NewGRPCMiddleware(cfg, globalLogger, metrics)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middleware.MetricsGRPCUnaryInterceptor),
	)

	//promMetrics.RunGRPCMetricsServer(":9003")

	fStorageGRPCServer := _fStorageServer.NewFStorageServerGRPC(grpcServer, fStorageUC)
	log.Info("file storage server started")
	err = fStorageGRPCServer.Start(":" + cfg.FileGPRCService.Port)
	if err != nil {
		log.Fatal(err)
	}
}
