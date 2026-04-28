DROP VIEW IF EXISTS vw_webhook_delivery_status;

CREATE VIEW vw_webhook_delivery_status AS
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

GRANT SELECT ON vw_webhook_delivery_status TO depay_merchant_readonly, depay_compliance;

DROP INDEX IF EXISTS idx_webhook_deliveries_retry_due;

UPDATE merchant_webhook_deliveries
SET status = 'failed'
WHERE status IN ('retry_scheduled', 'dead_letter');

ALTER TABLE merchant_webhook_deliveries
    DROP CONSTRAINT chk_webhook_delivery_status,
    ADD CONSTRAINT chk_webhook_delivery_status
        CHECK (status IN ('pending', 'delivered', 'failed'));

ALTER TABLE merchant_webhook_deliveries
    DROP COLUMN IF EXISTS next_attempt_at,
    DROP COLUMN IF EXISTS last_attempt_at;
