package middlewares

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"proxy-fileserver/common/log"
	"proxy-fileserver/helpers"
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
	tokenObject, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return p.publicKey, nil
	})
	switch err.(type) {
	case nil:
		if !tokenObject.Valid {
			return false, nil
		}
		return true, nil
	case *jwt.ValidationError:
		vErr := err.(*jwt.ValidationError)
		switch vErr.Errors {
		case jwt.ValidationErrorExpired:
			return false, nil
		default:
			log.Errorf("Error when validate token: %v", err)
		}
	default:
		return false, err
	}
	return true, nil
}
