package repositories

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"sync"

	"wallet-service/internal/models"
)

// WalletRepository описывает методы для работы с данными кошельков
type WalletRepository interface {
	GetAllWallets(ctx context.Context) ([]models.Wallet, error)
	Create(ctx context.Context, req models.CreateWalletRequest) (models.Wallet, error)
	GetByID(ctx context.Context, id string) (models.Wallet, error)
	Delete(ctx context.Context, id string) error
	GetBalances(ctx context.Context, walletID string) ([]models.WalletBalance, error)
}

type walletRepo struct {
	wallets []models.Wallet
	mu      sync.RWMutex
	nextID  int64
}

func NewWalletRepository() WalletRepository {
	return &walletRepo{
		nextID: 2,
		wallets: []models.Wallet{
			{
				ID:         "wallet1",
				UserID:     "user1",
				Name:       "Main Wallet",
				Blockchain: "ethereum",
				Address:    "0x1234567890abcdef",
				Addresses: map[string]string{
					"ethereum": "0x1234567890abcdef",
				},
			},
		},
	}
}

func (r *walletRepo) GetAllWallets(ctx context.Context) ([]models.Wallet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.wallets) == 0 {
		return nil, errors.New("no wallets found")
	}
	return r.wallets, nil
}

func (r *walletRepo) Create(ctx context.Context, req models.CreateWalletRequest) (models.Wallet, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if req.UserID == "" && req.StoreID == "" {
		return models.Wallet{}, errors.New("wallet owner is required")
	}
	blockchain := req.Blockchain
	if blockchain == "" {
		blockchain = "ethereum"
	}
	wallet := models.Wallet{
		ID:         "wallet" + strconv.FormatInt(r.nextID, 10),
		UserID:     req.UserID,
		StoreID:    req.StoreID,
		Name:       req.Name,
		Blockchain: blockchain,
		Address:    req.Address,
		Addresses: map[string]string{
			blockchain: req.Address,
		},
	}
	r.nextID++
	r.wallets = append(r.wallets, wallet)
	return wallet, nil
}

func (r *walletRepo) GetByID(ctx context.Context, id string) (models.Wallet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, wallet := range r.wallets {
		if wallet.ID == id {
			return wallet, nil
		}
	}
	return models.Wallet{}, errors.New("wallet not found")
}

func (r *walletRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, wallet := range r.wallets {
		if wallet.ID == id {
			r.wallets = append(r.wallets[:i], r.wallets[i+1:]...)
			return nil
		}
	}
	return errors.New("wallet not found")
}

func (r *walletRepo) GetBalances(ctx context.Context, walletID string) ([]models.WalletBalance, error) {
	if _, err := r.GetByID(ctx, walletID); err != nil {
		return nil, err
	}
	return []models.WalletBalance{
		{
			WalletID:    walletID,
			AssetSymbol: "ETH",
			Balance:     "1.00000000",
		},
	}, nil
}

type postgresWalletRepo struct {
	db *sql.DB
}

func NewPostgresWalletRepository(db *sql.DB) WalletRepository {
	return &postgresWalletRepo{db: db}
}

func (r *postgresWalletRepo) GetAllWallets(ctx context.Context) ([]models.Wallet, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT w.wallet_id::text, COALESCE(w.user_id::text, ''), COALESCE(w.store_id::text, ''),
		       COALESCE(w.wallet_label, ''), bc.chain_name, w.wallet_address
		FROM wallets w
		JOIN blockchains bc ON bc.chain_id = w.chain_id
		ORDER BY w.wallet_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []models.Wallet
	for rows.Next() {
		var wallet models.Wallet
		if err := rows.Scan(&wallet.ID, &wallet.UserID, &wallet.StoreID, &wallet.Name, &wallet.Blockchain, &wallet.Address); err != nil {
			return nil, err
		}
		wallet.Addresses = map[string]string{wallet.Blockchain: wallet.Address}
		wallets = append(wallets, wallet)
	}
	return wallets, rows.Err()
}

func (r *postgresWalletRepo) Create(ctx context.Context, req models.CreateWalletRequest) (models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO wallets(user_id, store_id, chain_id, wallet_address, public_key, is_store_wallet, wallet_label)
		VALUES (NULLIF($1, '')::bigint, NULLIF($2, '')::bigint, $3::bigint, $4, NULLIF($5, ''), $6, $7)
		RETURNING wallet_id::text
	`, req.UserID, req.StoreID, req.ChainID, req.Address, req.PublicKey, req.IsStore, req.Name).Scan(&wallet.ID)
	if err != nil {
		return models.Wallet{}, err
	}
	return r.GetByID(ctx, wallet.ID)
}

func (r *postgresWalletRepo) GetByID(ctx context.Context, id string) (models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.QueryRowContext(ctx, `
		SELECT w.wallet_id::text, COALESCE(w.user_id::text, ''), COALESCE(w.store_id::text, ''),
		       COALESCE(w.wallet_label, ''), bc.chain_name, w.wallet_address
		FROM wallets w
		JOIN blockchains bc ON bc.chain_id = w.chain_id
		WHERE w.wallet_id = $1
	`, id).Scan(&wallet.ID, &wallet.UserID, &wallet.StoreID, &wallet.Name, &wallet.Blockchain, &wallet.Address)
	if err != nil {
		return models.Wallet{}, err
	}
	wallet.Addresses = map[string]string{wallet.Blockchain: wallet.Address}
	return wallet, nil
}

func (r *postgresWalletRepo) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM wallets WHERE wallet_id = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("wallet not found")
	}
	return nil
}

func (r *postgresWalletRepo) GetBalances(ctx context.Context, walletID string) ([]models.WalletBalance, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT b.wallet_id::text, a.symbol, b.balance::text, b.updated_at::text
		FROM balances b
		JOIN assets a ON a.asset_id = b.asset_id
		WHERE b.wallet_id = $1
		ORDER BY a.symbol
	`, walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []models.WalletBalance
	for rows.Next() {
		var balance models.WalletBalance
		if err := rows.Scan(&balance.WalletID, &balance.AssetSymbol, &balance.Balance, &balance.UpdatedAt); err != nil {
			return nil, err
		}
		balances = append(balances, balance)
	}
	return balances, rows.Err()
}
