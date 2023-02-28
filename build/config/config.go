package config

import "time"

const (
	Port = "8001"

	CookieName = "MailBoxSession"
	CookieTTL  = time.Hour * 24 * 100 // 100 days
	CookieLen  = 32

	ContentTypeJSON = "application/json"
	RouteSignin     = "/api/signin"
	RouteSignup     = "/api/signup"
	RouteAuth       = "/api/auth"
	RouteLogout     = "/api/logout"
)
