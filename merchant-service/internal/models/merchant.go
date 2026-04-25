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
