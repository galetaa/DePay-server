package models

// BalanceRequest описывает запрос для получения баланса по адресу
type BalanceRequest struct {
	Address string `json:"address" binding:"required"`
}
