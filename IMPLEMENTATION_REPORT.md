# Отчет по реализованным работам DePay

Дата составления: 2026-04-26

Файл подготовлен как рабочий отчет по изменениям, которые были добавлены с начала текущей серии работ над DePay, и актуализирован после coursework stabilization.

## 1. Краткое резюме

Проект был доведен из состояния MVP/прототипа до coursework-ready состояния с полноценной PostgreSQL-схемой, миграциями, seed-данными, SQL-функциями, триггерами, SQL-тестами, backend API, admin API и web-интерфейсом для демонстрации.

После этого поверх coursework MVP были добавлены production-like доработки: PostgreSQL-backed repositories, merchant flow, webhook observability, pluggable provider integrations, Prometheus metrics, Docker Compose profiles, Kong/Vault/K8s конфигурации, kind/kubectl validation workflow и документация по запуску/проверке.

Главный результат: DePay теперь можно показывать как курсовой проект с базой данных, SQL-логикой, web demo-flow и воспроизводимым локальным окружением, а также развивать дальше как pet project с более production-like архитектурой.

## 2. Основные добавленные директории и файлы

Добавлены или существенно расширены:

- `database/` - миграции, seed data, SQL tests.
- `admin-service/` - backend API для admin/web интерфейса.
- `merchant-service/` - merchant registration/login/verification/invoices/terminals/webhooks.
- `apps/web/` - React + Vite + TypeScript web-интерфейс.
- `shared/auth`, `shared/db`, `shared/events`, `shared/validation`, `shared/observability` - общая инфраструктура.
- `observability/` - Prometheus config, Grafana provisioning и dashboard.
- `kong/kong.yml` - declarative Kong config.
- `k8s/` - Kubernetes manifests для сервисов и инфраструктуры.
- `docs/coursework/README.md` - инструкция по coursework/demo readiness.
- `.env.example` - пример env-конфигурации.
- `Makefile` - локальные workflow targets для тестов, миграций, Docker Compose и Kubernetes validation.
- `AGENTS.md` - актуальные инструкции по реальности репозитория и workflow.

## 3. Документация и спецификации

Добавлен пакет спецификаций в `specs/`:

- `specs/00-current-state-audit.md` - аудит текущего состояния.
- `specs/01-product-scope.md` - границы продукта.
- `specs/02-architecture.md` - целевая архитектура.
- `specs/03-repository-structure.md` - структура репозитория.
- `specs/04-database-design.md` - дизайн базы данных.
- `specs/05-coursework-sql-pack.md` - SQL pack для курсовой.
- `specs/06-backend-plan.md` - backend plan.
- `specs/07-web-plan.md` - web plan.
- `specs/08-api-contracts.md` - API contracts.
- `specs/09-security-compliance.md` - security/compliance.
- `specs/10-roadmap.md` - roadmap.
- `specs/11-feature-backlog.md` - backlog.
- `specs/12-testing-devops.md` - testing/devops.
- `specs/13-coursework-map.md` - соответствие требованиям курсовой.
- `specs/14-defense-demo-script.md` - demo script для защиты.
- `specs/15-coding-standards.md` - coding standards.
- `specs/README.md` - оглавление.

Обновлены:

- `docs/openAPI.yaml` - API контракт приведен ближе к фактическим MVP endpoints.
- `docs/coursework/README.md` - добавлен путь запуска и проверки coursework-ready версии.
- `README.md` - добавлено описание проекта, архитектуры, запуска и demo-flow.
- `docs/api.md` - добавлен human-readable список фактически реализованных endpoints.
- `docs/defense-demo-script.md` - добавлен сценарий защиты.
- `docs/screenshots.md` - добавлен список скриншотов для отчета.
- `verification-results.md` - добавлен журнал проверок стабилизации.
- `AGENTS.md` - добавлена реальная информация о multi-module Go monorepo, тестах по модулям, external systems и conventions.

## 4. PostgreSQL и coursework SQL pack

Добавлен каталог `database/` с миграциями `golang-migrate`:

- `000001_create_extensions` - extensions.
- `000002_create_types` - enum-типы.
- `000003_create_tables` - основная схема таблиц.
- `000004_create_indexes` - индексы.
- `000005_create_triggers` - trigger functions и triggers.
- `000006_create_functions` - SQL-функции и views.
- `000007_create_roles` - роли и grants.
- `000008_create_webhooks_observability` - webhook tables/views/grants.

### 4.1 Enum-типы

Добавлены enum-типы:

- `kyc_status_enum`
- `verification_status_enum`
- `transaction_status_enum`
- `invoice_status_enum`
- `nfc_session_status_enum`
- `terminal_status_enum`
- `rpc_node_status_enum`
- `risk_level_enum`
- `risk_alert_status_enum`
- `document_type_enum`

### 4.2 Таблицы

Добавлены основные таблицы:

- `users`
- `roles`
- `user_roles`
- `refresh_tokens`
- `stores`
- `kyc_applications`
- `kyc_documents`
- `merchant_verification_applications`
- `blockchains`
- `assets`
- `wallets`
- `balances`
- `balance_history`
- `payment_invoices`
- `payment_invoice_items`
- `terminals`
- `nfc_sessions`
- `payment_transactions`
- `rpc_nodes`
- `rpc_node_checks`
- `exchange_rates`
- `blacklisted_wallets`
- `risk_alerts`
- `audit_logs`
- `merchant_webhooks`
- `merchant_webhook_deliveries`

### 4.3 SQL-функции

Реализованы 12 PostgreSQL-функций:

- `get_user_kyc_wallet_summary()`
- `get_user_wallet_balances(p_user_id BIGINT)`
- `get_wallet_asset_distribution(p_wallet_id BIGINT)`
- `get_transaction_card(p_tx_id BIGINT)`
- `get_user_transaction_history(p_user_id BIGINT, p_date_from TIMESTAMPTZ, p_date_to TIMESTAMPTZ)`
- `get_store_transaction_history(p_store_id BIGINT, p_date_from TIMESTAMPTZ, p_date_to TIMESTAMPTZ)`
- `get_blockchain_asset_activity(p_chain_id BIGINT, p_date_from TIMESTAMPTZ, p_date_to TIMESTAMPTZ)`
- `get_rpc_nodes_activity(p_chain_id BIGINT, p_date_from TIMESTAMPTZ, p_date_to TIMESTAMPTZ)`
- `get_store_turnover(p_store_id BIGINT, p_date_from TIMESTAMPTZ, p_date_to TIMESTAMPTZ)`
- `get_store_success_rate(p_store_id BIGINT, p_date_from TIMESTAMPTZ, p_date_to TIMESTAMPTZ)`
- `get_unverified_active_users(p_min_tx_count BIGINT, p_min_amount_usdt NUMERIC, p_date_from TIMESTAMPTZ, p_date_to TIMESTAMPTZ)`
- `get_failed_transactions_analytics(p_date_from TIMESTAMPTZ, p_date_to TIMESTAMPTZ)`

### 4.4 Views

Добавлены views для admin/web/analytics:

- `vw_user_wallet_balances`
- `vw_store_transactions`
- `vw_failed_transactions`
- `vw_rpc_node_status`
- `vw_compliance_kyc_queue`
- `vw_webhook_delivery_status`

### 4.5 Триггеры

Добавлены 7 триггеров:

- `trg_wallet_owner_check`
- `trg_balance_non_negative`
- `trg_balance_history`
- `trg_transaction_status_flow`
- `trg_transaction_completed_at`
- `trg_audit_payment_transactions`
- `trg_risk_alert_large_unverified_payment`

Trigger functions проверяют:

- корректность владельца wallet;
- запрет отрицательных balances;
- запись истории balances;
- допустимый lifecycle статусов транзакции;
- автоматическое заполнение `completed_at`;
- audit log для payment transactions;
- risk alert для крупных платежей unverified users.

### 4.6 Роли и grants

Добавлены PostgreSQL roles:

- `depay_admin`
- `depay_app`
- `depay_compliance`
- `depay_merchant_readonly`
- `depay_user_readonly`

Для ролей настроены grants на таблицы, sequences, views и webhook observability tables.

### 4.7 Seed data и SQL-тесты

Добавлен `database/seeds/seed_coursework.sql` с демонстрационными данными, на которых SQL-функции возвращают результат.

Добавлены SQL-тесты:

- `database/tests/test_functions.sql`
- `database/tests/test_triggers.sql`
- `database/tests/test_webhooks.sql`

Проверяются:

- что функции возвращают строки;
- owner check;
- запрет negative balance;
- invalid transaction status flow;
- audit log;
- webhook delivery status.

## 5. Shared backend инфраструктура

Добавлены и расширены shared packages:

- `shared/db` - подключение к PostgreSQL через `DATABASE_URL`.
- `shared/auth` - access/refresh token pair, hash refresh token.
- `shared/errors` - единый формат API errors/data responses.
- `shared/events` - простая структура событий.
- `shared/validation` - EVM address и positive amount validation.
- `shared/observability` - Prometheus HTTP metrics middleware и `/metrics` handler.

Обновлены:

- `shared/middleware/auth.go` - JWT middleware и role middleware.
- `shared/utils/utils.go` - JWT secret из env вместо hardcoded-only подхода.
- `shared/config/config.go` - env-driven конфигурация.
- `shared/middleware/cors.go` и error handler продолжили использоваться как общий middleware stack.

## 6. Backend services

### 6.1 user-service

Реализовано:

- PostgreSQL-backed users.
- Roles в access token.
- Refresh tokens через hash storage.
- Profile endpoints.
- KYC submit/status.
- Backward-compatible legacy routes и новые `/api/user/*` routes.
- In-memory repository оставлен для тестов.

Основные endpoints:

- `POST /api/user/register`
- `POST /api/user/login`
- `POST /api/user/refresh-token`
- `GET /api/user/me`
- `PUT /api/user/me`
- `POST /api/user/kyc`
- `GET /api/user/kyc/status`
- `GET /health`
- `GET /metrics`

### 6.2 merchant-service

Добавлен новый сервис `merchant-service`.

Реализовано:

- merchant register/login/me;
- JWT token pair для merchant;
- merchant verification submit/status;
- invoices;
- terminals;
- запрет invoice/terminal creation для unverified merchant;
- merchant webhook endpoints;
- PostgreSQL repository;
- in-memory repository для тестов/локального fallback.

Основные endpoints:

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
- `GET /health`
- `GET /metrics`

### 6.3 wallet-service

Реализовано:

- PostgreSQL-backed wallets.
- In-memory repository для тестов.
- CRUD кошельков.
- Balances endpoint.
- Wallet sync endpoint.
- Redis как cache для balance lookup.
- RPC balance provider через `WALLET_BALANCE_RPC_URL` или `BLOCKCHAIN_RPC_URL`.
- Mock balance provider для local/dev/test.
- Поддержка user wallets и store wallets на уровне схемы/репозитория.

Основные endpoints:

- `POST /api/wallets`
- `GET /api/wallets`
- `GET /api/wallets/:wallet_id`
- `DELETE /api/wallets/:wallet_id`
- `GET /api/wallets/:wallet_id/balances`
- `POST /api/wallets/:wallet_id/sync`
- `GET /api/wallets/:wallet_id/balance`
- legacy `GET /wallet/export`
- legacy `POST /wallet/balance`
- `GET /health`
- `GET /metrics`

### 6.4 transaction-core-service

Реализовано:

- PostgreSQL-backed transaction repository.
- In-memory repository для тестов.
- Transaction lifecycle через статусы `created`, `submitted`, `validated`, `broadcasted`, `confirmed`, `failed`, `cancelled`.
- Status endpoint.
- Cancel endpoint.
- Broadcast endpoint.
- Optional RabbitMQ mode через `SKIP_RABBITMQ`.
- Blockchain broadcaster:
  - mock mode для local/dev;
  - JSON-RPC provider mode через `BLOCKCHAIN_RPC_URL`.
- Webhook dispatcher:
  - noop mode;
  - PostgreSQL-backed dispatcher для `merchant_webhooks` и `merchant_webhook_deliveries`;
  - delivery log/status.
- Events dispatch на transaction lifecycle события.

Основные endpoints:

- `POST /api/transaction/initiate`
- `POST /api/transaction/submit`
- `GET /api/transaction/:transaction_id/status`
- `POST /api/transaction/:transaction_id/cancel`
- `POST /api/transaction/:transaction_id/submit`
- `POST /api/transaction/:transaction_id/validate`
- `POST /api/transaction/:transaction_id/broadcast`
- `POST /api/transaction/:transaction_id/confirm`
- `GET /api/transaction/:transaction_id`
- legacy `POST /transaction/initiate`
- `GET /health`
- `GET /metrics`

### 6.5 transaction-validation-service

Реализовано:

- Stateless validation fallback.
- PostgreSQL-backed validation service.
- Проверки:
  - amount positive;
  - EVM sender/recipient address;
  - signature required;
  - transaction exists;
  - status is not terminal;
  - amount matches persisted transaction;
  - currency matches asset;
  - sender wallet matches user wallet;
  - recipient wallet matches merchant/store wallet;
  - user KYC approved;
  - merchant verification approved;
  - wallet not blacklisted;
  - no open high/critical risk alert;
  - sufficient balance;
  - user/store wallet ownership.

Endpoints:

- `POST /api/transaction/validate`
- legacy `POST /transaction/validate`
- `GET /health`
- `GET /metrics`

### 6.6 gas-info-service

Реализовано:

- Gas info endpoint.
- Gas history endpoint.
- Redis cache.
- Mock provider для dev/local.
- HTTP provider mode через env provider URL.
- Metrics endpoint.

Endpoints:

- `GET /api/transactions/gas-info`
- `GET /api/gas-info/history`
- legacy `GET /gas-info`
- `GET /health`
- `GET /metrics`

### 6.7 kyc-service

Реализовано:

- Pluggable KYC provider.
- Mock provider для dev/local.
- HTTP provider mode через `KYC_PROVIDER_URL`, `KYC_PROVIDER_API_KEY`, `KYC_PROVIDER_TIMEOUT_MS`.
- API-prefixed endpoint.
- Metrics endpoint.

Endpoints:

- `POST /api/kyc`
- legacy `POST /kyc`
- `GET /health`
- `GET /metrics`

### 6.8 admin-service

Добавлен новый `admin-service`.

Реализовано:

- browser для разрешенных таблиц;
- executor для разрешенных SQL-функций;
- audit/risk endpoints;
- analytics endpoints;
- demo endpoints для web demo-flow.

Endpoints:

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
- `GET /health`
- `GET /metrics`

## 7. Web-интерфейс

Добавлен `apps/web` на:

- React;
- Vite;
- TypeScript;
- TanStack Query;
- Recharts;
- Vitest/Testing Library.

Добавлены страницы:

- `/admin/tables` - просмотр таблиц.
- `/admin/functions` - запуск SQL-функций.
- `/admin/analytics` - графики и аналитика.
- `/admin/demo` - demo-flow invoice/payment.

Добавлены компоненты:

- `Layout`
- `DataTable`
- `FunctionRunner` реализован внутри страницы функций.
- `StatusBadge`
- `ErrorAlert`
- `ChartCard`
- `DemoFlow` реализован через `DemoPage`

Графики:

- bar turnover;
- pie transaction statuses;
- line failed transactions / time series style analytics;
- RPC health/status analytics.

Demo-flow:

- создание demo invoice;
- отправка demo payment;
- показ transaction status;
- обновление данных через admin API.

Добавлены web tests:

- `apps/web/src/api/client.test.ts`
- `apps/web/src/pages/DemoPage.test.tsx`

## 8. Observability

Добавлено:

- `shared/observability/metrics.go`
- `/metrics` endpoint во всех Go-сервисах.
- Gin middleware для HTTP metrics.
- `observability/prometheus.yml`
- Grafana provisioning:
  - datasource Prometheus;
  - dashboard `depay-overview.json`.

Сервисы теперь имеют:

- `/health`
- `/metrics`

## 9. DevOps и local workflow

### 9.1 Makefile

Добавлены targets:

- `up`
- `backend-up`
- `web-up`
- `gateway-up`
- `observability-up`
- `secrets-up`
- `full-up`
- `prod-like-up`
- `dev-ready`
- `down`
- `reset-local-data`
- `wait-db`
- `migrate-up`
- `migrate-down`
- `seed`
- `sql-test`
- `test`
- `test-go`
- `web-dev`
- `web-build`
- `web-test`
- `kind-install`
- `kind-up`
- `kind-down`
- `k8s-validate`
- `k8s-dry-run`

### 9.2 Docker Compose

`docker-compose.yml` расширен profiles:

- base infra: PostgreSQL, Redis, RabbitMQ;
- `backend` - Go services;
- `web` - React/Vite web;
- `gateway` - Kong;
- `observability` - Prometheus/Grafana;
- `secrets` - Vault dev mode.

### 9.3 Env configuration

Добавлен `.env.example`.

Секреты и external URLs вынесены в env:

- PostgreSQL;
- Redis;
- RabbitMQ;
- JWT;
- Blockchain RPC;
- Wallet balance RPC;
- Gas provider;
- KYC provider;
- Webhook delivery mode;
- Grafana/Vault local secrets.

### 9.4 Kong

Добавлен `kong/kong.yml` и обновлены K8s/Kong manifests.

Kong routes покрывают:

- user;
- merchant;
- wallet;
- transaction;
- validation;
- gas;
- kyc;
- admin.

### 9.5 Vault

Добавлена dev-конфигурация Vault:

- Docker Compose profile `secrets`;
- K8s manifest `k8s/vault.yaml`;
- обновленный helper script `server-config/setup_vault_and_kong.sh`.

### 9.6 Kubernetes

Добавлены/обновлены manifests:

- `k8s/backend-services.yaml`
- `k8s/kong.yaml`
- `k8s/observability.yaml`
- `k8s/postgres.yaml`
- `k8s/rabbitmq.yaml`
- `k8s/redis.yaml`
- `k8s/user-service.yaml`
- `k8s/vault.yaml`

Добавлен kind/kubectl workflow:

- `make kind-up` - устанавливает/поднимает local kind cluster и переключает context.
- `make k8s-validate` - выполняет `kubectl apply --dry-run=server --validate=strict -f k8s -o name`.
- `make kind-down` - удаляет local kind cluster.

Проблема с `kubectl` была в пустом kubeconfig/current-context. Для проверки был создан local kind cluster `depay-kubectl-check`; server-side dry-run по `k8s/` прошел успешно.

## 10. Тесты и проверки

Добавлены или обновлены Go tests:

- `user-service/tests/user_service_test.go`
- `wallet-service/tests/wallet_service_test.go`
- `transaction-core-service/tests/transaction_core_test.go`
- `transaction-validation-service/tests/validation_service_test.go`
- `gas-info-service/tests/gas_info_test.go`
- `kyc-service/tests/kyc_service_test.go`

Добавлены web tests:

- `apps/web/src/api/client.test.ts`
- `apps/web/src/pages/DemoPage.test.tsx`

Добавлены SQL tests:

- `database/tests/test_functions.sql`
- `database/tests/test_triggers.sql`
- `database/tests/test_webhooks.sql`

В ходе работы прогонялись:

- `make test-go`
- `cd shared && go test ./...`
- `make migrate-up`
- `make seed`
- `make sql-test`
- `make web-test`
- `make web-build`
- `make reset-local-data`
- `docker compose --profile backend --profile web --profile gateway --profile observability --profile secrets config --quiet`
- health checks backend сервисов;
- transaction lifecycle smoke: initiate -> submit -> validate -> broadcast -> confirm;
- webhook delivery log smoke;
- `make k8s-validate`

## 11. Demo-flow, который теперь можно показывать

Базовый coursework demo-flow:

1. Поднять PostgreSQL/Redis/RabbitMQ.
2. Прогнать миграции.
3. Загрузить seed data.
4. Поднять backend services.
5. Открыть web admin.
6. В `/admin/tables` показать таблицы.
7. В `/admin/functions` запустить SQL-функции.
8. В `/admin/analytics` показать графики.
9. В `/admin/demo` создать invoice и провести payment.
10. Показать изменение transaction status, таблиц, графиков и webhook delivery log.

Локальные команды:

```bash
make dev-ready
make web-up
```

Для честного project-local запуска с пустой БД используется:

```bash
make reset-local-data
make dev-ready
make web-up
```

Проверки:

```bash
make test-go
make sql-test
make web-test
make web-build
make k8s-validate
```

## 12. Production-like доработки

Реализованы вещи, которые изначально относились к production-later:

- merchant webhooks;
- webhook delivery audit/observability;
- pluggable KYC provider;
- pluggable gas provider;
- wallet balance RPC provider;
- transaction blockchain broadcaster;
- Prometheus metrics;
- Grafana dashboard;
- Kong gateway config;
- Vault dev config;
- Kubernetes manifests;
- kind/kubectl validation workflow;
- Docker Compose profiles для production-like local stack;
- env-driven secrets/configuration.

## 13. Важные caveats и что еще не является настоящим production

Проект стал production-like, но не production-ready в строгом смысле.

Остается важным:

- Go-сервисы в K8s manifests сейчас используют образы вида `myregistry/<service>:latest`; для реального запуска в Kubernetes нужно добавить Dockerfiles/build pipeline, загрузку images в kind или push в registry.
- Vault пока dev-mode/local setup, не полноценная production Vault схема с policies, auth methods и secret rotation.
- Blockchain/KYC/Gas integrations pluggable, но реальные провайдеры требуют настоящих credentials, SLA/error handling и provider-specific контракты.
- RabbitMQ optional mode есть, но нет полноценной consumer pipeline/worker topology.
- Webhook delivery хранится в PostgreSQL и вызывается в log/http modes, но production-grade retry/backoff/dead-letter policy можно развивать дальше.
- Observability покрывает базовые HTTP metrics; tracing/log correlation/alerting rules можно добавить позже.
- K8s manifests валидируются server-side dry-run, но end-to-end Kubernetes runtime зависит от сборки service images.
- API auth/roles есть на уровне middleware/shared, но для полного production нужно довести role enforcement по всем sensitive endpoints.

## 14. Коммиты, сделанные в рамках работ

Ключевая последовательность коммитов:

- `f2ac388 docs(specs): add coursework implementation plan`
- `30eb9e4 docs(agents): update repository guidance`
- `a0e8421 feat(database): add coursework sql pack`
- `63a2246 feat(shared): add env-driven infrastructure helpers`
- `5a9eb2a feat(user-service): add postgres-backed profile and kyc API`
- `d1e5b0c feat(wallet-service): add postgres-backed wallet API`
- `4d2402b feat(transaction-core): add postgres-backed lifecycle API`
- `7feaef4 feat(transaction-validation): expose API-prefixed validation`
- `4236fd3 feat(gas-info): add gas history API`
- `bf55d44 feat(merchant-service): add merchant flow MVP`
- `366d0b5 feat(admin-service): add coursework admin API`
- `3f8b5be feat(web): add admin demo UI`
- `28f9c5b docs(api): align contracts with MVP endpoints`
- `9bea6fb chore(devops): add local MVP workflow`
- `07746d6 chore(gitignore): keep specs tracked`
- `efbd369 feat(merchant): persist merchant flows in postgres`
- `5a5b7dc feat(validation): check transactions against postgres`
- `1a6433d feat(gas): add provider-backed gas history`
- `60af1c7 feat(kyc): support pluggable verification providers`
- `a86cfdf fix(transaction-core): resolve postgres transaction references`
- `a3126b6 chore(devops): add full local readiness targets`
- `78e9009 test(web): add demo smoke coverage`
- `52dab28 docs: refresh coursework readiness guide`
- `9de6667 feat(database): add merchant webhook observability`
- `f61b6ab feat(merchant): manage webhook endpoints`
- `273d04e feat(transaction): broadcast and notify payments`
- `4fdbf75 feat(wallet): sync balances from rpc provider`
- `95e2293 feat(observability): expose service metrics`
- `fd62d21 chore(infra): add production-like stack configs`
- `c65fe1b docs: document production readiness paths`
- `c97afd9 chore(k8s): add kind validation workflow`
- `7cc3198 fix(web): route admin pages by URL`
- `db7929c fix(transaction): enforce lifecycle in test repository`
- `11c1196 fix(webhooks): log demo payment deliveries`
- `0d5eb3b chore(devops): add local data reset target`

## 15. Текущий status на момент актуализации

На момент актуализации были обнаружены локальные изменения, не относящиеся к стабилизационным коммитам:

- `.idea/workspace.xml`
- `.idea/go.imports.xml`
- `.idea/inspectionProfiles/`
- `apps/web/package-lock.json`

Они не включались в code/doc коммиты стабилизации, если не были нужны для выполняемой задачи.
