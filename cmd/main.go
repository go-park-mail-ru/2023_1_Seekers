package main

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	_authHandler "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/delivery/http"
	_authRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/repository/inmemory"
	_authUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/usecase"
	_mailHandler "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/delivery"
	_mailRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/repository/inmemory"
	_mailUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/usecase"
	_middleware "github.com/go-park-mail-ru/2023_1_Seekers/internal/middleware"
	_userRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/repository/inmemory"
	_userUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/usecase"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	mailRepo := _mailRepo.NewMailRepository()
	userRepo := _userRepo.New()
	authRepo := _authRepo.New()

	usersUCase := _userUCase.New(userRepo)
	authUCase := _authUCase.New(authRepo)
	mailUC := _mailUCase.New(mailRepo)

	middleware := _middleware.New(authUCase)

	authH := _authHandler.New(authUCase, usersUCase)
	mailH := _mailHandler.New(mailUC)

	_authHandler.RegisterHTTPRoutes(router, authH, middleware)
	_mailHandler.RegisterHTTPRoutes(router, mailH, middleware)

	corsRouter := middleware.Cors(router)

	server := http.Server{
		Addr:    ":" + config.Port,
		Handler: corsRouter,
	}

	log.Info("server started")
	err := server.ListenAndServe()
	if err != nil {
		log.Errorf("server stopped %v", err)
	}
}
