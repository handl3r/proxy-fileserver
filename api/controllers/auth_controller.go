package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"proxy-fileserver/services"
)

type AuthController struct {
	AuthService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}

func (c *AuthController) GetToken(ctx *gin.Context) {
	token, err := c.AuthService.GenerateToken()
	if err != nil {
		ctx.AbortWithStatusJSON(err.GetCode(), err.GetMessage())
		return
	}
	ctx.JSON(http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}
