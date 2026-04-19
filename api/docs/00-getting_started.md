# Getting Started — Bank API

## 1. Overview

This document describes how to run the API locally from a clean environment.

The setup assumes:

- Docker Desktop installed and running
- Go installed (for API execution)
- migrate CLI installed

The database is treated as the **source of truth**, and must always be initialized before running the API.

---

## 2. Prerequisites

### 2.1 Environment variables

Before running the API, create the environment file:

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

### 2.2 Docker

Start Docker:

```bash
open -a Docker
```

Wait until Docker is ready:

```bash
docker info
```

---

## 3. Bootstrap (first run)

Initialize the full environment:

```bash
make setup
```

This will:

1. validate Docker availability
2. start PostgreSQL container
3. wait for database readiness
4. apply all migrations

---

## 4. Run the API

Start the API server:

```bash
make run
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
