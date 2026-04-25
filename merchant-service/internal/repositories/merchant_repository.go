package repositories

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"merchant-service/internal/models"
)

type MerchantRepository interface {
	Create(merchant models.Merchant) (models.Merchant, error)
	GetByEmail(email string) (models.Merchant, error)
	GetByID(id string) (models.Merchant, error)
	SubmitVerification(id string, req models.VerificationRequest) (models.Merchant, error)
	CreateInvoice(merchantID string, req models.CreateInvoiceRequest) (models.Invoice, error)
	ListInvoices(merchantID string) ([]models.Invoice, error)
	CreateTerminal(merchantID string, req models.CreateTerminalRequest) (models.Terminal, error)
	ListTerminals(merchantID string) ([]models.Terminal, error)
}

type memoryMerchantRepo struct {
	mu        sync.RWMutex
	merchants map[string]models.Merchant
	byEmail   map[string]string
	invoices  map[string][]models.Invoice
	terminals map[string][]models.Terminal
	nextID    int64
}

func NewMerchantRepository() MerchantRepository {
	return &memoryMerchantRepo{
		merchants: make(map[string]models.Merchant),
		byEmail:   make(map[string]string),
		invoices:  make(map[string][]models.Invoice),
		terminals: make(map[string][]models.Terminal),
		nextID:    1,
	}
}

func (r *memoryMerchantRepo) Create(merchant models.Merchant) (models.Merchant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byEmail[merchant.OwnerEmail]; exists {
		return models.Merchant{}, errors.New("merchant already exists")
	}
	merchant.ID = strconv.FormatInt(r.nextID, 10)
	merchant.CreatedAt = time.Now().UTC()
	merchant.VerificationStatus = "not_submitted"
	r.nextID++
	r.merchants[merchant.ID] = merchant
	r.byEmail[merchant.OwnerEmail] = merchant.ID
	return merchant, nil
}

func (r *memoryMerchantRepo) GetByEmail(email string) (models.Merchant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exists := r.byEmail[email]
	if !exists {
		return models.Merchant{}, errors.New("merchant not found")
	}
	return r.merchants[id], nil
}

func (r *memoryMerchantRepo) GetByID(id string) (models.Merchant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	merchant, exists := r.merchants[id]
	if !exists {
		return models.Merchant{}, errors.New("merchant not found")
	}
	return merchant, nil
}

func (r *memoryMerchantRepo) SubmitVerification(id string, req models.VerificationRequest) (models.Merchant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	merchant, exists := r.merchants[id]
	if !exists {
		return models.Merchant{}, errors.New("merchant not found")
	}
	merchant.LegalName = req.LegalName
	merchant.VerificationStatus = "pending"
	r.merchants[id] = merchant
	return merchant, nil
}

func (r *memoryMerchantRepo) CreateInvoice(merchantID string, req models.CreateInvoiceRequest) (models.Invoice, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	merchant, exists := r.merchants[merchantID]
	if !exists {
		return models.Invoice{}, errors.New("merchant not found")
	}
	if merchant.VerificationStatus != "approved" {
		return models.Invoice{}, errors.New("merchant verification is required")
	}
	invoice := models.Invoice{
		ID:              strconv.Itoa(len(r.invoices[merchantID]) + 1),
		MerchantID:      merchantID,
		ExternalOrderID: req.ExternalOrderID,
		AmountUSDT:      req.AmountUSDT,
		Status:          "issued",
		ExpiresAt:       time.Now().UTC().Add(30 * time.Minute),
	}
	r.invoices[merchantID] = append(r.invoices[merchantID], invoice)
	return invoice, nil
}

func (r *memoryMerchantRepo) ListInvoices(merchantID string) ([]models.Invoice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.invoices[merchantID], nil
}

func (r *memoryMerchantRepo) CreateTerminal(merchantID string, req models.CreateTerminalRequest) (models.Terminal, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	merchant, exists := r.merchants[merchantID]
	if !exists {
		return models.Terminal{}, errors.New("merchant not found")
	}
	if merchant.VerificationStatus != "approved" {
		return models.Terminal{}, errors.New("merchant verification is required")
	}
	terminal := models.Terminal{
		ID:           strconv.Itoa(len(r.terminals[merchantID]) + 1),
		MerchantID:   merchantID,
		SerialNumber: req.SerialNumber,
		Status:       "active",
		CreatedAt:    time.Now().UTC(),
	}
	r.terminals[merchantID] = append(r.terminals[merchantID], terminal)
	return terminal, nil
}

func (r *memoryMerchantRepo) ListTerminals(merchantID string) ([]models.Terminal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.terminals[merchantID], nil
}
