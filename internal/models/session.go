package models

//go:generate easyjson -all

type Session struct {
	UID       uint64 `json:"userId"`
	SessionID string `json:"sessionId"`
}
