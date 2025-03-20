package models

import "time"

// Transaction представляет транзакцию, инициированную терминалом.
// Все расчёты ведутся в единой базовой валюте (например, ETH).
type Transaction struct {
	TransactionID string    `json:"transaction_id" binding:"required"`
	StoreID       string    `json:"store_id" binding:"required"`
	Timestamp     time.Time `json:"timestamp" binding:"required"`
	Amount        string    `json:"amount" binding:"required"` // сумма в базовой валюте (например, wei)
	Currency      string    `json:"currency"`                  // поле не является обязательным, его устанавливает контроллер
}
