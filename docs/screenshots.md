# Screenshot Checklist

Capture these screenshots for the coursework report and portfolio README. The coursework set should stay focused on PostgreSQL, SQL functions, triggers and the web demo; the portfolio set can also show production-like extras.

## Portfolio Overview

- Architecture diagram from `README.md` or `docs/architecture.md`.
- DBeaver ER diagram for the full `mydatabase` schema.
- Admin tables view showing `payment_transactions`, `payment_invoices` or `merchant_webhook_deliveries`.
- SQL function runner with a successful function call.
- Analytics dashboard with charts visible.
- Login/persona selector at `/login`.
- User dashboard at `/user/dashboard`.
- Merchant dashboard at `/merchant/dashboard`.
- Merchant invoices at `/merchant/invoices`.
- Merchant webhooks and delivery attempts at `/merchant/webhooks`.
- Merchant API key create/list response via API client or terminal.
- Merchant terminals at `/merchant/terminals`.
- Merchant analytics at `/merchant/analytics`.
- Compliance queue at `/compliance/kyc`.
- Admin system health at `/admin/system-health`.
- Demo payment flow before creating an invoice.
- Demo payment flow after payment status becomes `confirmed`.
- Webhook delivery table or `vw_webhook_delivery_status`.
- Optional: Prometheus targets or Grafana dashboard.
- Optional: `admin-service` `/metrics` showing business metrics.

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
- Docker image mode with `make image-up`.
- Kong route config.
- `make k8s-validate` output.
