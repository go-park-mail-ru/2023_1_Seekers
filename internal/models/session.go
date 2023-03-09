package models

type Session struct {
	UID       uint64 `json:"user_id"`
	SessionID string `json:"session_id"`
}
