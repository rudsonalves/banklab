# REST API Documentation - Bank API

## Table of Contents

- [REST API Documentation - Bank API](#rest-api-documentation---bank-api)
  - [Table of Contents](#table-of-contents)
  - [1. Overview](#1-overview)
  - [2. Response Envelope](#2-response-envelope)
    - [2.1 Error Payload Examples (Standard)](#21-error-payload-examples-standard)
  - [3. Authentication Endpoints](#3-authentication-endpoints)
    - [3.1 Register User](#31-register-user)
    - [3.2 Login User](#32-login-user)
    - [3.3 Refresh Access Token](#33-refresh-access-token)
    - [3.4 Get Current User](#34-get-current-user)
    - [3.5 Approve User (Admin Only)](#35-approve-user-admin-only)
  - [4. Account Endpoints](#4-account-endpoints)
    - [4.1 List Accounts](#41-list-accounts)
    - [4.2 Create Account](#42-create-account)
    - [4.3 Deposit](#43-deposit)
    - [4.4 Withdraw](#44-withdraw)
    - [4.5 Internal Transfer Recipient Lookup](#45-internal-transfer-recipient-lookup)
    - [4.6 Internal Transfer](#46-internal-transfer)
    - [4.7 Transfer Receipt](#47-transfer-receipt)
    - [4.8 Get Balance](#48-get-balance)
    - [4.9 Get Statement](#49-get-statement)
  - [5. Customer Endpoints](#5-customer-endpoints)
    - [5.1 Get My Customer Profile](#51-get-my-customer-profile)
  - [6. Authorization Model](#6-authorization-model)
  - [7. Error Code Reference](#7-error-code-reference)
  - [8. Domain Notes for API Consumers](#8-domain-notes-for-api-consumers)
  - [9. Error Scenarios by Endpoint (with Payload)](#9-error-scenarios-by-endpoint-with-payload)
    - [9.1 POST /auth/register](#91-post-authregister)
    - [9.2 POST /auth/login](#92-post-authlogin)
    - [9.3 POST /auth/refresh](#93-post-authrefresh)
    - [9.4 GET /auth/me](#94-get-authme)
    - [9.5 GET /accounts](#95-get-accounts)
    - [9.6 POST /accounts](#96-post-accounts)
    - [9.7 POST /accounts/{id}/deposit](#97-post-accountsiddeposit)
    - [9.8 POST /accounts/{id}/withdraw](#98-post-accountsidwithdraw)
    - [9.9 GET /accounts/internal-transfers/recipients](#99-get-accountsinternal-transfersrecipients)
    - [9.10 POST /accounts/internal-transfers](#910-post-accountsinternal-transfers)
    - [9.11 GET /accounts/transfer/{transaction_reference}/receipt](#911-get-accountstransfertransaction_referencereceipt)
    - [9.12 GET /accounts/{id}/balance](#912-get-accountsidbalance)
    - [9.13 GET /accounts/{id}/statement](#913-get-accountsidstatement)
    - [9.14 GET /customers/me](#914-get-customersme)
  - [10. Postman Setup](#10-postman-setup)
    - [10.1 Files in Repository](#101-files-in-repository)
    - [10.2 Environment Variables](#102-environment-variables)
    - [10.3 How to Import and Configure](#103-how-to-import-and-configure)
    - [10.4 Recommended Execution Flow](#104-recommended-execution-flow)

## 1. Overview

This document describes the HTTP REST contract currently implemented by the service.

Base URL (local):
- http://localhost:8080

Content type:
- request: application/json
- response: application/json

Authentication:
- `POST /auth/register` and `POST /auth/login` require header `X-App-Token: <app_token>`
- `POST /auth/refresh`, `GET /auth/me`, all `/accounts` and `/accounts/*`, and all `/customers/*` require JWT Bearer token
- Send JWT in header `Authorization: Bearer <access_token>`

Access control summary:
- Auth entry routes: AppToken only
- Auth session and service routes: JWT only

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

### 2.1 Error Payload Examples (Standard)

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

### 3.1 Register User

- Method: POST
- Path: /auth/register
- Auth required: AppToken (`X-App-Token`)

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
- 401 INVALID_APP_TOKEN: missing or invalid `X-App-Token`
- 400 INVALID_REQUEST: invalid JSON body
- 400 INVALID_DATA: invalid email or password format
- 409 USER_ALREADY_EXISTS: duplicate email
- 409 (customer domain): duplicate CPF or email in customers table
- 500 INTERNAL_ERROR: unexpected internal error

### 3.2 Login User

- Method: POST
- Path: /auth/login
- Auth required: AppToken (`X-App-Token`)

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
    "refresh_token": "<opaque-token>",
    "user_id": "d3de5f8b-4892-42e8-9680-979cf3f37844",
    "email": "user@example.com",
    "role": "customer",
    "customer_id": "6f3ebf86-bf82-4b75-a2ce-cd261ca47ec3"
  },
  "error": null
}
```

`customer_id` is always populated for users with role `customer`. The JWT embeds this value for use in subsequent requests.

Every login issues a new refresh token and persists a corresponding server-side session. The refresh token is required to obtain a new access token via `POST /auth/refresh`.

Possible errors:
- 401 INVALID_APP_TOKEN: missing or invalid `X-App-Token`
- 400 INVALID_REQUEST: invalid JSON body or unknown fields
- 400 INVALID_DATA: invalid email or password input
- 401 INVALID_CREDENTIALS: invalid email/password
- 500 INTERNAL_ERROR: unexpected internal error

### 3.3 Refresh Access Token

- Method: POST
- Path: /auth/refresh
- Auth required: JWT Bearer token

Exchanges a valid refresh token for a new access token and a new refresh token (token rotation). Each refresh token is single-use — after a successful refresh the old token is immediately revoked and a new session is created atomically.

Request body:

```json
{
  "refresh_token": "<opaque-token>"
}
```

The `refresh_token` field is required and must not be blank. Unknown fields are rejected.

Success response (200):

```json
{
  "data": {
    "access_token": "<new-jwt>",
    "refresh_token": "<new-opaque-token>"
  },
  "error": null
}
```

Behaviour:
- The old refresh token is revoked and the new token is persisted in a single database transaction. If either step fails the entire rotation is rolled back and the original token remains valid.
- A refresh token that has been revoked, is expired, or does not correspond to any session returns `401 INVALID_TOKEN`.

Possible errors:
- 401 UNAUTHORIZED: authentication required
- 400 INVALID_REQUEST: missing or blank `refresh_token`, invalid JSON, or unknown fields
- 401 INVALID_TOKEN: token invalid, revoked, expired, or not found
- 500 INTERNAL_ERROR: unexpected internal error

### 3.4 Get Current User

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

### 3.5 Approve User (Admin Only)

- Method: POST
- Path: /admin/users/{id}/approve
- Auth required: JWT (admin role)

Approves a pending user, transitioning them from `pending` to `active` status. Also creates the associated account atomically.

Path parameters:

- `id`: UUID of the user to approve

Request body:

Empty object or no body required.

```json
{}
```

Success response (200):

```json
{
  "data": {
    "user_id": "d3de5f8b-4892-42e8-9680-979cf3f37844",
    "status": "active",
    "account_id": "a1b2c3d4-e5f6-4789-a012-b3c4d5e6f789"
  },
  "error": null
}
```

Response fields:

- `user_id`: UUID of the approved user
- `status`: new status (always "active" on success)
- `account_id`: UUID of the newly created account

Atomicity:

- User status update and account creation occur within a single database transaction
- If account creation fails, user status is not updated
- No partial state is possible

Possible errors:

- 401 UNAUTHORIZED: authentication required
- 401 INVALID_TOKEN: token invalid or expired
- 403 FORBIDDEN: authenticated user does not have admin role
- 404 USER_NOT_FOUND: user does not exist
- 404 CUSTOMER_NOT_FOUND: associated customer does not exist
- 409 USER_ALREADY_ACTIVE: user is already active (cannot approve active/blocked users)
- 500 INTERNAL_ERROR: unexpected internal error

## 4. Account Endpoints

All account routes are protected and require Authorization header with Bearer token.

Ownership is enforced automatically. A customer-role user can only access accounts that belong to their own `customer_id`. Admin-role users can access account-scoped operations, but `GET /accounts` is customer-context scoped and requires a non-nil `customer_id` in the authenticated principal.

### 4.1 List Accounts

- Method: GET
- Path: /accounts
- Auth required: yes

Returns the list of accounts that belong to the authenticated user.

Query params:
- none (any query param returns `400 INVALID_REQUEST`)

Success response (200):

```json
{
  "data": [
    {
      "id": "fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b",
      "customer_id": "6f3ebf86-bf82-4b75-a2ce-cd261ca47ec3",
      "number": "10000001",
      "branch": "0001",
      "status": "active"
    }
  ],
  "error": null
}
```

Notes:
- The account list is derived from the authenticated user's `customer_id`
- `balance` is intentionally omitted from this endpoint

Possible errors:
- 401 UNAUTHORIZED: authentication required
- 401 INVALID_TOKEN: token invalid, malformed, or expired
- 400 INVALID_REQUEST: unexpected query params
- 403 FORBIDDEN: authenticated user has no customer context
- 500 INTERNAL_ERROR: unexpected internal error

### 4.2 Create Account

- Method: POST
- Path: /accounts
- Auth required: yes

Creates a new account for the authenticated user. The user **must have status = active** to create an account.

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
- 403 FORBIDDEN: user is not active or access denied
- 404 CUSTOMER_NOT_FOUND: customer does not exist
- 500 INTERNAL_ERROR: unexpected internal error

### 4.3 Deposit

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
- 403 FORBIDDEN: access denied
- 404 ACCOUNT_NOT_FOUND: account does not exist
- 422 ACCOUNT_INACTIVE: account not active
- 500 INTERNAL_ERROR: unexpected internal error

### 4.4 Withdraw

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
- 403 FORBIDDEN: access denied
- 404 ACCOUNT_NOT_FOUND: account does not exist
- 422 INSUFFICIENT_FUNDS: insufficient funds
- 422 ACCOUNT_INACTIVE: account not active
- 500 INTERNAL_ERROR: unexpected internal error

### 4.5 Internal Transfer Recipient Lookup

- Method: GET
- Path: /accounts/internal-transfers/recipients
- Auth required: yes

This endpoint searches eligible recipient accounts for internal transfers only.
It is not a Pix, TED, DOC, or interbank account discovery endpoint.

Query modes:

1. By branch and account number:

```text
/accounts/internal-transfers/recipients?branch=0001&account_number=00067890
```

2. By CPF document:

```text
/accounts/internal-transfers/recipients?document=12345678901
```

Rules:
- `branch` and `account_number` must be provided together.
- `document` may be used alone and currently accepts CPF only.
- CNPJ and legal person account lookup are intentionally out of scope for this
  version.
- Requests with neither query mode are invalid.
- Mixed query modes are invalid unless explicitly supported by a future version.
- Branch, account number, and CPF are normalized before lookup.

Success response (200):

```json
{
  "data": {
    "accounts": [
      {
        "account_id": "fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b",
        "holder_name": "Maria Silva",
        "document": "***.456.789-**",
        "branch": "0001",
        "account_number": "00067890"
      }
    ]
  },
  "error": null
}
```

No results response (200):

```json
{
  "data": {
    "accounts": []
  },
  "error": null
}
```

Response rules:
- `account_id` is the account identifier used later as `to_account_id`.
- `document` is always masked.
- Lookup by branch + account number returns zero or one eligible account.
- Lookup by CPF may return zero, one, or many eligible accounts.
- If multiple accounts are returned, the client must require user selection.
- The response must not include customer ID, balance, full document, phone,
  e-mail, or unrelated internal fields.

Possible errors:
- 401 UNAUTHORIZED: authentication required
- 401 INVALID_TOKEN: token invalid, malformed, or expired
- 400 INVALID_DATA: invalid or unsupported query parameter combination
- 403 FORBIDDEN: access denied
- 500 INTERNAL_ERROR: unexpected internal error

### 4.6 Internal Transfer

- Method: POST
- Path: /accounts/internal-transfers
- Auth required: yes

Executes an internal transfer between two account IDs previously known by the
client. The destination account ID should come from recipient lookup or another
trusted internal account selection flow.

Request body:

```json
{
  "from_account_id": "fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b",
  "to_account_id": "7c8e75d0-60c9-4b2d-a9dc-605293d98e0c",
  "amount": 2500,
  "idempotency_key": "transfer-client-key",
  "description": "Aluguel de maio"
}
```

Notes:
- This endpoint is for internal transfers only.
- Transfer execution uses account IDs, not branch/account-number identifiers.
- The backend validates that the authenticated user can operate
  `from_account_id`.
- The backend validates that `to_account_id` exists and can receive internal
  transfers.
- `idempotency_key` is required and must be stable across retries of the same
  transfer attempt.
- `description` is optional. When omitted or blank, no description is stored.
- Idempotency is scoped to `from_account_id` and `idempotency_key`.
- Different source accounts may reuse the same `idempotency_key` independently.
- Replay responses return the historical transfer result from ledger data (not current account balances).
- Replay responses preserve the original transfer description; a different
  `description` in a retry does not overwrite the stored value.

Success response (200):

```json
{
  "data": {
    "from_account_id": "fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b",
    "transaction_reference": "2e3ef0c7-ef10-4f4e-a62b-56c71c3c5b31",
    "to_account_id": "7c8e75d0-60c9-4b2d-a9dc-605293d98e0c",
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
- 400 INVALID_DATA: missing or malformed account IDs, or missing
  `idempotency_key`
- 400 INVALID_AMOUNT: amount must be greater than zero
- 400 SAME_ACCOUNT_TRANSFER: source and destination are equal
- 403 FORBIDDEN: authenticated user cannot operate the source account
- 404 ACCOUNT_NOT_FOUND: source or destination account not found
- 422 INSUFFICIENT_FUNDS: source account has insufficient funds
- 422 ACCOUNT_INACTIVE: one account is inactive
- 500 INTERNAL_ERROR: unexpected internal error

### 4.7 Transfer Receipt

- Method: GET
- Path: /accounts/transfer/{transaction_reference}/receipt
- Auth required: yes

Returns persisted receipt details for an internal transfer.

Path params:
- `transaction_reference`: public transfer reference returned by `POST /accounts/internal-transfers`

Success response (200):

```json
{
  "data": {
    "operation_type": "transfer_out",
    "amount": 2500,
    "status": "completed",
    "transaction_reference": "2e3ef0c7-ef10-4f4e-a62b-56c71c3c5b31",
    "operation_date": "2026-05-06T12:30:00Z",
    "source_branch": "0001",
    "source_account_number": "00012345",
    "destination_branch": "0001",
    "destination_account_number": "00067890",
    "recipient_name": "Maria Silva",
    "description": "Aluguel de maio"
  },
  "error": null
}
```

Notes:
- `description` is omitted when the original transfer did not include one.
- Receipt data remains public/confirmation-oriented and does not expose internal
  customer IDs, balances, phone, e-mail, or full documents.

Status semantics (`data.status`):
- `completed`: transfer executed successfully (terminal success)
- `pending`: transfer accepted and still processing (intermediate)
- `failed`: transfer failed due to technical/system error (terminal failure)
- `cancelled`: transfer cancelled before completion (terminal failure)
- `rejected`: transfer rejected by validation/business rules (terminal failure)

Current behavior:
- The current backend implementation returns `completed` for persisted receipts.
- Additional status values are documented here for upcoming backend evolution.

Possible errors:
- 401 UNAUTHORIZED: authentication required
- 401 INVALID_TOKEN: token invalid, malformed, or expired
- 400 INVALID_DATA: missing or malformed transaction reference
- 403 FORBIDDEN: authenticated user cannot access this receipt
- 404 TRANSACTION_NOT_FOUND: transfer receipt does not exist
- 500 INTERNAL_ERROR: unexpected internal error

### 4.8 Get Balance

- Method: GET
- Path: /accounts/{id}/balance
- Auth required: yes

Query params:
- none (any query param returns `400 INVALID_DATA`)

Success response (200):

```json
{
  "data": {
    "account_id": "fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b",
    "balance": 12000
  },
  "error": null
}
```

Notes:
- `balance` is returned in cents

Possible errors:
- 401 UNAUTHORIZED: authentication required
- 401 INVALID_TOKEN: token invalid, malformed, or expired
- 400 INVALID_DATA: invalid account id or unexpected query params
- 403 FORBIDDEN: access denied
- 404 ACCOUNT_NOT_FOUND: account does not exist
- 500 INTERNAL_ERROR: unexpected internal error

### 4.9 Get Statement

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
- 403 FORBIDDEN: access denied
- 404 ACCOUNT_NOT_FOUND: account does not exist
- 500 INTERNAL_ERROR: unexpected internal error

## 5. Customer Endpoints

All customer routes are protected and require Authorization header with Bearer token.

### 5.1 Get My Customer Profile

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
- A user with role `admin` can access account/customer scoped resources when identifiers are provided
- `GET /accounts` is scoped to the authenticated principal's `customer_id` and returns `403 FORBIDDEN` when `customer_id` is absent
- The `customer_id` is never accepted from the client — it is always read from the JWT token
- Cross-customer access returns `403 FORBIDDEN`
- Any operation where the user has no `customer_id` returns `409 INVALID_USER_STATE`

This rule is enforced in the application layer via `CanAccessAccount` and `CanAccessCustomer` helpers, not in HTTP handlers.

## 7. Error Code Reference

Common error codes currently used by handlers:
- INVALID_APP_TOKEN
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
- TRANSACTION_NOT_FOUND
- INTERNAL_ERROR

`INVALID_APP_TOKEN` (HTTP 401) is returned when `POST /auth/register` or `POST /auth/login` is called without `X-App-Token` or with an invalid app token.

`INVALID_TOKEN` (HTTP 401) is returned for any of the following conditions on the `/auth/refresh` endpoint: token not found, already revoked, expired, or JWT signature invalid.

`INVALID_USER_STATE` (HTTP 409) indicates the system detected an invariant violation: a user with role `customer` has no linked `customer_id`. This should never occur under normal operation; it signals a data consistency bug.

## 8. Domain Notes for API Consumers

- Monetary values are represented as integer cents
- UUID is used for all resource identifiers
- Financial operations are synchronous and strongly consistent
- Transfer operation is atomic: debit and credit are committed together

## 9. Error Scenarios by Endpoint (with Payload)

This section lists common error situations and the expected payload shape.

### 9.1 POST /auth/register

Scenario: missing or invalid app token
- Status: 401
- Code: INVALID_APP_TOKEN

```json
{
  "data": null,
  "error": {
    "code": "INVALID_APP_TOKEN",
    "message": "invalid application token"
  }
}
```

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

Scenario: missing or invalid app token
- Status: 401
- Code: INVALID_APP_TOKEN

```json
{
  "data": null,
  "error": {
    "code": "INVALID_APP_TOKEN",
    "message": "invalid application token"
  }
}
```

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

### 9.3 POST /auth/refresh

Scenario: missing/invalid JWT authentication
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

### 9.4 GET /auth/me

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

### 9.5 GET /accounts

Scenario: authenticated user has no customer context
- Status: 403
- Code: FORBIDDEN

```json
{
  "data": null,
  "error": {
    "code": "FORBIDDEN",
    "message": "Access denied"
  }
}
```

Scenario: unexpected query params
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

### 9.6 POST /accounts

Scenario: authenticated user cannot create account for requested context
- Status: 403
- Code: FORBIDDEN

```json
{
  "data": null,
  "error": {
    "code": "FORBIDDEN",
    "message": "Access denied"
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

### 9.7 POST /accounts/{id}/deposit

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

### 9.8 POST /accounts/{id}/withdraw

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

### 9.9 GET /accounts/internal-transfers/recipients

Scenario: invalid query combination
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

Scenario: no recipients found
- Status: 200

```json
{
  "data": {
    "accounts": []
  },
  "error": null
}
```

Scenario: access denied
- Status: 403
- Code: FORBIDDEN

```json
{
  "data": null,
  "error": {
    "code": "FORBIDDEN",
    "message": "Access denied"
  }
}
```

### 9.10 POST /accounts/internal-transfers

Scenario: invalid JSON or unknown public request field
- Status: 400
- Code: INVALID_REQUEST

```json
{
  "data": null,
  "error": {
    "code": "INVALID_REQUEST",
    "message": "Invalid request"
  }
}
```

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

Scenario: missing or malformed account IDs or idempotency key
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
    "message": "Access denied"
  }
}
```

Scenario: source or destination account not found
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

Scenario: inactive account
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

### 9.11 GET /accounts/transfer/{transaction_reference}/receipt

Scenario: receipt found
- Status: 200
- `data.status` currently returned as `completed`
- Future-compatible values for `data.status`: `completed`, `pending`, `failed`, `cancelled`, `rejected`

Scenario: malformed transaction reference
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

Scenario: receipt not found
- Status: 404
- Code: TRANSACTION_NOT_FOUND

```json
{
  "data": null,
  "error": {
    "code": "TRANSACTION_NOT_FOUND",
    "message": "Transaction not found"
  }
}
```

Scenario: access denied to receipt
- Status: 403
- Code: FORBIDDEN

```json
{
  "data": null,
  "error": {
    "code": "FORBIDDEN",
    "message": "Access denied"
  }
}
```

### 9.12 GET /accounts/{id}/balance

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

### 9.13 GET /accounts/{id}/statement

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

### 9.14 GET /customers/me

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

## 10. Postman Setup

The repository includes a ready-to-use Postman collection and environment under `tools/postman`.

### 10.1 Files in Repository

- `tools/postman/Banklab_API.postman_collection.json`
- `tools/postman/Environment.postman_environment.json`
- `tools/postman/README.md`

### 10.2 Environment Variables

Use these variables when configuring the Postman environment:

- `base_url`: API base URL (default: `http://localhost:8080`)
- `app_token`: application token used by auth entry routes (`/auth/register` and `/auth/login`)
- `access_token`: JWT used for protected routes
- `refresh_token`: opaque refresh token used by `/auth/refresh`
- `account_id`: account UUID for account operations
- `from_account_id`: source account UUID used by internal transfer requests
- `to_account_id`: destination account UUID used by internal transfer requests
- `recipient_branch`: branch used for recipient lookup examples
- `recipient_account_number`: account number used for recipient lookup examples
- `recipient_document`: CPF used for recipient lookup examples
- `transaction_reference`: public reference returned by successful transfer requests
- `id`: user UUID used by admin approval route (`/admin/users/{id}/approve`)

### 10.3 How to Import and Configure

1. Import `tools/postman/Banklab_API.postman_collection.json` into Postman.
2. Import `tools/postman/Environment.postman_environment.json` into Postman.
3. Select the imported environment in Postman.
4. Adjust `base_url` if your API is not running on `http://localhost:8080`.
5. Confirm `app_token` matches the value configured in your local API environment.
6. Run auth requests to obtain tokens and update `access_token` / `refresh_token`.

### 10.4 Recommended Execution Flow

Use this flow to bootstrap test data and credentials quickly:

1. `Auth/Register`
2. `Auth/Login` (copy `access_token` and `refresh_token` from response)
3. `Auth/Me` (validate JWT)
4. `Account/User/Approve` (admin only, using `id`)
5. Account endpoints using `account_id` as needed
6. Recipient lookup endpoint using `recipient_branch` + `recipient_account_number`, or `recipient_document`
7. Internal transfer endpoint using `from_account_id` and `to_account_id`
8. Receipt endpoint using `transaction_reference` from a successful transfer response

Notes:
- Keep `X-App-Token` for register/login requests as documented in this file.
- For protected routes, send `Authorization: Bearer <access_token>`.
- If a request returns `401 INVALID_TOKEN`, run `Auth/Refresh` and update tokens.
