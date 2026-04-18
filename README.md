# banklab

banklab is a monorepo built around a simplified banking core with emphasis on **transactional consistency** and **explicit business invariants**.

The system is structured around the premise that the transaction is the central element — balance is a consequence of recorded movements, not a value maintained directly.

It consists of two applications:

- **API (Go)** — implements the banking core: customers, accounts, and financial operations with reliable transactional control
- **Mobile (Flutter)** — consumes the API and validates end-to-end flows across auth, account management, and transactions

The project is intentionally focused on correctness, consistency, and architecture decisions, not on feature volume.

## Repository structure

```text
banklab/
|-- api/            # Go backend (modular monolith)
|-- mobile/         # Flutter app (BankFlow)
|-- docs/           # Architecture and design docs
|-- infra/          # Docker and infrastructure scripts
|-- docker-compose.yml
`-- Makefile
```

## System scope

### Goal

Implement a simplified banking core capable of:

- managing customers
- maintaining bank accounts
- recording financial movements
- guaranteeing balance consistency

The focus is on **reliable transactional control**, not on peripheral features.

### Nature

> A balance-control system based on records of financial movements.

- the balance is a consequence
- the transaction is the central element

### In scope

| Domain    | Responsibilities                                |
| --------- | ----------------------------------------------- |
| Customers | creation, identification by CPF and email       |
| Accounts  | opening, balance query, status control          |
| Movements | deposit, withdraw, transfer, full operation log |
| Statement | transaction listing per account                 |

### Out of scope (at this stage)

- integration with external systems (Pix, TED, etc.)
- anti-fraud and risk analysis
- notifications (email, push)
- multi-currency
- bank reconciliation and external settlement

### System guarantees

- **Financial integrity** — no balance inconsistency; every movement is recorded
- **Atomicity** — critical operations (especially transfers) are indivisible
- **Traceability** — all operations are auditable; no balance change without a record
- **Consistency** — system state is always valid, even under concurrency
- **Synchronous model** — all operations complete at request time; no eventual consistency
- **Single source of truth** — the relational database is the only authority

---

## What is implemented

### API (Go)

- auth: register, login, current user (JWT)
- customer creation, identified by CPF and email
- account opening, balance query, status control
- financial operations: deposit, withdraw, transfer between accounts
- account statement with pagination
- transactional consistency enforced at the database level

### Mobile (Flutter)

- authentication flow with JWT
- account creation and management
- deposit, withdrawal, and transfer operations integrated with the API
- transaction history browsing

## Quick start

### Prerequisites

- Docker and Docker Compose
- Go 1.26.1+
- Flutter SDK (matching your local setup)
- golang-migrate CLI

Install golang-migrate (macOS/Homebrew):

```bash
brew install golang-migrate
```

### 1) Start infrastructure

```bash
make docker-up
```

This starts PostgreSQL 16 at localhost:5432.

### 2) Run database migrations

```bash
make api-migrate-up
```

### 3) Build and run API

```bash
make api-build
export JWT_SECRET=dev-change-me
./api/build/bank-api
```

API base URL: http://localhost:8080

### 4) Run mobile app

```bash
cd mobile
flutter pub get
flutter run
```

For emulator/device networking, point the app to the API URL that is reachable from your device.

## Main endpoints

```text
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

Routes other than register/login require JWT authentication.

## Development commands

```bash
make help

# monorepo
make build
make test

# api
make api-build
make api-migrate-up
make api-migrate-down
make api-test

# mobile
make mobile-test
make mobile-test-unit
make fclean
make fbuild

# docker
make docker-up
make docker-down
make docker-logs
```

## Project docs

### API (Go)

- [api/README.md](api/README.md) — API guide and setup
- [api/docs/ARCHITECTURE.md](api/docs/ARCHITECTURE.md) — API Architecture
- [api/docs/objetivos.md](api/docs/objetivos.md) — System Scope — Bank API
- [api/docs/01-domain_model.md](api/docs/01-domain_model.md) — Domain Model
- [api/docs/02-use_case_flows.md](api/docs/02-use_case_flows.md) — Use Case Flows
- [api/docs/03-application_model.md](api/docs/03-application_model.md) — Application Model
- [api/docs/04-consistency_and_concorrency.md](api/docs/04-consistency_and_concorrency.md) — Consistency and Concurrency Strategy
- [api/docs/05-error_and_response.md](api/docs/05-error_and_response.md) — Error and Response Standard
- [api/docs/06-implementation.md](api/docs/06-implementation.md) — Implementation Documentation
- [api/docs/07-api-rest.md](api/docs/07-api-rest.md) — REST API Documentation
- [api/docs/08-auth_implementation.md](api/docs/08-auth_implementation.md) — Auth & Authorization
- [api/docs/09-database.md](api/docs/09-database.md) — Database Documentation
- [api/docs/infra.md](api/docs/infra.md) — Infrastructure

### Mobile (Flutter)

- [mobile/README.md](mobile/README.md) — Mobile guide and setup
- [mobile/docs/ARCHITECTURE.md](mobile/docs/ARCHITECTURE.md) — Mobile Architecture

## License

MIT. See [LICENSE](LICENSE).
