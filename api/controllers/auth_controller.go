package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"proxy-fileserver/api/dtos"
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

func (c *AuthController) ValidateToken(ctx *gin.Context) {
	var token *dtos.Token
	err := ctx.ShouldBindJSON(&token)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	valid, err := c.AuthService.ValidateToken(token.Token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, struct {
			Message string
		}{
			Message: "System error. Please contact admin!",
		})
		return
	}
	if valid {
		ctx.JSON(http.StatusOK, nil)
	} else {
		ctx.JSON(http.StatusUnauthorized, nil)
	}
}
