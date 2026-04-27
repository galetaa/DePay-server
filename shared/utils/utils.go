package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GetJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}

// GenerateToken генерирует JWT для заданного userId с указанной длительностью
func GenerateToken(userId string, duration time.Duration) (string, error) {
	return GenerateTokenWithRoles(userId, []string{"user"}, duration)
}

func GenerateTokenWithRoles(userId string, roles []string, duration time.Duration) (string, error) {
	secret := GetJWTSecret()
	if secret == "" {
		return "", errors.New("JWT_SECRET is required")
	}

	claims := jwt.MapClaims{
		"sub":   userId,
		"roles": roles,
		"type":  "access",
		"exp":   time.Now().Add(duration).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
