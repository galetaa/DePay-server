# DePay Architecture

This document describes the current architecture that the portfolio roadmap builds on. The repository is intentionally kept as a top-level multi-module Go monorepo rather than being moved into a `services/` directory.

## Runtime View

```mermaid
flowchart LR
  Browser["Browser"] --> Web["React/Vite Web :5173"]
  Web --> Admin["admin-service :8090"]
  Web --> Gateway["Kong :8000 (optional)"]

  Gateway --> User["user-service :8080"]
  Gateway --> Merchant["merchant-service :8083"]
  Gateway --> Wallet["wallet-service :8084"]
  Gateway --> TxCore["transaction-core-service :8085"]
  Gateway --> TxValidation["transaction-validation-service :8081"]
  Gateway --> Gas["gas-info-service :8082"]
  Gateway --> KYC["kyc-service :8086"]

  User --> PG[("PostgreSQL")]
  Merchant --> PG
  Wallet --> PG
  TxCore --> PG
  TxValidation --> PG
  Admin --> PG

  Wallet --> Redis[("Redis")]
  Gas --> Redis
  TxCore --> Rabbit[("RabbitMQ optional")]

  TxCore --> Chain["EVM JSON-RPC or mock"]
  Wallet --> BalanceRPC["Balance RPC or mock"]
  Gas --> GasProvider["HTTP gas provider or mock"]
  KYC --> KYCProvider["HTTP KYC provider or mock"]
  TxCore --> Webhooks["Merchant webhook endpoints"]

  Prom["Prometheus"] --> User
  Prom --> Merchant
  Prom --> Wallet
  Prom --> TxCore
  Prom --> Admin
  Grafana["Grafana"] --> Prom
```

## Service Shape

Most Go services follow the same shape:

```mermaid
flowchart LR
  Main["cmd/<service>/main.go"] --> Controller["controller"]
  Controller --> Service["service"]
  Service --> Repository["repository"]
  Repository --> PG[("PostgreSQL")]
  Main --> Shared["shared middleware/config/logging"]
```

Where a repository exists, PostgreSQL is the production-like source of truth. Some services keep in-memory repositories as local or test fallbacks when `DATABASE_URL` is not available.

## Payment Lifecycle

The current persisted transaction flow is:

```mermaid
stateDiagram-v2
  [*] --> created
  created --> submitted
  submitted --> validated
  validated --> broadcasted
  broadcasted --> confirmed
  created --> cancelled
  submitted --> cancelled
  created --> failed
  submitted --> failed
  validated --> failed
  broadcasted --> failed
  confirmed --> [*]
  cancelled --> [*]
  failed --> [*]
```

The PostgreSQL trigger is a safety guard for persisted status transitions. The next roadmap phase moves more lifecycle policy into a centralized transaction-core state machine with idempotency and event creation.

## Data Ownership

PostgreSQL owns durable state:

- users, merchants, verification and KYC;
- wallets, assets and balances;
- invoices, payment sessions and transactions;
- webhook registrations and delivery attempts;
- RPC nodes, audit logs and risk alerts;
- SQL views and reporting functions.

Redis is used as cache/history support for gas and wallet paths. RabbitMQ is optional and should remain disabled in focused tests through `SKIP_RABBITMQ=true`.

## Provider Modes

Local development defaults to safe mock/dev behavior:

| Area | Default | Real/provider switch |
| --- | --- | --- |
| Blockchain broadcast | Mock transaction hash | `BLOCKCHAIN_RPC_URL` |
| Wallet balance sync | Deterministic mock balance | `WALLET_BALANCE_RPC_URL` or `BLOCKCHAIN_RPC_URL` |
| Gas info | Mock gas data | `GAS_PROVIDER_URL` |
| KYC | Fast mock provider | `KYC_PROVIDER_URL`, `KYC_PROVIDER_API_KEY` |
| Webhook delivery | PostgreSQL log mode | `WEBHOOK_DELIVERY_MODE=http` |

No user private keys are stored by the demo.

## Local Profiles

Docker Compose profiles group optional layers:

| Command | Purpose |
| --- | --- |
| `make up` | PostgreSQL, Redis and RabbitMQ |
| `make backend-up` | Go backend services |
| `make web-up` | Backend plus React/Vite web |
| `make gateway-up` | Kong gateway |
| `make observability-up` | Prometheus and Grafana |
| `make secrets-up` | Vault dev mode |
| `make prod-like-up` | Backend, web, gateway, observability and secrets |

The baseline green-state is documented in the root README and in `depay_next_dev_pack/README.md`.
