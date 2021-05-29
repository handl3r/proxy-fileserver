package configs

import (
	"proxy-fileserver/common/config"
	"time"
)

const (
	// NoTokenMode no token required
	NoTokenMode = 1
	// MediumTokenMode required token with no strict path
	MediumTokenMode = 2
	// HighTokenMode required token with strict path
	HighTokenMode = 3
)

var MapTokenMode = map[int]bool{
	NoTokenMode:     true,
	MediumTokenMode: true,
	HighTokenMode:   true,
}

type Config struct {
	Env                      string
	SharedRootFolder         string
	SharedRootFolderID       string
	SharedRootFolderLocal    string
	CacheTimeLocalFileSystem int

	CycleTimeCleaner int

	AuthPublicKeyLocation string
	PrivateKeyLocation    string
	ExpiredTimeToken      time.Duration

	MysqlFileInfoTable string
	MysqlUser          string
	MysqlPassword      string
	MysqlPort          string
	MysqlHost          string
	MysqlDatabase      string

	HttpPort  string
	TokenMode int

	InteractiveMode        bool
	GoogleDriveOAuthConfig GoogleDriveOAuth2Config
	TelegramBotConfig      TelegramBotConfig
}

type GoogleDriveOAuth2Config struct {
	CredentialFile string
	TokenFile      string
	Enable         bool
}

type TelegramBotConfig struct {
	BaseURL   string
	BotToken  string
	ChannelID string
}

var Common *Config

func Get() *Config {
	return Common
}

func LoadConfigs() {
	tokenMode := config.GetIntWithDefault("TOKEN_MODE", 1)
	if _, ok := MapTokenMode[tokenMode]; !ok {
		panic("Invalid TokenMode")
	}
	expiredTimeToken, err := config.GetTimeDuration("EXPIRED_TIME_TOKEN")
	if err != nil {
		panic(err)
	}
	gOAuth2Enable, err := config.GetBoolWithD("GOOGLE_OAUTH2_ENABLE", true)
	if err != nil {
		panic(err)
	}
	gOAuth2Config := GoogleDriveOAuth2Config{
		CredentialFile: config.GetString("CREDENTIAL_GOOGLE_OAUTH2_FILE"),
		TokenFile:      config.GetString("TOKEN_GOOGLE_OAUTH2_FILE"),
		Enable:         gOAuth2Enable,
	}
	interactiveMode, err := config.GetBoolWithD("INTERACTIVE_MODE", false)
	if err != nil {
		panic(nil)
	}
	telegramBaseURL := config.GetString("TELEGRAM_BASE_URL")
	telegramBotToken := config.GetString("TELEGRAM_BOT_TOKEN")
	telegramChannelID := config.GetString("TELEGRAM_CHANNEL_ID")

	Common = &Config{
		Env:                      config.GetString("PROXY_SERVER_ENV"),
		SharedRootFolder:         config.GetString("SHARED_ROOT_FOLDER"),
		SharedRootFolderID:       config.GetString("SHARED_ROOT_FOLDER_ID"),
		SharedRootFolderLocal:    config.GetString("SHARED_ROOT_FOLDER_LOCAL"),
		CacheTimeLocalFileSystem: config.GetInt("CACHE_TIME_LOCAL_FILE_SYSTEM"),

		AuthPublicKeyLocation: config.GetString("AUTH_PUBLIC_KEY"),
		PrivateKeyLocation:    config.GetString("PRIVATE_KEY_LOCATION"),
		ExpiredTimeToken:      expiredTimeToken,
		CycleTimeCleaner:      config.GetInt("CYCLE_TIME_CLEANER"),

		MysqlUser:     config.GetString("MYSQL_USER"),
		MysqlPassword: config.GetString("MYSQL_PASSWORD"),
		MysqlPort:     config.GetString("MYSQL_PORT"),
		MysqlHost:     config.GetString("MYSQL_HOST"),
		MysqlDatabase: config.GetString("MYSQL_DATABASE"),

		HttpPort:  config.GetString("HTTP_PORT"),
		TokenMode: tokenMode,

		InteractiveMode:        interactiveMode,
		GoogleDriveOAuthConfig: gOAuth2Config,
		TelegramBotConfig: TelegramBotConfig{
			BaseURL:   telegramBaseURL,
			BotToken:  telegramBotToken,
			ChannelID: telegramChannelID,
		},
	}
}
