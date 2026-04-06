# REST API Documentation - Bank API

## 1. Overview

This document describes the HTTP REST contract currently implemented by the service.

Base URL (local):
- http://localhost:8080

Content type:
- request: application/json
- response: application/json

Authentication:
- JWT Bearer token
- Send in header Authorization: Bearer <access_token>

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

## 3. Authentication Endpoints

## 3.1 Register User

- Method: POST
- Path: /auth/register
- Auth required: no

Request body:

```json
{
  "email": "user@example.com",
  "password": "P@ssword123"
}
```

Success response (201):

```json
{
  "data": {
    "id": "d3de5f8b-4892-42e8-9680-979cf3f37844",
    "email": "user@example.com",
    "role": "customer",
    "customer_id": null
  },
  "error": null
}
```

Possible errors:
- 400 INVALID_REQUEST: invalid JSON body
- 400 INVALID_DATA: invalid email or password
- 409 USER_ALREADY_EXISTS: duplicate email
- 500 INTERNAL_ERROR: unexpected internal error

## 3.2 Login User

- Method: POST
- Path: /auth/login
- Auth required: no

Request body:

```json
{
  "email": "user@example.com",
  "password": "P@ssword123"
}
```

Success response (200):

```json
{
  "data": {
    "access_token": "<jwt>",
    "user_id": "d3de5f8b-4892-42e8-9680-979cf3f37844",
    "email": "user@example.com",
    "role": "customer",
    "customer_id": null
  },
  "error": null
}
```

Possible errors:
- 400 INVALID_REQUEST: invalid JSON body
- 400 INVALID_DATA: invalid email or password input
- 401 INVALID_CREDENTIALS: invalid email/password
- 500 INTERNAL_ERROR: unexpected internal error

## 3.3 Get Current User

- Method: GET
- Path: /auth/me
- Auth required: yes

Success response (200):

```json
{
  "data": {
    "id": "d3de5f8b-4892-42e8-9680-979cf3f37844",
    "email": "user@example.com",
    "role": "customer",
    "customer_id": "6f3ebf86-bf82-4b75-a2ce-cd261ca47ec3"
  },
  "error": null
}
```

Possible errors:
- 401 UNAUTHORIZED: authentication required
- 401 INVALID_TOKEN: token invalid, malformed, or expired
- 500 INTERNAL_ERROR: unexpected internal error

## 4. Customer Endpoint

## 4.1 Create Customer

- Method: POST
- Path: /customers
- Auth required: no

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
- 400 INVALID_DATA: missing/invalid fields
- 409 USER_ALREADY_EXISTS: duplicate CPF or email
- 500 INTERNAL_ERROR: unexpected internal error

## 5. Account Endpoints

All account routes are protected and require Authorization header with Bearer token.

## 5.1 Create Account

- Method: POST
- Path: /accounts
- Auth required: yes

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
- 401 UNAUTHORIZED: authentication required
- 401 INVALID_TOKEN: token invalid, malformed, or expired
- 400 INVALID_REQUEST: invalid JSON body
- 400 INVALID_DATA: customer_id invalid
- 403 FORBIDDEN: access denied to account
- 404 CUSTOMER_NOT_FOUND: customer does not exist
- 500 INTERNAL_ERROR: unexpected internal error

## 5.2 Deposit

- Method: POST
- Path: /accounts/{id}/deposit
- Auth required: yes

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
- 401 UNAUTHORIZED: authentication required
- 401 INVALID_TOKEN: token invalid, malformed, or expired
- 400 INVALID_DATA: invalid account id
- 400 INVALID_REQUEST: invalid JSON body
- 400 INVALID_AMOUNT: amount must be greater than zero
- 403 FORBIDDEN: access denied to account
- 404 ACCOUNT_NOT_FOUND: account does not exist
- 422 ACCOUNT_INACTIVE: account not active
- 500 INTERNAL_ERROR: unexpected internal error

## 5.3 Withdraw

- Method: POST
- Path: /accounts/{id}/withdraw
- Auth required: yes

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
- 401 UNAUTHORIZED: authentication required
- 401 INVALID_TOKEN: token invalid, malformed, or expired
- 400 INVALID_DATA: invalid account id
- 400 INVALID_REQUEST: invalid JSON body
- 400 INVALID_AMOUNT: amount must be greater than zero
- 403 FORBIDDEN: access denied to account
- 404 ACCOUNT_NOT_FOUND: account does not exist
- 422 INSUFFICIENT_FUNDS: insufficient funds
- 422 ACCOUNT_INACTIVE: account not active
- 500 INTERNAL_ERROR: unexpected internal error

## 5.4 Transfer

- Method: POST
- Path: /accounts/transfer
- Auth required: yes

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
- 401 UNAUTHORIZED: authentication required
- 401 INVALID_TOKEN: token invalid, malformed, or expired
- 400 INVALID_REQUEST: invalid JSON body
- 400 INVALID_DATA: invalid UUID data
- 400 INVALID_AMOUNT: amount must be greater than zero
- 400 SAME_ACCOUNT_TRANSFER: source and destination are equal
- 403 FORBIDDEN: access denied to account
- 404 ACCOUNT_NOT_FOUND: source or destination account not found
- 422 INSUFFICIENT_FUNDS: source account has insufficient funds
- 422 ACCOUNT_INACTIVE: one account is inactive
- 500 INTERNAL_ERROR: unexpected internal error

## 5.5 Get Statement

- Method: GET
- Path: /accounts/{id}/statement
- Auth required: yes

Query params (optional):
- limit: integer, default 50, max 100
- cursor: RFC3339 datetime
- cursor_id: UUID
- from: RFC3339 datetime
- to: RFC3339 datetime

Notes:
- cursor and cursor_id must be provided together
- items are returned in descending order by created_at and id

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
      }
    ],
    "next_cursor": null
  },
  "error": null
}
```

When there are more results, `next_cursor` is an object — pass both fields as query params for the next page:

```json
{
  "next_cursor": {
    "created_at": "2026-04-02T11:59:00Z",
    "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6"
  }
}
```

Possible errors:
- 401 UNAUTHORIZED: authentication required
- 401 INVALID_TOKEN: token invalid, malformed, or expired
- 400 INVALID_DATA: invalid path/query value or cursor/cursor_id mismatch
- 403 FORBIDDEN: access denied to account
- 404 ACCOUNT_NOT_FOUND: account does not exist
- 500 INTERNAL_ERROR: unexpected internal error

## 6. Error Code Reference

Common error codes currently used by handlers:
- INVALID_REQUEST
- INVALID_DATA
- INVALID_AMOUNT
- USER_ALREADY_EXISTS
- CUSTOMER_NOT_FOUND
- INVALID_CREDENTIALS
- UNAUTHORIZED
- INVALID_TOKEN
- FORBIDDEN
- ACCOUNT_NOT_FOUND
- ACCOUNT_INACTIVE
- INSUFFICIENT_FUNDS
- SAME_ACCOUNT_TRANSFER
- INTERNAL_ERROR

## 7. Domain Notes for API Consumers

- Monetary values are represented as integer cents
- UUID is used for all resource identifiers
- Financial operations are synchronous and strongly consistent
- Transfer operation is atomic: debit and credit are committed together
