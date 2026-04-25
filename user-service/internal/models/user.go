package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PhoneNumber  string    `json:"phone_number,omitempty"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	KYCStatus    string    `json:"kyc_status"`
	Roles        []string  `json:"roles,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	PasswordHash string    `json:"-"`
}
