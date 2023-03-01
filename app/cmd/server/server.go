package server

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/cmd/router"
	"net/http"
)

func Run(port string, route *router.Router) error {
	server := http.Server{
		Addr:    ":" + port,
		Handler: route,
	}
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("server stopped %v", err)
	}
	return nil
}
