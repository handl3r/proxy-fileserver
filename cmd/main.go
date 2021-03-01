package main

import (
	"errors"
	"proxy-fileserver/common/config"
	"proxy-fileserver/common/log"
	"proxy-fileserver/configs"
)

func initLogger() log.Logging {
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
	_ = initLogger()
	log.Infof("Test me")
	log.Infof("Testme: %d, %s", 1, "thai123")
	log.Error("Oh shit, this is a error")
	log.Errorf("Oh shit, this is a error: %v", errors.New("ANSASASAJIS ERR"))
}
