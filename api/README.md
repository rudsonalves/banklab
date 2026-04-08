# Bank API

Bank API is a modular Go service that implements a simplified banking core, designed with a strong emphasis on transactional consistency and system correctness.

Rather than a feature-complete product, this project serves as an engineering-focused environment to explore backend design decisions in financial systems.

## Purpose

This project is part of a broader effort to build and validate a complete system (backend + mobile client), focusing on:

- correctness of financial operations
- consistency guarantees under concurrent access
- clear API contracts between backend and client
- explicit modeling of business rules and invariants

## Core Capabilities

- customer creation
- user registration and authentication (JWT)
- account creation and ownership enforcement
- financial operations (deposit, withdraw, transfer)
- account statement with pagination and filtering

## Architectural Approach

The system follows a layered modular monolith structure:

- Delivery (HTTP layer)
- Application (use cases and transaction orchestration)
- Domain (business rules and invariants)
- Infrastructure (PostgreSQL and repositories)

Dependency direction:

- Delivery → Application → Domain  
- Infrastructure → Domain  

This structure is intentionally enforced to keep business logic isolated and predictable.

## Design Focus

### Transactional Consistency

Financial operations are treated as critical sections:

- all balance-changing operations run inside explicit database transactions
- row-level locking (`SELECT FOR UPDATE`) is used to prevent race conditions
- deterministic lock ordering is applied in transfers to reduce deadlock risk

### Ledger Integrity

- all balance changes are recorded as immutable ledger entries
- account balance is derived from controlled updates, not arbitrary mutations
- operations are designed to be traceable and auditable

### Domain Invariants

- business rules are enforced at domain and application levels
- invalid state transitions are rejected early
- consistency is prioritized over convenience

### API Design

- consistent response envelope (`data` / `error`)
- explicit error codes and messages
- endpoints designed to reflect real-world operations rather than CRUD abstractions

## Technology Stack

- Go 1.26.1
- PostgreSQL 16
- pgx/v5
- net/http (standard library)

## Implemented Endpoints

Authentication:
- POST /auth/register
- POST /auth/login
- GET /auth/me

Customer:
- POST /customers

Accounts:
- POST /accounts
- POST /accounts/{id}/deposit
- POST /accounts/{id}/withdraw
- POST /accounts/transfer
- GET /accounts/{id}/statement

Protected routes require JWT authentication.

## Persistence Model

- relational schema in PostgreSQL
- balances stored in minor units (BIGINT)
- ledger entries stored in `account_transactions`
- explicit transaction boundaries for all critical operations

## Local Development

**Start database:**

```
make docker-up
make api-migrate-up
```

This starts PostgreSQL 16 in a container and runs all pending database migrations.

**Run tests:**

```
make api-test
```

Runs all tests with coverage report.

**Build and run API:**

```
make api-build
export JWT_SECRET=dev-change-me
./api/build/bank-api
```

Compiles the binary into the `api/build/` directory and runs the server on `http://localhost:8080`.

## Project Structure

- cmd/api: application bootstrap and route registration
- internal/customer: customer module (domain, application, delivery, infra)
- internal/account: account and transaction logic
- internal/auth: authentication and authorization
- internal/database: database initialization
- migrations: schema evolution
- docs: architectural and technical documentation

## Documentation

Detailed design decisions are documented in `../docs/api/`:

- [Architecture](../docs/api/00-arquitetura.md)
- [Domain Model](../docs/api/01-modelo_de_dominio.md)
- [Use Case Flows](../docs/api/02-fluxos_de_caso_de_uso.md)
- [Data Model](../docs/api/03-modelo_de_dados.md)
- [Consistency and Concurrency Strategy](../docs/api/04-estrategia_de_consistencia_e_concorrencia.md)
- [Error Handling](../docs/api/05-padrao_de_erros_e_respostas.md)
- [Implementation Details](../docs/api/06-implementation.md)
- [API REST Design](../docs/api/07-api-rest.md)
- [Authentication and Authorization](../docs/api/08-auth_implementation.md)

## Future Work

Planned extensions include:

- Zero Trust Architecture (context-aware request validation)
- transaction-level authorization (transaction password / step-up auth)
- onboarding flows aligned with financial systems
- improved observability and audit capabilities

## Notes

- this project prioritizes correctness and design clarity over feature breadth
- it is intended as an engineering exploration, not a production-ready system
