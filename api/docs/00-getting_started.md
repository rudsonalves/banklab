# Getting Started — Bank API

## Index

- [Getting Started — Bank API](#getting-started--bank-api)
  - [Index](#index)
  - [1. Overview](#1-overview)
  - [2. Prerequisites](#2-prerequisites)
    - [2.1 Environment variables](#21-environment-variables)
    - [Description](#description)
    - [2.2 Docker engine](#22-docker-engine)
  - [3. Bootstrap (first run)](#3-bootstrap-first-run)
    - [3.1 Bootstrap the first admin user](#31-bootstrap-the-first-admin-user)
  - [4. Run the API](#4-run-the-api)
  - [5. Reset environment](#5-reset-environment)
  - [6. Notes](#6-notes)
  - [7. Troubleshooting](#7-troubleshooting)
    - [Docker not running](#docker-not-running)
    - [Database not ready](#database-not-ready)
    - [Missing environment variables](#missing-environment-variables)

---

## 1. Overview

This document describes how to run the API locally from a clean environment.

The setup assumes:

- a running Docker engine, such as Colima
- Go installed (for API execution)
- migrate CLI installed

The database is treated as the **source of truth**, and must always be initialized before running the API.

The local database runs on a custom PostgreSQL 17 image with `pg_cron` enabled.
`pg_cron` is loaded through `shared_preload_libraries` and is configured to run
inside the `bank` database using the `America/Sao_Paulo` timezone.

## 2. Prerequisites

### 2.1 Environment variables

Before running the API, initialize environment files with:

```bash
make env-init
```

This command creates independent API and mobile environment files:

```bash
./api/dev.env
./api/staging.env
./mobile/dev.env
./mobile/staging.env
```

Existing files and secrets are preserved. Newly created mobile files receive the
endpoint, mode, and application token from the corresponding API environment.

The API never discovers an environment file implicitly. `ENV_FILE` selects one
file explicitly, while the Make targets perform that selection for you.

Each API environment contains:

```env
APP_ENV=dev
PUBLIC_BASE_URL=http://localhost:8080
EXPOSE_DEBUG_VERIFICATION_TOKEN=true

APP_TOKEN=your_app_token_here
JWT_SECRET=your_jwt_secret_here
JWT_ACCESS_TOKEN_DURATION=15m
JWT_REFRESH_TOKEN_DURATION=168h
TRANSACTION_PASSWORD_PEPPER=your_transaction_password_pepper_here

SERVER_HOST=0.0.0.0
SERVER_PORT=8080

DB_HOST=127.0.0.1
DB_PORT=5432
DB_NAME=bank
DB_USER=postgres
DB_PASSWORD=postgres
```

Generated defaults use separate databases and Compose projects:

- development: `bank`, port `5432`, project `banklab-dev`
- staging: `bank_staging`, port `5433`, project `banklab-staging`

Only PostgreSQL runs in Docker. The API runs directly on the host and connects
to the port published by the selected PostgreSQL container:

- development: `127.0.0.1:5432`
- staging: `127.0.0.1:5433`

Run only the API process:

```bash
make api-run-dev
make api-run-staging
```

Start PostgreSQL, apply migrations, and run the selected API:

```bash
make dev
make staging
```

On Linux Mint, `make staging` starts the staging PostgreSQL container, applies
migrations, and runs the Go API on `127.0.0.1:8080`. Run
`cloudflared tunnel run banklab` in a separate terminal. The tunnel ingress
should point to `http://localhost:8080`.

To generate a strong pepper value:

```bash
openssl rand -base64 32
```

Requirements:

- `TRANSACTION_PASSWORD_PEPPER` must be at least 32 characters
- `TRANSACTION_PASSWORD_PEPPER` must be different from `APP_TOKEN` and `JWT_SECRET`
- `JWT_ACCESS_TOKEN_DURATION` and `JWT_REFRESH_TOKEN_DURATION` use Go duration syntax (`15m`, `168h`)
- never commit the real value to Git

### Description

- **APP_TOKEN**
  - Protects onboarding endpoints (`/auth/register`, `/auth/login`)
  - Must be sent in header: `X-App-Token`

- **JWT_SECRET**
  - Used to sign and validate JWT tokens
  - Must remain stable between application restarts

- **JWT_ACCESS_TOKEN_DURATION**
  - Controls the JWT access token expiration
  - Optional; defaults to `15m` when omitted

- **JWT_REFRESH_TOKEN_DURATION**
  - Controls the persisted refresh session expiration
  - Optional; defaults to `168h` when omitted

- **TRANSACTION_PASSWORD_PEPPER**
  - Used as an API-side secret to derive transaction password input before bcrypt
  - Is never stored in the database
  - Rotating this value invalidates existing transaction password hashes unless a migration strategy is implemented

- **EXPOSE_DEBUG_VERIFICATION_TOKEN**
  - Includes the contact verification code in API responses for controlled development/staging use

- **DB_HOST**, **DB_PORT**, **DB_NAME**, **DB_USER**, **DB_PASSWORD**
  - Configure the host API connection to PostgreSQL in Docker
  - Are also used by Make targets for Docker Compose, migrations, reset, readiness checks, and schema export
  - Use `DB_HOST=127.0.0.1`
  - Use `DB_PORT=5432` for development and `DB_PORT=5433` for staging

---

### 2.2 Docker engine

You need a running Docker engine. If you use Colima, `make docker-up` will start it automatically when needed.

If you prefer to start Docker manually, wait until the engine is ready:

```bash
docker info
```

---

## 3. Bootstrap (first run)

On the first run, you can bootstrap the whole stack with a single command:

```bash
make bootstrap
```

This flow will:

1. create the API and mobile `.env` files if they do not exist
2. start Colima if it is available and not already running
3. create and start the PostgreSQL container if it does not already exist
4. apply all migrations
5. start the API server

There is no extra Make target needed to create the container; `make docker-up`
runs the selected PostgreSQL service.
Docker Compose is invoked with `--env-file api/dev.env` by default, so PostgreSQL user,
password, database name, and host port follow the API environment file.
On a clean environment Docker builds the local PostgreSQL image from `infra/docker/postgres/Dockerfile`, which installs the `postgresql-17-cron` package required by `pg_cron`.

If you prefer the explicit sequence, the bootstrap is equivalent to:

```bash
make env-init
make docker-up
make migrate-up
make api-run
```

For staging, use `make staging`; for daily development, use `make dev`.

### 3.1 Bootstrap the first admin user

After the first user is registered, the system still needs an administrator to approve new accounts.

At the moment there is no admin UI for this first setup step. Promote the first user directly in PostgreSQL by changing their `role` to `admin`:

```bash
docker compose --env-file api/dev.env -p banklab-dev exec -T postgres \
  psql -U postgres -d bank \
  -c "UPDATE users SET role = 'admin', updated_at = NOW() WHERE email = 'admin@example.com';"
```

Replace `admin@example.com` with the email of the user that should become the initial administrator.

This bootstrap admin user does not need a bank account. The admin only needs to log in through the web client or Bruno, obtain a JWT, and call the approval endpoint for newly registered users:

```http
POST {{base_url}}/admin/users/{{id}}/approve
Authorization: Bearer {{access_token}}
```

The `id` path parameter is the UUID of the pending user returned by `POST /auth/register`.

Only users whose account approval flow has been completed can use the customer-facing parts of the application. In practice, a newly registered customer starts as `pending`; an admin must call `POST /admin/users/{id}/approve`, which changes the user to `active` and creates the associated account atomically.

---

## 4. Run the API

Start the API server:

```bash
make run
```

Or run only the API process:

```bash
make api-run-dev
```

`make dev` updates only `mobile/dev.env` with the current host LAN IP before
starting development. It never changes the staging configuration.

On Linux Mint, start staging with:

```bash
make staging
```

Then start the tunnel in another terminal:

```bash
cloudflared tunnel run banklab
```

To stop the API process listening on port 8080:

```bash
make api-stop
```

The server will be available at:

```
http://localhost:8080
```

---

## 5. Reset environment

To fully reset the system (including database):

```bash
make reset
```

This performs:

1. container removal (including volumes)
2. database recreation
3. migration reapplication

This guarantees a clean and deterministic state.

---

## 6. Notes

* Migrations are safe to re-run
* The database must always be ready before API startup
* Partial resets are discouraged due to the transactional model of the system

---

## 7. Troubleshooting

### Docker not running

```bash
make docker-check
```

If it fails:

```bash
sudo systemctl start docker
```

---

### Database not ready

Check container logs:

```bash
make docker-logs
```

---

### Missing environment variables

The application will fail at startup with:

```
missing required environment variable: APP_TOKEN
```

or

```
missing required environment variable: JWT_SECRET
```

or

```
TRANSACTION_PASSWORD_PEPPER environment variable is required
```

Ensure `api/dev.env` or `api/staging.env` exists and is correctly populated.
The API fails fast when any required `APP_TOKEN`, `JWT_SECRET`,
`TRANSACTION_PASSWORD_PEPPER`, or `DB_*` value is missing.
JWT duration variables are optional, but if provided they must be valid positive
Go durations such as `15m` or `168h`.

If files are missing, regenerate them safely with:

```bash
make env-init
```
