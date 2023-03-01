package main

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	_authHandler "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/delivery/http"
	_authRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/repository/inmemory"
	_authUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/auth/usecase"
	_middleware "github.com/go-park-mail-ru/2023_1_Seekers/internal/middleware"
	_userRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/repository/inmemory"
	_userUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/user/usecase"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func RegisterRoutes(r *mux.Router) {
	userRepo := _userRepo.New()
	authRepo := _authRepo.New()

	usersUCase := _userUCase.New(userRepo)
	authUCase := _authUCase.New(authRepo)

	middleware := _middleware.New(authUCase)

	authH := _authHandler.New(authUCase, usersUCase)

	_authHandler.RegisterHTTPRoutes(r, authH, middleware)
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
