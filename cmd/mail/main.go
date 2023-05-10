package main

import (
	"flag"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	_fStorageClient "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage/client"
	_mailRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/repository/postgres"
	_mailServer "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/server"
	_mailUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/usecase"
	_userClient "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/client"
	_middleware "github.com/go-park-mail-ru/2023_1_Seekers/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/connectors"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/logger"
	promMetrics "github.com/go-park-mail-ru/2023_1_Seekers/pkg/metrics/prometheus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var configFile string

	flag.StringVar(&configFile, "config", "cmd/config/debug.yml", "-config=./cmd/config/debug.yml")
	flag.Parse()

	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}

	globalLogger := logger.Init(log.InfoLevel, *cfg.Logger.LogsUseStdOut, cfg.Logger.LogsMailFileName, cfg.Logger.LogsTimeFormat, cfg.Project.ProjectBaseDir, cfg.Logger.LogsDir)

	var connStr = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		cfg.DB.DBUser, cfg.DB.DBPassword, cfg.DB.DBHost, cfg.DB.DBPort, cfg.DB.DBName, cfg.DB.DBSSLMode)

	tablePrefix := cfg.DB.DBSchemaName + "."

	db, err := connectors.NewGormDb(connStr, tablePrefix)
	if err != nil {
		log.Fatalf("db connection error %v", err)
	}

	userServiceCon, err := grpc.Dial(
		cfg.UserGRPCService.Host+":"+cfg.UserGRPCService.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("failed connect to file microservice", err)
	}

	userServiceClient := _userClient.NewUserClientGRPC(userServiceCon)

	fServiceCon, err := grpc.Dial(
		cfg.FileGPRCService.Host+":"+cfg.FileGPRCService.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("failed connect to file microservice", err)
	}

	fStorageClient := _fStorageClient.NewFstorageClientGRPC(fServiceCon)

	mailRepo := _mailRepo.New(cfg, db)
	mailUC := _mailUCase.New(cfg, mailRepo, userServiceClient, fStorageClient)

	metrics, err := promMetrics.NewMetricsGRPCServer(cfg.MailGRPCService.MetricsName)
	if err != nil {
		log.Fatal("mail - failed create metrics server", err)
	}
	middleware := _middleware.NewGRPCMiddleware(cfg, globalLogger, metrics)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middleware.LoggerGRPCUnaryInterceptor, middleware.MetricsGRPCUnaryInterceptor),
	)

	go func() {
		if err = metrics.RunGRPCMetricsServer(":" + cfg.MailGRPCService.MetricsPort); err != nil {
			log.Fatal("mail - failed run metrics server", err)
		}
	}()

	userGRPCServer := _mailServer.NewAuthServerGRPC(grpcServer, mailUC)

	log.Info("mail server started")
	err = userGRPCServer.Start(":" + cfg.MailGRPCService.Port)
	if err != nil {
		log.Fatal(err)
	}
}
