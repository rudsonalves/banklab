# REST API Documentation - Bank API

## 1. Overview

This document describes the current HTTP REST contract implemented by the service.

Base URL (local):
- http://localhost:8080

Content type:
- request: application/json
- response: application/json

## 2. Response Envelope

All endpoints return a standard envelope.

Success:

```json
{
  "data": {},
  "error": null
}
```

Error:

```json
{
  "data": null,
  "error": {
    "code": "ERROR_CODE",
    "message": "human readable message",
    "details": {}
  }
}
```

## 3. Endpoints

## 3.1 Create Customer

- Method: POST
- Path: /customers

Request body:

```json
{
  "name": "Maria Silva",
  "cpf": "12345678901",
  "email": "maria@example.com"
}
```

Success response (201):

```json
{
  "data": {
    "id": "f992f5f5-c4c5-4e25-95b8-3f3f1ee4ae39",
    "name": "Maria Silva",
    "cpf": "12345678901",
    "email": "maria@example.com",
    "created_at": "2026-04-02T12:00:00Z"
  },
  "error": null
}
```

Possible errors:
- 400 INVALID_REQUEST: invalid JSON body
- 400 INVALID_DATA: missing/invalid fields (may include details.field)
- 409 CUSTOMER_ALREADY_EXISTS: duplicate CPF or email
- 500 INTERNAL_ERROR: unexpected internal error

## 3.2 Create Account

- Method: POST
- Path: /accounts

Request body:

```json
{
  "customer_id": "6f3ebf86-bf82-4b75-a2ce-cd261ca47ec3"
}
```

Success response (201):

```json
{
  "data": {
    "id": "fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b",
    "customer_id": "6f3ebf86-bf82-4b75-a2ce-cd261ca47ec3",
    "number": "10000001",
    "branch": "0001",
    "balance": 0,
    "status": "active"
  },
  "error": null
}
```

Possible errors:
- 400 INVALID_REQUEST: invalid JSON body
- 400 INVALID_DATA: customer_id invalid
- 404 CUSTOMER_NOT_FOUND: customer does not exist
- 500 INTERNAL_ERROR: unexpected internal error

## 3.3 Deposit

- Method: POST
- Path: /accounts/{id}/deposit
- Path param:
  - id: account UUID

Request body:

```json
{
  "amount": 5000
}
```

Success response (200):

```json
{
  "data": {
    "id": "fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b",
    "balance": 15000
  },
  "error": null
}
```

Possible errors:
- 400 INVALID_DATA: invalid account id
- 400 INVALID_REQUEST: invalid JSON body
- 400 INVALID_AMOUNT: amount must be greater than zero
- 404 ACCOUNT_NOT_FOUND: account does not exist
- 422 ACCOUNT_INACTIVE: account not active
- 500 INTERNAL_ERROR: unexpected internal error

## 3.4 Withdraw

- Method: POST
- Path: /accounts/{id}/withdraw
- Path param:
  - id: account UUID

Request body:

```json
{
  "amount": 3000
}
```

Success response (200):

```json
{
  "data": {
    "id": "fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b",
    "balance": 12000
  },
  "error": null
}
```

Possible errors:
- 400 INVALID_DATA: invalid account id
- 400 INVALID_REQUEST: invalid JSON body
- 400 INVALID_AMOUNT: amount must be greater than zero
- 404 ACCOUNT_NOT_FOUND: account does not exist
- 422 INSUFFICIENT_BALANCE: insufficient funds
- 422 ACCOUNT_INACTIVE: account not active
- 500 INTERNAL_ERROR: unexpected internal error

## 3.5 Transfer

- Method: POST
- Path: /accounts/transfer

Request body:

```json
{
  "from_account_id": "7e2a56a4-b56c-44aa-9204-5e6c2df659d5",
  "to_account_id": "f2ec464e-dd1d-4b89-9f29-bf45dcbf16ff",
  "amount": 2500
}
```

Success response (200):

```json
{
  "data": {
    "from_account_id": "7e2a56a4-b56c-44aa-9204-5e6c2df659d5",
    "to_account_id": "f2ec464e-dd1d-4b89-9f29-bf45dcbf16ff",
    "amount": 2500,
    "from_balance": 97500,
    "to_balance": 32500
  },
  "error": null
}
```

Possible errors:
- 400 INVALID_REQUEST: invalid JSON body
- 400 INVALID_DATA: invalid UUID data
- 400 INVALID_AMOUNT: amount must be greater than zero
- 400 SAME_ACCOUNT_TRANSFER: source and destination are equal
- 404 ACCOUNT_NOT_FOUND: source or destination account not found
- 422 INSUFFICIENT_BALANCE: source account has insufficient funds
- 422 ACCOUNT_INACTIVE: one account is inactive
- 500 INTERNAL_ERROR: unexpected internal error

## 3.6 Get Statement

- Method: GET
- Path: /accounts/{id}/statement
- Path param:
  - id: account UUID

Query params (optional):
- limit: integer, default 50, max 100
- cursor: RFC3339 datetime
- cursor_id: UUID
- from: RFC3339 datetime
- to: RFC3339 datetime

Notes:
- cursor and cursor_id must be provided together.
- items are returned in descending order by created_at and id.

Example request:
- GET /accounts/{id}/statement?limit=2&from=2026-04-01T00:00:00Z&to=2026-04-02T23:59:59Z

Success response (200):

```json
{
  "data": {
    "account_id": "fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b",
    "items": [
      {
        "transaction_id": "0fd87d49-d94e-4449-bde4-0c0808f7645f",
        "type": "deposit",
        "amount": 5000,
        "balance_after": 15000,
        "reference_id": null,
        "created_at": "2026-04-02T12:00:00Z"
      },
      {
        "transaction_id": "6b667f90-c94b-4f81-a8c2-b6ca38a4e45d",
        "type": "transfer_in",
        "amount": 2500,
        "balance_after": 10000,
        "reference_id": "7a8da098-5351-469a-95be-b59d14fe43a8",
        "created_at": "2026-04-02T11:55:00Z"
      }
    ],
    "next_cursor": {
      "created_at": "2026-04-02T11:55:00Z",
      "id": "6b667f90-c94b-4f81-a8c2-b6ca38a4e45d"
    }
  },
  "error": null
}
```

Possible errors:
- 400 INVALID_DATA: invalid path/query value or cursor/cursor_id mismatch
- 404 ACCOUNT_NOT_FOUND: account does not exist
- 500 INTERNAL_ERROR: unexpected internal error

## 4. Error Code Reference

Common error codes currently used by handlers:
- INVALID_REQUEST
- INVALID_DATA
- INVALID_AMOUNT
- CUSTOMER_ALREADY_EXISTS
- CUSTOMER_NOT_FOUND
- ACCOUNT_NOT_FOUND
- ACCOUNT_INACTIVE
- INSUFFICIENT_BALANCE
- SAME_ACCOUNT_TRANSFER
- INTERNAL_ERROR

## 5. Domain Notes for API Consumers

- Monetary values are represented as integer cents.
- UUID is used for all resource identifiers.
- Financial operations are synchronous and strongly consistent.
- Transfer operation is atomic: debit and credit are committed together.
