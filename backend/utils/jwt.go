package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret = []byte("SparkPass_Dev-Key-456$%^789XYZ") // Cambia esto en producción!

func GenerateJWT(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(JWTSecret)
}

func ParseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})
}