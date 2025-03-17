package models

// Wallet описывает структуру кошелька
type Wallet struct {
	ID         string            `json:"id"`
	UserID     string            `json:"user_id"`
	Name       string            `json:"name"`
	Blockchain string            `json:"blockchain"`
	Addresses  map[string]string `json:"addresses"`
}
