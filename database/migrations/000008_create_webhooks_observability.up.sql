CREATE TABLE merchant_webhooks (
    webhook_id BIGSERIAL PRIMARY KEY,
    store_id BIGINT NOT NULL REFERENCES stores(store_id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    event_types TEXT[] NOT NULL DEFAULT ARRAY['transaction.created', 'transaction.submitted', 'transaction.validated', 'transaction.broadcasted', 'transaction.confirmed', 'transaction.failed'],
    secret_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    failure_count INT NOT NULL DEFAULT 0,
    last_success_at TIMESTAMPTZ,
    last_failure_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (store_id, url),
    CONSTRAINT chk_merchant_webhooks_url CHECK (url ~ '^https?://'),
    CONSTRAINT chk_merchant_webhooks_failure_count CHECK (failure_count >= 0)
);

CREATE TABLE merchant_webhook_deliveries (
    webhook_delivery_id BIGSERIAL PRIMARY KEY,
    webhook_id BIGINT REFERENCES merchant_webhooks(webhook_id) ON DELETE SET NULL,
    store_id BIGINT REFERENCES stores(store_id) ON DELETE SET NULL,
    transaction_id BIGINT REFERENCES payment_transactions(transaction_id) ON DELETE SET NULL,
    event_type VARCHAR(120) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(40) NOT NULL DEFAULT 'pending',
    attempts INT NOT NULL DEFAULT 0,
    response_status INT,
    response_body TEXT,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    delivered_at TIMESTAMPTZ,
    CONSTRAINT chk_webhook_delivery_status CHECK (status IN ('pending', 'delivered', 'failed')),
    CONSTRAINT chk_webhook_delivery_attempts CHECK (attempts >= 0),
    CONSTRAINT chk_webhook_delivery_response_status CHECK (response_status IS NULL OR response_status BETWEEN 100 AND 599)
);

CREATE INDEX idx_merchant_webhooks_store_id ON merchant_webhooks(store_id);
CREATE INDEX idx_merchant_webhooks_active ON merchant_webhooks(is_active);
CREATE INDEX idx_webhook_deliveries_store_id ON merchant_webhook_deliveries(store_id);
CREATE INDEX idx_webhook_deliveries_transaction_id ON merchant_webhook_deliveries(transaction_id);
CREATE INDEX idx_webhook_deliveries_status_created ON merchant_webhook_deliveries(status, created_at DESC);

CREATE OR REPLACE VIEW vw_webhook_delivery_status AS
SELECT
    d.webhook_delivery_id,
    d.webhook_id,
    d.store_id,
    s.store_name,
    d.transaction_id,
    d.event_type,
    d.status,
    d.attempts,
    d.response_status,
    d.error_message,
    d.created_at,
    d.delivered_at
FROM merchant_webhook_deliveries d
LEFT JOIN stores s ON s.store_id = d.store_id;

GRANT ALL PRIVILEGES ON merchant_webhooks, merchant_webhook_deliveries TO depay_admin;
GRANT ALL PRIVILEGES ON merchant_webhooks_webhook_id_seq, merchant_webhook_deliveries_webhook_delivery_id_seq TO depay_admin;
GRANT SELECT, INSERT, UPDATE, DELETE ON merchant_webhooks, merchant_webhook_deliveries TO depay_app;
GRANT USAGE, SELECT ON merchant_webhooks_webhook_id_seq, merchant_webhook_deliveries_webhook_delivery_id_seq TO depay_app;
GRANT SELECT ON merchant_webhooks, merchant_webhook_deliveries, vw_webhook_delivery_status TO depay_merchant_readonly, depay_compliance;
