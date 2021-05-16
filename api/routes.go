package api

import (
	"github.com/gin-gonic/gin"
	"proxy-fileserver/api/controllers"
	"proxy-fileserver/api/middlewares"
	"proxy-fileserver/configs"
)

func NewRouterWithMiddleware(controllerProvider controllers.ControllerProvider, middlewareProvider middlewares.MiddlewareProvider,
	tokenMode int) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.NoCache)
	router.RouterGroup.POST("/auth", controllerProvider.GetAuthController().GetToken)
	router.RouterGroup.POST("/verify", controllerProvider.GetAuthController().ValidateToken)
	router.RouterGroup.GET("/", controllerProvider.GetAuthController().Home)
	router.NoRoute(controllerProvider.GetStreamFileController().GetFile)
	switch tokenMode {
	case configs.MediumTokenMode:
		router.Use(middlewareProvider.GetAuthorizationProcessor().ValidateRequestWithToken)
	case configs.HighTokenMode:
		router.Use(middlewareProvider.GetAuthorizationProcessor().ValidateRequestWithStrictToken)
	}
	return router
}
