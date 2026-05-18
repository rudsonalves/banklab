# BankLab

BankLab is an open source lab for studying and practicing engineering applied to financial systems, bringing together a **Go backend**, a **Flutter mobile app**, **transactional consistency**, **onboarding**, and future evolution toward **Zero Trust Architecture**.

![BankLab](docs/images/Capa_BankLab_LinkedIn_v3.png)

The project came from a practical gap: after studying Zero Trust in financial applications, I needed a dedicated environment where research, articles, and architectural decisions could become real implementation, with enough room to make mistakes, review, document, and mature the design.

A banking domain, even a simplified one, is a strong foundation for this kind of study because it forces the project to deal with problems that appear in real systems: identity, authentication, user approval, financial movements, traceability, API contracts, mobile/backend integration, and onboarding flows that cannot be treated as simple forms.

BankLab's central idea is to treat financial movements as the source of truth. An account balance is not just a loose number: it must be the consequence of ledger records persisted in the `transactions` table.

The project currently has two main applications:

- **Go API**: implements the banking core, including authentication, customers, accounts, and financial movements.
- **Flutter mobile app**: consumes the API and validates end-to-end flows for authentication, accounts, and transactions.

This repository is still evolving. The goal is for it to serve as a practical lab for people who want to study backend, mobile, tests, financial products, technical documentation, and collaboration in a project that feels like a real system.

## Goal

Build, incrementally, a simplified banking foundation capable of:

- registering users and customers;
- approving users before full access;
- opening and querying accounts;
- recording financial movements;
- querying balance and statements;
- performing internal transfers;
- keeping data traceability and consistency.

The focus is not to ship many features quickly. The focus is to build a smaller surface with care: clear rules, tests, documentation, and well-explained technical decisions.

## Project Principles

- **Financial integrity**: every relevant change must be recorded.
- **Atomicity**: critical operations, such as transfers, must either complete fully or not happen.
- **Traceability**: balance changes must be auditable.
- **Ledger as authority**: records in `transactions` are the source of truth for financial movements.
- **Consistency before volume**: the project prioritizes reliability over the number of screens or endpoints.
- **Explicit architecture**: layers and responsibilities should be easy to understand.

## Current Scope

| Domain         | Responsibilities                                                  |
| -------------- | ----------------------------------------------------------------- |
| Authentication | registration, login, refresh token, current user, and JWT control |
| Users          | registration, approval state, and administrative flow             |
| Customers      | automatic creation from user, CPF, and email                      |
| Accounts       | opening, listing, balance query, and lifecycle                    |
| Ledger         | deposit, withdraw, internal transfer, and movement history        |
| Statement      | paginated transaction listing by account                          |
| Mobile         | Flutter experience integrated with the API                        |

## Out of Scope for Now

- Pix, TED, or external banking integrations;
- antifraud and risk analysis;
- email, push, or SMS notifications;
- multiple currencies;
- external bank reconciliation;
- asynchronous processing of financial transactions.

These topics may appear in the future roadmap, but the foundation needs to become solid first.

## Stack

### API

- Go 1.26.1
- PostgreSQL 16
- `net/http`
- `pgx/v5`
- JWT
- migrations with `golang-migrate`

### Mobile

- Flutter
- Dart SDK ^3.11.4
- `dio`
- `go_router`
- `flutter_secure_storage`
- `auto_injector`

### Local Infra

- Docker
- Docker Compose
- Makefile for development commands

## Repository Structure

```text
banklab/
├── CHANGELOG.md       # Project changelog
├── CONTRIBUTING.md    # Contribution guidelines
├── LICENSE            # Project license
├── Makefile           # Development shortcuts
├── README.md          # Main documentation in Portuguese
├── README_en.md       # Main documentation in English
├── api/               # Go backend
├── docker-compose.yml # Local PostgreSQL configuration
├── docs/              # Roadmap, backlogs, decisions and reports
├── infra/             # Infrastructure scripts and configurations
├── mobile/            # Flutter BankFlow app
├── packages/          # Future repository for helper packages
├── templates/         # Pandoc templates
└── tools/postman/     # Collections and support for testing the API
```

## What Already Exists

### Go API

- registration, login, refresh token, and current user;
- automatic customer creation during registration;
- administrative approval flow for pending users;
- account opening and listing;
- balance query;
- internal transfers;
- transfer receipt;
- paginated statement;
- append-only persistence of movements in `transactions`;
- automated tests in important modules.

### Flutter Mobile App

- JWT authentication flow;
- handling of users pending approval;
- account visualization and creation;
- transfer between accounts;
- transaction history;
- integration with the local API.

## Running Locally

### Prerequisites

- Docker and Docker Compose;
- Go 1.26.1 or higher;
- Flutter configured on your machine;
- `golang-migrate`.

On macOS, `golang-migrate` can be installed with:

```bash
brew install golang-migrate
```

### 1. Start the API

From the repository root:

```bash
make run
```

This command validates Docker, starts PostgreSQL, waits for the database to become ready, applies migrations, and starts the API.

Default API URL:

```text
http://localhost:8080
```

### 2. Run the Mobile App

In another terminal:

```bash
cd mobile
flutter pub get
flutter run --dart-define-from-file=dev.env
```

Detailed guides are available at:

- [api/docs/00-getting_started.md](api/docs/00-getting_started.md)
- [mobile/docs/00-getting_started.md](mobile/docs/00-getting_started.md)

## Main Endpoints

```text
POST   /auth/register
POST   /auth/login
POST   /auth/refresh
GET    /auth/me

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

Registration and login routes require `X-App-Token`. The other routes require JWT authentication. Administrative routes also require an administrator user role.

## Useful Commands

```bash
make help

# local environment
make setup
make run
make reset

# API
make api-build
make api-tests

# Mobile
make mobile-tests
make mobile-test-unit

# All tests
make tests

# Docker
make docker-up
make docker-down
make docker-logs
make docker-clean
```

## Contributing

BankLab is open to collaboration from people who want to practice, learn, and grow together with an applied engineering project.

Good areas to contribute:

- **Go backend**: new use cases, tests, handlers, domain rules, and consistency improvements.
- **Flutter**: screens, navigation, error states, components, and user experience.
- **Tests**: unit tests, integration tests, concurrency scenarios, and critical flow coverage.
- **Documentation**: setup guides, architecture explanations, diagrams, and usage examples.
- **Financial product**: flow modeling, business rules, and roadmap refinement.
- **Infra and DevEx**: automation, scripts, CI, local environment, and Postman collections.

To organize the work, use the pattern described in [CONTRIBUTING.md](CONTRIBUTING.md). The idea is to classify each task by type, area, and priority, making it easier for new contributors to join.

## Suggested First Contributions

- review the setup guide and report confusing points;
- improve mobile error messages;
- add tests for transfer scenarios;
- create payload examples for main endpoints;
- document flows with simple diagrams;
- improve statement and receipt screens;
- propose small issues with well-defined scope.

## Documentation

### General

- [CONTRIBUTING.md](CONTRIBUTING.md)
- [docs/README.md](docs/README.md)
- [docs/ROADMAP.md](docs/ROADMAP.md)
- [docs/backlogs/README.md](docs/backlogs/README.md)
- [docs/backlogs/api/000 - pre-onboarding.md](<docs/backlogs/api/000 - pre-onboarding.md>)
- [docs/backlogs/api/001 - onboarding.md](<docs/backlogs/api/001 - onboarding.md>)
- [CHANGELOG.md](CHANGELOG.md)
- [tools/postman/README.md](tools/postman/README.md)

### API

- [api/README.md](api/README.md)
- [api/docs/00-getting_started.md](api/docs/00-getting_started.md)
- [api/docs/ARCHITECTURE.md](api/docs/ARCHITECTURE.md)
- [api/docs/01-domain_model.md](api/docs/01-domain_model.md)
- [api/docs/02-use_case_flows.md](api/docs/02-use_case_flows.md)
- [api/docs/03-application_model.md](api/docs/03-application_model.md)
- [api/docs/04-consistency_and_concorrency.md](api/docs/04-consistency_and_concorrency.md)
- [api/docs/05-error_and_response.md](api/docs/05-error_and_response.md)
- [api/docs/06-implementation.md](api/docs/06-implementation.md)
- [api/docs/07-api-rest.md](api/docs/07-api-rest.md)
- [api/docs/08-auth_implementation.md](api/docs/08-auth_implementation.md)
- [api/docs/09-database.md](api/docs/09-database.md)
- [api/docs/infra.md](api/docs/infra.md)

### Mobile

- [mobile/README.md](mobile/README.md)
- [mobile/docs/00-getting_started.md](mobile/docs/00-getting_started.md)
- [mobile/docs/ARCHITECTURE.md](mobile/docs/ARCHITECTURE.md)
- [mobile/docs/01-implemented-features.md](mobile/docs/01-implemented-features.md)

## Status

Active development. The foundation already supports studying real flows of a simplified banking application, but there is still room to evolve the product, tests, documentation, and mobile experience.

## License

MIT. See [LICENSE](LICENSE).
