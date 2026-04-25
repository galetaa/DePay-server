package models

// Wallet описывает структуру кошелька
type Wallet struct {
	ID         string            `json:"id"`
	UserID     string            `json:"user_id"`
	StoreID    string            `json:"store_id,omitempty"`
	Name       string            `json:"name"`
	Blockchain string            `json:"blockchain"`
	Address    string            `json:"address,omitempty"`
	Addresses  map[string]string `json:"addresses"`
}

type CreateWalletRequest struct {
	UserID     string `json:"user_id"`
	StoreID    string `json:"store_id"`
	ChainID    string `json:"chain_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Address    string `json:"address" binding:"required"`
	PublicKey  string `json:"public_key"`
	IsStore    bool   `json:"is_store_wallet"`
	Blockchain string `json:"blockchain"`
}

type WalletBalance struct {
	WalletID    string `json:"wallet_id"`
	AssetSymbol string `json:"asset_symbol"`
	Balance     string `json:"balance"`
	UpdatedAt   string `json:"updated_at"`
}
