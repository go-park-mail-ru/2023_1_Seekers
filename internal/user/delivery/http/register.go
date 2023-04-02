package http

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/user"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHTTPRoutes(r *mux.Router, h user.HandlersI, m *middleware.Middleware) {
	r.HandleFunc(config.RouteUser, m.CheckAuth(m.CheckCSRF(h.Delete))).Methods(http.MethodDelete)
	r.HandleFunc(config.RouteUserInfo, m.CheckAuth(m.CheckCSRF(h.EditInfo))).Methods(http.MethodPost)
	r.HandleFunc(config.RouteUserPw, m.CheckAuth(m.CheckCSRF(h.EditPw))).Methods(http.MethodPost)
	r.HandleFunc(config.RouteUserAvatar, m.CheckAuth(m.CheckCSRF(h.EditAvatar))).Methods(http.MethodPost)
	r.HandleFunc(config.RouteUserAvatar, h.GetAvatar).
		Methods(http.MethodGet).Queries(config.RouteUserAvatarQueryEmail, "{email}")
	r.HandleFunc(config.RouteUserInfo, h.GetInfo).
		Methods(http.MethodGet).Queries(config.RouteUserInfoQueryEmail, "{email}")
	r.HandleFunc(config.RouteUserInfo, m.CheckAuth(h.GetPersonalInfo)).Methods(http.MethodGet)
}
