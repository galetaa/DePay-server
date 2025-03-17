package models

type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}
