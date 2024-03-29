package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"time"
)

type Config struct {
	Api struct {
		Port              string `yaml:"port" env-default:"8001"`
		MetricsPort       string `yaml:"metrics_port" env-default:"9001"`
		MetricsName       string `yaml:"metrics_name" env-default:"api"`
		MailHubCtx        string `yaml:"mail_hub_ctx" env-default:"hub"`
		MailTplDir        string `yaml:"mail_tpl_dir" env-default:"./internal/api/http/templates/"`
		Host              string `yaml:"host" env-default:"https://mailbx.ru"`
		MailPreviewMaxLen int    `yaml:"mail_preview_max_len" env-default:"60"`
	} `yaml:"api"`

	Project struct {
		ProjectBaseDir string `yaml:"project_base_dir" env-default:"2023_1_Seekers"`
	} `yaml:"project"`

	Logger struct {
		LogsDir                 string `yaml:"logs_dir" env-default:"logs/app/"`
		LogsApiFileName         string `yaml:"logs_api_file_name"`
		LogsAuthFileName        string `yaml:"logs_auth_file_name"`
		LogsFileServiceFileName string `yaml:"logs_file_service_file_name"`
		LogsUserFileName        string `yaml:"logs_user_file_name"`
		LogsMailFileName        string `yaml:"logs_mail_file_name"`
		LogsUseStdOut           *bool  `yaml:"logs_use_std_out" env-default:"true"`
		LogsTimeFormat          string `yaml:"logs_time_format" env-default:"2006-01-02_15:04:05_MST"`
	} `yaml:"logger"`

	FileGPRCService struct {
		Port        string `yaml:"port" env-default:"8005"`
		Host        string `yaml:"host" env-default:"127.0.0.1"`
		MetricsPort string `yaml:"metrics_port" env-default:"9003"`
		MetricsName string `yaml:"metrics_name" env-default:"file_service"`
	} `yaml:"file_grpc_service"`

	UserGRPCService struct {
		Port        string `yaml:"port" env-default:"8006"`
		Host        string `yaml:"host" env-default:"127.0.0.1"`
		MetricsPort string `yaml:"metrics_port" env-default:"9005"`
		MetricsName string `yaml:"metrics_name" env-default:"user"`
	} `yaml:"user_grpc_service"`

	AuthGRPCService struct {
		Port        string `yaml:"port" env-default:"8007"`
		Host        string `yaml:"host" env-default:"127.0.0.1"`
		MetricsPort string `yaml:"metrics_port" env-default:"9002"`
		MetricsName string `yaml:"metrics_name" env-default:"auth"`
	} `yaml:"auth_grpc_service"`

	MailGRPCService struct {
		Port        string `yaml:"port" env-default:"8008"`
		Host        string `yaml:"host" env-default:"127.0.0.1"`
		MetricsPort string `yaml:"metrics_port" env-default:"9004"`
		MetricsName string `yaml:"metrics_name" env-default:"mail"`
	} `yaml:"mail_grpc_service"`

	SmtpServer struct {
		Port               string        `yaml:"port" env-default:"25"`
		Domain             string        `yaml:"domain" env-default:"mailbx"`
		ReadTimeout        time.Duration `yaml:"read_timeout" env-default:"10s"`
		WriteTimeout       time.Duration `yaml:"write_timeout_timeout" env-default:"10s"`
		MaxMessageBytes    int           `yaml:"max_message_bytes" env-default:"104857600"`
		MaxRecipients      int           `yaml:"max_recipients" env-default:"50"`
		AllowInsecureAuth  bool          `yaml:"allow_insecure_auth" env-default:"false"`
		CertFile           string        `yaml:"cert_file"`
		KeyFile            string        `yaml:"key_file"`
		DkimPrivateKeyFile string        `yaml:"dkim_private_key_file"`
		SecretPassword     string        `env:"SMTP_SECRET_PASSWORD"`
	} `yaml:"smtp_server"`

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
		S3AttachBucket string `yaml:"s3_attach_bucket"`
	} `yaml:"s3"`

	Sessions struct {
		CookieName string        `yaml:"cookie_name" env-default:"MailBxSession"`
		CookieTTL  time.Duration `yaml:"cookie_ttl" env-default:"2400h0m0s"` //time.Hour * 24 * 100 // 100 days
		CookiePath string        `yaml:"cookie_path" env-default:"/"`
		CSRFHeader string        `yaml:"csrf_header" env-default:"Csrf-Token"`
		CookieLen  int           `yaml:"cookie_len" env-deault:"32"`
	} `yaml:"sessions"`

	Routes struct {
		RoutePrefix string `yaml:"route_prefix"`

		//AuthRoutes
		RouteSignin string `yaml:"route_signin" env-default:"/signin"`
		RouteSignup string `yaml:"route_signup" env-default:"/signup"`
		RouteLogout string `yaml:"route_logout" env-default:"/logout"`
		RouteAuth   string `yaml:"route_auth" env-default:"/auth"`
		RouteCSRF   string `yaml:"route_csrf" env-default:"/csrf"`

		// MailRoutes
		RouteMessage                       string `yaml:"route_message" env-default:"/message/{id:[0-9]+}"`
		RouteAttach                        string `yaml:"route_attach" env-default:"/attach/{id:[0-9]+}"`
		RouteAttachB64                     string `yaml:"route_attach_b64" env-default:"/attach/{id:[0-9]+}/b64"`
		RouteMessageAttaches               string `yaml:"route_message_attaches" env-default:"/message/{id:[0-9]+}/attaches"`
		RouteExternalAttach                string `yaml:"route_external_attach" env-default:"/external/attach/{id:[0-9]+}"`
		RoutePreviewAttach                 string `yaml:"route_preview_attach" env-default:"/attach/{id:[0-9]+}/preview"`
		RouteSendMessage                   string `yaml:"route_send_message" env-default:"/message/send"`
		RouteSaveDraftMessage              string `yaml:"route_save_draft_message" env-default:"/message/save"`
		RouteReadMessage                   string `yaml:"route_read_message" env-default:"/message/{id:[0-9]+}/read"`
		RouteUnreadMessage                 string `yaml:"route_unread_message" env-default:"/message/{id:[0-9]+}/unread"`
		RouteMoveToFolder                  string `yaml:"route_move_to_folder" env-default:"/message/{id:[0-9]+}/move"`
		RouteQueryFromFolderSlug           string `yaml:"route_move_to_folder_query_from_folder" env-default:"fromFolder"`
		RouteMoveToFolderQueryToFolderSlug string `yaml:"route_move_to_folder_query_to_folder" env-default:"toFolder"`
		RouteGetFolders                    string `yaml:"route_get_folders" env-default:"/folders"`
		RouteFolder                        string `yaml:"route_folder" env-default:"/folder/{slug}"`
		RouteSearch                        string `yaml:"route_search" env-default:"/messages/search"`
		RouteRecipients                    string `yaml:"route_recipients" env-default:"/recipients/search"`
		RouteCreateFolder                  string `yaml:"route_create_folder" env-default:"/folder/create"`
		RouteEditFolder                    string `yaml:"route_edit_folder" env-default:"/folder/{slug}/edit"`
		RouteWS                            string `yaml:"route_ws" env-default:"/ws"`
		RouteWsQueryEmail                  string `yaml:"route_ws_query_email" env-default:"email"`
		RouteGetFoldersIsCustom            string `yaml:"route_get_folders_is_custom" env-default:"custom"`
		QueryAccessKey                     string `yaml:"query_access_key" env-default:"accessKey"`

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
		ExternalUserPassword  string `env:"EXTERNAL_USER_PASSWORD"`
		AvatarTTFPath         string `yaml:"avatar_ttf_path" env-default:"./cmd/config/wqy-zenhei.ttf"`
	} `yaml:"user_service"`

	Password struct {
		PasswordMinLen  int `yaml:"password_min_len" env-default:"5"`
		PasswordSaltLen int `yaml:"password_salt_len" env-default:"10"`
	} `yaml:"password"`

	Mail struct {
		PostDomain   string `yaml:"post_domain" env-default:"mailbx.ru"`
		PostAtDomain string `yaml:"post_at_domain" env-default:"@mailbx.ru"`
	} `yaml:"mail"`

	Cors struct {
		AllowedHeaders []string `yaml:"allowed_headers"`
		AllowedOrigins []string `yaml:"allowed_origins"`
		AllowedMethods []string `yaml:"allowed_methods"`
	} `yaml:"cors"`
}

func Parse(path string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, errors.Wrap(err, "parse config")
	}

	return &cfg, nil
}
