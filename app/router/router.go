package router

import (
	_authHandler "github.com/go-park-mail-ru/2023_1_Seekers/app/auth/delivery/http"
	_authRepo "github.com/go-park-mail-ru/2023_1_Seekers/app/auth/repository/inmemory"
	_authUCase "github.com/go-park-mail-ru/2023_1_Seekers/app/auth/usecase"
	_sessionRepo "github.com/go-park-mail-ru/2023_1_Seekers/app/session/repository/inmemory"
	_sessionUcase "github.com/go-park-mail-ru/2023_1_Seekers/app/session/usecase"
	_userRepo "github.com/go-park-mail-ru/2023_1_Seekers/app/user/repository/inmemory"
	_userUCase "github.com/go-park-mail-ru/2023_1_Seekers/app/user/usecase"
	"github.com/gorilla/mux"
)

func Register(r *mux.Router) {
	userRepo := _userRepo.New()
	sessionRepo := _sessionRepo.New()
	authRepo := _authRepo.New()

	usersUCase := _userUCase.New(userRepo)
	sessionUCase := _sessionUcase.New(sessionRepo)
	authUCase := _authUCase.New(authRepo)

	authH := _authHandler.New(authUCase, sessionUCase, usersUCase)

	_authHandler.RegisterHTTPRoutes(r, authH)
}

func New() *mux.Router {
	r := mux.NewRouter()
	// TODO корсы
	return r
}
