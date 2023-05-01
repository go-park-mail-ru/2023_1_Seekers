package main

import (
	"flag"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	_mailRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/repository/postgres"
	_mailServer "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/server"
	_mailUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/usecase"
	_userClient "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/client"
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

	userServiceCon, err := grpc.Dial(
		cfg.UserGRPCService.Host+":"+cfg.UserGRPCService.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("failed connect to file microservice", err)
	}

	userServiceClient := _userClient.NewUserClientGRPC(userServiceCon)

	mailRepo := _mailRepo.New(cfg, db)
	mailUC := _mailUCase.New(cfg, mailRepo, userServiceClient)

	grpcServer := grpc.NewServer()
	userGRPCServer := _mailServer.NewAuthServerGRPC(grpcServer, mailUC)

	log.Info("mail server started")
	err = userGRPCServer.Start(":" + cfg.MailGRPCService.Port)
	if err != nil {
		log.Fatal(err)
	}
}
