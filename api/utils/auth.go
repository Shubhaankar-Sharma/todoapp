package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// EncodeAuthToken signs authentication token
func EncodeAuthToken(uid string, secret string) (string, error) {
	claims := jwt.MapClaims{}
	claims["uid"] = uid
	claims["IssuedAt"] = time.Now().Unix()
	claims["ExpiresAt"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	return token.SignedString([]byte(secret))
}