package configs

import (
	"proxy-fileserver/common/config"
	"time"
)

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

	HttpPort      string
	RequiredToken bool

	InteractiveMode        bool
	GoogleDriveOAuthConfig GoogleDriveOAuth2Config
}

type GoogleDriveOAuth2Config struct {
	CredentialFile string
	TokenFile      string
	Enable         bool
}

var Common *Config

func Get() *Config {
	return Common
}

func LoadConfigs() {
	//cacheTimeLocalFileSystem, err := config.GetTimeDuration("CACHE_TIME_LOCAL_FILE_SYSTEM")
	//if err != nil {
	//	panic(err)
	//}
	requiredToken, err := config.GetBoolWithD("REQUIRED_TOKEN", true)
	if err != nil {
		panic(err)
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

		HttpPort:      config.GetString("HTTP_PORT"),
		RequiredToken: requiredToken,

		InteractiveMode:        interactiveMode,
		GoogleDriveOAuthConfig: gOAuth2Config,
	}
}
