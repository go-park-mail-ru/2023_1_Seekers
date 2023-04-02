package delivery

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHTTPRoutes(r *mux.Router, h mail.HandlersI, m *middleware.Middleware) {
	r.HandleFunc(config.RouteGetFolderMessages, m.CheckAuth(m.CheckCSRF(h.GetFolderMessages))).Methods(http.MethodGet)
	r.HandleFunc(config.RouteGetFolders, m.CheckAuth(m.CheckCSRF(h.GetFolders))).Methods(http.MethodGet)
	r.HandleFunc(config.RouteGetMessage, m.CheckAuth(m.CheckCSRF(h.GetMessage))).Methods(http.MethodGet)
	r.HandleFunc(config.RouteSendMessage, m.CheckAuth(m.CheckCSRF(h.SendMessage))).Methods(http.MethodPost)
	r.HandleFunc(config.RouteReadMessage, m.CheckAuth(m.CheckCSRF(h.ReadMessage))).Methods(http.MethodPost)
	r.HandleFunc(config.RouteUnreadMessage, m.CheckAuth(m.CheckCSRF(h.UnreadMessage))).Methods(http.MethodPost)
}
