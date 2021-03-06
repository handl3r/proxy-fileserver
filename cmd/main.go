package main

import (
	"context"
	"proxy-fileserver/api"
	"proxy-fileserver/bootstrap"
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
	_initLogger()
	ctx := context.Background()
	appContext := bootstrap.InitService(ctx)
	router := api.NewRouterWithMiddleware(appContext.AppContext.ControllerProvider, appContext.AppContext.MiddlewareProvider)
	_ = router.Run(":8080")

}

//log.Infof("Test me")
//log.Infof("Testme: %d, %s", 1, "thai123")
//log.Error("Oh shit, this is a error")
//log.Errorf("Oh shit, this is a error: %v", errors.New("ANSASASAJIS ERR"))
