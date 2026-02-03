package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("SUPER_SECRET")

//Token Generating
func generate(userID uint, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID ,
		"exp": time.Now().Add(ttl).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

// AccessToken
func GenerateAccess( userID uint) (string, error) {
	return generate(userID, 15*time.Minute)
}

//RefreshToken
func GenerateRefresh(userID uint) (string, error) {
	return generate(userID, 7*24*time.Hour)
}

//Parse
func Parse(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
}