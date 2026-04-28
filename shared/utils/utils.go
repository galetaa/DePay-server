package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JwtSecret используется для подписи JWT (замените на более безопасное хранение в продакшене)
var JwtSecret = []byte(GetJWTSecret())

func GetJWTSecret() string {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return secret
	}
	return "dev_secret_change_me"
}

// GenerateToken генерирует JWT для заданного userId с указанной длительностью
func GenerateToken(userId string, duration time.Duration) (string, error) {
	return GenerateTokenWithRoles(userId, []string{"user"}, duration)
}

func GenerateTokenWithRoles(userId string, roles []string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userId,
		"roles": roles,
		"type":  "access",
		"exp":   time.Now().Add(duration).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(GetJWTSecret()))
}
