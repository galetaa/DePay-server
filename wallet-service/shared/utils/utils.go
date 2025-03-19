package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JwtSecret используется для подписи JWT (замените на более безопасное хранение в продакшене)
var JwtSecret = []byte("your_secret_key")

// GenerateToken генерирует JWT для заданного userId с указанной длительностью
func GenerateToken(userId string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}
