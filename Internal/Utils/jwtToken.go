package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("SUPER_SECRET")

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Token Generating
func generate(userID uint, role string, ttl time.Duration) (string, error) {
	// claims := jwt.MapClaims{
	// 	"userID": userID ,
	// 	"exp": time.Now().Add(ttl).Unix(),
	// }
	// return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	exp := time.Now().Add(ttl)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	return signed, err
}

// AccessToken
func GenerateAccess(userID uint, role string) (string, error) {
	return generate(userID, role, 15*time.Minute)
}

// RefreshToken
func GenerateRefresh(userID uint, role string) (string, error) {
	return generate(userID, role, 7*24*time.Hour)
}

// Parse
func Parse(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
}

func VerifyToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
