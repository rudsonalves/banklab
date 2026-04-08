# banklab

**banklab** is a monorepo combining a Go banking API and a Flutter mobile client, built as an engineering-focused environment to explore system design decisions across backend and mobile layers, with a strong emphasis on transactional consistency, API contract design, and end-to-end correctness.

> This project prioritizes correctness and design clarity over feature breadth. It is intended as an engineering exploration, not a production-ready system.

---

## Repository Structure

```
banklab/
├── api/        # Go banking API (modular monolith)
├── mobile/     # Flutter mobile client (BankFlow)
├── infra/      # Docker and infrastructure scripts
├── docs/       # Architecture and design decisions
└── Makefile    # Monorepo task runner
```

---

## Purpose

This project is part of a broader effort to build and validate a complete system (backend + mobile client), focusing on:

- Correctness of financial operations
- Consistency guarantees under concurrent access
- Clear API contracts between backend and client
- Explicit modeling of business rules and invariants
- End-to-end validation of financial workflows

---

## API — Bank API (Go)

A modular Go service implementing a simplified banking core.

### Core Capabilities

- Customer creation
- User registration and authentication (JWT)
- Account creation and ownership enforcement
- Financial operations: deposit, withdraw, transfer
- Account statement with pagination and filtering

### Architectural Approach

Layered modular monolith:

- **Delivery** — HTTP layer
- **Application** — use cases and transaction orchestration
- **Domain** — business rules and invariants
- **Infrastructure** — PostgreSQL and repositories

Dependency direction: `Delivery → Application → Domain` / `Infrastructure → Domain`

### Design Focus

**Transactional Consistency**
- All balance-changing operations run inside explicit database transactions
- Row-level locking (`SELECT FOR UPDATE`) prevents race conditions
- Deterministic lock ordering in transfers reduces deadlock risk

**Ledger Integrity**
- All balance changes recorded as immutable ledger entries
- Operations designed to be traceable and auditable

**API Design**
- Consistent response envelope (`data` / `error`)
- Explicit error codes and messages
- Endpoints reflect real-world operations, not CRUD abstractions

### Technology Stack

- Go 1.26.1
- PostgreSQL 16
- pgx/v5
- net/http (standard library)

### Implemented Endpoints

```
POST   /auth/register
POST   /auth/login
GET    /auth/me

POST   /customers

POST   /accounts
POST   /accounts/{id}/deposit
POST   /accounts/{id}/withdraw
POST   /accounts/transfer
GET    /accounts/{id}/statement
```

Protected routes require JWT authentication.

### API Project Structure

```
api/
├── cmd/api/            # Bootstrap and route registration
├── internal/
│   ├── customer/       # Customer module
│   ├── account/        # Account and transaction logic
│   ├── auth/           # Authentication and authorization
│   └── database/       # Database initialization
├── migrations/         # Schema evolution
└── docs/               # Design documentation
```

### API Documentation

- [Architecture](docs/api/00-arquitetura.md)
- [Domain Model](docs/api/01-modelo_de_dominio.md)
- [Use Case Flows](docs/api/02-fluxos_de_caso_de_uso.md)
- [Data Model](docs/api/03-modelo_de_dados.md)
- [Consistency and Concurrency Strategy](docs/api/04-estrategia_de_consistencia_e_concorrencia.md)
- [Error Handling](docs/api/05-padrao_de_erros_e_respostas.md)
- [Implementation Details](docs/api/06-implementation.md)
- [API REST Design](docs/api/07-api-rest.md)
- [Authentication and Authorization](docs/api/08-auth_implementation.md)

---

## Mobile — BankFlow (Flutter)

A Flutter mobile client designed to validate and exercise the banking API, acting as a controlled integration environment rather than a feature-driven product.

### Scope

- JWT-based authentication
- Account creation and lifecycle management
- Financial operations: deposit, withdraw, transfer
- Transaction history with cursor-based pagination

### Architectural Role

- Validates backend assumptions through real usage flows
- Exposes inconsistencies in API design and data contracts
- Ensures alignment between user interaction and backend behavior
- Structured to reflect production concerns: clear separation between UI, state, and business logic

---

## Local Development

### Prerequisites

- Go 1.26.1+
- Flutter SDK
- Docker & Docker Compose
- [golang-migrate](https://github.com/golang-migrate/migrate)

### Start infrastructure

```bash
make docker-up
make api-migrate-up
```

### Run the API

```bash
make api-build
export JWT_SECRET=dev-change-me
./api/build/bank-api
```

API runs at `http://localhost:8080`.

### Run tests

```bash
make test              # API + Mobile tests
make api-test          # Go tests with coverage
make mobile-test       # Flutter tests
make mobile-test-unit  # Flutter unit tests (test/core)
```

---

## Makefile Reference

```
make help              List all available commands

make build             Build API binary
make test              Run API and Mobile tests

make api-build         Build API binary into api/build/
make api-migrate-up    Run database migrations
make api-migrate-down  Rollback last migration
make api-test          Run API tests with coverage

make mobile-test       Run all Flutter tests
make mobile-test-unit  Run Flutter unit tests

make docker-up         Start Docker containers
make docker-down       Stop Docker containers
make docker-logs       Follow Docker logs

make commit            Commit using ~/commit.md message
make diff              Show staged diff and line count
make push [branch=x]   Push branch to origin
make pull [branch=x]   Pull branch from origin
make gitlog            Show git log (one line)
```

---

## Future Work

- Zero Trust Architecture (context-aware request validation)
- Transaction-level authorization (transaction password / step-up auth)
- Onboarding flows aligned with financial systems
- Improved observability and audit capabilities

---

## License

MIT License — see the LICENSE file for details.
