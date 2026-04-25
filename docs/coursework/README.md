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

5. Start admin API and web:

```bash
docker compose --profile backend up admin-service
cd apps/web && npm install && npm run dev
```

## Defense Screens

- DBeaver schema from the migrated PostgreSQL database.
- `users`, `stores`, `wallets`, `balances`, `payment_transactions`.
- `/admin/functions` with `get_store_turnover`.
- Trigger errors from `database/tests/test_triggers.sql`.
- Analytics charts in `/admin/analytics`.
- Demo payment in `/admin/demo`.
