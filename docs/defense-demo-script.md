# Defense Demo Script

Use this script for a 7-10 minute coursework defense.

## 1. Project Overview

DePay is a crypto payment system prototype for users, merchants, wallets, invoices, NFC-style terminals, transactions, KYC, risk monitoring and analytics.

Stack:

- Go microservices;
- PostgreSQL as source of truth;
- Redis and RabbitMQ as supporting systems;
- React/Vite/TypeScript admin web UI;
- Docker Compose for local reproducible launch.

## 2. Database

Open DBeaver and show the schema.

Mention:

- normalized entities: users, roles, stores, wallets, balances, invoices, terminals, transactions, KYC, risk alerts, audit logs, webhooks;
- foreign keys between users/stores/wallets/assets/invoices/transactions;
- enum types for statuses;
- indexes for common lookup fields;
- PostgreSQL roles and grants.

## 3. SQL Functions

Open `/admin/functions`.

Run:

- `get_user_kyc_wallet_summary`
- `get_store_turnover`
- `get_failed_transactions_analytics`

Explain that these functions use joins, grouping and aggregates over seed data.

## 4. Triggers

Run or show output from:

```bash
make sql-test
```

Explain trigger coverage:

- wallet owner validation;
- negative balance blocking;
- balance history;
- transaction lifecycle;
- audit log;
- risk alert for large unverified payment.

## 5. Admin Web UI

Open:

- `/admin/tables`
- `/admin/functions`
- `/admin/analytics`
- `/admin/demo`

Show that direct URLs work and each page renders the expected screen.

## 6. Analytics

Open `/admin/analytics`.

Show at least three charts:

- store turnover;
- transaction statuses;
- failed transactions;
- RPC latency.

## 7. Demo Payment

Open `/admin/demo`.

1. Click `Create`.
2. Show created invoice id.
3. Click `Submit`.
4. Show transaction id and `confirmed` status.
5. Open `/admin/tables`.
6. Inspect `payment_transactions`, `audit_logs`, `merchant_webhook_deliveries` or `vw_webhook_delivery_status`.

## 8. Testing

Show:

```bash
make sql-test
make test-go
make web-test
make web-build
```

Mention that provider integrations are mock/dev by default and can be switched through env variables.

## 9. Honest Boundaries

Say explicitly:

- The coursework demo does not require real blockchain, real KYC or real secrets.
- Kubernetes, Vault and observability configs are production-like extensions, not the main defense path.
- Real production would still need image build/push pipeline, stronger auth enforcement for every admin endpoint, secret rotation, retries/dead-letter queues and real provider contracts.

