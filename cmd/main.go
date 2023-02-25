package main

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/app/server"
	"github.com/go-park-mail-ru/2023_1_Seekers/config"
)

func main() {
	server.Run(config.Port)
}
