# 06. Backend Plan

## Общие правила

Слои:

```text
controller -> service -> repository -> database
```

Все методы сервисов и репозиториев принимают `context.Context`. In-memory repositories оставить только для тестов.

## Shared packages

Добавить:

- `shared/db` — PostgreSQL connection;
- `shared/auth` — JWT, roles;
- `shared/errors` — типизированные ошибки;
- `shared/events` — RabbitMQ events;
- `shared/validation` — общие проверки.

## User Service

Исправить:

- заменить map repository на PostgreSQL;
- заменить fixed ID на UUID/BIGSERIAL;
- реализовать refresh tokens;
- добавить `/user/me`;
- добавить KYC endpoints.

Endpoints:

```http
POST /api/user/register
POST /api/user/login
POST /api/user/refresh-token
GET  /api/user/me
PUT  /api/user/me
POST /api/user/kyc
GET  /api/user/kyc/status
```

## Merchant Service

Новый сервис.

Endpoints:

```http
POST /api/merchant/register
POST /api/merchant/login
GET  /api/merchant/me
POST /api/merchant/verification
POST /api/merchant/invoices
GET  /api/merchant/invoices
POST /api/merchant/terminals
GET  /api/merchant/terminals
```

## Wallet Service

Исправить:

- убрать статический список;
- хранить wallets и balances в PostgreSQL;
- Redis оставить только как cache;
- добавить создание и удаление кошелька.

## Transaction Core Service

Исправить:

- PostgreSQL repository;
- связь с invoice/NFC session;
- status lifecycle;
- status endpoint;
- идемпотентность;
- RabbitMQ events.

## Transaction Validation Service

Добавить проверки:

- сумма > 0;
- адрес корректен;
- актив существует;
- кошелек принадлежит пользователю;
- магазинный кошелек принадлежит магазину;
- баланс достаточен;
- KYC для крупных сумм;
- магазин verified;
- blacklist.

## Gas Info Service

Добавить:

- запись истории gas;
- `rpc_node_checks`;
- endpoint истории;
- mock provider оставить для dev.

## Admin/Analytics API

Нужен для курсовой:

```http
GET  /api/admin/tables
GET  /api/admin/tables/{table_name}
POST /api/admin/functions/{function_name}/execute
GET  /api/analytics/store-turnover
GET  /api/analytics/transaction-statuses
GET  /api/analytics/failed-transactions
GET  /api/analytics/rpc-health
```
