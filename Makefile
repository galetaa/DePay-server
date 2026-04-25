DATABASE_URL ?= postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable
SERVICES := user-service wallet-service transaction-core-service transaction-validation-service gas-info-service kyc-service merchant-service admin-service

.PHONY: up down migrate-up migrate-down seed sql-test test test-go web-dev

up:
	docker compose up -d postgres redis rabbitmq

down:
	docker compose down

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
