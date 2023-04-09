package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	_ "github.com/go-park-mail-ru/2023_1_Seekers/docs"
	_authHandler "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/delivery/http"
	_sessionRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/repository/redis"
	_authUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/usecase"
	_fStorageRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/file_storage/repository"
	_fStorageUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/file_storage/usecase"
	_mailHandler "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/delivery"
	_mailRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/repository/postgres"
	_mailUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/usecase"
	_middleware "github.com/go-park-mail-ru/2023_1_Seekers/internal/middleware"
	_userHandler "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/delivery/http"
	_userRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/repository/postgres"
	_userUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/usecase"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/connectors"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var connStr = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
	os.Getenv(config.DBUserEnv),
	os.Getenv(config.DBPasswordEnv),
	os.Getenv(config.DBHostEnv),
	os.Getenv(config.DBPortEnv),
	os.Getenv(config.DBNameEnv),
	os.Getenv(config.DBSSLModeEnv),
)

// @title MailBox Swagger API
// @version 1.0
// @host localhost:8001
// @BasePath	/api/v1
func main() {
	pkg.InitLogger()
	logger := pkg.GetLogger()
	router := mux.NewRouter()

	tablePrefix := os.Getenv(config.DBSchemaNameEnv) + "."
	db, err := connectors.NewGormDb(connStr, tablePrefix)
	if err != nil {
		logger.Fatalf("db connection error %v", err)
	}

	redisAddr := os.Getenv(config.RedisHostEnv) + ":" + os.Getenv(config.RedisPortEnv)
	redisPw := os.Getenv(config.RedisPasswordEnv)
	rdb, err := connectors.NewRedisClient(redisAddr, redisPw)
	if err != nil {
		log.Fatalf("failed connect to redis : %v", err)
	}

	endpoint := aws.String(config.S3Endpoint)
	region := aws.String(config.S3Region)
	disableSSL := aws.Bool(true)
	s3ForcePathStyle := aws.Bool(true)
	creds := credentials.NewStaticCredentials(
		os.Getenv(config.S3AccessKeyEnv),
		os.Getenv(config.S3ASecretKeyEnv),
		"",
	)
	s3Session, err := connectors.NewS3(endpoint, region, disableSSL, s3ForcePathStyle, creds)
	if err != nil {
		log.Fatalf("Failed create S3 session : %v", err)
	}

	userRepo := _userRepo.New(db)
	sessionRepo := _sessionRepo.NewSessionRepo(rdb)
	mailRepo := _mailRepo.New(db)
	fStorageRepo := _fStorageRepo.New(s3Session)

	fStorageUC := _fStorageUCase.New(fStorageRepo)
	usersUC := _userUCase.New(userRepo, fStorageUC)
	mailUC := _mailUCase.New(mailRepo, userRepo)
	sessionUC := _authUCase.NewSessionUC(sessionRepo)
	authUC := _authUCase.NewAuthUC(sessionUC, mailUC, usersUC)

	middleware := _middleware.New(sessionUC, logger)

	authH := _authHandler.New(authUC)
	mailH := _mailHandler.New(mailUC)
	userH := _userHandler.New(usersUC)

	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	_authHandler.RegisterHTTPRoutes(router, authH, middleware)
	_mailHandler.RegisterHTTPRoutes(router, mailH, middleware)
	_userHandler.RegisterHTTPRoutes(router, userH, middleware)

	router.Use(middleware.HandlerLogger)
	corsRouter := middleware.Cors(router)

	server := http.Server{
		Addr:         ":" + config.Port,
		Handler:      corsRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("server started")
		if err := server.ListenAndServe(); err != nil {
			logger.Fatalf("server stopped %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Kill, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("failed to gracefully shutdown server")
	}

}
