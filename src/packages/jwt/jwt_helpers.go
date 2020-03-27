// Package jwt is used for authorization
package jwt

import (
	"crypto/rsa"
	"fmt"
	"os"

	jwtgo "github.com/dgrijalva/jwt-go"
)

//nolint: gochecknoglobals
var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

// PublicKey returns the public key
func PublicKey() *rsa.PublicKey {
	return publicKey
}

// PrivateKey return the priate key
func PrivateKey() *rsa.PrivateKey {
	return privateKey
}

// LoadKeys load JWT key pair
func LoadKeys() error {
	jwtKey := []byte(os.Getenv("JWT_KEY"))
	jwtPub := []byte(os.Getenv("JWT_KEY_PUB"))

	if len(jwtKey) == 0 {
		return fmt.Errorf("jwt private key not found in JWT_KEY")
	}
	if len(jwtPub) == 0 {
		return fmt.Errorf("jwt public key not found in JWT_KEY_PUB")
	}

	return ParseKeys(jwtKey, jwtPub)
}

// ParseKeys parse priv pub key pair as bytes
func ParseKeys(jwtKey, jwtPub []byte) error {
	keyParsed, err := jwtgo.ParseRSAPrivateKeyFromPEM(jwtKey)
	if err != nil {
		return fmt.Errorf("failed to load jwt private key error: %s", err)
	}
	privateKey = keyParsed

	pubParsed, err := jwtgo.ParseRSAPublicKeyFromPEM(jwtPub)
	if err != nil {
		return fmt.Errorf("failed to load jwt public key error: %s", err)
	}
	publicKey = pubParsed
	return nil
}

// LoadPublicKey loads just the jwt public key set in the JWT_KEY_PUB variable, unlike LoadKeys it returns an error if the key cannot be found instead of creating a new one
func LoadPublicKey() error {
	jwtPub := []byte(os.Getenv("JWT_KEY_PUB"))
	if len(jwtPub) == 0 {
		return fmt.Errorf("jwt key not found in environment variable JWT_KEY_PUB")
	}

	pubParsed, err := jwtgo.ParseRSAPublicKeyFromPEM(jwtPub)
	if err != nil {
		return fmt.Errorf("failed to load jwt public key error: %s", err)
	}
	publicKey = pubParsed
	return nil
}
