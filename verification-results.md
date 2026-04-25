# Verification Results

Date: 2026-04-26

This file records the latest stabilization checks for the coursework-ready DePay setup.

## Environment Notes

- Docker Desktop was initially stopped. It was started before running the clean setup.
- `docker volume prune -f` did not remove Compose named volumes such as `depay-server_pgdata` on this machine.
- `make reset-local-data` was added for a project-scoped clean reset with `docker compose down -v`.
- After removing Go/NPM cache volumes, backend and web containers may need several minutes for cold `go run`/`npm install` startup before `/health` and Vite respond.
- Kubernetes validation uses a local kind cluster only for manifest validation; it is not part of the main coursework demo.

## Final Clean Reset Check

| Check | Result | Notes |
| --- | --- | --- |
| Docker Desktop | Passed after manual start | Initial `make down` failed while Docker daemon was stopped. |
| `make down && docker volume prune -f && make dev-ready` | Passed | Confirmed ordinary prune does not remove Compose named `pgdata`. |
| `make reset-local-data` | Passed | Removed project containers and named volumes, including `depay-server_pgdata`. |
| `make dev-ready` after reset | Passed | Migrations ran from `1/u` through `8/u`, seed loaded, backend started. |
| `make web-up` after reset | Passed | Web container started; cold `npm install` took extra time after volume reset. |
| `make sql-test` | Passed | 12 SQL functions returned rows; trigger tests passed; webhook view test inserts and reads a delivery. |
| `make test-go` | Passed | All service modules passed. |
| `make web-test` | Passed | 3 files, 6 tests. |
| `make web-build` | Passed | Vite build succeeded; chunk-size warning only. |
| Backend health checks | Passed | Ports 8080, 8081, 8082, 8083, 8084, 8085, 8086, 8090 returned `{"status":"ok"}` after cold start. |
| Web direct routes | Passed | `/admin/tables`, `/admin/functions`, `/admin/analytics`, `/admin/demo` render the correct page. |
| UI demo-flow | Passed | Create invoice -> submit payment -> status `confirmed`. |
| Webhook delivery log | Passed | Demo confirmed payment inserted `merchant_webhook_deliveries` row with `delivered` and response status `204`. |
| `docs/openAPI.yaml` YAML parse | Passed | Ruby YAML parser loaded the file successfully. |

## Final Command Sequence

```bash
make reset-local-data
make dev-ready
make web-up
make sql-test
make test-go
make web-test
make web-build
```

Additional checks:

```bash
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health
curl http://localhost:8084/health
curl http://localhost:8085/health
curl http://localhost:8086/health
curl http://localhost:8090/health
```

Browser verification was performed for:

- `http://localhost:5173/admin/tables`
- `http://localhost:5173/admin/functions`
- `http://localhost:5173/admin/analytics`
- `http://localhost:5173/admin/demo`
