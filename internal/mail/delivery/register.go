package delivery

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHTTPRoutes(r *mux.Router, h mail.DeliveryI, m *middleware.Middleware) {
	r.HandleFunc(config.RouteInboxMessages, m.CheckAuth(h.GetInboxMessages)).Methods(http.MethodGet)
	r.HandleFunc(config.RouteOutboxMessages, m.CheckAuth(h.GetOutboxMessages)).Methods(http.MethodGet)
	r.HandleFunc(config.RouteFolderMessages, m.CheckAuth(h.GetFolderMessages)).Methods(http.MethodGet)
}
