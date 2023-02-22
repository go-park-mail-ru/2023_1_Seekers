package server

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/app/router"
)

func Run(port string) {
	r := router.New()
	router.Register(r)
	r.Logger.Fatal(r.Start("127.0.0.1:" + port))
	//TODO to config
}
