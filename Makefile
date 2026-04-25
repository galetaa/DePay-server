DATABASE_URL ?= postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable
KIND ?= $(shell command -v kind 2>/dev/null || printf "%s/bin/kind" "$$(go env GOPATH)")
KIND_CLUSTER ?= depay-kubectl-check
SERVICES := user-service wallet-service transaction-core-service transaction-validation-service gas-info-service kyc-service merchant-service admin-service

.PHONY: up backend-up web-up gateway-up observability-up secrets-up full-up prod-like-up dev-ready down reset-local-data wait-db migrate-up migrate-down seed sql-test test test-go web-dev web-build web-test kind-install kind-up kind-down k8s-validate k8s-dry-run

up:
	docker compose up -d postgres redis rabbitmq

backend-up:
	docker compose --profile backend up -d

gateway-up:
	docker compose --profile gateway up -d kong

observability-up:
	docker compose --profile backend --profile observability up -d prometheus grafana

secrets-up:
	docker compose --profile secrets up -d vault

web-up:
	docker compose --profile backend --profile web up -d web

full-up:
	docker compose --profile backend --profile web up -d

prod-like-up:
	docker compose --profile backend --profile web --profile gateway --profile observability --profile secrets up -d

dev-ready: up wait-db migrate-up seed backend-up

down:
	docker compose --profile backend --profile web --profile gateway --profile observability --profile secrets down

reset-local-data:
	docker compose --profile backend --profile web --profile gateway --profile observability --profile secrets down -v --remove-orphans

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
	psql "$(DATABASE_URL)" -f database/tests/test_webhooks.sql

test: test-go

test-go:
	@set -e; \
	for service in $(SERVICES); do \
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

kind-install:
	@if ! command -v kind >/dev/null 2>&1 && [ ! -x "$$(go env GOPATH)/bin/kind" ]; then \
		go install sigs.k8s.io/kind@latest; \
	fi

kind-up: kind-install
	@"$(KIND)" get clusters | grep -qx "$(KIND_CLUSTER)" || "$(KIND)" create cluster --name "$(KIND_CLUSTER)"
	kubectl config use-context "kind-$(KIND_CLUSTER)"

kind-down:
	@"$(KIND)" delete cluster --name "$(KIND_CLUSTER)"

k8s-validate:
	kubectl apply --dry-run=server --validate=strict -f k8s -o name

k8s-dry-run: k8s-validate
