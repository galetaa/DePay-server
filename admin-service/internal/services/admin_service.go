package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"admin-service/internal/models"
)

var allowedTables = map[string]string{
	"users":                       "users",
	"stores":                      "stores",
	"wallets":                     "wallets",
	"balances":                    "balances",
	"payment_invoices":            "payment_invoices",
	"payment_transactions":        "payment_transactions",
	"rpc_nodes":                   "rpc_nodes",
	"rpc_node_checks":             "rpc_node_checks",
	"risk_alerts":                 "risk_alerts",
	"audit_logs":                  "audit_logs",
	"merchant_webhooks":           "merchant_webhooks",
	"merchant_webhook_deliveries": "merchant_webhook_deliveries",
	"vw_user_wallet_balances":     "vw_user_wallet_balances",
	"vw_store_transactions":       "vw_store_transactions",
	"vw_failed_transactions":      "vw_failed_transactions",
	"vw_rpc_node_status":          "vw_rpc_node_status",
	"vw_compliance_kyc_queue":     "vw_compliance_kyc_queue",
	"vw_webhook_delivery_status":  "vw_webhook_delivery_status",
}

var allowedFunctions = map[string]int{
	"get_user_kyc_wallet_summary":       0,
	"get_user_wallet_balances":          1,
	"get_wallet_asset_distribution":     1,
	"get_transaction_card":              1,
	"get_user_transaction_history":      3,
	"get_store_transaction_history":     3,
	"get_blockchain_asset_activity":     3,
	"get_rpc_nodes_activity":            3,
	"get_store_turnover":                3,
	"get_store_success_rate":            3,
	"get_unverified_active_users":       4,
	"get_failed_transactions_analytics": 2,
}

type AdminService interface {
	ListTables() []string
	GetTableRows(ctx context.Context, tableName string, limit int) (models.TableRowsResponse, error)
	ExecuteFunction(ctx context.Context, functionName string, params []string) (models.TableRowsResponse, error)
	AuditLogs(ctx context.Context, limit int) (models.TableRowsResponse, error)
	RiskAlerts(ctx context.Context, limit int) (models.TableRowsResponse, error)
	CreateDemoInvoice(ctx context.Context) (map[string]any, error)
	SubmitDemoPayment(ctx context.Context, req models.DemoPaymentRequest) (map[string]any, error)
}

type adminService struct {
	db *sql.DB
}

func NewAdminService(db *sql.DB) AdminService {
	return &adminService{db: db}
}

func (s *adminService) ListTables() []string {
	tables := make([]string, 0, len(allowedTables))
	for table := range allowedTables {
		tables = append(tables, table)
	}
	return tables
}

func (s *adminService) GetTableRows(ctx context.Context, tableName string, limit int) (models.TableRowsResponse, error) {
	table, ok := allowedTables[tableName]
	if !ok {
		return models.TableRowsResponse{}, errors.New("table is not allowed")
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	return s.queryRows(ctx, fmt.Sprintf("SELECT * FROM %s LIMIT %d", table, limit))
}

func (s *adminService) ExecuteFunction(ctx context.Context, functionName string, params []string) (models.TableRowsResponse, error) {
	expectedParams, ok := allowedFunctions[functionName]
	if !ok {
		return models.TableRowsResponse{}, errors.New("function is not allowed")
	}
	if len(params) != expectedParams {
		return models.TableRowsResponse{}, fmt.Errorf("function expects %d params", expectedParams)
	}

	placeholders := make([]string, len(params))
	args := make([]any, len(params))
	for i, param := range params {
		placeholders[i] = "$" + strconv.Itoa(i+1)
		args[i] = param
	}

	query := fmt.Sprintf("SELECT * FROM %s(%s)", functionName, strings.Join(placeholders, ", "))
	return s.queryRows(ctx, query, args...)
}

func (s *adminService) AuditLogs(ctx context.Context, limit int) (models.TableRowsResponse, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	return s.queryRows(ctx, fmt.Sprintf("SELECT * FROM audit_logs ORDER BY created_at DESC LIMIT %d", limit))
}

func (s *adminService) RiskAlerts(ctx context.Context, limit int) (models.TableRowsResponse, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	return s.queryRows(ctx, fmt.Sprintf("SELECT * FROM risk_alerts ORDER BY created_at DESC LIMIT %d", limit))
}

func (s *adminService) CreateDemoInvoice(ctx context.Context) (map[string]any, error) {
	var invoiceID string
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO payment_invoices(store_id, user_id, external_order_id, amount_usdt, status, expires_at)
		VALUES (1, 1, 'DEMO-' || replace(gen_random_uuid()::text, '-', ''), 42.50, 'issued', now() + interval '30 minutes')
		RETURNING invoice_id::text
	`).Scan(&invoiceID)
	if err != nil {
		return nil, err
	}
	return map[string]any{"invoice_id": invoiceID, "status": "issued"}, nil
}

func (s *adminService) SubmitDemoPayment(ctx context.Context, req models.DemoPaymentRequest) (map[string]any, error) {
	if req.InvoiceID == "" {
		return nil, errors.New("invoice_id is required")
	}
	if req.UserID == "" {
		req.UserID = "1"
	}
	if req.WalletID == "" {
		req.WalletID = "1"
	}

	var transactionID string
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO payment_transactions(
			invoice_id, user_id, store_id, asset_id, user_wallet_id, store_wallet_id,
			amount, amount_in_usdt, status, created_at
		)
		SELECT $1::bigint, $2::bigint, store_id, 2, $3::bigint, 16, amount_usdt, amount_usdt, 'created', now()
		FROM payment_invoices
		WHERE invoice_id = $1::bigint
		RETURNING transaction_id::text
	`, req.InvoiceID, req.UserID, req.WalletID).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	for _, status := range []string{"submitted", "validated", "broadcasted", "confirmed"} {
		if _, err := s.db.ExecContext(ctx, `UPDATE payment_transactions SET status = $2 WHERE transaction_id = $1`, transactionID, status); err != nil {
			return nil, err
		}
	}
	_, _ = s.db.ExecContext(ctx, `UPDATE payment_invoices SET status = 'paid', paid_at = now() WHERE invoice_id = $1`, req.InvoiceID)
	_ = s.recordDemoWebhookDeliveries(ctx, transactionID)

	return map[string]any{"transaction_id": transactionID, "status": "confirmed"}, nil
}

func (s *adminService) recordDemoWebhookDeliveries(ctx context.Context, transactionID string) error {
	_, err := s.db.ExecContext(ctx, `
		WITH tx AS (
			SELECT transaction_id, external_transaction_id, store_id, user_id, invoice_id, amount, amount_in_usdt, status
			FROM payment_transactions
			WHERE transaction_id = $1::bigint
		),
		inserted AS (
			INSERT INTO merchant_webhook_deliveries(
				webhook_id,
				store_id,
				transaction_id,
				event_type,
				payload,
				status,
				attempts,
				response_status,
				response_body,
				delivered_at
			)
			SELECT
				mw.webhook_id,
				mw.store_id,
				tx.transaction_id,
				'transaction.confirmed',
				jsonb_build_object(
					'type', 'transaction.confirmed',
					'payload', jsonb_build_object(
						'transaction_id', COALESCE(tx.external_transaction_id, tx.transaction_id::text),
						'store_id', tx.store_id::text,
						'user_id', tx.user_id::text,
						'invoice_id', COALESCE(tx.invoice_id::text, ''),
						'amount', tx.amount::text,
						'amount_usdt', tx.amount_in_usdt::text,
						'status', tx.status::text
					)
				),
				'delivered',
				1,
				204,
				'admin demo webhook delivery log mode',
				now()
			FROM merchant_webhooks mw
			JOIN tx ON tx.store_id = mw.store_id
			WHERE mw.is_active = true
			  AND mw.event_types @> ARRAY['transaction.confirmed']::text[]
			RETURNING webhook_id
		)
		UPDATE merchant_webhooks
		SET last_success_at = now(),
		    updated_at = now()
		WHERE webhook_id IN (SELECT webhook_id FROM inserted)
	`, transactionID)
	return err
}

func (s *adminService) queryRows(ctx context.Context, query string, args ...any) (models.TableRowsResponse, error) {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return models.TableRowsResponse{}, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return models.TableRowsResponse{}, err
	}

	result := models.TableRowsResponse{Columns: columns, Rows: []map[string]any{}}
	for rows.Next() {
		values := make([]any, len(columns))
		dest := make([]any, len(columns))
		for i := range values {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return models.TableRowsResponse{}, err
		}

		row := make(map[string]any, len(columns))
		for i, column := range columns {
			switch value := values[i].(type) {
			case []byte:
				row[column] = string(value)
			default:
				row[column] = value
			}
		}
		result.Rows = append(result.Rows, row)
	}
	return result, rows.Err()
}
