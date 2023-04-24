package main

import (
	"flag"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	_authRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/repository/redis"
	_authServer "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/server"
	_authUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/usecase"
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

	redisAddr := cfg.Redis.RedisHost + ":" + cfg.Redis.RedisPort
	redisPw := cfg.Redis.RedisPassword

	rdb, err := connectors.NewRedisClient(redisAddr, redisPw)
	if err != nil {
		log.Fatalf("failed connect to redis : %v", err)
	}

	sessionRepo := _authRepo.NewSessionRepo(cfg, rdb)

	userServiceCon, err := grpc.Dial(
		cfg.UserGRPCService.Host+":"+cfg.UserGRPCService.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("failed connect to file microservice", err)
	}

	userServiceClient := _userClient.NewUserClientGRPC(userServiceCon)

	authUC := _authUCase.NewAuthUC(cfg, userServiceClient, sessionRepo)

	grpcServer := grpc.NewServer()
	authGRPCServer := _authServer.NewAuthServerGRPC(grpcServer, authUC)

	log.Info("auth server started")
	err = authGRPCServer.Start(":" + cfg.AuthGRPCService.Port)
	if err != nil {
		log.Fatal(err)
	}
}
