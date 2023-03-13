package main

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	_ "github.com/go-park-mail-ru/2023_1_Seekers/docs"
	_authHandler "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/delivery/http"
	_authRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/repository/inmemory"
	_authUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/usecase"
	_mailHandler "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/delivery"
	_mailRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/repository/inmemory"
	_mailUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/usecase"
	_middleware "github.com/go-park-mail-ru/2023_1_Seekers/internal/middleware"
	_userRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/repository/inmemory"
	_userUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/usecase"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var connStr = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
	config.DBUser,
	config.DBPassword,
	config.DBHost,
	config.DBPort,
	config.DBName,
)

// @title MailBox Swagger API
// @version 1.0
// @host localhost:8001
// @BasePath	/api/v1
func main() {
	pkg.InitLogger()
	logger := pkg.GetLogger()
	router := mux.NewRouter()

	_, err := gorm.Open(postgres.New(postgres.Config{DSN: connStr}), &gorm.Config{})
	if err != nil {
		logger.Fatalf("db connection error %v", err)
		return
	}

	userRepo := _userRepo.New()
	authRepo := _authRepo.New()
	mailRepo := _mailRepo.New(userRepo)

	usersUC := _userUCase.New(userRepo)
	authUC := _authUCase.New(authRepo, usersUC)
	mailUC := _mailUCase.New(mailRepo)

	middleware := _middleware.New(authUC, logger)

	authH := _authHandler.New(authUC, usersUC, mailUC)
	mailH := _mailHandler.New(mailUC)

	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	_authHandler.RegisterHTTPRoutes(router, authH, middleware)
	_mailHandler.RegisterHTTPRoutes(router, mailH, middleware)

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
