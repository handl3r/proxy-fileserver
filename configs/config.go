package configs

import "proxy-fileserver/common/config"

type Config struct {
	Env string
}

var Common *Config

func Get() *Config {
	return Common
}

func LoadConfigs() {
	Common = &Config{
		Env: config.GetString("PROXY_SERVER_ENV"),
	}
}
