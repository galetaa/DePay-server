package repositories

import (
	"database/sql"
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
	_, err := r.db.Exec(`
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
			created_at
		)
		VALUES (
			$1,
			NULLIF($2, '')::bigint,
			NULLIF($3, '')::bigint,
			NULLIF($4, '')::bigint,
			$5::bigint,
			1,
			1,
			16,
			$6::numeric,
			$7::numeric,
			$8,
			$9
		)
	`, tx.TransactionID, tx.InvoiceID, tx.NFCSessionID, tx.UserID, tx.StoreID, tx.Amount, amountUSDT, status, tx.Timestamp)
	return err
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
			COALESCE(failure_reason, '')
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
