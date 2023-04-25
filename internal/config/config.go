package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"time"
)

type Config struct {
	Api struct {
		Port string `yaml:"port" env-default:"8001"`
	} `yaml:"api"`

	Project struct {
		ProjectBaseDir string `yaml:"project_base_dir" env-default:"2023_1_Seekers"`
	} `yaml:"project"`

	Logger struct {
		LogsDir        string `yaml:"logs_dir" env-default:"logs/app/"`
		LogsFileName   string `yaml:"logs_file_name" env-default:"server_"`
		LogsTimeFormat string `yaml:"logs_time_format" env-default:"2006-01-02_15:04:05_MST"`
	} `yaml:"logger"`

	FileGPRCService struct {
		Port string `yaml:"port" env-default:"8005"`
		Host string `yaml:"host" env-default:"127.0.0.1"`
	} `yaml:"file_grpc_service"`

	UserGRPCService struct {
		Port string `yaml:"port" env-default:"8006"`
		Host string `yaml:"host" env-default:"127.0.0.1"`
	} `yaml:"user_grpc_service"`

	AuthGRPCService struct {
		Port string `yaml:"port" env-default:"8007"`
		Host string `yaml:"host" env-default:"127.0.0.1"`
	} `yaml:"auth_grpc_service"`

	MailGRPCService struct {
		Port string `yaml:"port" env-default:"8008"`
		Host string `yaml:"host" env-default:"127.0.0.1"`
	} `yaml:"mail_grpc_service"`

	DB struct {
		DBUser       string `env:"POSTGRES_USER"`
		DBPassword   string `env:"POSTGRES_PASSWORD"`
		DBHost       string `env:"POSTGRES_HOST"`
		DBPort       string `env:"POSTGRES_PORT"`
		DBName       string `env:"POSTGRES_DB"`
		DBSchemaName string `env:"POSTGRES_SCHEMA"`
		DBSSLMode    string `env:"POSTGRES_SSL_MODE"`
	} `yaml:"db"`

	Redis struct {
		RedisHost     string `env:"REDIS_HOST"`
		RedisPort     string `env:"REDIS_PORT"`
		RedisPassword string `env:"REDIS_PASSWORD"`
	} `yaml:"redis"`

	S3 struct {
		S3AccessKey  string `env:"S3_ACCESS_KEY"`
		S3ASecretKey string `env:"S3_SECRET_KEY"`
		S3Region     string `yaml:"s3_region" env-default:"eu-west-2"`
		//-----VK cloud solutions--------
		S3Endpoint     string `yaml:"s3_endpoint"`
		S3AvatarBucket string `yaml:"s3_avatar_bucket"`
	} `yaml:"s3"`

	Sessions struct {
		CookieName string        `yaml:"cookie_name" env-default:"MailBoxSession"`
		CookieTTL  time.Duration `yaml:"cookie_ttl" env-default:"2400h0m0s"` //time.Hour * 24 * 100 // 100 days
		CookiePath string        `yaml:"cookie_path" env-default:"/"`
		CSRFHeader string        `yaml:"csrf_header" env-default:"Csrf-Token"`
		CookieLen  int           `yaml:"cookie_len" env-deault:"32"`
	} `yaml:"sessions"`

	Routes struct {
		RoutePrefix string `yaml:"route_prefix" env-default:"/api/v1"`

		//AuthRoutes
		RouteSignin string `yaml:"route_signin" env-default:"/signin"`
		RouteSignup string `yaml:"route_signup" env-default:"/signup"`
		RouteLogout string `yaml:"route_logout" env-default:"/logout"`
		RouteAuth   string `yaml:"route_auth" env-default:"/auth"`
		RouteCSRF   string `yaml:"route_csrf" env-default:"/csrf"`

		// MailRoutes
		RouteMessage                     string `yaml:"route_message" env-default:"/message/{id:[0-9]+}"`
		RouteSendMessage                 string `yaml:"route_send_message" env-default:"/message/send"`
		RouteSaveDraftMessage            string `yaml:"route_save_draft_message" env-default:"/message/save"`
		RouteReadMessage                 string `yaml:"route_read_message" env-default:"/message/{id:[0-9]+}/read"`
		RouteUnreadMessage               string `yaml:"route_unread_message" env-default:"/message/{id:[0-9]+}/unread"`
		RouteMoveToFolder                string `yaml:"route_move_to_folder" env-default:"/message/{id:[0-9]+}/move"`
		RouteMoveToFolderQueryFolderSlug string `yaml:"route_move_to_folder_query_folder_slug" env-default:"folderSlug"`
		RouteGetFolders                  string `yaml:"route_get_folders" env-default:"/folders"`
		RouteFolder                      string `yaml:"route_folder" env-default:"/folder/{slug}"`
		RouteCreateFolder                string `yaml:"route_create_folder" env-default:"/folder/create"`
		RouteEditFolder                  string `yaml:"route_edit_folder" env-default:"/folder/{slug}/edit"`
		RouteGetFoldersIsCustom          string `yaml:"route_get_folders_is_custom" env-default:"custom"`

		// UserRoutes
		RouteUser                 string `yaml:"route_user" env-default:"/user"`
		RouteUserInfo             string `yaml:"route_user_info" env-default:"/user/info"`
		RouteUserAvatar           string `yaml:"route_user_avatar" env-default:"/user/avatar"`
		RoutePw                   string `yaml:"route_pw" env-default:"/user/pw"`
		RouteUserAvatarQueryEmail string `yaml:"route_user_avatar_query_email" env-default:"email"`
		RouteUserInfoQueryEmail   string `yaml:"route_user_info_query_email" env-default:"email"`
	} `yaml:"routes"`

	UserService struct {
		MaxImageSize          int    `yaml:"max_image_size" env-default:"2147483648"` //  2 << 30
		UserFormNewAvatar     string `yaml:"user_form_new_avatar" env-default:"avatar"`
		UserDefaultAvatarSize int    `yaml:"user_default_avatar_size" env-default:"46"`
		DefaultAvatar         string `yaml:"default_avatar" env-default:"default_avatar.png"`
		AvatarTTFPath         string `yaml:"avatar_ttf_path" env-default:"./cmd/config/wqy-zenhei.ttf"`
	} `yaml:"user_service"`

	Password struct {
		PasswordMinLen  int `yaml:"password_min_len" env-default:"5"`
		PasswordSaltLen int `yaml:"password_salt_len" env-default:"10"`
	} `yaml:"password"`

	Mail struct {
		PostDomain   string `yaml:"post_domain" env-default:"mailbox.ru"`
		PostAtDomain string `yaml:"post_at_domain" env-default:"@mailbox.ru"`
	} `yaml:"mail"`

	Cors struct {
		AllowedHeaders []string `yaml:"allowed_headers" env-default:"Content-Type,Content-Length,X-Csrf-Token,application/json,text/xml"`
		AllowedOrigins []string `yaml:"allowed_origins" env-default:"http://127.0.0.1:8002,http://localhost:8002,http://localhost,http://127.0.0.1,http://89.208.197.150:8002,https://mailbx.ru,https://www.mailbx.ru"`
		AllowedMethods []string `yaml:"allowed_methods" env-default:"POST,GET,PUT,DELETE"`
	} `yaml:"cors"`
}

func Parse(path string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, errors.Wrap(err, "parse config")
	}

	return &cfg, nil
}
