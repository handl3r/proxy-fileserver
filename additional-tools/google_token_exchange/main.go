package main

import (
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"io/ioutil"
	"proxy-fileserver/adapter"
	"proxy-fileserver/common/config"
	"proxy-fileserver/common/log"
	"proxy-fileserver/configs"
)

func _initLogger() log.Logging {
	logger, err := log.NewLogger()
	if err != nil {
		panic("Error when init logger")
	}
	log.RegisterGlobal(logger)
	return logger
}

func main() {
	config.LoadEnvironments()
	configs.LoadConfigs()
	conf := configs.Get()
	_initLogger()
	credentials, err := ioutil.ReadFile(conf.GoogleDriveOAuthConfig.CredentialFile)
	if err != nil {
		panic(err)
	}
	gConfig, err := google.ConfigFromJSON(credentials, drive.DriveReadonlyScope, drive.DriveMetadataScope)
	if err != nil {
		panic(err)
	}
	_ = adapter.GetDriveClient(gConfig, conf.GoogleDriveOAuthConfig.TokenFile, true)
	log.Infof("You can find your access token in token.json file at %v.\n"+
		"Please delete token file before run this tool if you want to generate new new token", conf.GoogleDriveOAuthConfig.TokenFile)
}
