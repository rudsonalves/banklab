# banklab

banklab is a monorepo with two main applications:

- a Go API that implements a simplified banking core
- a Flutter mobile app used to validate end-to-end flows against the API

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

## What is implemented

### API (Go)

- JWT authentication (register, login, me)
- customer creation
- account creation
- money operations: deposit, withdraw, transfer
- account statement with pagination

### Mobile (Flutter)

- authentication flow with JWT
- account creation and management flows
- transaction operations integrated with API
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

- API guide: [api/README.md](api/README.md)
- Mobile guide: [mobile/README.md](mobile/README.md)
- Mobile architecture: [docs/mobile/ARCHITECTURE.md](docs/mobile/ARCHITECTURE.md)
- Architecture and design docs: [docs/api](docs/api)

## License

MIT. See [LICENSE](LICENSE).
