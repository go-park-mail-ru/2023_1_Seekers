package config

import "time"

const (
	Port           = "8001"
	ProjectBaseDir = "2023_1_Seekers"
	LogsDir        = "logs/"
	LogsFileName   = "server_"
	LogsTimeFormat = "2006-01-02_15:04:05_MST"

	// Postgres
	DBUserEnv       = "POSTGRES_USER"
	DBPasswordEnv   = "POSTGRES_PASSWORD"
	DBHostEnv       = "POSTGRES_HOST"
	DBPortEnv       = "POSTGRES_PORT"
	DBNameEnv       = "POSTGRES_DB" //
	DBSchemaNameEnv = "POSTGRES_SCHEMA"

	// Redis
	RedisHostEnv     = "REDIS_HOST"
	RedisPortEnv     = "REDIS_PORT"
	RedisPasswordEnv = "REDIS_PASSWORD"

	//MINIO S3
	S3AccessKeyEnv  = "MINIO_ACCESS_KEY"
	S3ASecretKeyEnv = "MINIO_SECRET_KEY"
	S3Endpoint      = "http://172.28.0.5:9000"
	S3Region        = "eu-west-2" // В случае с minio - не играет роли, нь указывается обязательно
	S3AvatarBucket  = "avatars"

	// Sessions
	CookieName = "MailBoxSession"
	CookieTTL  = time.Hour * 24 * 100 // 100 days
	CookieLen  = 32
	CookiePath = "/"

	//AuthRoutes
	RouteSignin = "/api/v1/signin"
	RouteSignup = "/api/v1/signup"
	RouteLogout = "/api/v1/logout"

	// MailRoutes
	RouteGetFolderMessages = "/api/v1/folder/{slug}"
	RouteGetFolders        = "/api/v1/folders"
	RouteGetMessage        = "/api/v1/message/{id:[0-9]+}"
	RouteSendMessage       = "/api/v1/message/send"
	RouteReadMessage       = "/api/v1/message/{id:[0-9]+}/read"
	RouteUnreadMessage     = "/api/v1/message/{id:[0-9]+}/unread"

	// UserService
	MaxImageSize      = 32 << 20
	UserFormNewAvatar = "newAvatar"

	// UserRoutes
	RouteUser                 = "/api/v1/user"
	RouteUserInfo             = "/api/v1/user/info"
	RouteUserPw               = "/api/v1/user/pw"
	RouteUserAvatar           = "/api/v1/user/avatar"
	RouteUserAvatarQueryEmail = "email"
	RouteUserInfoQueryEmail   = "email"

	PasswordMinLen   = 5
	DefaultAvatar    = "default_avatar.png"
	DefaultAvatarDir = "./cmd/config/static/"
	PostDomain       = "mailbox.ru"
	PostAtDomain     = "@" + PostDomain
)

var (
	AllowedHeaders = []string{"Content-Type", "Content-Length", "X-Csrf-Token"}
	AllowedOrigins = []string{"http://127.0.0.1:8002", "http://localhost:8002", "http://89.208.197.150:8002"}
	AllowedMethods = []string{"POST", "GET", "PUT"}
)
