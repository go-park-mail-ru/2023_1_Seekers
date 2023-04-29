package main

import (
	"flag"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	_authRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/repository/redis"
	_authServer "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/server"
	_authUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/usecase"
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

	flag.StringVar(&configFile, "config", "cmd/config/debug.yml", "-config=./cmd/configs/debug.yml")
	flag.Parse()

	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}

	globalLogger := logger.Init(log.InfoLevel, *cfg.Logger.LogsUseStdOut, cfg.Logger.LogsAuthFileName, cfg.Logger.LogsTimeFormat, cfg.Project.ProjectBaseDir, cfg.Logger.LogsDir)

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

	metrics, err := promMetrics.NewMetricsGRPCServer("auth")
	if err != nil {
		log.Fatal("auth - failed create metrics server", err)
	}
	middleware := _middleware.NewGRPCMiddleware(cfg, globalLogger, metrics)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middleware.LoggerGRPCUnaryInterceptor, middleware.MetricsGRPCUnaryInterceptor),
	)
	authGRPCServer := _authServer.NewAuthServerGRPC(grpcServer, authUC)

	// TODO to conf
	//go func() {
	//	if err := promMetrics.RunGRPCMetricsServer(":9002"); err != nil {
	//		log.Fatal("auth - failed run metrics server", err)
	//	}
	//
	//}()

	log.Info("auth server started")
	err = authGRPCServer.Start(":" + cfg.AuthGRPCService.Port)
	if err != nil {
		log.Fatal(err)
	}
}
