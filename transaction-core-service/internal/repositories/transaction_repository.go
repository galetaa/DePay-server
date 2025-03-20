package repositories

import (
	"errors"
	"sync"
	"transaction-core-service/internal/models"
)

// TransactionRepository описывает методы для работы с транзакциями.
// Здесь используется синхронный in-memory репозиторий для демонстрации.
type TransactionRepository interface {
	Save(tx models.Transaction) error
	Get(transactionID string) (*models.Transaction, error)
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
