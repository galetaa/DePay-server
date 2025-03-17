package models

// KYCResponse представляет ответ по KYC проверке
type KYCResponse struct {
	UserID    string `json:"user_id"`
	KYCStatus string `json:"kyc_status"` // Например, "pending", "verified", "rejected"
	Message   string `json:"message"`
}
