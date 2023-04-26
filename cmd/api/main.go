package main

import (
	"context"
	"flag"
	_ "github.com/go-park-mail-ru/2023_1_Seekers/docs"
	_api "github.com/go-park-mail-ru/2023_1_Seekers/internal/api/http"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	_authClient "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/client"
	_mailClient "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/client"
	_userClient "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/client"
	_middleware "github.com/go-park-mail-ru/2023_1_Seekers/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/logger"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// @title MailBox Swagger API
// @version 1.0
// @host localhost:8001
// @BasePath	/api/v1
func main() {
	var configFile string

	flag.StringVar(&configFile, "config", "cmd/config/debug.yml", "-config=./cmd/configs/debug.yml")
	flag.Parse()

	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}

	globalLogger := logger.Init(log.InfoLevel, *cfg.Logger.LogsUseStdOut, cfg.Logger.LogsApiFileName, cfg.Logger.LogsTimeFormat, cfg.Project.ProjectBaseDir, cfg.Logger.LogsDir)
	router := mux.NewRouter()

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

	authH := _api.NewAuthHandlers(cfg, authServiceClient, mailServiceClient, userServiceClient)
	mailH := _api.NewMailHandlers(cfg, mailServiceClient)
	userH := _api.NewUserHandlers(cfg, userServiceClient)
	middleware := _middleware.New(cfg, authServiceClient, globalLogger)

	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	_api.RegisterHTTPRoutes(router, cfg, authH, userH, mailH, middleware)

	router.Use(middleware.HandlerLogger)
	corsRouter := middleware.Cors(router)

	server := http.Server{
		Addr:         ":" + cfg.Api.Port,
		Handler:      corsRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		globalLogger.Info("server started")
		if err := server.ListenAndServe(); err != nil {
			globalLogger.Fatalf("server stopped %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Kill, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	if err = server.Shutdown(ctx); err != nil {
		globalLogger.Errorf("failed to gracefully shutdown server")
	}

}
