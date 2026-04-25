package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sync"
	"transaction-core-service/internal/models"
)

// TransactionRepository описывает методы для работы с транзакциями.
// Здесь используется синхронный in-memory репозиторий для демонстрации.
type TransactionRepository interface {
	Save(tx models.Transaction) error
	Get(transactionID string) (*models.Transaction, error)
	UpdateStatus(transactionID string, status string, failureReason string) error
	MarkBroadcasted(transactionID string, blockchainTxHash string) error
}

type transactionRepo struct {
	data map[string]models.Transaction
	mu   sync.RWMutex
}

func NewTransactionRepository() TransactionRepository {
	return &transactionRepo{
		data: make(map[string]models.Transaction),
	}
}

func (r *transactionRepo) Save(tx models.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.data[tx.TransactionID]; exists {
		return errors.New("transaction already exists")
	}
	if tx.Status == "" {
		tx.Status = "created"
	}
	r.data[tx.TransactionID] = tx
	return nil
}

func (r *transactionRepo) Get(transactionID string) (*models.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tx, exists := r.data[transactionID]
	if !exists {
		return nil, errors.New("transaction not found")
	}
	return &tx, nil
}

func (r *transactionRepo) UpdateStatus(transactionID string, status string, failureReason string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, exists := r.data[transactionID]
	if !exists {
		return errors.New("transaction not found")
	}
	tx.Status = status
	tx.FailureReason = failureReason
	r.data[transactionID] = tx
	return nil
}

func (r *transactionRepo) MarkBroadcasted(transactionID string, blockchainTxHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, exists := r.data[transactionID]
	if !exists {
		return errors.New("transaction not found")
	}
	tx.Status = "broadcasted"
	tx.BlockchainTxHash = blockchainTxHash
	r.data[transactionID] = tx
	return nil
}

type postgresTransactionRepo struct {
	db *sql.DB
}

func NewPostgresTransactionRepository(db *sql.DB) TransactionRepository {
	return &postgresTransactionRepo{db: db}
}

func (r *postgresTransactionRepo) Save(tx models.Transaction) error {
	status := tx.Status
	if status == "" {
		status = "created"
	}
	amountUSDT := tx.AmountUSDT
	if amountUSDT == "" {
		amountUSDT = tx.Amount
	}

	refs, err := r.resolveReferences(context.Background(), tx, amountUSDT)
	if err != nil {
		return err
	}

	var signedPayload any
	if tx.SignedTransaction != "" {
		signedPayload = map[string]string{"signed_transaction": tx.SignedTransaction}
	}
	signedPayloadBytes, err := json.Marshal(signedPayload)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		INSERT INTO payment_transactions(
			external_transaction_id,
			invoice_id,
			nfc_session_id,
			user_id,
			store_id,
			asset_id,
			user_wallet_id,
			store_wallet_id,
			amount,
			amount_in_usdt,
			status,
			signed_payload,
			created_at
		)
		VALUES (
			$1,
			NULLIF($2, '')::bigint,
			NULLIF($3, '')::bigint,
			$4::bigint,
			$5::bigint,
			$6::bigint,
			$7::bigint,
			$8::bigint,
			$9::numeric,
			$10::numeric,
			$11,
			NULLIF($12, 'null')::jsonb,
			$13
		)
	`, tx.TransactionID, refs.InvoiceID, tx.NFCSessionID, refs.UserID, refs.StoreID, refs.AssetID, refs.UserWalletID, refs.StoreWalletID, tx.Amount, refs.AmountUSDT, status, string(signedPayloadBytes), tx.Timestamp)
	return err
}

type transactionReferences struct {
	InvoiceID     string
	UserID        string
	StoreID       string
	AmountUSDT    string
	AssetID       int64
	UserWalletID  int64
	StoreWalletID int64
}

func (r *postgresTransactionRepo) resolveReferences(ctx context.Context, tx models.Transaction, amountUSDT string) (transactionReferences, error) {
	refs := transactionReferences{
		InvoiceID:  tx.InvoiceID,
		UserID:     tx.UserID,
		StoreID:    tx.StoreID,
		AmountUSDT: amountUSDT,
	}

	if refs.InvoiceID != "" {
		var invoiceUserID sql.NullString
		var invoiceStoreID string
		var invoiceAmountUSDT string
		err := r.db.QueryRowContext(ctx, `
			SELECT COALESCE(user_id::text, ''), store_id::text, amount_usdt::text
			FROM payment_invoices
			WHERE invoice_id = $1::bigint
		`, refs.InvoiceID).Scan(&invoiceUserID, &invoiceStoreID, &invoiceAmountUSDT)
		if err != nil {
			return transactionReferences{}, err
		}
		if refs.UserID == "" && invoiceUserID.Valid {
			refs.UserID = invoiceUserID.String
		}
		if refs.StoreID == "" {
			refs.StoreID = invoiceStoreID
		}
		if refs.AmountUSDT == "" {
			refs.AmountUSDT = invoiceAmountUSDT
		}
	}

	if refs.UserID == "" {
		return transactionReferences{}, errors.New("user_id is required")
	}
	if refs.StoreID == "" {
		return transactionReferences{}, errors.New("store_id is required")
	}

	currency := tx.Currency
	if currency == "" {
		currency = "ETH"
	}

	err := r.db.QueryRowContext(ctx, `
		SELECT uw.wallet_id, sw.wallet_id, a.asset_id
		FROM wallets uw
		JOIN wallets sw
		  ON sw.store_id = $2::bigint
		 AND sw.is_store_wallet = true
		 AND sw.chain_id = uw.chain_id
		JOIN assets a
		  ON a.chain_id = uw.chain_id
		 AND upper(a.symbol) = upper($3)
		WHERE uw.user_id = $1::bigint
		  AND uw.is_store_wallet = false
		ORDER BY uw.wallet_id, sw.wallet_id, a.asset_id
		LIMIT 1
	`, refs.UserID, refs.StoreID, currency).Scan(&refs.UserWalletID, &refs.StoreWalletID, &refs.AssetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return transactionReferences{}, errors.New("matching user/store wallets and asset were not found")
		}
		return transactionReferences{}, err
	}

	return refs, nil
}

func (r *postgresTransactionRepo) Get(transactionID string) (*models.Transaction, error) {
	tx := &models.Transaction{}
	err := r.db.QueryRow(`
		SELECT
			COALESCE(external_transaction_id, transaction_id::text),
			store_id::text,
			user_id::text,
			COALESCE(invoice_id::text, ''),
			COALESCE(nfc_session_id::text, ''),
			created_at,
			amount::text,
			amount_in_usdt::text,
			status::text,
			COALESCE(failure_reason, ''),
			COALESCE(signed_payload->>'signed_transaction', ''),
			COALESCE(blockchain_tx_hash, '')
		FROM payment_transactions
		WHERE transaction_id::text = $1
		   OR external_transaction_id = $1
	`, transactionID).Scan(
		&tx.TransactionID,
		&tx.StoreID,
		&tx.UserID,
		&tx.InvoiceID,
		&tx.NFCSessionID,
		&tx.Timestamp,
		&tx.Amount,
		&tx.AmountUSDT,
		&tx.Status,
		&tx.FailureReason,
		&tx.SignedTransaction,
		&tx.BlockchainTxHash,
	)
	if err != nil {
		return nil, err
	}
	tx.Currency = "ETH"
	return tx, nil
}

func (r *postgresTransactionRepo) UpdateStatus(transactionID string, status string, failureReason string) error {
	result, err := r.db.Exec(`
		UPDATE payment_transactions
		SET status = $2,
		    failure_reason = NULLIF($3, '')
		WHERE transaction_id::text = $1
		   OR external_transaction_id = $1
	`, transactionID, status, failureReason)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("transaction not found")
	}
	return nil
}

func (r *postgresTransactionRepo) MarkBroadcasted(transactionID string, blockchainTxHash string) error {
	result, err := r.db.Exec(`
		UPDATE payment_transactions
		SET status = 'broadcasted',
		    blockchain_tx_hash = $2
		WHERE transaction_id::text = $1
		   OR external_transaction_id = $1
	`, transactionID, blockchainTxHash)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("transaction not found")
	}
	return nil
}
