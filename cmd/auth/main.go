package main

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	_authRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/repository/redis"
	_authServer "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/server"
	_authUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/usecase"
	_userClient "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/client"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/connectors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

func main() {
	redisAddr := os.Getenv(config.RedisHostEnv) + ":" + os.Getenv(config.RedisPortEnv)
	redisPw := os.Getenv(config.RedisPasswordEnv)
	rdb, err := connectors.NewRedisClient(redisAddr, redisPw)
	if err != nil {
		log.Fatalf("failed connect to redis : %v", err)
	}

	sessionRepo := _authRepo.NewSessionRepo(rdb)

	userServiceCon, err := grpc.Dial(
		config.UserGRPCHost+":"+config.UserGRPCPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("failed connect to file microservice", err)
	}

	userServiceClient := _userClient.NewUserClientGRPC(userServiceCon)

	authUC := _authUCase.NewAuthUC(userServiceClient, sessionRepo)

	grpcServer := grpc.NewServer()
	authGRPCServer := _authServer.NewAuthServerGRPC(grpcServer, authUC)

	log.Info("auth server started")
	err = authGRPCServer.Start(":" + config.AuthGRPCPort)
	if err != nil {
		log.Fatal(err)
	}
}
