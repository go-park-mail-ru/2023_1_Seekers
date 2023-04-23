package http

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHTTPRoutes(r *mux.Router, authH AuthHandlersI, userH UserHandlersI, mailH MailHandlersI, m *middleware.Middleware) {
	// Auth
	r.HandleFunc(config.RouteSignin, authH.SignIn).Methods(http.MethodPost)
	r.HandleFunc(config.RouteSignup, authH.SignUp).Methods(http.MethodPost)
	r.HandleFunc(config.RouteLogout, m.CheckAuth(m.CheckCSRF(authH.Logout))).Methods(http.MethodDelete)
	r.HandleFunc(config.RouteAuth, m.CheckAuth(authH.Auth)).Methods(http.MethodGet)
	r.HandleFunc(config.RouteCSRF, authH.GetCSRF).Methods(http.MethodGet)

	// User
	r.HandleFunc(config.RouteUser, m.CheckAuth(m.CheckCSRF(userH.Delete))).Methods(http.MethodDelete)
	r.HandleFunc(config.RouteUserInfo, m.CheckAuth(m.CheckCSRF(userH.EditInfo))).Methods(http.MethodPut)
	r.HandleFunc(config.RouteUserAvatar, m.CheckAuth(m.CheckCSRF(userH.EditAvatar))).Methods(http.MethodPut)
	r.HandleFunc(config.RoutePw, m.CheckAuth(m.CheckCSRF(userH.EditPw))).Methods(http.MethodPut)
	r.HandleFunc(config.RouteUserAvatar, userH.GetAvatar).
		Methods(http.MethodGet).Queries(config.RouteUserAvatarQueryEmail, "{email}")
	r.HandleFunc(config.RouteUserInfo, userH.GetInfo).
		Methods(http.MethodGet).Queries(config.RouteUserInfoQueryEmail, "{email}")
	r.HandleFunc(config.RouteUserInfo, m.CheckAuth(userH.GetPersonalInfo)).Methods(http.MethodGet)

	// Mail
	r.HandleFunc(config.RouteGetFolderMessages, m.CheckAuth(m.CheckCSRF(mailH.GetFolderMessages))).Methods(http.MethodGet)
	r.HandleFunc(config.RouteGetFolders, m.CheckAuth(m.CheckCSRF(mailH.GetFolders))).Methods(http.MethodGet)
	r.HandleFunc(config.RouteGetMessage, m.CheckAuth(m.CheckCSRF(mailH.GetMessage))).Methods(http.MethodGet)
	r.HandleFunc(config.RouteSendMessage, m.CheckAuth(m.CheckCSRF(mailH.SendMessage))).Methods(http.MethodPost)
	r.HandleFunc(config.RouteReadMessage, m.CheckAuth(m.CheckCSRF(mailH.ReadMessage))).Methods(http.MethodPost)
	r.HandleFunc(config.RouteUnreadMessage, m.CheckAuth(m.CheckCSRF(mailH.UnreadMessage))).Methods(http.MethodPost)
}
