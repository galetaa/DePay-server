# 08. API Contracts

## Общий формат

Все endpoints через `/api`.

Ошибка:

```json
{
  "error": {
    "code": "INVALID_INPUT",
    "message": "Amount must be greater than zero",
    "details": {}
  }
}
```

Успех:

```json
{
  "data": {}
}
```

## User API

```http
POST /api/user/register
POST /api/user/login
POST /api/user/refresh-token
GET  /api/user/me
PUT  /api/user/me
POST /api/user/kyc
GET  /api/user/kyc/status
```

## Merchant API

```http
POST /api/merchant/register
POST /api/merchant/login
GET  /api/merchant/me
POST /api/merchant/verification
GET  /api/merchant/verification/status
POST /api/merchant/invoices
GET  /api/merchant/invoices
GET  /api/merchant/invoices/{invoice_id}
POST /api/merchant/terminals
GET  /api/merchant/terminals
```

## Wallet API

```http
POST   /api/wallets
GET    /api/wallets
GET    /api/wallets/{wallet_id}
DELETE /api/wallets/{wallet_id}
GET    /api/wallets/{wallet_id}/balances
POST   /api/wallets/{wallet_id}/sync
GET    /api/wallets/{address}/balance
```

## Terminal API

```http
POST /api/terminal/register
POST /api/terminal/refresh-token
POST /api/terminal/sessions
GET  /api/terminal/sessions/{session_id}
```

## Transaction API

```http
POST /api/transaction/initiate
POST /api/transaction/submit
POST /api/transaction/validate
GET  /api/transaction/{transaction_id}
GET  /api/transaction/{transaction_id}/status
POST /api/transaction/{transaction_id}/cancel
```

## Admin/Analytics API

```http
GET  /api/admin/tables
GET  /api/admin/tables/{table_name}
POST /api/admin/functions/{function_name}/execute
GET  /api/admin/audit-logs
GET  /api/admin/risk-alerts
GET  /api/analytics/store-turnover
GET  /api/analytics/transaction-statuses
GET  /api/analytics/failed-transactions
GET  /api/analytics/rpc-health
```

## Что исправить в текущем API

- унифицировать `/wallet/export` и `/wallets/export`;
- унифицировать `/gas-info` и `/transactions/gas-info`;
- добавить `/api` prefix;
- добавить merchant/invoice/terminal/admin endpoints;
- синхронизировать `openAPI.yaml` с кодом.
