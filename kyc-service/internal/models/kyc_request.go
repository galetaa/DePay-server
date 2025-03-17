package models

// KYCRequest представляет данные для запроса KYC проверки
type KYCRequest struct {
	UserID       string `json:"user_id" binding:"required"`
	DocumentType string `json:"document_type" binding:"required"`
	DocumentData string `json:"document_data" binding:"required"`
}
