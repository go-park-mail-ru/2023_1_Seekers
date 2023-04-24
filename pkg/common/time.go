package common

import (
	"time"
)

func GetCurrentTime(layout string) string {
	return time.Now().Format(layout)
}
