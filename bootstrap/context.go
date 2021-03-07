package bootstrap

import (
	"context"
	"proxy-fileserver/adapter"
	"proxy-fileserver/api/controllers"
	"proxy-fileserver/api/middlewares"
	"proxy-fileserver/repository"
	"proxy-fileserver/services"
)

type Context struct {
	CommonContext context.Context
	AppContext    *AppContext
}

type AppContext struct {
	AdapterProvider    adapter.ProviderAdapter
	RepoProvider       repository.ProviderRepository
	ServiceProvider    services.ServiceProvider
	ControllerProvider controllers.ControllerProvider
	MiddlewareProvider middlewares.MiddlewareProvider
}
