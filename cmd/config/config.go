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

	RouteSignin = "/api/v1/signin"
	RouteSignup = "/api/v1/signup"
	RouteLogout = "/api/v1/logout"

	RouteInboxMessages  = "/api/v1/inbox/"
	RouteOutboxMessages = "/api/v1/outbox/"
	RouteFolderMessages = "/api/v1/folder/{id:[0-9]+}"

	PasswordMinLen = 5
)
