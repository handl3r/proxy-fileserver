package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"proxy-fileserver/api/dtos"
	"proxy-fileserver/api/validation"
	"proxy-fileserver/common/log"
	"proxy-fileserver/enums"
	"proxy-fileserver/services"
)

type AuthController struct {
	AuthService *services.AuthService
	Validator   *validation.Validator
}

func NewAuthController(authService *services.AuthService, validator *validation.Validator) *AuthController {
	return &AuthController{
		AuthService: authService,
		Validator:   validator,
	}
}

func (c *AuthController) Home(ctx *gin.Context) {
	log.Errorf("test-error: %s", "thai")
	ctx.JSON(http.StatusOK, `Hide on thesis degree! Get me on https://handl3r.netlify.app or https://github.com/handl3r or https://www.linkedin.com/in/thaibuixuan`)
}

func (c *AuthController) GetToken(ctx *gin.Context) {
	var createTokenRequest *dtos.CreateTokenRequest
	err := ctx.ShouldBindJSON(&createTokenRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, struct {
			Message string `json:"message"`
		}{
			Message: "Invalid request body",
		})
		return
	}

	if err = c.Validator.ValidateStruct(createTokenRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, struct {
			Message string `json:"message"`
		}{
			Message: "invalid request: " + err.Error(),
		})
		return
	}
	var token string
	var errRes enums.Response
	switch createTokenRequest.Type {
	case enums.MediumTokenType:
		token, errRes = c.AuthService.GenerateToken()
	case enums.HighTokenType:
		token, errRes = c.AuthService.GenerateTokenWithPath(createTokenRequest.Path)
	}
	if errRes != nil {
		ctx.AbortWithStatusJSON(errRes.GetCode(), errRes.GetMessage())
		return
	}
	ctx.JSON(http.StatusOK, dtos.CreateTokenResponse{
		Token: token,
	})
}

func (c *AuthController) ValidateToken(ctx *gin.Context) {
	var validateTokenRequest *dtos.ValidateTokenRequest
	err := ctx.ShouldBindJSON(&validateTokenRequest)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err = c.Validator.ValidateStruct(validateTokenRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, struct {
			Message string
		}{
			Message: "invalid request: " + err.Error(),
		})
		return
	}
	valid, errRes := c.AuthService.FullValidateToken(validateTokenRequest.Token, validateTokenRequest.Path)
	if errRes != nil {
		ctx.AbortWithStatusJSON(errRes.GetCode(), errRes.GetMessage())
		return
	}

	if !valid {
		ctx.JSON(http.StatusUnauthorized, nil)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}
