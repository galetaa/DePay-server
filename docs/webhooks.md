# Webhooks

Merchant webhooks notify stores about invoice and transaction lifecycle events.

## Events

- `invoice.created`
- `invoice.paid`
- `invoice.expired`
- `transaction.created`
- `transaction.submitted`
- `transaction.validated`
- `transaction.broadcasted`
- `transaction.confirmed`
- `transaction.failed`
- `transaction.cancelled`

## Delivery Headers

Transaction-core sends HTTP callbacks with these headers when `WEBHOOK_DELIVERY_MODE=http`:

- `X-DePay-Event`
- `X-DePay-Delivery`
- `X-DePay-Timestamp`
- `X-DePay-Signature`

The signature is `sha256=<hex hmac>` over:

```text
<timestamp>.<raw_body>
```

## Retry Policy

Each delivery starts as `pending`. Retryable failures are network errors, timeouts, `429`, `408`, and `5xx` responses.

```text
attempt 1: immediate
attempt 2: +30s
attempt 3: +2m
attempt 4: +10m
attempt 5: +30m
then dead_letter
```

Delivery statuses:

- `delivered`
- `failed`
- `retry_scheduled`
- `dead_letter`

`vw_webhook_delivery_status` exposes `payload`, `attempts`, response metadata, `last_attempt_at`, and `next_attempt_at` for the merchant dashboard.
