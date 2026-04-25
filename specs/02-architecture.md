# 02. Target Architecture

## Архитектурный стиль

Монорепозиторий с несколькими backend-сервисами, общей PostgreSQL-БД, Redis-кэшем, RabbitMQ-событиями и веб-приложением.

```text
apps/web
  |
api-gateway или web-api
  |
  |-- user-service
  |-- merchant-service
  |-- wallet-service
  |-- transaction-core-service
  |-- transaction-validation-service
  |-- gas-info-service
  |-- analytics-service
  |-- admin-service
  |
PostgreSQL
Redis
RabbitMQ
```

## Правило MVP

Для курсовой можно временно объединить часть admin/analytics API в один `web-api`, но доменную архитектуру сохранить в документации и структуре.

## Сервисы

### `user-service`

Регистрация, login, refresh token, профиль, KYC-заявка.

### `merchant-service`

Магазины, verification, invoices, terminals, merchant wallets.

### `wallet-service`

Кошельки, балансы, история балансов, синхронизация.

### `transaction-core-service`

Создание транзакций, статусы, связь с invoice/NFC session, RabbitMQ events.

### `transaction-validation-service`

Проверка адресов, баланса, KYC, магазина, blacklist, risk scoring.

### `gas-info-service`

Gas price, RPC nodes, Redis cache, история проверок.

### `analytics-service`

SQL-функции и данные для графиков.

### `admin-service`

Таблицы БД, запуск функций, аудит, risk alerts, справочники.

## Shared packages

- `shared/config`;
- `shared/db`;
- `shared/logging`;
- `shared/middleware`;
- `shared/auth`;
- `shared/errors`;
- `shared/events`;
- `shared/validation`.

## RabbitMQ events

- `transaction.created`;
- `transaction.submitted`;
- `transaction.validated`;
- `transaction.broadcasted`;
- `transaction.confirmed`;
- `transaction.failed`;
- `balance.updated`;
- `risk_alert.created`.

## Хранилища

### PostgreSQL

Источник истины: users, stores, wallets, balances, invoices, transactions, KYC, risk, audit.

### Redis

Кэш gas info, balance info, rate limit state, temporary session state.

### RabbitMQ

Асинхронная обработка транзакций. Для MVP должен быть optional mode.
