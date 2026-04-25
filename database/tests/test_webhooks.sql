\echo 'Webhook observability tests'

BEGIN;

SELECT 'merchant_webhooks_seeded' AS test_name, count(*) AS row_count
FROM merchant_webhooks;

WITH target AS (
    SELECT mw.webhook_id, mw.store_id, pt.transaction_id
    FROM merchant_webhooks mw
    JOIN payment_transactions pt ON pt.store_id = mw.store_id
    WHERE mw.is_active = true
      AND mw.event_types @> ARRAY['transaction.confirmed']::text[]
    LIMIT 1
)
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
    webhook_id,
    store_id,
    transaction_id,
    'transaction.confirmed',
    jsonb_build_object('type', 'transaction.confirmed', 'payload', jsonb_build_object('transaction_id', transaction_id)),
    'delivered',
    1,
    204,
    'sql test delivery',
    now()
FROM target;

SELECT 'vw_webhook_delivery_status_queryable' AS test_name, count(*) AS row_count
FROM vw_webhook_delivery_status
WHERE event_type = 'transaction.confirmed'
  AND status = 'delivered';

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM vw_webhook_delivery_status
        WHERE event_type = 'transaction.confirmed'
          AND status = 'delivered'
          AND response_status = 204
    ) THEN
        RAISE EXCEPTION 'webhook delivery view did not expose inserted delivery';
    END IF;
END;
$$;

ROLLBACK;
