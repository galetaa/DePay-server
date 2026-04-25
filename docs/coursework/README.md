# DePay Coursework Pack

## Run Order

1. Start infrastructure:

```bash
make up
```

2. Apply migrations:

```bash
make migrate-up
```

3. Load coursework seed data:

```bash
make seed
```

4. Run SQL checks:

```bash
make sql-test
```

5. Or prepare seeded backend in one command:

```bash
make dev-ready
```

6. Start all backend services:

```bash
make backend-up
```

7. Start backend plus web:

```bash
make full-up
```

8. Start production-like local stack:

```bash
make prod-like-up
```

9. Run app checks:

```bash
make test-go
make web-test
make web-build
```

## Defense Screens

- DBeaver schema from the migrated PostgreSQL database.
- `users`, `stores`, `wallets`, `balances`, `payment_transactions`.
- `/admin/functions` with `get_store_turnover`.
- Trigger errors from `database/tests/test_triggers.sql`.
- Analytics charts in `/admin/analytics`.
- Demo payment in `/admin/demo`.

## Provider Modes

- Gas data uses a mock provider by default and can call `GAS_PROVIDER_URL`.
- KYC uses a mock provider by default and can call `KYC_PROVIDER_URL` with `KYC_PROVIDER_API_KEY`.
- Transaction broadcasting uses mock hashes by default and can call EVM JSON-RPC via `BLOCKCHAIN_RPC_URL`.
- Wallet sync uses deterministic mock balances by default and can call EVM JSON-RPC via `WALLET_BALANCE_RPC_URL` or `BLOCKCHAIN_RPC_URL`.
- Merchant webhooks are registered under `/api/merchant/webhooks`; delivery is logged by default and switches to HTTP with `WEBHOOK_DELIVERY_MODE=http`.
- Each backend exposes `/metrics`; Prometheus and Grafana run with `make observability-up`.
- Kong runs with `make gateway-up`; Vault dev bootstrap runs with `make secrets-up`.
- Redis/RabbitMQ remain optional supporting systems; PostgreSQL is the source of truth.
