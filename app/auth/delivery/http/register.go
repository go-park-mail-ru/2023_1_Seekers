package http

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/app/auth"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHTTPRoutes(r *mux.Router, h auth.Handlers) {
	r.HandleFunc("/api/signin", h.SignIn).Methods(http.MethodPost)
	r.HandleFunc("/api/signup", h.SignUp).Methods(http.MethodPost)
	r.HandleFunc("/api/logout", h.Logout).Methods(http.MethodGet)
	r.HandleFunc("/api/auth", h.Auth).Methods(http.MethodGet)
}
