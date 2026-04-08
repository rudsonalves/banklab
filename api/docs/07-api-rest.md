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
    "message": "human readable message"
  }
}
```

Notes:
- Current implementation returns `error.code` and `error.message`.
- `error.details` is not currently populated by handlers.

## 2.1 Error Payload Examples (Standard)

Example - 400 INVALID_REQUEST:

```json
{
  "data": null,
  "error": {
    "code": "INVALID_REQUEST",
    "message": "Invalid request body"
  }
}
```

Example - 401 UNAUTHORIZED:

```json
{
  "data": null,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Authentication required"
  }
}
```

Example - 401 INVALID_TOKEN:

```json
{
  "data": null,
  "error": {
    "code": "INVALID_TOKEN",
    "message": "Invalid token"
  }
}
```

Example - 500 INTERNAL_ERROR:

```json
{
  "data": null,
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "Internal server error"
  }
}
```

## 3. Authentication Endpoints

## 3.1 Register User

- Method: POST
- Path: /auth/register
- Auth required: no

This endpoint creates a User and an associated Customer atomically in a single transaction. The Customer is created automatically — the client never needs to call a separate customer creation endpoint.

Request body:

```json
{
  "email": "user@example.com",
  "password": "P@ssword123",
  "name": "Maria Silva",
  "cpf": "12345678901"
}
```

All four fields are required.

Success response (201):

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

`customer_id` is always populated for users with role `customer`.

Possible errors:
- 400 INVALID_REQUEST: invalid JSON body
- 400 INVALID_DATA: invalid email or password format
- 409 USER_ALREADY_EXISTS: duplicate email
- 409 (customer domain): duplicate CPF or email in customers table
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
    "customer_id": "6f3ebf86-bf82-4b75-a2ce-cd261ca47ec3"
  },
  "error": null
}
```

`customer_id` is always populated for users with role `customer`. The JWT embeds this value for use in subsequent requests.

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

## 4. Account Endpoints

All account routes are protected and require Authorization header with Bearer token.

Ownership is enforced automatically. A customer-role user can only access accounts that belong to their own `customer_id`. Admin-role users can access any account.

## 4.1 Create Account

- Method: POST
- Path: /accounts
- Auth required: yes

The `customer_id` is derived automatically from the authenticated user's JWT token. The client MUST NOT send a `customer_id` in the request body.

Request body:

```json
{}
```

Body can also be empty. Any extra fields are rejected (400 INVALID_REQUEST).

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
- 403 FORBIDDEN: access denied to account
- 404 CUSTOMER_NOT_FOUND: customer does not exist
- 500 INTERNAL_ERROR: unexpected internal error

## 4.2 Deposit

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

## 4.3 Withdraw

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

## 4.4 Transfer

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

## 4.5 Get Statement

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

## 5. Customer Endpoints

All customer routes are protected and require Authorization header with Bearer token.

## 5.1 Get My Customer Profile

- Method: GET
- Path: /customers/me
- Auth required: yes

Returns the customer profile linked to the authenticated user. No path or query parameters required.

Success response (200):

```json
{
  "data": {
    "id": "6f3ebf86-bf82-4b75-a2ce-cd261ca47ec3",
    "name": "Maria Silva",
    "cpf": "12345678901",
    "email": "user@example.com",
    "created_at": "2026-04-07T10:00:00Z"
  },
  "error": null
}
```

Possible errors:
- 401 UNAUTHORIZED: authentication required
- 401 INVALID_TOKEN: token invalid, malformed, or expired
- 409 INVALID_USER_STATE: authenticated user has no associated customer (inconsistent state)
- 404 CUSTOMER_NOT_FOUND: customer record not found
- 500 INTERNAL_ERROR: unexpected internal error

## 6. Authorization Model

All account and customer operations enforce ownership based on the authenticated user's context.

Rules:
- A user with role `customer` can only access resources where `resource.customer_id == user.customer_id`
- A user with role `admin` can access any resource
- The `customer_id` is never accepted from the client — it is always read from the JWT token
- Cross-customer access returns `403 FORBIDDEN`
- Any operation where the user has no `customer_id` returns `409 INVALID_USER_STATE`

This rule is enforced in the application layer via `CanAccessAccount` and `CanAccessCustomer` helpers, not in HTTP handlers.

## 7. Error Code Reference

Common error codes currently used by handlers:
- INVALID_REQUEST
- INVALID_DATA
- INVALID_AMOUNT
- INVALID_USER_STATE
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

`INVALID_USER_STATE` (HTTP 409) indicates the system detected an invariant violation: a user with role `customer` has no linked `customer_id`. This should never occur under normal operation; it signals a data consistency bug.

## 8. Domain Notes for API Consumers

- Monetary values are represented as integer cents
- UUID is used for all resource identifiers
- Financial operations are synchronous and strongly consistent
- Transfer operation is atomic: debit and credit are committed together

## 9. Error Scenarios by Endpoint (with Payload)

This section lists common error situations and the expected payload shape.

### 9.1 POST /auth/register

Scenario: malformed JSON
- Status: 400
- Code: INVALID_REQUEST

```json
{
  "data": null,
  "error": {
    "code": "INVALID_REQUEST",
    "message": "Invalid request body"
  }
}
```

Scenario: duplicate email/CPF
- Status: 409
- Code: USER_ALREADY_EXISTS

```json
{
  "data": null,
  "error": {
    "code": "USER_ALREADY_EXISTS",
    "message": "User already exists"
  }
}
```

### 9.2 POST /auth/login

Scenario: invalid credentials
- Status: 401
- Code: INVALID_CREDENTIALS

```json
{
  "data": null,
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Invalid credentials"
  }
}
```

### 9.3 GET /auth/me

Scenario: missing/invalid authentication
- Status: 401
- Code: UNAUTHORIZED or INVALID_TOKEN

```json
{
  "data": null,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Authentication required"
  }
}
```

### 9.4 POST /accounts

Scenario: authenticated user cannot create account for requested context
- Status: 403
- Code: FORBIDDEN

```json
{
  "data": null,
  "error": {
    "code": "FORBIDDEN",
    "message": "Access denied to account"
  }
}
```

Scenario: customer does not exist
- Status: 404
- Code: CUSTOMER_NOT_FOUND

```json
{
  "data": null,
  "error": {
    "code": "CUSTOMER_NOT_FOUND",
    "message": "Customer not found"
  }
}
```

### 9.5 POST /accounts/{id}/deposit

Scenario: invalid amount
- Status: 400
- Code: INVALID_AMOUNT

```json
{
  "data": null,
  "error": {
    "code": "INVALID_AMOUNT",
    "message": "Invalid amount"
  }
}
```

Scenario: account not found
- Status: 404
- Code: ACCOUNT_NOT_FOUND

```json
{
  "data": null,
  "error": {
    "code": "ACCOUNT_NOT_FOUND",
    "message": "Account not found"
  }
}
```

Scenario: account inactive
- Status: 422
- Code: ACCOUNT_INACTIVE

```json
{
  "data": null,
  "error": {
    "code": "ACCOUNT_INACTIVE",
    "message": "Account is not active"
  }
}
```

### 9.6 POST /accounts/{id}/withdraw

Scenario: insufficient funds
- Status: 422
- Code: INSUFFICIENT_FUNDS

```json
{
  "data": null,
  "error": {
    "code": "INSUFFICIENT_FUNDS",
    "message": "Insufficient balance"
  }
}
```

### 9.7 POST /accounts/transfer

Scenario: source and destination are the same
- Status: 400
- Code: SAME_ACCOUNT_TRANSFER

```json
{
  "data": null,
  "error": {
    "code": "SAME_ACCOUNT_TRANSFER",
    "message": "Source and destination accounts must be different"
  }
}
```

Scenario: access denied to source account
- Status: 403
- Code: FORBIDDEN

```json
{
  "data": null,
  "error": {
    "code": "FORBIDDEN",
    "message": "Access denied to account"
  }
}
```

### 9.8 GET /accounts/{id}/statement

Scenario: invalid query/path data
- Status: 400
- Code: INVALID_DATA

```json
{
  "data": null,
  "error": {
    "code": "INVALID_DATA",
    "message": "Invalid data"
  }
}
```

### 9.9 GET /customers/me

Scenario: user has inconsistent state (customer role without customer_id)
- Status: 409
- Code: INVALID_USER_STATE

```json
{
  "data": null,
  "error": {
    "code": "INVALID_USER_STATE",
    "message": "Invalid user state"
  }
}
```

Scenario: customer not found
- Status: 404
- Code: CUSTOMER_NOT_FOUND

```json
{
  "data": null,
  "error": {
    "code": "CUSTOMER_NOT_FOUND",
    "message": "Customer not found"
  }
}
```
