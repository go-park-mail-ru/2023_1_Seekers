package config

import "time"

const (
	Port        = "8001"
	CookieName  = "MailBoxSession"
	CookieTTL   = time.Hour * 24
	CookieLen   = 32
	NetTypeJSON = "application/json"
)
