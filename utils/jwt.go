package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateAccessToken(userID string, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"roles":   roles,
		"exp":     time.Now().Add(time.Minute * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken(userID string, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"roles":   roles,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}
