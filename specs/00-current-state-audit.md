# 00. Аудит текущего состояния

## Общий вывод

Текущий `DePay-server` — хороший API-прототип, но еще не полноценная информационная система. Уже есть микросервисная структура, Gin-сервисы, тесты, OpenAPI-документация, Redis/RabbitMQ-заготовки и инфраструктурные идеи. Для курсовой и сильного pet project нужно добавить PostgreSQL, веб, миграции, SQL-функции, триггеры, роли и полноценный платежный flow.

## Что уже есть

### `user-service`

Готово:

- регистрация;
- login;
- refresh-token endpoint;
- bcrypt-хэширование пароля;
- JWT-ответ;
- тесты.

Проблемы:

- данные в `map`, а не в PostgreSQL;
- фиксированный `generateID`;
- refresh-token фактически заглушка;
- нет ролей и middleware авторизации;
- KYC описан в API, но не реализован полноценно.

### `wallet-service`

Готово:

- экспорт кошельков;
- получение баланса;
- Redis-кэш;
- mock balance provider;
- тесты.

Проблемы:

- кошельки статические;
- баланс фиксированный;
- нет таблиц `wallets`, `balances`, `assets`, `blockchains`;
- нет добавления кошельков;
- нет магазинных кошельков.

### `transaction-core-service`

Готово:

- инициализация транзакции;
- in-memory repository;
- RabbitMQ publish;
- `SKIP_RABBITMQ=true` для тестов;
- тесты.

Проблемы:

- нет PostgreSQL;
- нет invoice/NFC session;
- нет статусов жизненного цикла;
- нет status endpoint;
- `Currency` ставится в `ETH` принудительно.

### `transaction-validation-service`

Готово:

- endpoint валидации;
- проверка amount != 0;
- проверка длины EVM-адресов;
- тесты.

Проблемы:

- нет проверки подписи;
- нет проверки баланса в БД;
- нет проверки KYC;
- нет merchant verification;
- нет blacklist/risk alerts.

### `gas-info-service`

Готово:

- gas-info endpoint;
- Redis-кэш;
- mock external provider.

Проблемы:

- нет истории газа;
- нет `rpc_nodes` и `rpc_node_checks`;
- нет графиков состояния сети.

## Главный план исправления

1. Сначала БД и SQL-часть.
2. Потом PostgreSQL repositories.
3. Потом web MVP.
4. Потом merchant/invoice/terminal flow.
5. Потом production-like фичи: webhooks, API keys, real blockchain provider, Kong/Vault/K8s.
