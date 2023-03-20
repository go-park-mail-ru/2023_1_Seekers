package config

import "time"

const (
	Port           = "8001"
	ProjectBaseDir = "2023_1_Seekers"
	LogsDir        = "logs/"
	LogsFileName   = "server_"
	LogsTimeFormat = "2006-01-02_15:04:05_MST"
	// TODO to env!
	DBUser       = "postgres"
	DBPassword   = "41219"
	DBHost       = "localhost"
	DBPort       = "5432"
	DBName       = "postgres" //
	DBSchemaName = "mail"

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
	RouteFolderMessages = "/api/v1/folder/{slug}"
	RouteFolders        = "/api/v1/folders"

	// UserService
	MaxImageSize      = 32 << 20
	UserFormNewAvatar = "newAvatar"

	// UserRoutes
	RouteUser                 = "/api/v1/user"
	RouteUserInfo             = "/api/v1/user/info"
	RouteUserPw               = "/api/v1/user/pw"
	RouteUserAvatar           = "/api/v1/user/avatar"
	RouteUserAvatarQueryEmail = "email"
	RouteUserInfoQueryEmail

	PasswordMinLen   = 5
	DefaultAvatar    = "default_avatar.png"
	DefaultAvatarDir = "./cmd/config/static/"
	PostDomain       = "mailbox.ru"
	PostAtDomain     = "@" + PostDomain

	S3AccessKeyID     = "ACCESS_TO_ENV!"
	S3AccessKeySecret = "SECRET_TO_ENV!"
	S3Endpoint        = "http://127.0.0.1:9000"
	S3Region          = "eu-west-2" // В случае с minio - не играет роли, нь указывается обязательно
	S3AvatarBucket    = "avatars"
)

var (
	AllowedHeaders = []string{"Content-Type", "Content-Length", "X-Csrf-Token"}
	AllowedOrigins = []string{"http://127.0.0.1:8002", "http://localhost:8002", "http://89.208.197.150:8002"}
	AllowedMethods = []string{"POST", "GET", "PUT"}
)
