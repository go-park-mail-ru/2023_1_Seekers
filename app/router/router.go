package router

import (
	_authHandler "github.com/go-park-mail-ru/2023_1_Seekers/app/auth/delivery/http"
	_authUCase "github.com/go-park-mail-ru/2023_1_Seekers/app/auth/usecase"
	_sessionRepo "github.com/go-park-mail-ru/2023_1_Seekers/app/session/repository/inmemory"
	_sessionUcase "github.com/go-park-mail-ru/2023_1_Seekers/app/session/usecase"
	_userRepo "github.com/go-park-mail-ru/2023_1_Seekers/app/user/repository/inmemory"
	_userUCase "github.com/go-park-mail-ru/2023_1_Seekers/app/user/usecase"
	"github.com/gorilla/mux"
	"net/http"
)

func Register(r *mux.Router) {
	userRepo := _userRepo.New()
	sessionRepo := _sessionRepo.New()

	usersUCase := _userUCase.New(userRepo)
	sessionUCase := _sessionUcase.New(sessionRepo)
	authUCase := _authUCase.New(sessionUCase, usersUCase)

	authH := _authHandler.New(authUCase)
	r.HandleFunc("/api/signin", authH.SignIn).Methods(http.MethodPost)
	r.HandleFunc("/api/signup", authH.SignUp).Methods(http.MethodPost)
	r.HandleFunc("/api/logout", authH.Logout).Methods(http.MethodGet)
	r.HandleFunc("/api/auth", authH.Auth).Methods(http.MethodGet)
}

func New() *mux.Router {
	r := mux.NewRouter()
	// TODO корсы
	return r
}
