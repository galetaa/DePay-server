package repositories

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/lib/pq"

	"merchant-service/internal/models"
	"shared/auth"
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
	CreateWebhook(merchantID string, req models.CreateWebhookRequest) (models.Webhook, error)
	ListWebhooks(merchantID string) ([]models.Webhook, error)
	DeleteWebhook(merchantID string, webhookID string) error
}

type memoryMerchantRepo struct {
	mu            sync.RWMutex
	merchants     map[string]models.Merchant
	byEmail       map[string]string
	invoices      map[string][]models.Invoice
	terminals     map[string][]models.Terminal
	webhooks      map[string][]models.Webhook
	nextID        int64
	nextWebhookID int64
}

func NewMerchantRepository() MerchantRepository {
	return &memoryMerchantRepo{
		merchants:     make(map[string]models.Merchant),
		byEmail:       make(map[string]string),
		invoices:      make(map[string][]models.Invoice),
		terminals:     make(map[string][]models.Terminal),
		webhooks:      make(map[string][]models.Webhook),
		nextID:        1,
		nextWebhookID: 1,
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

func (r *memoryMerchantRepo) CreateWebhook(merchantID string, req models.CreateWebhookRequest) (models.Webhook, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	merchant, exists := r.merchants[merchantID]
	if !exists {
		return models.Webhook{}, errors.New("merchant not found")
	}
	if merchant.VerificationStatus != "approved" {
		return models.Webhook{}, errors.New("merchant verification is required")
	}
	webhook := models.Webhook{
		ID:           strconv.FormatInt(r.nextWebhookID, 10),
		MerchantID:   merchantID,
		URL:          req.URL,
		EventTypes:   defaultWebhookEvents(req.EventTypes),
		IsActive:     true,
		FailureCount: 0,
		CreatedAt:    time.Now().UTC(),
	}
	r.nextWebhookID++
	r.webhooks[merchantID] = append(r.webhooks[merchantID], webhook)
	return webhook, nil
}

func (r *memoryMerchantRepo) ListWebhooks(merchantID string) ([]models.Webhook, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.webhooks[merchantID], nil
}

func (r *memoryMerchantRepo) DeleteWebhook(merchantID string, webhookID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	webhooks := r.webhooks[merchantID]
	for i, webhook := range webhooks {
		if webhook.ID == webhookID {
			r.webhooks[merchantID] = append(webhooks[:i], webhooks[i+1:]...)
			return nil
		}
	}
	return errors.New("webhook not found")
}

type postgresMerchantRepo struct {
	db *sql.DB
}

func NewPostgresMerchantRepository(db *sql.DB) MerchantRepository {
	return &postgresMerchantRepo{db: db}
}

func (r *postgresMerchantRepo) Create(merchant models.Merchant) (models.Merchant, error) {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Merchant{}, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `INSERT INTO roles(role_name) VALUES ('merchant') ON CONFLICT DO NOTHING`); err != nil {
		return models.Merchant{}, err
	}

	var userID int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO users(username, email, password_hash, kyc_status)
		VALUES ($1, $1, $2, 'not_submitted')
		RETURNING user_id
	`, merchant.OwnerEmail, merchant.PasswordHash).Scan(&userID)
	if err != nil {
		return models.Merchant{}, err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO user_roles(user_id, role_id)
		SELECT $1, role_id FROM roles WHERE role_name = 'merchant'
		ON CONFLICT DO NOTHING
	`, userID); err != nil {
		return models.Merchant{}, err
	}

	var storeID int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO stores(owner_id, store_name, legal_name, contact_email, verification_status)
		VALUES ($1, $2, NULLIF($3, ''), $4, 'not_submitted')
		RETURNING store_id, created_at
	`, userID, merchant.StoreName, merchant.LegalName, merchant.OwnerEmail).Scan(&storeID, &merchant.CreatedAt)
	if err != nil {
		return models.Merchant{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Merchant{}, err
	}

	merchant.ID = strconv.FormatInt(storeID, 10)
	merchant.VerificationStatus = "not_submitted"
	return merchant, nil
}

func (r *postgresMerchantRepo) GetByEmail(email string) (models.Merchant, error) {
	return r.getOne(context.Background(), "u.email = $1", email)
}

func (r *postgresMerchantRepo) GetByID(id string) (models.Merchant, error) {
	return r.getOne(context.Background(), "s.store_id = $1", id)
}

func (r *postgresMerchantRepo) getOne(ctx context.Context, where string, arg any) (models.Merchant, error) {
	query := `
		SELECT
			s.store_id::text,
			u.email,
			s.store_name,
			COALESCE(s.legal_name, ''),
			s.verification_status::text,
			s.created_at,
			u.password_hash
		FROM stores s
		JOIN users u ON u.user_id = s.owner_id
		WHERE ` + where + `
		ORDER BY s.store_id
		LIMIT 1
	`

	var merchant models.Merchant
	if err := r.db.QueryRowContext(ctx, query, arg).Scan(
		&merchant.ID,
		&merchant.OwnerEmail,
		&merchant.StoreName,
		&merchant.LegalName,
		&merchant.VerificationStatus,
		&merchant.CreatedAt,
		&merchant.PasswordHash,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Merchant{}, errors.New("merchant not found")
		}
		return models.Merchant{}, err
	}
	return merchant, nil
}

func (r *postgresMerchantRepo) SubmitVerification(id string, req models.VerificationRequest) (models.Merchant, error) {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Merchant{}, err
	}
	defer tx.Rollback()

	var storeID int64
	err = tx.QueryRowContext(ctx, `
		UPDATE stores
		SET legal_name = $2,
		    contact_email = $3,
		    store_address = $4,
		    verification_status = 'pending',
		    updated_at = now()
		WHERE store_id = $1
		RETURNING store_id
	`, id, req.LegalName, req.ContactEmail, req.Address).Scan(&storeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Merchant{}, errors.New("merchant not found")
		}
		return models.Merchant{}, err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO merchant_verification_applications(store_id, status)
		VALUES ($1, 'pending')
	`, storeID); err != nil {
		return models.Merchant{}, err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO audit_logs(action, entity_type, entity_id, new_values)
		VALUES ('merchant.verification_submitted', 'store', $1, jsonb_build_object('document_url', $2, 'legal_name', $3))
	`, id, req.DocumentURL, req.LegalName); err != nil {
		return models.Merchant{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Merchant{}, err
	}
	return r.GetByID(id)
}

func (r *postgresMerchantRepo) CreateInvoice(merchantID string, req models.CreateInvoiceRequest) (models.Invoice, error) {
	if err := r.ensureApprovedMerchant(context.Background(), merchantID); err != nil {
		return models.Invoice{}, err
	}

	var invoice models.Invoice
	invoice.MerchantID = merchantID
	err := r.db.QueryRowContext(context.Background(), `
		INSERT INTO payment_invoices(store_id, external_order_id, amount_usdt, status, expires_at)
		VALUES ($1, $2, $3::numeric, 'issued', now() + interval '30 minutes')
		RETURNING invoice_id::text, external_order_id, amount_usdt::text, status::text, expires_at
	`, merchantID, req.ExternalOrderID, req.AmountUSDT).Scan(
		&invoice.ID,
		&invoice.ExternalOrderID,
		&invoice.AmountUSDT,
		&invoice.Status,
		&invoice.ExpiresAt,
	)
	if err != nil {
		return models.Invoice{}, err
	}
	return invoice, nil
}

func (r *postgresMerchantRepo) ListInvoices(merchantID string) ([]models.Invoice, error) {
	rows, err := r.db.QueryContext(context.Background(), `
		SELECT invoice_id::text, store_id::text, COALESCE(external_order_id, ''), amount_usdt::text, status::text, expires_at
		FROM payment_invoices
		WHERE store_id = $1
		ORDER BY created_at DESC, invoice_id DESC
	`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := make([]models.Invoice, 0)
	for rows.Next() {
		var invoice models.Invoice
		if err := rows.Scan(&invoice.ID, &invoice.MerchantID, &invoice.ExternalOrderID, &invoice.AmountUSDT, &invoice.Status, &invoice.ExpiresAt); err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}
	return invoices, rows.Err()
}

func (r *postgresMerchantRepo) CreateTerminal(merchantID string, req models.CreateTerminalRequest) (models.Terminal, error) {
	if err := r.ensureApprovedMerchant(context.Background(), merchantID); err != nil {
		return models.Terminal{}, err
	}

	secretHash := auth.HashToken(req.SerialNumber + ":" + merchantID)
	var terminal models.Terminal
	terminal.MerchantID = merchantID
	err := r.db.QueryRowContext(context.Background(), `
		INSERT INTO terminals(store_id, serial_number, secret_hash, status)
		VALUES ($1, $2, $3, 'active')
		RETURNING terminal_id::text, serial_number, status::text, created_at
	`, merchantID, req.SerialNumber, secretHash).Scan(
		&terminal.ID,
		&terminal.SerialNumber,
		&terminal.Status,
		&terminal.CreatedAt,
	)
	if err != nil {
		return models.Terminal{}, err
	}
	return terminal, nil
}

func (r *postgresMerchantRepo) ListTerminals(merchantID string) ([]models.Terminal, error) {
	rows, err := r.db.QueryContext(context.Background(), `
		SELECT terminal_id::text, store_id::text, serial_number, status::text, created_at
		FROM terminals
		WHERE store_id = $1
		ORDER BY created_at DESC, terminal_id DESC
	`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	terminals := make([]models.Terminal, 0)
	for rows.Next() {
		var terminal models.Terminal
		if err := rows.Scan(&terminal.ID, &terminal.MerchantID, &terminal.SerialNumber, &terminal.Status, &terminal.CreatedAt); err != nil {
			return nil, err
		}
		terminals = append(terminals, terminal)
	}
	return terminals, rows.Err()
}

func (r *postgresMerchantRepo) CreateWebhook(merchantID string, req models.CreateWebhookRequest) (models.Webhook, error) {
	if err := r.ensureApprovedMerchant(context.Background(), merchantID); err != nil {
		return models.Webhook{}, err
	}

	secretHash := req.Secret
	events := defaultWebhookEvents(req.EventTypes)
	var webhook models.Webhook
	var lastSuccessAt sql.NullTime
	var lastFailureAt sql.NullTime
	webhook.MerchantID = merchantID
	err := r.db.QueryRowContext(context.Background(), `
		INSERT INTO merchant_webhooks(store_id, url, event_types, secret_hash, is_active)
		VALUES ($1, $2, $3, $4, true)
		RETURNING webhook_id::text, url, event_types, is_active, failure_count, last_success_at, last_failure_at, created_at
	`, merchantID, req.URL, pq.Array(events), secretHash).Scan(
		&webhook.ID,
		&webhook.URL,
		(*pq.StringArray)(&webhook.EventTypes),
		&webhook.IsActive,
		&webhook.FailureCount,
		&lastSuccessAt,
		&lastFailureAt,
		&webhook.CreatedAt,
	)
	if err != nil {
		return models.Webhook{}, err
	}
	if lastSuccessAt.Valid {
		webhook.LastSuccessAt = &lastSuccessAt.Time
	}
	if lastFailureAt.Valid {
		webhook.LastFailureAt = &lastFailureAt.Time
	}
	return webhook, nil
}

func (r *postgresMerchantRepo) ListWebhooks(merchantID string) ([]models.Webhook, error) {
	rows, err := r.db.QueryContext(context.Background(), `
		SELECT webhook_id::text, store_id::text, url, event_types, is_active, failure_count, last_success_at, last_failure_at, created_at
		FROM merchant_webhooks
		WHERE store_id = $1
		ORDER BY created_at DESC, webhook_id DESC
	`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	webhooks := make([]models.Webhook, 0)
	for rows.Next() {
		var webhook models.Webhook
		var lastSuccessAt sql.NullTime
		var lastFailureAt sql.NullTime
		if err := rows.Scan(
			&webhook.ID,
			&webhook.MerchantID,
			&webhook.URL,
			(*pq.StringArray)(&webhook.EventTypes),
			&webhook.IsActive,
			&webhook.FailureCount,
			&lastSuccessAt,
			&lastFailureAt,
			&webhook.CreatedAt,
		); err != nil {
			return nil, err
		}
		if lastSuccessAt.Valid {
			webhook.LastSuccessAt = &lastSuccessAt.Time
		}
		if lastFailureAt.Valid {
			webhook.LastFailureAt = &lastFailureAt.Time
		}
		webhooks = append(webhooks, webhook)
	}
	return webhooks, rows.Err()
}

func (r *postgresMerchantRepo) DeleteWebhook(merchantID string, webhookID string) error {
	result, err := r.db.ExecContext(context.Background(), `
		DELETE FROM merchant_webhooks
		WHERE store_id = $1
		  AND webhook_id = $2
	`, merchantID, webhookID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("webhook not found")
	}
	return nil
}

func (r *postgresMerchantRepo) ensureApprovedMerchant(ctx context.Context, merchantID string) error {
	var status string
	err := r.db.QueryRowContext(ctx, `
		SELECT verification_status::text
		FROM stores
		WHERE store_id = $1
	`, merchantID).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("merchant not found")
		}
		return err
	}
	if status != "approved" {
		return errors.New("merchant verification is required")
	}
	return nil
}

func defaultWebhookEvents(events []string) []string {
	if len(events) > 0 {
		return events
	}
	return []string{
		"transaction.created",
		"transaction.submitted",
		"transaction.validated",
		"transaction.broadcasted",
		"transaction.confirmed",
		"transaction.failed",
	}
}
