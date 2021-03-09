package services

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"proxy-fileserver/common/log"
	"proxy-fileserver/enums"
	"proxy-fileserver/helpers"
	"time"
)

type AuthService struct {
	privateKey       *rsa.PrivateKey
	tokenExpiredTime time.Duration
}

func NewAuthService(privateKeyLocation string, expiredTime time.Duration) *AuthService {
	key, err := helpers.LoadPrivateKey(privateKeyLocation)
	if err != nil {
		panic(err)
	}
	return &AuthService{
		privateKey:       key,
		tokenExpiredTime: expiredTime,
	}
}

func (s *AuthService) GenerateToken() (string, enums.Response) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	claims := jwt.MapClaims{}
	now := time.Now()
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(s.tokenExpiredTime).Unix()
	t.Claims = claims
	token, err := t.SignedString(s.privateKey)
	if err != nil {
		log.Errorf("Failure when generate token at: %v", now)
		return "", enums.ErrorSystem
	}
	return token, nil
}
