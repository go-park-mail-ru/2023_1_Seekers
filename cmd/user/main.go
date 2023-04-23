package main

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	_fStorageClient "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage/client"
	_userRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/repository/postgres"
	_userServer "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/server"
	_userUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/usecase"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/connectors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

var connStr = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
	os.Getenv(config.DBUserEnv),
	os.Getenv(config.DBPasswordEnv),
	os.Getenv(config.DBHostEnv),
	os.Getenv(config.DBPortEnv),
	os.Getenv(config.DBNameEnv),
	os.Getenv(config.DBSSLModeEnv),
)

func main() {
	tablePrefix := os.Getenv(config.DBSchemaNameEnv) + "."
	db, err := connectors.NewGormDb(connStr, tablePrefix)
	if err != nil {
		log.Fatalf("db connection error %v", err)
	}

	fServiceCon, err := grpc.Dial(
		config.FileServiceGRPCHost+":"+config.FileServiceGRPCPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("failed connect to file microservice", err)
	}

	fStorageClient := _fStorageClient.NewFstorageClientGRPC(fServiceCon)

	userRepo := _userRepo.New(db)
	usersUC := _userUCase.New(userRepo, fStorageClient)

	grpcServer := grpc.NewServer()
	userGRPCServer := _userServer.NewUserServerGRPC(grpcServer, usersUC)
	log.Info("user server started")
	err = userGRPCServer.Start(":" + config.UserGRPCPort)
	if err != nil {
		log.Fatal(err)
	}
}
