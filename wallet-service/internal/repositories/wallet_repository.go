package repositories

import (
	"errors"
	"wallet-service/internal/models"
)

// WalletRepository описывает методы для работы с данными кошельков
type WalletRepository interface {
	GetAllWallets() ([]models.Wallet, error)
}

// walletRepo – пример реализации с использованием статического списка (в продакшене использовать базу данных)
type walletRepo struct {
	wallets []models.Wallet
}

func NewWalletRepository() WalletRepository {
	return &walletRepo{
		wallets: []models.Wallet{
			{
				ID:         "wallet1",
				UserID:     "user1",
				Name:       "Main Wallet",
				Blockchain: "ethereum",
				Addresses: map[string]string{
					"ethereum": "0x1234567890abcdef",
				},
			},
		},
	}
}

func (r *walletRepo) GetAllWallets() ([]models.Wallet, error) {
	if len(r.wallets) == 0 {
		return nil, errors.New("no wallets found")
	}
	return r.wallets, nil
}
