package api

import (
	"github.com/gin-gonic/gin"
	"proxy-fileserver/api/controllers"
	"proxy-fileserver/api/middlewares"
)

func NewRouterWithMiddleware(controllerProvider controllers.ControllerProvider, middlewareProvider middlewares.MiddlewareProvider) *gin.Engine {
	router := gin.Default()
	router.NoRoute(controllerProvider.GetStreamFileController().GetFile)
	router.Use(middlewareProvider.GetAuthorizationProcessor().ValidateRequest)
	return router
}
