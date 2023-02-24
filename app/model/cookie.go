package model

import "time"

type Cookie struct {
	UId     int       `json:"user_id"`
	Session string    `json:"session"`
	Expire  time.Time `json:"expire"`
}
