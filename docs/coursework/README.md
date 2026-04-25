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

8. Run app checks:

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
- Redis/RabbitMQ remain optional supporting systems; PostgreSQL is the source of truth.
