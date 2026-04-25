DATABASE_URL ?= postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable
SERVICES := user-service wallet-service transaction-core-service transaction-validation-service gas-info-service kyc-service merchant-service admin-service

.PHONY: up backend-up web-up full-up dev-ready down wait-db migrate-up migrate-down seed sql-test test test-go web-dev web-build web-test

up:
	docker compose up -d postgres redis rabbitmq

backend-up:
	docker compose --profile backend up -d

web-up:
	docker compose --profile backend --profile web up -d web

full-up:
	docker compose --profile backend --profile web up -d

dev-ready: up wait-db migrate-up seed backend-up

down:
	docker compose --profile backend --profile web down

wait-db:
	@until pg_isready -d "$(DATABASE_URL)" >/dev/null 2>&1; do \
		echo "waiting for postgres"; \
		sleep 1; \
	done

migrate-up:
	migrate -path database/migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path database/migrations -database "$(DATABASE_URL)" down

seed:
	psql "$(DATABASE_URL)" -f database/seeds/seed_coursework.sql

sql-test:
	psql "$(DATABASE_URL)" -f database/tests/test_functions.sql
	psql "$(DATABASE_URL)" -f database/tests/test_triggers.sql

test: test-go

test-go:
	@for service in $(SERVICES); do \
		if [ -f "$$service/go.mod" ]; then \
			echo "==> $$service"; \
			(cd "$$service" && go test ./...); \
		fi; \
	done

web-dev:
	cd apps/web && npm run dev -- --host 0.0.0.0

web-build:
	cd apps/web && npm run build

web-test:
	cd apps/web && npm test
