# DePay Server

DePay is a coursework-ready cryptocurrency payment system prototype. It models users, merchants, wallets, invoices, NFC-style payment sessions, transactions, KYC, risk monitoring, webhook delivery logs and analytics.

The repository is a multi-module Go monorepo. Each backend service has its own `go.mod`; shared cross-cutting packages live in `shared/`.

## Current Scope

The stable demo path is Docker Compose + PostgreSQL + Go services + React web UI.

Implemented for coursework:

- PostgreSQL schema with normalized tables, indexes, enum types and grants.
- Seed data for users, merchants, wallets, balances, invoices, transactions, RPC checks, risk alerts and webhooks.
- 12 SQL functions with joins, grouping and aggregates.
- 7 PostgreSQL triggers covering ownership checks, negative balances, balance history, transaction lifecycle, audit and risk alerts.
- Admin API for table browsing, SQL function execution, analytics and demo flow.
- React/Vite web UI for `/admin/tables`, `/admin/functions`, `/admin/analytics`, `/admin/demo`.
- At least 3 charts in the analytics page.
- Go, SQL and web tests.

Production-like but not the main coursework demo:

- Kong config.
- Vault dev-mode config.
- Prometheus/Grafana config.
- Kubernetes manifests and kind validation.
- Mock/dev providers for blockchain, KYC, gas and wallet balances, with env-based HTTP/RPC provider switches.

## Services

- `user-service` on `8080`
- `transaction-validation-service` on `8081`
- `gas-info-service` on `8082`
- `merchant-service` on `8083`
- `wallet-service` on `8084`
- `transaction-core-service` on `8085`
- `kyc-service` on `8086`
- `admin-service` on `8090`
- `apps/web` on `5173`

## Quick Start

For the usual local coursework run:

```bash
make dev-ready
make web-up
```

Open:

```text
http://localhost:5173/admin/tables
http://localhost:5173/admin/functions
http://localhost:5173/admin/analytics
http://localhost:5173/admin/demo
```

For a true project-local reset that removes PostgreSQL data and Compose volumes:

```bash
make reset-local-data
make dev-ready
make web-up
```

`docker volume prune` may not remove Compose named volumes such as `depay-server_pgdata`; use `make reset-local-data` when you need a real empty local database.

## Verification

Run:

```bash
make sql-test
make test-go
make web-test
make web-build
```

Optional Kubernetes manifest validation:

```bash
make kind-up
make k8s-validate
```

This validates manifests with `kubectl apply --dry-run=server`. It does not prove runtime Kubernetes deployment, because service images are still referenced as `myregistry/<service>:latest` and need a real image build/push/load workflow.

## Demo Flow

1. Start the local stack with `make dev-ready && make web-up`.
2. Open `http://localhost:5173/admin/tables` and show seeded database tables/views.
3. Open `/admin/functions` and run SQL functions such as `get_store_turnover`.
4. Open `/admin/analytics` and show charts for turnover, statuses, failed transactions and RPC latency.
5. Open `/admin/demo`.
6. Click `Create` to create a demo invoice.
7. Click `Submit` to create and progress a demo payment to `confirmed`.
8. Check `payment_transactions`, `payment_invoices`, `audit_logs` and `merchant_webhook_deliveries`.

## Provider Modes

By default the project uses mock/dev providers:

- Gas provider: mock unless `GAS_PROVIDER_URL` is set.
- KYC provider: mock unless `KYC_PROVIDER_URL` is set.
- Blockchain broadcaster: mock unless `BLOCKCHAIN_RPC_URL` is set.
- Wallet balance provider: mock unless `WALLET_BALANCE_RPC_URL` or `BLOCKCHAIN_RPC_URL` is set.
- Webhook delivery: PostgreSQL log mode unless `WEBHOOK_DELIVERY_MODE=http` is set.

No user private keys are stored by the demo.

## Documentation

- Coursework guide: `docs/coursework/README.md`
- API summary: `docs/api.md`
- OpenAPI contract: `docs/openAPI.yaml`
- Defense script: `docs/defense-demo-script.md`
- Screenshot checklist: `docs/screenshots.md`
- Implementation report: `IMPLEMENTATION_REPORT.md`
- Latest verification notes: `verification-results.md`

