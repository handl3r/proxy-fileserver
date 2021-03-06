package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthorizationProcessor struct {
	publicKey string
}

func NewAuthorizationProcessor(publicKey string) *AuthorizationProcessor {
	return &AuthorizationProcessor{
		publicKey: publicKey,
	}
}

func (p *AuthorizationProcessor) ValidateRequest(c *gin.Context) {
	rawQuery := c.Request.URL.Query()
	listToken, ok := rawQuery["token"]
	if !ok || len(listToken) == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// TODO validate here
	c.Next()
}
