package main

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	_authHandler "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/delivery/http"
	_authRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/repository/inmemory"
	_authUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/usecase"
	_mailHandler "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/delivery"
	_mailRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/reporsitory/inmemory"
	_mailUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/mail/usecase"
	_middleware "github.com/go-park-mail-ru/2023_1_Seekers/internal/middleware"
	_userRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/repository/inmemory"
	_userUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/usecase"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func RegisterRoutes(r *mux.Router) {
	mailRepo := _mailRepo.NewMailRepository()
	userRepo := _userRepo.New()
	authRepo := _authRepo.New()

	usersUCase := _userUCase.New(userRepo)
	authUCase := _authUCase.New(authRepo)
	mailUC := _mailUCase.New(mailRepo)

	middleware := _middleware.New(authUCase)

	authH := _authHandler.New(authUCase, usersUCase)
	mailH := _mailHandler.New(mailUC)

	_authHandler.RegisterHTTPRoutes(r, authH, middleware)
	_mailHandler.RegisterHTTPRoutes(r, mailH, middleware)
}

func main() {
	router := mux.NewRouter()
	RegisterRoutes(router)
	server := http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}
	log.Info("server started")
	err := server.ListenAndServe()
	if err != nil {
		log.Errorf("server stopped %v", err)
	}
}
