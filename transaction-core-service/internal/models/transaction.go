package models

import "time"

// Transaction представляет транзакцию, инициированную терминалом.
// Все расчёты ведутся в единой базовой валюте (например, ETH).
type Transaction struct {
	TransactionID     string    `json:"transaction_id" binding:"required"`
	StoreID           string    `json:"store_id" binding:"required"`
	UserID            string    `json:"user_id"`
	InvoiceID         string    `json:"invoice_id,omitempty"`
	NFCSessionID      string    `json:"nfc_session_id,omitempty"`
	Timestamp         time.Time `json:"timestamp" binding:"required"`
	Amount            string    `json:"amount" binding:"required"` // сумма в базовой валюте (например, wei)
	AmountUSDT        string    `json:"amount_usdt,omitempty"`
	Currency          string    `json:"currency"` // поле не является обязательным, его устанавливает контроллер
	Status            string    `json:"status,omitempty"`
	FailureReason     string    `json:"failure_reason,omitempty"`
	SignedTransaction string    `json:"signed_transaction,omitempty"`
	BlockchainTxHash  string    `json:"blockchain_tx_hash,omitempty"`
}

type TransactionStatusResponse struct {
	TransactionID    string `json:"transaction_id"`
	Status           string `json:"status"`
	FailureReason    string `json:"failure_reason,omitempty"`
	BlockchainTxHash string `json:"blockchain_tx_hash,omitempty"`
}
