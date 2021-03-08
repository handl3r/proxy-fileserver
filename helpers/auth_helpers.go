package helpers

import (
	"crypto/rsa"
	"encoding/pem"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"os"
)

func LoadPublicKey(location string) (*rsa.PublicKey, []byte, error) {
	if _, err := os.Stat(location); err != nil {
		panic(err)
	}
	pub, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, nil, err
	}
	pubPem, _ := pem.Decode(pub)
	if pubPem == nil {
		return nil, nil, errors.New("decode public key fail")
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		return nil, nil, err
	}
	return publicKey, pub, nil
}
