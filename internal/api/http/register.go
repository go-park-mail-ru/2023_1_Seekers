package http

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHTTPRoutes(r *mux.Router, cfg *config.Config, authH AuthHandlersI, userH UserHandlersI, mailH MailHandlersI, m *middleware.Middleware) {
	// Auth
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteSignin, authH.SignIn).Methods(http.MethodPost)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteSignup, authH.SignUp).Methods(http.MethodPost)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteLogout, m.CheckAuth(m.CheckCSRF(authH.Logout))).Methods(http.MethodDelete)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteAuth, m.CheckAuth(authH.Auth)).Methods(http.MethodGet)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteCSRF, authH.GetCSRF).Methods(http.MethodGet)

	// User
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteUser, m.CheckAuth(m.CheckCSRF(userH.Delete))).Methods(http.MethodDelete)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteUserInfo, m.CheckAuth(m.CheckCSRF(userH.EditInfo))).Methods(http.MethodPut)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteUserAvatar, m.CheckAuth(m.CheckCSRF(userH.EditAvatar))).Methods(http.MethodPut)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RoutePw, m.CheckAuth(m.CheckCSRF(userH.EditPw))).Methods(http.MethodPut)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteUserAvatar, userH.GetAvatar).
		Methods(http.MethodGet).Queries(cfg.Routes.RoutePrefix+cfg.Routes.RouteUserAvatarQueryEmail, "{email}")
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteUserInfo, userH.GetInfo).
		Methods(http.MethodGet).Queries(cfg.Routes.RoutePrefix+cfg.Routes.RouteUserInfoQueryEmail, "{email}")
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteUserInfo, m.CheckAuth(userH.GetPersonalInfo)).Methods(http.MethodGet)

	// Mail
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteGetFolderMessages, m.CheckAuth(m.CheckCSRF(mailH.GetFolderMessages))).Methods(http.MethodGet)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteGetFolders, m.CheckAuth(m.CheckCSRF(mailH.GetFolders))).Methods(http.MethodGet)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteGetMessage, m.CheckAuth(m.CheckCSRF(mailH.GetMessage))).Methods(http.MethodGet)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteSendMessage, m.CheckAuth(m.CheckCSRF(mailH.SendMessage))).Methods(http.MethodPost)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteReadMessage, m.CheckAuth(m.CheckCSRF(mailH.ReadMessage))).Methods(http.MethodPost)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteUnreadMessage, m.CheckAuth(m.CheckCSRF(mailH.UnreadMessage))).Methods(http.MethodPost)
}
