package main

import (
	"flag"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	_authClient "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/client"
	_mailClient "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/client"
	_userClient "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/client"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/smtp/server"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// this server must be running on port 25 ( Root is needed )
// sudo go run cmd/smtp/main.go -config=./cmd/configs/debug.yml

func main() {
	var configFile string

	flag.StringVar(&configFile, "config", "./cmd/config/debug.yml", "-config=./cmd/config/debug.yml")
	flag.Parse()

	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}

	authServiceCon, err := grpc.Dial(
		cfg.AuthGRPCService.Host+":"+cfg.AuthGRPCService.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("failed connect to file microservice", err)
	}

	authServiceClient := _authClient.NewAuthClientGRPC(authServiceCon)

	mailServiceCon, err := grpc.Dial(
		cfg.MailGRPCService.Host+":"+cfg.MailGRPCService.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("failed connect to file microservice", err)
	}

	mailServiceClient := _mailClient.NewMailClientGRPC(mailServiceCon)

	userServiceCon, err := grpc.Dial(
		cfg.UserGRPCService.Host+":"+cfg.UserGRPCService.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("failed connect to file microservice", err)
	}

	userServiceClient := _userClient.NewUserClientGRPC(userServiceCon)

	if err = server.RunSmtpServer(cfg, mailServiceClient, userServiceClient, authServiceClient); err != nil {
		log.Fatal("smtp server stopped", err)
	}
}
