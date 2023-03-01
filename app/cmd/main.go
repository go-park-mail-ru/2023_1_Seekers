package main

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/cmd/router"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/cmd/server"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/internal/mail/delivery"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/internal/mail/reporsitory/inmemory"
	"github.com/go-park-mail-ru/2023_1_Seekers/app/internal/mail/usecase"
)

func main() {
	mailRepo := inmemory.NewMailRepository()
	mailUC := usecase.New(mailRepo)
	mailDeliver := delivery.New(mailUC)

	r := router.New(mailDeliver)

	fmt.Println("run server")

	server.Run(config.Port, r)
}
