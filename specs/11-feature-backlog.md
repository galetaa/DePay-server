# 11. Feature Backlog

## Critical

### DB-001 Migrations

Создать миграции типов, таблиц, индексов, функций, триггеров и ролей.

### DB-002 Seed data

Нужно минимум:

- 10 users;
- 5 stores;
- 4 blockchains;
- 8 assets;
- 20 wallets;
- 80 transactions;
- failed transactions;
- users without KYC;
- rpc checks.

### SQL-001 12 analytical functions

Все функции должны возвращать данные на seed.

### SQL-002 7 triggers

Проверить ошибки и сделать скриншоты.

### WEB-001 Admin table viewer

Выбор таблицы и просмотр строк.

### WEB-002 Function runner

Выбор функции, параметры, результат.

### WEB-003 Charts

3+ графика разных типов.

## High

### USER-001 PostgreSQL repository

Заменить map.

### USER-002 Refresh tokens

Реальное хранение refresh token hash.

### MERCHANT-001 Merchant service

Регистрация, login, verification.

### INVOICE-001 Payment invoices

Создание и оплата invoice.

### TERMINAL-001 NFC sessions

Терминал создает временную сессию.

### TX-001 Transaction lifecycle

Статусы, validation, audit, risk.

## Medium

- RPC node checks;
- merchant webhooks;
- API keys;
- OpenAPI cleanup;
- README screenshots;
- integration tests.

## Later

- real testnet provider;
- EVM signature verification;
- Kubernetes;
- Kong;
- Vault;
- Prometheus/Grafana.
