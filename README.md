# Bank API

Bank API is a modular Go service that implements a simplified banking core with strong transactional consistency.

It provides:
- customer creation
- user registration and login (JWT)
- account creation
- deposit, withdraw, and transfer operations
- account statement with pagination and filters
- authentication and account ownership authorization

The project is organized as a layered modular monolith, with clear separation between HTTP delivery, application use cases, domain rules, and PostgreSQL infrastructure.

## Main Goals

- keep financial operations correct and predictable
- guarantee ledger traceability for balance changes
- enforce business invariants at domain and application levels
- support safe concurrent operations through database transactions and row locking

## Architecture

Current dependency direction:

- Delivery -> Application -> Domain
- Infrastructure -> Domain

Entrypoint:
- cmd/api/main.go

Main modules:
- internal/customer
- internal/account
- internal/auth
- internal/database

## Stack

- Go 1.26.1
- PostgreSQL 16
- pgx/v5
- net/http (standard library)

## Implemented Endpoints

- POST /auth/register
- POST /auth/login
- GET /auth/me
- POST /customers
- POST /accounts
- POST /accounts/{id}/deposit
- POST /accounts/{id}/withdraw
- POST /accounts/transfer
- GET /accounts/{id}/statement

Protected routes (JWT required):
- GET /auth/me
- /accounts/*

## Response Pattern

All endpoints use a consistent envelope:

- success: data populated and error null
- failure: data null and error populated with code/message/details

## Persistence and Consistency

- relational model in PostgreSQL
- account balances stored in cents using BIGINT
- ledger entries persisted in account_transactions
- explicit transactions for all balance-changing operations
- SELECT FOR UPDATE for critical row locking
- deterministic lock ordering in transfer to reduce deadlock risk

## Local Development

1. Start PostgreSQL

	docker compose up -d

2. Run tests

	go test ./...

3. Run API

	export JWT_SECRET=dev-change-me
	go run ./cmd/api

Server default address:
- http://localhost:8080

## Project Structure

- cmd/api: process bootstrap and route registration
- internal/customer: customer domain, use case, HTTP, repository
- internal/account: account and ledger domain, use cases, HTTP, repository
- internal/auth: auth domain, use cases, middleware, and repository
- internal/database: DB pool creation
- db: SQL schema
- migrations: schema evolution files
- docs: architecture and technical documentation

## Documentation

- docs/00-arquitetura.md
- docs/01-modelo_de_dominio.md
- docs/02-fluxos_de_caso_de_uso.md
- docs/03-modelo_de_dados.md
- docs/04-estrategia_de_consistencia_e_concorrencia.md
- docs/05-padrao_de_erros_e_respostas.md
- docs/06-implementation.md
- docs/07-api-rest.md
- docs/08-auth_implementation.md

## Current Notes

- README now summarizes the implemented system and operational flow.
- Connection string is currently hard-coded in internal/database/db.go.
- One-account-per-customer rule exists as optional logic scaffold in create account use case.
