package bootstrap

import (
	"context"
	"proxy-fileserver/adapter"
	"proxy-fileserver/api/controllers"
	"proxy-fileserver/api/middlewares"
	"proxy-fileserver/configs"
	"proxy-fileserver/services"
)

func InitService(ctx context.Context) *Context {
	conf := configs.Get()
	adapterProvider, err := adapter.NewProviderAdapter(ctx, conf)
	if err != nil {
		panic(err)
	}
	serviceProvider := services.NewServiceProvider(adapterProvider)
	controllerProvider := controllers.NewControllerProvider(ctx, serviceProvider.GeFileSystemService())
	middlewareProvider := middlewares.NewMiddlewareProvider(conf.AuthPublicKey)
	context := &Context{
		CommonContext: ctx,
		AppContext: &AppContext{
			AdapterProvider:    adapterProvider,
			ServiceProvider:    serviceProvider,
			ControllerProvider: controllerProvider,
			MiddlewareProvider: middlewareProvider,
		},
	}
	return context

}
