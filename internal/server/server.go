package server

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/router"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Run(port string) error {
	mux := router.New()
	router.Register(mux)
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Info("server started")
	err := server.ListenAndServe()
	if err != nil {
		log.Errorf("server stopped %v", err)
		return fmt.Errorf("server stopped %v", err)
	}
	return nil
}
