package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"todo-list/config"
)

var secretKey = config.SecretKey

func GenerateToken(userID string, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
