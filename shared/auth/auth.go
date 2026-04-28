package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"shared/utils"
)

type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	RefreshTokenHash string
	RefreshExpiresAt time.Time
}

func NewTokenPair(userID string, roles []string, accessTTL time.Duration, refreshTTL time.Duration) (TokenPair, error) {
	accessToken, err := utils.GenerateTokenWithRoles(userID, roles, accessTTL)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := randomToken(32)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		RefreshTokenHash: HashToken(refreshToken),
		RefreshExpiresAt: time.Now().UTC().Add(refreshTTL),
	}, nil
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func NewOpaqueToken(size int) (string, error) {
	return randomToken(size)
}

func randomToken(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
