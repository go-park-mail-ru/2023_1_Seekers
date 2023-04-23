package main

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	_mailRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/repository/postgres"
	_mailServer "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/server"
	_mailUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/usecase"
	_userClient "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/client"
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

	userServiceCon, err := grpc.Dial(
		config.UserGRPCHost+":"+config.UserGRPCPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("failed connect to file microservice", err)
	}

	userServiceClient := _userClient.NewUserClientGRPC(userServiceCon)

	mailRepo := _mailRepo.New(db)
	mailUC := _mailUCase.New(mailRepo, userServiceClient)

	grpcServer := grpc.NewServer()
	userGRPCServer := _mailServer.NewAuthServerGRPC(grpcServer, mailUC)

	log.Info("mail server started")
	err = userGRPCServer.Start(":" + config.MailGRPCPort)
	if err != nil {
		log.Fatal(err)
	}
}
