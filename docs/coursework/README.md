# DePay Coursework Guide

This is the stable coursework path. Keep the defense focused on PostgreSQL, SQL logic, triggers, roles, table/function UI, charts and demo-flow. Kong, Vault, Prometheus/Grafana and Kubernetes are optional pet-project/production-like extras, not required for the main coursework demo.

## Clean Local Setup

For a normal run on an already prepared local database:

```bash
make dev-ready
make web-up
```

For a real project-local clean setup that removes the Compose PostgreSQL volume:

```bash
make reset-local-data
make dev-ready
make web-up
```

`docker volume prune` may leave Compose named volumes such as `depay-server_pgdata`, so `make reset-local-data` is the preferred reproducible reset command for this project.

## Required Checks

Run these before the defense:

```bash
make sql-test
make test-go
make web-test
make web-build
```

Expected result:

- `make sql-test` runs all SQL function, trigger and webhook view checks.
- `make test-go` runs Go tests per service module.
- `make web-test` runs Vitest tests for API client, routing and demo flow.
- `make web-build` produces a Vite production build. A large chunk warning is acceptable for the current MVP.

## What To Show In DBeaver

Connect to:

```text
postgres://myuser:mypassword@localhost:5432/mydatabase
```

Show:

- tables: `users`, `stores`, `wallets`, `balances`, `payment_invoices`, `payment_transactions`, `merchant_webhooks`, `merchant_webhook_deliveries`, `audit_logs`, `risk_alerts`;
- enum types for KYC, verification, transactions, invoices, terminals, RPC status and risk;
- foreign keys between users, stores, wallets, assets, invoices and transactions;
- views: `vw_user_wallet_balances`, `vw_store_transactions`, `vw_failed_transactions`, `vw_rpc_node_status`, `vw_compliance_kyc_queue`, `vw_webhook_delivery_status`;
- roles: `depay_admin`, `depay_app`, `depay_compliance`, `depay_merchant_readonly`, `depay_user_readonly`.

## SQL Functions To Demonstrate

Use `/admin/functions` or DBeaver:

- `get_user_kyc_wallet_summary()`
- `get_user_wallet_balances(1)`
- `get_wallet_asset_distribution(1)`
- `get_transaction_card(1)`
- `get_user_transaction_history(1, '2020-01-01', '2100-01-01')`
- `get_store_transaction_history(1, '2020-01-01', '2100-01-01')`
- `get_blockchain_asset_activity(1, '2020-01-01', '2100-01-01')`
- `get_rpc_nodes_activity(1, '2020-01-01', '2100-01-01')`
- `get_store_turnover(1, '2020-01-01', '2100-01-01')`
- `get_store_success_rate(1, '2020-01-01', '2100-01-01')`
- `get_unverified_active_users(1, 1, '2020-01-01', '2100-01-01')`
- `get_failed_transactions_analytics('2020-01-01', '2100-01-01')`

## Triggers To Demonstrate

`database/tests/test_triggers.sql` exercises:

- wallet owner check;
- negative balance rejection;
- balance history write;
- transaction lifecycle validation;
- transaction audit log.

The migration also includes risk alert creation for large unverified payments and timestamp updates for transaction status changes.

## Web Screens

Open:

```text
http://localhost:5173/admin/tables
http://localhost:5173/admin/functions
http://localhost:5173/admin/analytics
http://localhost:5173/admin/demo
```

Show:

- `/admin/tables` table/view browser;
- `/admin/functions` SQL function runner;
- `/admin/analytics` charts for store turnover, transaction statuses, failed transactions and RPC latency;
- `/admin/demo` invoice + payment flow.

## Demo Flow

1. Open `/admin/demo`.
2. Click `Create`; an invoice id should appear.
3. Click `Submit`; a transaction id should appear.
4. Status should become `confirmed`.
5. Open `/admin/tables` and inspect:
   - `payment_invoices`
   - `payment_transactions`
   - `audit_logs`
   - `merchant_webhook_deliveries`
   - `vw_webhook_delivery_status`

The demo uses PostgreSQL and mock/dev providers. It does not store user private keys.

## Optional Pet Project Checks

These are useful but should not be the centerpiece of the coursework defense:

```bash
make gateway-up
make observability-up
make secrets-up
make kind-up
make k8s-validate
```

Notes:

- Prometheus/Grafana are configured for local observability.
- Vault runs in dev mode.
- Kubernetes manifests validate with `kubectl` server-side dry-run, but runtime deployment still needs real service images.
