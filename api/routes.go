package api

import (
	"github.com/gin-gonic/gin"
	"proxy-fileserver/api/controllers"
	"proxy-fileserver/api/middlewares"
)

func NewRouterWithMiddleware(controllerProvider controllers.ControllerProvider, middlewareProvider middlewares.MiddlewareProvider,
	requiredToken bool) *gin.Engine {
	router := gin.Default()
	router.RouterGroup.POST("/auth", controllerProvider.GetAuthController().GetToken)
	router.RouterGroup.POST("/verify", controllerProvider.GetAuthController().ValidateToken)
	router.NoRoute(controllerProvider.GetStreamFileController().GetFile)
	if requiredToken {
		router.Use(middlewareProvider.GetAuthorizationProcessor().ValidateRequest)
	}
	return router
}
