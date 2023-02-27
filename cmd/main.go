package main

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/server"
)

func main() {
	server.Run(config.Port)
}
