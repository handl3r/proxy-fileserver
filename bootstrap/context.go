package bootstrap

import (
	"context"
	"proxy-fileserver/adapter"
	"proxy-fileserver/api/controllers"
	"proxy-fileserver/api/middlewares"
	"proxy-fileserver/services"
)

type Context struct {
	CommonContext context.Context
	AppContext *AppContext
}

type AppContext struct {
	AdapterProvider adapter.ProviderAdapter
	serviceProvider services.ServiceProvider
	controllerProvider controllers.ControllerProvider
	middlewareProvider middlewares.MiddlewareProvider
}
