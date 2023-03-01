package config

import "time"

const (
	Port = "8001"

	CookieName      = "MailBoxSession"
	CookieTTL       = time.Hour * 24 * 100 // 100 days
	CookieLen       = 32
	CookiePath      = "/"
	ContextUser     = "user-ctx"
	ContentTypeJSON = "application/json"

	RouteSignin = "/api/signin"
	RouteSignup = "/api/signup"
	RouteLogout = "/api/logout"

	RouteInboxMessages  = "/api/inbox/"
	RouteOutboxMessages = "/api/outbox/"
	RouteFolderMessages = "/api/folder/{id:[0-9]+}"
)
