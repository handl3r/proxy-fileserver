package configs

import (
	"proxy-fileserver/common/config"
)

type Config struct {
	Env                      string
	SharedRootFolder         string
	SharedRootFolderID       string
	SharedRootFolderLocal    string
	CacheTimeLocalFileSystem int

	CycleTimeCleaner int

	AuthPublicKey string

	MysqlFileInfoTable string
	MysqlUser          string
	MysqlPassword      string
	MysqlPort          string
	MysqlHost          string
	MysqlDatabase      string
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

	Common = &Config{
		Env:                      config.GetString("PROXY_SERVER_ENV"),
		SharedRootFolder:         config.GetString("SHARED_ROOT_FOLDER"),
		SharedRootFolderID:       config.GetString("SHARED_ROOT_FOLDER_ID"),
		SharedRootFolderLocal:    config.GetString("SHARED_ROOT_FOLDER_LOCAL"),
		CacheTimeLocalFileSystem: config.GetInt("CACHE_TIME_LOCAL_FILE_SYSTEM"),

		AuthPublicKey: config.GetString("AUTH_PUBLIC_KEY"),
		CycleTimeCleaner: config.GetInt("CYCLE_TIME_CLEANER"),

		MysqlFileInfoTable: config.GetString("MYSQL_FILE_INFO_TABLE_NAME"),
		MysqlUser:          config.GetString("MYSQL_USER"),
		MysqlPassword:      config.GetString("MYSQL_PASSWORD"),
		MysqlPort:          config.GetString("MYSQL_PORT"),
		MysqlHost:          config.GetString("MYSQL_HOST"),
		MysqlDatabase:      config.GetString("MYSQL_DATABASE"),
	}
}
