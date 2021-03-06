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
	CacheTimeLocalFileSystem time.Duration

	AuthPublicKey string

	MysqlFileInfoTable string
}

var Common *Config

func Get() *Config {
	return Common
}

func LoadConfigs() {
	cacheTimeLocalFileSystem, err := config.GetTimeDuration("CACHE_TIME_LOCAL_FILE_SYSTEM")
	if err != nil {
		panic(err)
	}
	Common = &Config{
		Env:                      config.GetString("PROXY_SERVER_ENV"),
		SharedRootFolder:         config.GetString("SHARED_ROOT_FOLDER"),
		SharedRootFolderID:       config.GetString("SHARED_ROOT_FOLDER_ID"),
		SharedRootFolderLocal:    config.GetString("SHARED_ROOT_FOLDER_LOCAL"),
		CacheTimeLocalFileSystem: cacheTimeLocalFileSystem,

		AuthPublicKey: config.GetString("AUTH_PUBLIC_KEY"),

		MysqlFileInfoTable: config.GetString("MYSQL_FILE_INFO_TABLE_NAME"),
	}
}
