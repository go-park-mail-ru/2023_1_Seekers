package pkg

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"time"
)

func GetCurrentTime() string {
	return time.Now().Format(config.LogsTimeFormat)
}
