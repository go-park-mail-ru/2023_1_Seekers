package config

import "time"

const (
	Port = "8001"

	CookieName = "MailBoxSession"
	CookieTTL  = time.Hour * 24
	CookieLen  = 32

	ContentTypeJSON = "application/json"
	RouteSignin     = "/api/signin"
	RouteSignup     = "/api/signup"
	RouteAuth       = "/api/auth"
	RouteLogout     = "/api/logout"
)
