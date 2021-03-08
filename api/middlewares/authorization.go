package middlewares

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"proxy-fileserver/common/log"
	"proxy-fileserver/helpers"
	"strings"
)

type AuthorizationProcessor struct {
	publicKey         *rsa.PublicKey
	publicKeyLocation string
}

func NewAuthorizationProcessor(publicKeyLocation string) *AuthorizationProcessor {
	publicKey, err := helpers.LoadPublicKey(publicKeyLocation)
	if err != nil {
		panic(err)
	}
	return &AuthorizationProcessor{
		publicKey:         publicKey,
		publicKeyLocation: publicKeyLocation,
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
	valid, err := p.validateToken(listToken[0])
	if err != nil {
		log.Errorf("Can not validate token with error: %s", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Next()
}

func (p *AuthorizationProcessor) validateToken(token string) (bool, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false, nil
	}
	err := jwt.SigningMethodRS256.Verify(strings.Join(parts[0:2], "."), parts[2], p.publicKey)
	if err != nil {
		return false, err
	}
	return true, nil
}
