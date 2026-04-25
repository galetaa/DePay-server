# AGENTS.md

## Repository Reality Check (read first)

- This is a multi-module Go monorepo. Each service has its own `go.mod` and uses `replace shared => ../shared`.
- The repository root is not a Go module. Do not run `go test ./...` from the root.
- `specs/*.md` describe the target architecture and coursework scope. The current code is a coursework-ready MVP plus production-like local/dev paths for providers, webhooks, observability, Kong, Vault, and Kubernetes manifests.
- Keep the existing top-level service directories. Do not move services into `services/` unless a task explicitly asks for that refactor.

## Current Modules and Apps

- Go modules:
  - `shared`
  - `user-service`
  - `wallet-service`
  - `transaction-core-service`
  - `transaction-validation-service`
  - `gas-info-service`
  - `kyc-service`
  - `merchant-service`
  - `admin-service`
- Frontend app:
  - `apps/web` is a React + Vite + TypeScript app for admin tables, SQL function execution, analytics, and demo flow.
- Database assets:
  - `database/migrations` contains `golang-migrate` migrations.
  - `database/seeds/seed_coursework.sql` contains demo data for coursework defense.
  - `database/tests/test_functions.sql` and `database/tests/test_triggers.sql` are SQL smoke/acceptance tests.

## Big Picture Architecture

- Service shape is consistent: `cmd/<service>/main.go` wires `controller -> service -> repository` where a repository exists.
- Shared cross-cutting code lives in `shared/`:
  - `config`, `logging`, `middleware`, `utils`
  - `auth`, `db`, `errors`, `events`, `observability`, `validation`
- PostgreSQL is the intended source of truth for the coursework MVP.
- `user-service`, `wallet-service`, `transaction-core-service`, and `merchant-service` use PostgreSQL repositories when `DATABASE_URL` is set and reachable, with in-memory repositories as local/test fallback.
- `wallet-service` can sync balances from EVM JSON-RPC through `WALLET_BALANCE_RPC_URL` or `BLOCKCHAIN_RPC_URL`; without those env vars it uses deterministic mock balances.
- `admin-service` is PostgreSQL-backed and is the main API for the coursework web/admin UI.
- `merchant-service` implements merchant registration/login, verification, invoices, terminals, and webhook registration on top of PostgreSQL.
- `transaction-core-service` supports validated -> broadcasted -> confirmed lifecycle endpoints, EVM `eth_sendRawTransaction` via `BLOCKCHAIN_RPC_URL`, mock broadcast fallback, and merchant webhook delivery logging/HTTP delivery via `WEBHOOK_DELIVERY_MODE`.
- `transaction-validation-service` uses stateless validation by default and enables PostgreSQL-backed transaction/KYC/ownership/balance/merchant/blacklist/risk checks when `DATABASE_URL` is set.
- `gas-info-service` uses Redis as optional cache/history storage and supports either mock gas data or an HTTP gas provider via `GAS_PROVIDER_URL`.
- `kyc-service` uses a fast mock provider by default and can call an HTTP KYC provider via `KYC_PROVIDER_URL` and `KYC_PROVIDER_API_KEY`.
- Every Go service exposes `/health` and Prometheus-compatible `/metrics`.

## Database and SQL Pack

- Migrations are ordered:
  - extensions
  - enum types
  - tables
  - indexes
  - triggers
  - functions/views
  - roles/grants
- Seed data is expected to make all coursework SQL functions return rows.
- SQL tests should verify both function output and trigger behavior.
- When changing schema, update migrations, seed data, SQL tests, affected repositories, and `docs/openAPI.yaml` where API behavior changes.

## External Systems and Environment

- Local defaults live in `.env.example`.
- Important environment variables:
  - `DATABASE_URL`
  - `JWT_SECRET`
  - `REDIS_ADDR`
  - `REDIS_PASSWORD`
  - `RABBITMQ_URL`
  - `SKIP_RABBITMQ`
  - `BLOCKCHAIN_RPC_URL`
  - `WALLET_BALANCE_RPC_URL`
  - `WEBHOOK_DELIVERY_MODE`
  - `GAS_PROVIDER_URL`
  - `KYC_PROVIDER_URL`
  - `KYC_PROVIDER_API_KEY`
  - `PORT`
- JWT signing uses `JWT_SECRET` through shared helpers. Do not reintroduce hardcoded secrets.
- RabbitMQ publishing in `transaction-core-service` is optional and should remain disabled in tests via `SKIP_RABBITMQ=true`.
- Redis and RabbitMQ are supporting systems, not the source of truth.
- `docker-compose.yml` provides local infrastructure and `backend`, `web`, `gateway`, `observability`, and `secrets` profiles.

## Developer Workflows

- Start local infrastructure:

```bash
make up
```

- Run migrations and seed:

```bash
make migrate-up
make seed
```

- Prepare a seeded local backend in one command:

```bash
make dev-ready
```

- Start all backend services:

```bash
make backend-up
```

- Start backend plus web:

```bash
make full-up
```

- Start production-like local stack:

```bash
make prod-like-up
```

- Run SQL tests:

```bash
make sql-test
```

- Run Go tests across modules:

```bash
make test-go
```

- Or run Go tests per module:

```bash
cd user-service && go test ./...
cd wallet-service && go test ./...
cd transaction-core-service && go test ./...
cd transaction-validation-service && go test ./...
cd gas-info-service && go test ./...
cd kyc-service && go test ./...
cd merchant-service && go test ./...
cd admin-service && go test ./...
```

- Run the frontend:

```bash
cd apps/web && npm run dev
```

- Build the frontend:

```bash
cd apps/web && npm run build
```

- Run frontend smoke tests:

```bash
cd apps/web && npm test
```

## API and Routing Notes

- New public/backend endpoints should use the `/api` prefix.
- Existing legacy endpoints are still present for compatibility in some services. Do not treat them as the preferred contract for new work.
- Important current API groups:
  - `user-service`: `/api/user/*`
  - `wallet-service`: `/api/wallets/*`
  - `transaction-core-service`: `/api/transaction/*`
  - `transaction-validation-service`: `/api/transaction/validate`
  - `gas-info-service`: `/api/transactions/gas-info`, `/api/gas-info/history`
  - `merchant-service`: `/api/merchant/*`
  - `admin-service`: `/api/admin/*`, `/api/analytics/*`
- When adding or changing endpoints, update both the service router and `docs/openAPI.yaml` in the same change set.

## Project-Specific Conventions to Preserve

- Keep business logic in services. Controllers should bind/validate HTTP input and map responses.
- Reuse the shared middleware stack in service mains:
  - `gin.Recovery()`
  - `CORSMiddleware()`
  - `ErrorHandlerMiddleware()`
- Preserve `/health` on each service.
- Preserve `/metrics` on each service.
- Keep Redis and RabbitMQ optional for local development and tests.
- Use in-memory repositories for tests where that keeps the test focused, but do not add new production storage paths that bypass PostgreSQL.
- Prefer small, local changes that follow the existing controller/service/repository shape.
- Do not update generated or IDE files unless the task explicitly requires it.

## Production-Like Paths

- Blockchain calls are provider-driven through env (`BLOCKCHAIN_RPC_URL`, `WALLET_BALANCE_RPC_URL`) and must keep mock fallbacks for local/tests.
- Merchant webhooks are registered through `/api/merchant/webhooks`; transaction-core records delivery attempts and can deliver HTTP callbacks when `WEBHOOK_DELIVERY_MODE=http`.
- Observability is available through `/metrics`, the `observability` compose profile, and `k8s/observability.yaml`.
- Kong declarative routing lives in `kong/kong.yml` and `k8s/kong.yaml`.
- Vault bootstrap lives in `server-config/setup_vault_and_kong.sh` and `k8s/vault.yaml`.
- Keep frontend smoke coverage in Vitest; add Playwright only when the UI flow needs browser-level regression checks.
