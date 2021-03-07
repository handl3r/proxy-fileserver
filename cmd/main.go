package main

import (
	"context"
	"github.com/robfig/cron"
	"proxy-fileserver/api"
	"proxy-fileserver/bootstrap"
	"proxy-fileserver/common/config"
	"proxy-fileserver/common/lock"
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
	lock.InitMapLock()

	dbConnection := bootstrap.InitDBConnection(conf.MysqlUser, conf.MysqlPassword, conf.MysqlHost, conf.MysqlPort, conf.MysqlDatabase)

	ctx := context.Background()
	appContext := bootstrap.InitService(ctx, dbConnection)

	c := cron.New()
	cleaner := api.NewCleaner(appContext.AppContext.RepoProvider.GetFileInfoRepository(), configs.Get().CacheTimeLocalFileSystem,
		appContext.AppContext.AdapterProvider.GetLocalFileSystem())
	_ = c.AddFunc("@every 1m", cleaner.Run)
	c.Start()

	router := api.NewRouterWithMiddleware(appContext.AppContext.ControllerProvider, appContext.AppContext.MiddlewareProvider)
	_ = router.Run(":8080")
	//http.HandleFunc("/", appContext.AppContext.ControllerProvider.GetStreamFileController().GetFileBasicHttp)
	//_ = http.ListenAndServe(":8080", nil)
}
