package router

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/app/internal/mail/delivery"
	"github.com/gorilla/mux"
	"net/http"
)

type Router struct {
	*mux.Router
	mailD delivery.DeliveryI
}

func New(mailD delivery.DeliveryI) *Router {
	r := &Router{
		Router: mux.NewRouter(),
		mailD:  mailD,
	}

	r.HandleFunc("/inbox/", mailD.GetInboxMessages).Methods(http.MethodGet)
	r.HandleFunc("/outbox/", mailD.GetOutboxMessages).Methods(http.MethodGet)
	r.HandleFunc("/folder/{id:[0-9]+}", mailD.GetFolderMessages).Methods(http.MethodGet)

	return r
}
