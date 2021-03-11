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
	publicKey        *rsa.PublicKey
	tokenExpiredTime time.Duration
}

func NewAuthService(privateKeyLocation, publicKeyLocation string, expiredTime time.Duration) *AuthService {
	privateKey, err := helpers.LoadPrivateKey(privateKeyLocation)
	if err != nil {
		panic(err)
	}
	publicKey, err := helpers.LoadPublicKey(publicKeyLocation)
	if err != nil {
		panic(err)
	}
	return &AuthService{
		privateKey:       privateKey,
		publicKey:        publicKey,
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

func (s *AuthService) ValidateToken(token string) (bool, error) {
	tokenObject, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return s.publicKey, nil
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
			return false, nil
		}
	default:
		return false, err
	}
}
