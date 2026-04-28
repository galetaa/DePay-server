# Testing And Quality

Baseline commands:

```bash
make sql-test
make test-go
make web-test
make web-build
```

Use the full reset flow before portfolio demos or larger database changes:

```bash
make reset-local-data
make dev-ready
make web-up
```

## Coverage Map

- SQL: functions, triggers, constraints, seed sanity, webhook delivery views.
- Go: controllers, services, repositories, providers, validation, auth middleware.
- Web: API client, route smoke tests, dashboard rendering, error/loading states.
- Infra: Docker image build, compose health smoke, Kubernetes manifest validation.

## Definition Of Done

- Tests pass.
- Build passes.
- Migrations apply and down migrations exist where practical.
- Docs and OpenAPI are updated for endpoint changes.
- Secrets are not added to code paths.
- Mock and real-provider modes are documented.
- Demo flow still works after the change.
