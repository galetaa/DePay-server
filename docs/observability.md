# Observability

All Go services expose:

- `GET /health`
- `GET /metrics`

## Request IDs

`RequestIDMiddleware` reads `X-Request-ID` or generates one, stores it in the Gin context, and returns it in the response header. HTTP request logs include:

- `service`
- `request_id`
- `method`
- `route`
- `status`
- `duration`

## HTTP Metrics

Every service records:

- `depay_http_requests_total{service,method,path,status}`
- `depay_http_request_duration_seconds_sum{service,method,path,status}`
- `depay_process_uptime_seconds`

## Business Metrics

`admin-service` enriches `/metrics` with PostgreSQL-backed business gauges:

- `depay_transactions_total{status}`
- `depay_invoices_total{status}`
- `depay_webhook_deliveries_total{status}`
- `depay_rpc_node_latency_ms{node,chain}`
- `depay_risk_alerts_total{level,status}`
- `depay_kyc_applications_total{status}`

## Dashboards

Prometheus scrape config lives in `observability/prometheus.yml`. The starter Grafana dashboard is `observability/grafana/dashboards/depay-overview.json`.
