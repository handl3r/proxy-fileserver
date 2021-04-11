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

// GenerateToken generate medium level token
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

// GenerateTokenWithPath generate high level token
func (s *AuthService) GenerateTokenWithPath(path string) (string, enums.Response) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	claims := jwt.MapClaims{}
	now := time.Now()
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(s.tokenExpiredTime).Unix()
	claims["path"] = path
	t.Claims = claims
	token, err := t.SignedString(s.privateKey)
	if err != nil {
		log.Errorf("Failure when generate strict token with path: %s, error: %s", path, err)
		return "", enums.ErrorSystem
	}
	return token, nil
}

// ValidateToken validate every token with no care for path
func (s *AuthService) ValidateToken(token string) (bool, error) {
	_, valid, err := s.validateToken(token)
	return valid, err
}

// FullValidateToken validate for both case medium level token and high level token
func (s *AuthService) FullValidateToken(token, path string) (bool, enums.Response) {
	tokenObj, valid, err := s.validateToken(token)
	if err != nil {
		return false, enums.ErrorSystem
	}
	if !valid {
		return false, nil
	}
	claims, ok := tokenObj.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil
	}
	pathClaim, ok := claims["path"]
	if !ok {
		return valid, nil
	}
	pathClaimStr, _ := pathClaim.(string)
	if path != pathClaimStr {
		return false, nil
	}
	return true, nil
}

// ValidateTokenWithPath validate token with a specify path
func (s *AuthService) ValidateTokenWithPath(token, path string) (bool, error) {
	tokenObj, valid, err := s.validateToken(token)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, nil
	}
	claims, ok := tokenObj.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil
	}
	pathClaim, ok := claims["path"]
	if !ok {
		return false, nil
	}
	pathClaimStr, _ := pathClaim.(string)
	if path != pathClaimStr {
		return false, nil
	}
	return true, nil
}

func (s *AuthService) validateToken(token string) (*jwt.Token, bool, error) {
	tokenObject, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return s.publicKey, nil
	})
	switch err.(type) {
	case nil:
		if !tokenObject.Valid {
			return nil, false, nil
		}
		return tokenObject, true, nil
	case *jwt.ValidationError:
		vErr := err.(*jwt.ValidationError)
		switch vErr.Errors {
		case jwt.ValidationErrorExpired:
			return nil, false, nil
		default:
			return nil, false, nil
		}
	default:
		return nil, false, err
	}
}
