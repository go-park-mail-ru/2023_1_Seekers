package main

import (
	"flag"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	_fStorageClient "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage/client"
	_userRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/repository/postgres"
	_userServer "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/server"
	_userUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/usecase"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/connectors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var configFile string

	flag.StringVar(&configFile, "config", "cmd/config/debug.yml", "-config=./cmd/configs/debug.yml")
	flag.Parse()

	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var connStr = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		cfg.DB.DBUser, cfg.DB.DBPassword, cfg.DB.DBHost, cfg.DB.DBPort, cfg.DB.DBName, cfg.DB.DBSSLMode)

	tablePrefix := cfg.DB.DBSchemaName + "."

	db, err := connectors.NewGormDb(connStr, tablePrefix)
	if err != nil {
		log.Fatalf("db connection error %v", err)
	}

	fServiceCon, err := grpc.Dial(
		cfg.FileGPRCService.Host+":"+cfg.FileGPRCService.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("failed connect to file microservice", err)
	}

	fStorageClient := _fStorageClient.NewFstorageClientGRPC(fServiceCon)

	userRepo := _userRepo.New(cfg, db)
	usersUC := _userUCase.New(cfg, userRepo, fStorageClient)

	grpcServer := grpc.NewServer()
	userGRPCServer := _userServer.NewUserServerGRPC(grpcServer, usersUC)
	log.Info("user server started")
	err = userGRPCServer.Start(":" + cfg.UserGRPCService.Port)
	if err != nil {
		log.Fatal(err)
	}
}
