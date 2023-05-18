package http

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHTTPRoutes(r *mux.Router, cfg *config.Config, authH AuthHandlersI, userH UserHandlersI, mailH MailHandlersI, m *middleware.HttpMiddleware) {
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
		Methods(http.MethodGet).Queries(cfg.Routes.RouteUserAvatarQueryEmail, "{email}")
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteUserInfo, userH.GetInfo).
		Methods(http.MethodGet).Queries(cfg.Routes.RouteUserInfoQueryEmail, "{email}")
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteUserInfo, m.CheckAuth(userH.GetPersonalInfo)).Methods(http.MethodGet)

	// Mail
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteFolder, m.CheckAuth(m.CheckCSRF(mailH.GetFolderMessages))).Methods(http.MethodGet)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteSearch, m.CheckAuth(m.CheckCSRF(mailH.SearchMessages))).Methods(http.MethodGet)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteRecipients, m.CheckAuth(m.CheckCSRF(mailH.SearchRecipients))).Methods(http.MethodGet)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteGetFolders, m.CheckAuth(m.CheckCSRF(mailH.GetFolders))).Methods(http.MethodGet)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteMessage, m.CheckAuth(m.CheckCSRF(mailH.GetMessage))).Methods(http.MethodGet)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteMessage, m.CheckAuth(m.CheckCSRF(mailH.DeleteMessage))).Methods(http.MethodDelete)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteSendMessage, m.CheckAuth(m.CheckCSRF(mailH.SendMessage))).Methods(http.MethodPost)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteReadMessage, m.CheckAuth(m.CheckCSRF(mailH.ReadMessage))).Methods(http.MethodPost)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteUnreadMessage, m.CheckAuth(m.CheckCSRF(mailH.UnreadMessage))).Methods(http.MethodPost)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteCreateFolder, m.CheckAuth(m.CheckCSRF(mailH.CreateFolder))).Methods(http.MethodPost)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteFolder, m.CheckAuth(m.CheckCSRF(mailH.DeleteFolder))).Methods(http.MethodDelete)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteEditFolder, m.CheckAuth(m.CheckCSRF(mailH.EditFolder))).Methods(http.MethodPut)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteMoveToFolder, m.CheckAuth(m.CheckCSRF(mailH.MoveToFolder))).Methods(http.MethodPut)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteSaveDraftMessage, m.CheckAuth(m.CheckCSRF(mailH.SaveDraft))).Methods(http.MethodPost)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteMessage, m.CheckAuth(m.CheckCSRF(mailH.EditDraft))).Methods(http.MethodPut)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteAttach, m.CheckAuth(m.CheckCSRF(mailH.DownloadAttach))).Methods(http.MethodGet)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RoutePreviewAttach, m.CheckAuth(m.CheckCSRF(mailH.PreviewAttach))).Methods(http.MethodGet)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteExternalAttach, mailH.GetAttach).Methods(http.MethodGet)
	r.HandleFunc(cfg.Routes.RoutePrefix+cfg.Routes.RouteWS, m.CheckAuth(m.CheckCSRF(mailH.WSMessageHandler))).Methods(http.MethodGet)
	//r.HandleFunc(cfg.Routes.RoutePrefix+"/chat", mailH.File)
}
