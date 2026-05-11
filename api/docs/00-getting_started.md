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

## 2. Prerequisites

### 2.1 Environment variables

Before running the API, initialize environment files with:

```bash
make env-init
```

This command creates the files below only if they do not exist (existing files are preserved):

```bash
./api/.env
./mobile/dev.env
./mobile/staging.env
./mobile/prod.env
```

If you prefer manual setup, create the API environment file:

```bash
touch api/.env
```

Add the following variables:

```env
APP_TOKEN=your_app_token_here
JWT_SECRET=your_jwt_secret_here
```

Example:

```env
APP_TOKEN=a3f5905dc26977e9408b3eca832869c2d49e4f7cf6d2026cff234075fd703ad5
JWT_SECRET=b03ff724fc843ace8ea69f2e00bdb6192e342f90038a8532d55bae3d42427d2d
```

### Description

- **APP_TOKEN**
  - Protects onboarding endpoints (`/auth/register`, `/auth/login`)
  - Must be sent in header: `X-App-Token`

- **JWT_SECRET**
  - Used to sign and validate JWT tokens
  - Must remain stable between application restarts

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

There is no extra Make target needed to create the container; `make docker-up` already runs `docker compose up -d --no-recreate`.

If you prefer the explicit sequence, the bootstrap is equivalent to:

```bash
make env-init
make docker-up
make migrate-up
make api-run
```

If you only want the one-step bootstrap for the database portion, `make setup` still performs the Docker check, container startup, wait, and migrations.

### 3.1 Bootstrap the first admin user

After the first user is registered, the system still needs an administrator to approve new accounts.

At the moment there is no admin UI for this first setup step. Promote the first user directly in PostgreSQL by changing their `role` to `admin`:

```bash
docker exec -i bank-postgres psql -U postgres -d bank \
  -c "UPDATE users SET role = 'admin', updated_at = NOW() WHERE email = 'admin@example.com';"
```

Replace `admin@example.com` with the email of the user that should become the initial administrator.

This bootstrap admin user does not need a bank account. The admin only needs to log in through the web client or Postman, obtain a JWT, and call the approval endpoint for newly registered users:

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
make api-run
```

`make api-run` automatically executes `make mobile-sync-ip` first, updating `BASE_URL` in mobile `.env` files with the current host LAN IP.

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
open -a Docker
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

Ensure the file `api/.env` exists and is correctly populated.

If files are missing, regenerate them safely with:

```bash
make env-init
```
