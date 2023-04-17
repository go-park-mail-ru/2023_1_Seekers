package config

import "time"

const (
	Port           = "8001"
	ProjectBaseDir = "2023_1_Seekers"
	LogsDir        = "logs/app/"
	LogsFileName   = "server_"
	LogsTimeFormat = "2006-01-02_15:04:05_MST"

	// Postgres
	DBUserEnv       = "POSTGRES_USER"
	DBPasswordEnv   = "POSTGRES_PASSWORD"
	DBHostEnv       = "POSTGRES_HOST"
	DBPortEnv       = "POSTGRES_PORT"
	DBNameEnv       = "POSTGRES_DB" //
	DBSchemaNameEnv = "POSTGRES_SCHEMA"
	DBSSLModeEnv    = "POSTGRES_SSL_MODE"

	// Redis
	RedisHostEnv     = "REDIS_HOST"
	RedisPortEnv     = "REDIS_PORT"
	RedisPasswordEnv = "REDIS_PASSWORD"

	//MINIO S3
	S3AccessKeyEnv  = "S3_ACCESS_KEY"
	S3ASecretKeyEnv = "S3_SECRET_KEY"
	S3Region        = "eu-west-2"
	//-----VK cloud solutions--------
	S3Endpoint     = "https://hb.bizmrg.com"
	S3AvatarBucket = "avatars_mailbox_vkcloud"
	//-----MinioS3------------
	//S3Endpoint     = "http://172.28.0.5:9000"
	//S3AvatarBucket = "avatars"

	// Sessions
	CookieName = "MailBoxSession"
	CookieTTL  = time.Hour * 24 * 100 // 100 days
	CookiePath = "/"
	CSRFHeader = "Csrf-Token"

	RoutePrefix = "/api/v1"

	//AuthRoutes
	RouteSignin = RoutePrefix + "/signin"
	RouteSignup = RoutePrefix + "/signup"
	RouteLogout = RoutePrefix + "/logout"
	RouteAuth   = RoutePrefix + "/auth"
	RouteCSRF   = RoutePrefix + "/csrf"
	RoutePw     = RoutePrefix + "/user/pw"

	// MailRoutes
	RouteGetFolderMessages = RoutePrefix + "/folder/{slug}"
	RouteGetFolders        = RoutePrefix + "/folders"
	RouteGetMessage        = RoutePrefix + "/message/{id:[0-9]+}"
	RouteSendMessage       = RoutePrefix + "/message/send"
	RouteReadMessage       = RoutePrefix + "/message/{id:[0-9]+}/read"
	RouteUnreadMessage     = RoutePrefix + "/message/{id:[0-9]+}/unread"

	// UserService
	MaxImageSize          = 32 << 20
	UserFormNewAvatar     = "avatar"
	UserDefaultAvatarSize = 46
	AvatarTTFPath         = "./cmd/config/wqy-zenhei.ttf"

	// UserRoutes
	RouteUser                 = RoutePrefix + "/user"
	RouteUserInfo             = RouteUser + "/info"
	RouteUserAvatar           = RouteUser + "/avatar"
	RouteUserAvatarQueryEmail = "email"
	RouteUserInfoQueryEmail   = "email"

	PasswordMinLen = 5
	//PasswordSaltLen  = 10
	DefaultAvatar = "default_avatar.png"
	PostDomain    = "mailbox.ru"
	PostAtDomain  = "@" + PostDomain
)

var (
	CookieLen       = 32
	PasswordSaltLen = 10
	AllowedHeaders  = []string{"Content-Type", "Content-Length", "X-Csrf-Token", "application/json", "text/xml"}
	AllowedOrigins  = []string{"http://127.0.0.1:8002", "http://localhost:8002", "http://localhost", "http://127.0.0.1",
		"http://89.208.197.150:8002", "https://mailbx.ru", "https://www.mailbx.ru"}
	AllowedMethods = []string{"POST", "GET", "PUT", "DELETE"}
)
