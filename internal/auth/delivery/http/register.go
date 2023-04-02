package http

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHTTPRoutes(r *mux.Router, h auth.HandlersI, m *middleware.Middleware) {
	r.HandleFunc(config.RouteSignin, h.SignIn).Methods(http.MethodPost)
	r.HandleFunc(config.RouteSignup, h.SignUp).Methods(http.MethodPost)
	r.HandleFunc(config.RouteLogout, m.CheckAuth(m.CheckCSRF(h.Logout))).Methods(http.MethodGet)
	r.HandleFunc(config.RoutePw, m.CheckAuth(m.CheckCSRF(h.EditPw))).Methods(http.MethodPost)
	r.HandleFunc(config.RouteCSRF, h.GetCSRF).Methods(http.MethodGet)
}
