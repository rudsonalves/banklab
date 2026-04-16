# API Architecture

## Overview

Bank API is implemented as a modular monolith with layered boundaries. The architecture prioritizes transactional correctness for financial operations, explicit domain rules, and low accidental coupling.

The service is organized by business modules (account, auth, customer), and each module keeps its own layered split:

- delivery: HTTP handlers, request parsing, response mapping
- application: use case orchestration and transaction boundaries
- domain: entities, invariants, value objects, contracts, domain errors
- infrastructure: PostgreSQL repositories, hashing, JWT, transactors

## Architectural style

The current style is:

> Modular monolith + layered architecture with clean dependency direction.

Dependency rule:

```text
delivery -> application -> domain
infrastructure -> domain
```

No domain code depends on HTTP, persistence libraries, or framework-specific concerns.

## Runtime composition

Application wiring happens in [api/cmd/api/main.go](../../api/cmd/api/main.go).

At startup, the process:

1. Initializes bootstrap routines
2. Creates the PostgreSQL pool
3. Instantiates repositories and infrastructure adapters
4. Instantiates use cases
5. Registers HTTP routes and auth middleware

This keeps composition centralized and explicit.

## Module map

Current top-level modules under [api/internal](../../api/internal):

- account
- auth
- customer
- database (shared DB helpers)
- shared (cross-module utilities)
- bootstrap (startup initialization)

Each business module follows a similar internal shape:

```text
internal/<module>/
|-- application/
|-- delivery/
|-- domain/
`-- infrastructure/
```

## Layer responsibilities

### Domain

Defines core business concepts and invariants:

- account ownership and access constraints
- transfer validity rules
- balance constraints and operation semantics
- typed domain errors
- repository and service contracts

Domain is pure business logic and is the most stable layer.

### Application

Coordinates use cases and consistency guarantees:

- validates use-case-level input
- applies authorization checks with authenticated user context
- manages transaction scope through repository/transactor contracts
- orchestrates multi-step operations (for example transfers)

Examples include create account, deposit, withdraw, transfer, login, register, and refresh token.

### Infrastructure

Implements technical details behind domain contracts:

- PostgreSQL repositories (pgx)
- SQL statements and mapping
- password hashing (bcrypt)
- JWT token generation/parsing
- database transaction execution helpers

Infrastructure may depend on drivers and libraries, but it should not own business policy.

### Delivery

Exposes HTTP endpoints and translates between transport and application:

- input decoding and validation
- request-to-use-case mapping
- use-case result/error to HTTP response mapping
- authentication middleware integration

Delivery does not implement core business rules.

## Request lifecycle

Typical flow:

```text
HTTP Request
  -> Delivery handler
  -> Application use case
  -> Domain validation/rules
  -> Infrastructure persistence/integration
  -> Application result
  -> Delivery response
  -> HTTP Response
```

## Transaction and consistency model

Transaction boundaries are defined at the application/use-case level and executed through repository/transactor abstractions.

Important implemented patterns:

- row locking with `SELECT ... FOR UPDATE` on critical account reads
- deterministic lock ordering for transfer operations to reduce deadlock risk
- explicit commit/rollback behavior around multi-step operations
- idempotency support for transfer requests via idempotency key + operation table
- immutable ledger-style account transactions for auditable balance changes

These patterns are central to preventing race conditions in balance updates.

## Authentication and session architecture

Authentication is JWT-based with access and refresh tokens.

Current flow:

1. User logs in with credentials
2. API returns access token + refresh token
3. Refresh token hash is stored in sessions storage
4. Refresh endpoint validates token, user, revocation state, and expiration
5. Refresh rotates session token inside a DB transaction

This provides token rotation with server-side session control.

## API surface (current)

Registered routes include:

- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/refresh`
- `GET /auth/me`
- `GET /customers/me`
- `POST /accounts`
- `POST /accounts/{id}/deposit`
- `POST /accounts/{id}/withdraw`
- `POST /accounts/transfer`
- `GET /accounts/{id}/statement`

Protected routes are guarded by JWT middleware.

## Design decisions

Intentional decisions in the current phase:

- keep a modular monolith instead of distributed services
- avoid premature adoption of event-driven architecture
- avoid CQRS/event sourcing until complexity justifies it
- prioritize deterministic behavior and consistency over throughput-oriented complexity

## Known trade-offs

Current architecture is intentionally conservative and does not optimize for:

- large-scale horizontal distribution
- asynchronous workflows/background processing
- independent deployability per module

These trade-offs are acceptable for the current stage and goals.

## Evolution path

Expected evolution options:

- stronger feature/module boundaries as scope grows
- selective extraction of services only when operational pressure requires it
- asynchronous integration points where eventual consistency is acceptable
- richer observability and operational metrics around transaction-critical paths

The current architecture is designed to evolve incrementally without breaking core domain boundaries.
