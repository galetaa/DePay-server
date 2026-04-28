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

WITH target AS (
    SELECT mw.webhook_id, mw.store_id, pt.transaction_id
    FROM merchant_webhooks mw
    JOIN payment_transactions pt ON pt.store_id = mw.store_id
    WHERE mw.is_active = true
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
    error_message,
    last_attempt_at,
    next_attempt_at
)
SELECT
    webhook_id,
    store_id,
    transaction_id,
    'transaction.confirmed',
    jsonb_build_object('type', 'transaction.confirmed', 'payload', jsonb_build_object('transaction_id', transaction_id)),
    'retry_scheduled',
    1,
    500,
    'sql test retry',
    now(),
    now() + interval '30 seconds'
FROM target;

SELECT 'webhook_retry_status_queryable' AS test_name, count(*) AS row_count
FROM vw_webhook_delivery_status
WHERE status = 'retry_scheduled'
  AND next_attempt_at IS NOT NULL
  AND payload IS NOT NULL;

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

    IF NOT EXISTS (
        SELECT 1
        FROM vw_webhook_delivery_status
        WHERE status = 'retry_scheduled'
          AND attempts = 1
          AND response_status = 500
          AND next_attempt_at IS NOT NULL
    ) THEN
        RAISE EXCEPTION 'webhook retry status was not accepted or exposed';
    END IF;
END;
$$;

ROLLBACK;
