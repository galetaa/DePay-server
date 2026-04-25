# Screenshot Checklist

Capture these screenshots for the coursework report.

## Database

- DBeaver ER diagram for the full `mydatabase` schema.
- DBeaver table list showing core tables.
- DBeaver enum/types list.
- DBeaver roles/grants view or SQL output for roles.

## SQL

- `make sql-test` terminal output.
- Result of `get_store_turnover`.
- Result of `get_user_kyc_wallet_summary`.
- Result of `get_failed_transactions_analytics`.
- Trigger test output showing successful checks.

## Web UI

- `/admin/tables` with `users` or `payment_transactions`.
- `/admin/functions` with a completed function run.
- `/admin/analytics` showing charts.
- `/admin/demo` before running the flow.
- `/admin/demo` after payment status becomes `confirmed`.
- `/admin/tables` showing `merchant_webhook_deliveries` or `vw_webhook_delivery_status`.

## Backend

- Health checks for backend services.
- `make test-go` successful output.
- `make web-test` successful output.
- `make web-build` successful output.

## Optional Pet Project Screens

- Prometheus targets or metrics endpoint.
- Grafana overview dashboard.
- Kong route config.
- `make k8s-validate` output.

