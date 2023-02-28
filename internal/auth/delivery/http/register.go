package http

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/build/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHTTPRoutes(r *mux.Router, h auth.Handlers) {
	r.HandleFunc(config.RouteSignin, h.SignIn).Methods(http.MethodPost)
	r.HandleFunc(config.RouteSignup, h.SignUp).Methods(http.MethodPost)
	r.HandleFunc(config.RouteLogout, h.Logout).Methods(http.MethodGet)
	r.HandleFunc(config.RouteAuth, h.Auth).Methods(http.MethodGet)
}
