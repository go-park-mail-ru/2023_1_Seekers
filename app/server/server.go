package server

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/router"
	"github.com/labstack/gommon/log"
	"net/http"
)

func Run(port string) error {
	mux := router.New()
	router.Register(mux)
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Errorf("failed to run server %v", err)
		return fmt.Errorf("failed to run server %v", err)
	}
	return nil
}
