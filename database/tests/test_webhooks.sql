\echo 'Webhook observability tests'

SELECT 'merchant_webhooks_seeded' AS test_name, count(*) AS row_count
FROM merchant_webhooks;

SELECT 'vw_webhook_delivery_status_queryable' AS test_name, count(*) AS row_count
FROM vw_webhook_delivery_status;
