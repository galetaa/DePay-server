# 03. Repository Structure

## Целевая структура

```text
DePay/
  README.md
  Makefile
  docker-compose.yml
  .env.example

  docs/
    architecture/
    api/
    coursework/
    decisions/

  database/
    migrations/
    seeds/
    scripts/

  services/
    user-service/
    merchant-service/
    wallet-service/
    transaction-core-service/
    transaction-validation-service/
    gas-info-service/
    analytics-service/
    admin-service/

  apps/
    web/

  shared/
    go/
      config/
      db/
      logging/
      middleware/
      auth/
      errors/
      events/

  deployments/
    docker/
    k8s/
    kong/
    vault/
```

## Мягкая миграция без поломки

Сначала можно оставить текущие сервисы в корне и просто добавить:

```text
database/
web/
docs/coursework/
docker-compose.yml
Makefile
.env.example
```

Потом постепенно перенести сервисы в `services/`.

## Database

```text
database/
  migrations/
    000001_create_extensions.up.sql
    000002_create_types.up.sql
    000003_create_tables.up.sql
    000004_create_indexes.up.sql
    000005_create_triggers.up.sql
    000006_create_functions.up.sql
    000007_create_roles.up.sql
  seeds/
    seed_dev.sql
    seed_coursework.sql
```

## Web

```text
apps/web или web/
  package.json
  src/
    api/
    components/
    pages/
    routes/
    charts/
```

## Makefile минимум

```makefile
up:
	docker compose up -d

down:
	docker compose down

migrate-up:
	migrate -path database/migrations -database "$$DATABASE_URL" up

seed:
	psql "$$DATABASE_URL" -f database/seeds/seed_coursework.sql

test:
	go test ./...
```
