package services

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"shared/events"
	"shared/logging"
	"transaction-core-service/internal/models"

	"go.uber.org/zap"
)

type WebhookDispatcher interface {
	Dispatch(ctx context.Context, eventType string, tx models.Transaction)
}

type noopWebhookDispatcher struct{}

func (noopWebhookDispatcher) Dispatch(context.Context, string, models.Transaction) {}

type PostgresWebhookDispatcher struct {
	db     *sql.DB
	mode   string
	client *http.Client
}

func NewPostgresWebhookDispatcher(db *sql.DB) WebhookDispatcher {
	mode := os.Getenv("WEBHOOK_DELIVERY_MODE")
	if mode == "" {
		mode = "log"
	}
	return &PostgresWebhookDispatcher{
		db:     db,
		mode:   mode,
		client: &http.Client{Timeout: envDuration("WEBHOOK_DELIVERY_TIMEOUT_MS", 5*time.Second)},
	}
}

func (d *PostgresWebhookDispatcher) Dispatch(ctx context.Context, eventType string, tx models.Transaction) {
	if d == nil || d.db == nil || eventType == "" {
		return
	}
	if d.mode == "disabled" {
		return
	}

	hooks, err := d.loadWebhooks(ctx, tx.StoreID, eventType)
	if err != nil {
		logging.Logger.Warn("Failed to load merchant webhooks", zap.Error(err), zap.String("event_type", eventType))
		return
	}

	for _, hook := range hooks {
		if err := d.deliver(ctx, hook, eventType, tx); err != nil {
			logging.Logger.Warn("Failed to deliver merchant webhook", zap.Error(err), zap.String("webhook_id", hook.ID), zap.String("event_type", eventType))
		}
	}
}

func (d *PostgresWebhookDispatcher) loadWebhooks(ctx context.Context, storeID string, eventType string) ([]webhookTarget, error) {
	rows, err := d.db.QueryContext(ctx, `
		SELECT webhook_id::text, store_id::text, url, secret_hash
		FROM merchant_webhooks
		WHERE store_id = $1::bigint
		  AND is_active = true
		  AND event_types @> ARRAY[$2]::text[]
		ORDER BY webhook_id
	`, storeID, eventType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	targets := make([]webhookTarget, 0)
	for rows.Next() {
		var target webhookTarget
		if err := rows.Scan(&target.ID, &target.StoreID, &target.URL, &target.SecretHash); err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}
	return targets, rows.Err()
}

func (d *PostgresWebhookDispatcher) deliver(ctx context.Context, target webhookTarget, eventType string, tx models.Transaction) error {
	event := events.New(eventType, map[string]any{
		"transaction_id":     tx.TransactionID,
		"store_id":           tx.StoreID,
		"user_id":            tx.UserID,
		"invoice_id":         tx.InvoiceID,
		"amount":             tx.Amount,
		"amount_usdt":        tx.AmountUSDT,
		"currency":           tx.Currency,
		"status":             tx.Status,
		"blockchain_tx_hash": tx.BlockchainTxHash,
		"failure_reason":     tx.FailureReason,
	})
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	deliveryID, err := d.insertDelivery(ctx, target, tx.TransactionID, eventType, body)
	if err != nil {
		return err
	}

	if d.mode != "http" {
		return d.markDelivered(ctx, target, deliveryID, 204, "webhook delivery log mode")
	}
	if err := validateWebhookTarget(ctx, target.URL); err != nil {
		_ = d.markFailed(ctx, target, deliveryID, 0, "", err.Error())
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target.URL, bytes.NewReader(body))
	if err != nil {
		_ = d.markFailed(ctx, target, deliveryID, 0, "", err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DePay-Event", eventType)
	req.Header.Set("X-DePay-Delivery", deliveryID)
	req.Header.Set("X-DePay-Signature", signPayload(body, target.SecretHash))

	resp, err := d.client.Do(req)
	if err != nil {
		_ = d.markFailed(ctx, target, deliveryID, 0, "", err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := fmt.Sprintf("webhook endpoint returned status %d", resp.StatusCode)
		_ = d.markFailed(ctx, target, deliveryID, resp.StatusCode, "", msg)
		return fmt.Errorf("webhook endpoint returned status %d", resp.StatusCode)
	}
	return d.markDelivered(ctx, target, deliveryID, resp.StatusCode, "")
}

func (d *PostgresWebhookDispatcher) insertDelivery(ctx context.Context, target webhookTarget, transactionID string, eventType string, payload []byte) (string, error) {
	var deliveryID string
	err := d.db.QueryRowContext(ctx, `
		INSERT INTO merchant_webhook_deliveries(webhook_id, store_id, transaction_id, event_type, payload)
		VALUES (
			$1::bigint,
			$2::bigint,
			(SELECT transaction_id FROM payment_transactions WHERE transaction_id::text = $3 OR external_transaction_id = $3 LIMIT 1),
			$4,
			$5::jsonb
		)
		RETURNING webhook_delivery_id::text
	`, target.ID, target.StoreID, transactionID, eventType, string(payload)).Scan(&deliveryID)
	return deliveryID, err
}

func (d *PostgresWebhookDispatcher) markDelivered(ctx context.Context, target webhookTarget, deliveryID string, status int, body string) error {
	if _, err := d.db.ExecContext(ctx, `
		UPDATE merchant_webhook_deliveries
		SET status = 'delivered',
		    attempts = attempts + 1,
		    response_status = NULLIF($2, 0),
		    response_body = NULLIF($3, ''),
		    delivered_at = now()
		WHERE webhook_delivery_id = $1::bigint
	`, deliveryID, status, body); err != nil {
		return err
	}
	_, err := d.db.ExecContext(ctx, `
		UPDATE merchant_webhooks
		SET last_success_at = now(),
		    updated_at = now()
		WHERE webhook_id = $1::bigint
	`, target.ID)
	return err
}

func (d *PostgresWebhookDispatcher) markFailed(ctx context.Context, target webhookTarget, deliveryID string, status int, body string, message string) error {
	if _, err := d.db.ExecContext(ctx, `
		UPDATE merchant_webhook_deliveries
		SET status = 'failed',
		    attempts = attempts + 1,
		    response_status = NULLIF($2, 0),
		    response_body = NULLIF($3, ''),
		    error_message = NULLIF($4, '')
		WHERE webhook_delivery_id = $1::bigint
	`, deliveryID, status, body, message); err != nil {
		return err
	}
	_, err := d.db.ExecContext(ctx, `
		UPDATE merchant_webhooks
		SET failure_count = failure_count + 1,
		    last_failure_at = now(),
		    updated_at = now()
		WHERE webhook_id = $1::bigint
	`, target.ID)
	return err
}

func signPayload(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

type webhookTarget struct {
	ID         string
	StoreID    string
	URL        string
	SecretHash string
}

func validateWebhookTarget(ctx context.Context, rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid webhook url: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("webhook url scheme must be http or https")
	}
	host := parsed.Hostname()
	if host == "" {
		return errors.New("webhook url host is required")
	}
	lowerHost := strings.ToLower(host)
	if lowerHost == "localhost" || strings.HasSuffix(lowerHost, ".localhost") {
		return errors.New("webhook url host is not allowed")
	}

	resolver := net.DefaultResolver
	ips, err := resolver.LookupIP(ctx, "ip", host)
	if err != nil {
		return fmt.Errorf("unable to resolve webhook host: %w", err)
	}
	for _, ip := range ips {
		if isRestrictedIP(ip) {
			return fmt.Errorf("webhook target resolves to restricted address %s", ip.String())
		}
	}
	return nil
}

func isRestrictedIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() {
		return true
	}
	if ipv4 := ip.To4(); ipv4 != nil && ipv4[0] == 169 && ipv4[1] == 254 {
		return true
	}
	return false
}

