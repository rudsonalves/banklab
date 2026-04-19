# banklab

banklab is a monorepo built around a simplified banking core with emphasis on **transactional consistency** and **explicit business invariants**.

The system is structured around the premise that financial movements are the central element: balances are derived from ledger records, not treated as the primary source of truth.

The authoritative ledger is persisted in `account_transactions` (append-only). The legacy `transactions` table has been consolidated and is no longer part of the active model.

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
- ledger entries in `account_transactions` are the source of truth

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
- **Ledger authority** — financial records are persisted in `account_transactions` only
- **Consistency** — system state is always valid, even under concurrency
- **Synchronous model** — all operations complete at request time; no eventual consistency
- **Single source of truth** — the relational database is the only authority

---

## What is implemented

### API (Go)

- auth: register, login, refresh, current user (JWT)
- user registration with automatic customer creation (CPF and email)
- admin approval flow for pending users
- account opening, balance query, status control
- financial operations: deposit, withdraw, transfer between accounts
- account statement with pagination
- ledger persistence in `account_transactions` (append-only)
- transactional consistency enforced at the database level

### Mobile (Flutter)

- authentication flow with JWT
- account creation and management
- deposit, withdrawal, and transfer operations integrated with the API
- transaction history browsing

## Quick start

For the most up-to-date setup instructions, use:

- [api/docs/00-getting_started.md](api/docs/00-getting_started.md)
- [mobile/docs/00-getting_started.md](mobile/docs/00-getting_started.md)

### Prerequisites

- Docker and Docker Compose
- Go 1.26.1+
- Flutter SDK (matching your local setup)
- golang-migrate CLI

Install golang-migrate (macOS/Homebrew):

```bash
brew install golang-migrate
```

### 1) Start API stack (recommended)

```bash
make run
```

This validates Docker, starts PostgreSQL, applies migrations, and starts the API.

### 2) Run mobile app

```bash
cd mobile
flutter pub get
flutter run --dart-define-from-file=dev.env
```

For environment variables and detailed bootstrap/reset instructions, see the getting started guides above.

## Main endpoints

```text
POST   /auth/register
POST   /auth/login
POST   /auth/refresh
GET    /auth/me

POST   /admin/users/{id}/approve
GET    /customers/me

POST   /accounts
POST   /accounts/{id}/deposit
POST   /accounts/{id}/withdraw
POST   /accounts/transfer
GET    /accounts/{id}/statement
```

`/auth/register` and `/auth/login` require `X-App-Token`.
All other routes require JWT authentication.

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
- [api/docs/00-getting_started.md](api/docs/00-getting_started.md) — API getting started
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
- [mobile/docs/00-getting_started.md](mobile/docs/00-getting_started.md) — Mobile getting started
- [mobile/docs/ARCHITECTURE.md](mobile/docs/ARCHITECTURE.md) — Mobile Architecture

## License

MIT. See [LICENSE](LICENSE).
