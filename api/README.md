# Bank API

Bank API is a Go service that implements a simplified banking core with emphasis on transactional consistency and explicit business invariants.

This service is part of the banklab monorepo and is designed to be consumed by the Flutter mobile app in the same repository.

## Stack

- Go 1.26.1
- PostgreSQL 16
- pgx/v5
- net/http

## Architecture

Modular monolith with layered boundaries:

- Delivery: HTTP handlers and request/response mapping
- Application: use cases and transaction orchestration
- Domain: entities, value objects, invariants, domain errors
- Infrastructure: repositories and database integration

Dependency direction:

- Delivery -> Application -> Domain
- Infrastructure -> Domain

## Features

- auth: register, login, current user
- customer creation
- account creation
- balance-changing operations: deposit, withdraw, transfer
- account statement listing with pagination support

## API routes

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

All routes except register/login require JWT authentication.

## Local setup

The recommended flow is from repository root.

1. Start database:

```bash
make docker-up
```

2. Run migrations:

```bash
make api-migrate-up
```

3. Build API:

```bash
make api-build
```

4. Run API:

```bash
export JWT_SECRET=dev-change-me
./api/build/bank-api
```

Default URL: http://localhost:8080

## Tests

From repository root:

```bash
make api-test
```

Or directly from api directory:

```bash
cd api
go test -cover ./...
```

## Directory map

```text
api/
|-- cmd/api/            # application bootstrap
|-- internal/
|   |-- account/
|   |-- auth/
|   |-- customer/
|   |-- database/
|   `-- shared/
|-- migrations/
`-- README.md
```

## Design documents

- [Getting started](docs/00-getting_started.md)
- [Architecture](docs/ARCHITECTURE.md)
- [Domain model](docs/01-domain_model.md)
- [Use case flows](docs/02-use_case_flows.md)
- [Application model](docs/03-application_model.md)
- [Consistency and concurrency strategy](docs/04-consistency_and_concorrency.md)
- [Error and response standard](docs/05-error_and_response.md)
- [Implementation notes](docs/06-implementation.md)
- [REST API design](docs/07-api-rest.md)
- [Auth implementation](docs/08-auth_implementation.md)
- [Database](docs/09-database.md)
- [Infrastructure](docs/infra.md)

## Related docs

- Monorepo overview: [../README.md](../README.md)
