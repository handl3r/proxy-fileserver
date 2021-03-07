package bootstrap

import (
	"context"
	"database/sql"
	"proxy-fileserver/adapter"
	"proxy-fileserver/api/controllers"
	"proxy-fileserver/api/middlewares"
	"proxy-fileserver/configs"
	"proxy-fileserver/repository"
	"proxy-fileserver/services"
)

func InitService(ctx context.Context, db *sql.DB) *Context {
	conf := configs.Get()
	adapterProvider, err := adapter.NewProviderAdapter(ctx, conf)
	if err != nil {
		panic(err)
	}
	repoProvider := repository.NewProviderRepository(db, conf)
	serviceProvider := services.NewServiceProvider(adapterProvider)
	controllerProvider := controllers.NewControllerProvider(ctx, serviceProvider.GeFileSystemService())
	middlewareProvider := middlewares.NewMiddlewareProvider(conf.AuthPublicKey)
	context := &Context{
		CommonContext: ctx,
		AppContext: &AppContext{
			AdapterProvider:    adapterProvider,
			RepoProvider:       repoProvider,
			ServiceProvider:    serviceProvider,
			ControllerProvider: controllerProvider,
			MiddlewareProvider: middlewareProvider,
		},
	}
	return context

}
