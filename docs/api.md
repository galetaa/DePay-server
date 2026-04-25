# DePay API Summary

This document lists the endpoints that currently exist in the repository. For machine-readable details, see `docs/openAPI.yaml`.

## Common

- `GET /health`
- `GET /metrics`

`/metrics` exists on backend Go services. The web app is served separately by Vite in local demo mode.

## User Service `:8080`

- `POST /api/user/register`
- `POST /api/user/login`
- `POST /api/user/refresh-token`
- `GET /api/user/me`
- `PUT /api/user/me`
- `POST /api/user/kyc`
- `GET /api/user/kyc/status`

Legacy routes still present:

- `POST /user/register`
- `POST /user/login`
- `POST /user/refresh-token`

Protected `/api/user/*` profile/KYC routes use JWT middleware.

## Merchant Service `:8083`

- `POST /api/merchant/register`
- `POST /api/merchant/login`
- `GET /api/merchant/me`
- `POST /api/merchant/verification`
- `GET /api/merchant/verification/status`
- `POST /api/merchant/invoices`
- `GET /api/merchant/invoices`
- `POST /api/merchant/terminals`
- `GET /api/merchant/terminals`
- `POST /api/merchant/webhooks`
- `GET /api/merchant/webhooks`
- `DELETE /api/merchant/webhooks/:webhook_id`

Protected merchant routes use JWT middleware. Invoice, terminal and webhook creation require merchant verification in service logic.

## Wallet Service `:8084`

- `POST /api/wallets`
- `GET /api/wallets`
- `GET /api/wallets/:wallet_id`
- `DELETE /api/wallets/:wallet_id`
- `GET /api/wallets/:wallet_id/balances`
- `POST /api/wallets/:wallet_id/sync`
- `GET /api/wallets/:wallet_id/balance`

Legacy routes still present:

- `GET /wallet/export`
- `POST /wallet/balance`

Wallet balance lookup uses Redis as cache and mock/RPC provider as source.

## Transaction Core Service `:8085`

- `POST /api/transaction/initiate`
- `POST /api/transaction/submit`
- `GET /api/transaction/:transaction_id/status`
- `GET /api/transaction/:transaction_id`
- `POST /api/transaction/:transaction_id/submit`
- `POST /api/transaction/:transaction_id/validate`
- `POST /api/transaction/:transaction_id/broadcast`
- `POST /api/transaction/:transaction_id/confirm`
- `POST /api/transaction/:transaction_id/cancel`

Legacy route still present:

- `POST /transaction/initiate`

Supported lifecycle:

- `created -> submitted -> validated -> broadcasted -> confirmed`
- `created/submitted/validated/broadcasted -> failed` where service logic marks failure
- `created/submitted -> cancelled`

The PostgreSQL trigger enforces this lifecycle for persisted transactions.

## Transaction Validation Service `:8081`

- `POST /api/transaction/validate`

Legacy route still present:

- `POST /transaction/validate`

Validation checks amount, EVM addresses, signature, transaction status, asset, wallet ownership, KYC, merchant verification, blacklist, risk alerts and balance when PostgreSQL is configured.

## Gas Info Service `:8082`

- `GET /api/transactions/gas-info`
- `GET /api/gas-info/history`

Legacy route still present:

- `GET /gas-info`

## KYC Service `:8086`

- `POST /api/kyc`

Legacy route still present:

- `POST /kyc`

Uses a mock provider by default; HTTP provider mode is controlled by env.

## Admin Service `:8090`

These endpoints are for coursework/admin/demo usage. They are intentionally treated as local demo/admin surface, not public production API.

- `GET /api/admin/tables`
- `GET /api/admin/tables/:table_name`
- `POST /api/admin/functions/:function_name/execute`
- `GET /api/admin/audit-logs`
- `GET /api/admin/risk-alerts`
- `GET /api/analytics/store-turnover`
- `GET /api/analytics/transaction-statuses`
- `GET /api/analytics/failed-transactions`
- `GET /api/analytics/rpc-health`
- `POST /api/admin/demo/invoices`
- `POST /api/admin/demo/payments`

## Web UI `:5173`

- `/admin/tables`
- `/admin/functions`
- `/admin/analytics`
- `/admin/demo`

