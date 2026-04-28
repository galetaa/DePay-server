package models

import "time"

type Merchant struct {
	ID                 string    `json:"id"`
	OwnerEmail         string    `json:"owner_email"`
	StoreName          string    `json:"store_name"`
	LegalName          string    `json:"legal_name"`
	VerificationStatus string    `json:"verification_status"`
	CreatedAt          time.Time `json:"created_at"`
	PasswordHash       string    `json:"-"`
}

type RegisterRequest struct {
	OwnerEmail string `json:"owner_email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	StoreName  string `json:"store_name" binding:"required"`
	LegalName  string `json:"legal_name"`
}

type LoginRequest struct {
	OwnerEmail string `json:"owner_email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
}

type VerificationRequest struct {
	LegalName    string `json:"legal_name" binding:"required"`
	ContactEmail string `json:"contact_email" binding:"required,email"`
	Address      string `json:"address" binding:"required"`
	DocumentURL  string `json:"document_url" binding:"required"`
}

type Invoice struct {
	ID              string    `json:"id"`
	MerchantID      string    `json:"merchant_id"`
	ExternalOrderID string    `json:"external_order_id"`
	AmountUSDT      string    `json:"amount_usdt"`
	Status          string    `json:"status"`
	ExpiresAt       time.Time `json:"expires_at"`
}

type CreateInvoiceRequest struct {
	ExternalOrderID string `json:"external_order_id" binding:"required"`
	AmountUSDT      string `json:"amount_usdt" binding:"required"`
}

type Terminal struct {
	ID           string    `json:"id"`
	MerchantID   string    `json:"merchant_id"`
	SerialNumber string    `json:"serial_number"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateTerminalRequest struct {
	SerialNumber string `json:"serial_number" binding:"required"`
}

type Webhook struct {
	ID            string     `json:"id"`
	MerchantID    string     `json:"merchant_id"`
	URL           string     `json:"url"`
	EventTypes    []string   `json:"event_types"`
	IsActive      bool       `json:"is_active"`
	FailureCount  int        `json:"failure_count"`
	LastSuccessAt *time.Time `json:"last_success_at,omitempty"`
	LastFailureAt *time.Time `json:"last_failure_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type WebhookDelivery struct {
	ID             string     `json:"id"`
	WebhookID      string     `json:"webhook_id"`
	MerchantID     string     `json:"merchant_id"`
	EventType      string     `json:"event_type"`
	Status         string     `json:"status"`
	Attempts       int        `json:"attempts"`
	ResponseStatus int        `json:"response_status,omitempty"`
	ResponseBody   string     `json:"response_body,omitempty"`
	ErrorMessage   string     `json:"error_message,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
}

type APIKey struct {
	ID         string     `json:"id"`
	MerchantID string     `json:"merchant_id"`
	Name       string     `json:"name"`
	KeyPrefix  string     `json:"key_prefix"`
	Secret     string     `json:"secret,omitempty"`
	Scopes     []string   `json:"scopes"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

type CreateAPIKeyRequest struct {
	Name   string   `json:"name" binding:"required"`
	Scopes []string `json:"scopes" binding:"required"`
}

type CreateWebhookRequest struct {
	URL        string   `json:"url" binding:"required,url"`
	EventTypes []string `json:"event_types"`
	Secret     string   `json:"secret" binding:"required"`
}
