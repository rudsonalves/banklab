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

Start database:
