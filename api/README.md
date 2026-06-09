# Bank API

Bank API is a Go service that implements a simplified banking core with emphasis on transactional consistency and explicit business invariants.

This service is part of the banklab monorepo and is designed to be consumed by the Flutter mobile app in the same repository.

## Stack

- Go 1.26.1
- PostgreSQL 17 with pg_cron
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
- transaction password setup and step-up authorization
- customer self-profile lookup
- account creation/provisioning
- account balance lookup
- balance-changing operations: deposit, withdraw, transfer
- account statement listing with pagination support

Note: account creation is an admin provisioning capability exposed through
`POST /admin/customers/{customer_id}/accounts`. Direct `deposit` and `withdraw`
terminal routes are intentionally not registered while a real terminal channel
is outside the project scope.

## API routes

```text
POST   /auth/register
POST   /auth/login
POST   /auth/refresh
GET    /auth/me

POST   /security/transaction-password
POST   /security/step-up/authorize

POST   /admin/users/{id}/approve
POST   /admin/customers/{customer_id}/accounts

GET    /customers/me

GET    /accounts
GET    /accounts/{id}/balance
GET    /accounts/internal-transfers/recipients
POST   /accounts/internal-transfers
GET    /accounts/transfer/{transaction_reference}/receipt
GET    /accounts/{id}/statement
```

All routes except register/login require JWT authentication. `POST /accounts/internal-transfers` also requires an `X-Step-Up-Token` issued by `POST /security/step-up/authorize` for the public operation `method=POST` and `path=/accounts/internal-transfers`.

`POST /terminal/accounts/{id}/deposit` and
`POST /terminal/accounts/{id}/withdraw` are intentionally not registered while a
real terminal channel is outside the project scope.

## Local setup

The recommended flow is from repository root.
The API runtime environment lives in `api/.env`. Local Make targets use that same
file for Docker Compose, migrations, database reset, readiness checks, and API
startup.

Initialize missing environment files:

```bash
make env-init
```

1. Start database:

```bash
make docker-up
```

The local PostgreSQL container is built from `infra/docker/postgres/Dockerfile`.
It installs `pg_cron` and starts PostgreSQL with:

- `shared_preload_libraries=pg_cron`
- `cron.database_name=bank`
- `cron.timezone=America/Sao_Paulo`

2. Run migrations:

```bash
make migrate-up
```

3. Build API:

```bash
make api-build
```

4. Build and run the API container in the selected environment:

```bash
make api-run-dev
make api-run-staging
make api-run-prod
```

Use a dedicated random value for `TRANSACTION_PASSWORD_PEPPER` (for example `openssl rand -base64 32`) and do not reuse `APP_TOKEN` or `JWT_SECRET`.
JWT token lifetimes can be configured with `JWT_ACCESS_TOKEN_DURATION` and
`JWT_REFRESH_TOKEN_DURATION`; when omitted, they default to `15m` and `168h`.

The API loads only the file explicitly selected through `ENV_FILE`. The Make
targets select `api/dev.env`, `api/staging.env`, or `api/prod.env`.

Development URL: http://localhost:8080
Staging public URL: https://api.rralves.dev.br

For staging, `make staging` starts PostgreSQL, applies migrations, and starts the
API container. Then run `cloudflared tunnel run banklab` on the Docker host. The
tunnel should forward `api.rralves.dev.br` to `http://localhost:8080`.

## Tests

From repository root:

```bash
make api-tests
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
- [Architecture presentation material](docs/presentation-api-architecture.md)

## Related docs

- Monorepo overview: [../README.md](../README.md)
