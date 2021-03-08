package bootstrap

import (
	"context"
	"gorm.io/gorm"
	"proxy-fileserver/adapter"
	"proxy-fileserver/api/controllers"
	"proxy-fileserver/api/middlewares"
	"proxy-fileserver/configs"
	"proxy-fileserver/repository"
	"proxy-fileserver/services"
)

func InitService(ctx context.Context, db *gorm.DB) *Context {
	conf := configs.Get()
	adapterProvider, err := adapter.NewProviderAdapter(ctx, conf)
	if err != nil {
		panic(err)
	}
	repoProvider := repository.NewProviderRepository(db)
	serviceProvider := services.NewServiceProvider(adapterProvider, repoProvider)
	controllerProvider := controllers.NewControllerProvider(ctx, serviceProvider.GeFileSystemService())
	middlewareProvider := middlewares.NewMiddlewareProvider(conf.AuthPublicKeyLocation)
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
