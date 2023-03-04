package main

import (
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
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

// @title MailBox Swagger API
// @version 1.0
// @host localhost:8001
// @BasePath	/api/v1

func main() {
	router := mux.NewRouter()

	userRepo := _userRepo.New()
	authRepo := _authRepo.New()
	mailRepo := _mailRepo.New(userRepo)

	usersUC := _userUCase.New(userRepo)
	authUC := _authUCase.New(authRepo, usersUC)
	mailUC := _mailUCase.New(mailRepo)

	middleware := _middleware.New(authUC)

	authH := _authHandler.New(authUC, usersUC)
	mailH := _mailHandler.New(mailUC)

	//router
	//
	//r.Get("/swagger/*", httpSwagger.Handler(
	//	httpSwagger.URL("http://localhost:8001/swagger/doc.json"), //The url pointing to API definition
	//))

	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	_authHandler.RegisterHTTPRoutes(router, authH, middleware)
	_mailHandler.RegisterHTTPRoutes(router, mailH, middleware)

	corsRouter := middleware.Cors(router)

	server := http.Server{
		Addr:    ":" + config.Port,
		Handler: corsRouter,
	}

	log.Info("server started")
	err := server.ListenAndServe()
	//err := server.ListenAndServeTLS("go-server.crt", "go-server.key")
	if err != nil {
		log.Errorf("server stopped %v", err)
	}
}
