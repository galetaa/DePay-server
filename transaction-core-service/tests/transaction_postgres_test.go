package tests

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"transaction-core-service/internal/models"
	"transaction-core-service/internal/repositories"
	"transaction-core-service/internal/services"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestPostgresLifecycleConfirmUpdatesInvoiceAndWebhook(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is not set")
	}
	t.Setenv("SKIP_RABBITMQ", "true")
	t.Setenv("BLOCKCHAIN_RPC_URL", "")
	t.Setenv("WEBHOOK_DELIVERY_MODE", "log")

	db, err := sql.Open("postgres", databaseURL)
	require.NoError(t, err)
	defer db.Close()
	require.NoError(t, db.Ping())

	ctx := context.Background()
	externalOrderID := fmt.Sprintf("TEST-LIFECYCLE-%d", time.Now().UnixNano())
	transactionID := fmt.Sprintf("test-lifecycle-%d", time.Now().UnixNano())

	var invoiceID string
	err = db.QueryRowContext(ctx, `
		INSERT INTO payment_invoices(store_id, user_id, external_order_id, amount_usdt, status, expires_at)
		VALUES (1, 1, $1, 42.50, 'issued', now() + interval '1 hour')
		RETURNING invoice_id::text
	`, externalOrderID).Scan(&invoiceID)
	require.NoError(t, err)

	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `
			DELETE FROM merchant_webhook_deliveries
			WHERE transaction_id IN (
				SELECT transaction_id
				FROM payment_transactions
				WHERE external_transaction_id = $1
			)
		`, transactionID)
		_, _ = db.ExecContext(ctx, `DELETE FROM payment_transactions WHERE external_transaction_id = $1`, transactionID)
		_, _ = db.ExecContext(ctx, `DELETE FROM payment_invoices WHERE invoice_id = $1::bigint`, invoiceID)
	})

	repo := repositories.NewPostgresTransactionRepository(db)
	svc := services.NewTransactionService(repo, services.WithWebhookDispatcher(services.NewPostgresWebhookDispatcher(db)))

	err = svc.Initiate(models.Transaction{
		TransactionID: transactionID,
		StoreID:       "1",
		UserID:        "1",
		InvoiceID:     invoiceID,
		Timestamp:     time.Now().UTC(),
		Amount:        "0.01",
		Currency:      "ETH",
	})
	require.NoError(t, err)
	require.NoError(t, svc.UpdateStatus(transactionID, "submitted", ""))
	require.NoError(t, svc.UpdateStatus(transactionID, "validated", ""))
	_, err = svc.Broadcast(transactionID)
	require.NoError(t, err)
	require.NoError(t, svc.UpdateStatus(transactionID, "confirmed", ""))

	var invoiceStatus string
	var paidAt sql.NullTime
	err = db.QueryRowContext(ctx, `
		SELECT status::text, paid_at
		FROM payment_invoices
		WHERE invoice_id = $1::bigint
	`, invoiceID).Scan(&invoiceStatus, &paidAt)
	require.NoError(t, err)
	require.Equal(t, "paid", invoiceStatus)
	require.True(t, paidAt.Valid)

	var deliveryStatus string
	var attempts int
	var responseStatus int
	err = db.QueryRowContext(ctx, `
		SELECT status, attempts, response_status
		FROM merchant_webhook_deliveries
		WHERE transaction_id = (
			SELECT transaction_id
			FROM payment_transactions
			WHERE external_transaction_id = $1
		)
		  AND event_type = 'transaction.confirmed'
		ORDER BY webhook_delivery_id DESC
		LIMIT 1
	`, transactionID).Scan(&deliveryStatus, &attempts, &responseStatus)
	require.NoError(t, err)
	require.Equal(t, "delivered", deliveryStatus)
	require.Equal(t, 1, attempts)
	require.Equal(t, 204, responseStatus)
}
