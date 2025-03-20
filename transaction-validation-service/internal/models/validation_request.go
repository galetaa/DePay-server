package models

// ValidationRequest содержит данные подписанной транзакции, необходимые для проверки
type ValidationRequest struct {
	TransactionID    string `json:"transaction_id" binding:"required"`
	SignedData       string `json:"signed_data" binding:"required"`
	SenderAddress    string `json:"sender_address" binding:"required"`
	RecipientAddress string `json:"recipient_address" binding:"required"`
	Amount           string `json:"amount" binding:"required"`
	Currency         string `json:"currency" binding:"required"`
}
