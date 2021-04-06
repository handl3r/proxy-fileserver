package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"proxy-fileserver/api/dtos"
	"proxy-fileserver/services"
)

type AuthController struct {
	AuthService     *services.AuthService
	StrictTokenMode bool
}

func NewAuthController(authService *services.AuthService, strictTokenMode bool) *AuthController {
	return &AuthController{
		AuthService:     authService,
		StrictTokenMode: strictTokenMode,
	}
}

func (c *AuthController) GetToken(ctx *gin.Context) {
	if c.StrictTokenMode {
		var createTokenRequest *dtos.CreateTokenRequest
		err := ctx.ShouldBindJSON(&createTokenRequest)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, struct {
				Message string
			}{
				Message: "Invalid request body",
			})
			return
		}
		token, errRes := c.AuthService.GenerateTokenWithPath(createTokenRequest.Path)
		if errRes != nil {
			ctx.AbortWithStatusJSON(errRes.GetCode(), errRes.GetMessage())
			return
		}
		ctx.JSON(http.StatusOK, struct {
			Token string `json:"token"`
		}{
			Token: token,
		})
		return
	}

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

	if c.StrictTokenMode {
		valid, err := c.AuthService.ValidateTokenWithPath(token.Token, token.Path)
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
			return
		} else {
			ctx.JSON(http.StatusUnauthorized, nil)
			return
		}
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
