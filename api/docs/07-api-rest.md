# REST API Documentation - Bank API

## Table of Contents

- [REST API Documentation - Bank API](#rest-api-documentation---bank-api)
  - [Table of Contents](#table-of-contents)
  - [1. Overview](#1-overview)
  - [2. Response Envelope](#2-response-envelope)
    - [2.1 Error Payload Examples (Standard)](#21-error-payload-examples-standard)
  - [3. Authentication Endpoints](#3-authentication-endpoints)
    - [3.0 Check CPF](#30-check-cpf)
    - [3.0.1 Request Contact Verification](#301-request-contact-verification)
    - [3.0.2 Confirm Contact Verification](#302-confirm-contact-verification)
    - [3.1 Register User](#31-register-user)
    - [3.2 Login User](#32-login-user)
    - [3.3 Refresh Access Token](#33-refresh-access-token)
    - [3.4 Get Auth Session](#34-get-auth-session)
    - [3.4.1 Get Current User](#341-get-current-user)
    - [3.5 Approve User (Admin Only)](#35-approve-user-admin-only)
    - [3.6 Create Customer Account (Admin Only)](#36-create-customer-account-admin-only)
    - [3.7 Create Transaction Password](#37-create-transaction-password)
    - [3.8 Authorize Step-Up](#38-authorize-step-up)
    - [3.9 Register Installation](#39-register-installation)
    - [3.10 List Installations](#310-list-installations)
    - [3.11 Revoke Installation](#311-revoke-installation)
  - [4. Account Endpoints](#4-account-endpoints)
    - [4.1 List Accounts](#41-list-accounts)
    - [4.2 Customer Account Creation Removed](#42-customer-account-creation-removed)
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
    - [9.0 POST /auth/cpf-check](#90-post-authcpf-check)
    - [9.1 POST /auth/register](#91-post-authregister)
    - [9.2 POST /auth/login](#92-post-authlogin)
    - [9.3 POST /auth/refresh](#93-post-authrefresh)
    - [9.4 GET /auth/session](#94-get-authsession)
    - [9.4.1 GET /auth/me](#941-get-authme)
    - [9.5 GET /accounts](#95-get-accounts)
    - [9.6 POST /admin/customers/{customer\_id}/accounts](#96-post-admincustomerscustomer_idaccounts)
    - [9.7 POST /terminal/accounts/{id}/deposit](#97-post-terminalaccountsiddeposit)
    - [9.8 POST /terminal/accounts/{id}/withdraw](#98-post-terminalaccountsidwithdraw)
    - [9.9 GET /accounts/internal-transfers/recipients](#99-get-accountsinternal-transfersrecipients)
    - [9.10 POST /accounts/internal-transfers](#910-post-accountsinternal-transfers)
    - [9.11 GET /accounts/transfer/{transaction\_reference}/receipt](#911-get-accountstransfertransaction_referencereceipt)
    - [9.12 GET /accounts/{id}/balance](#912-get-accountsidbalance)
    - [9.13 GET /accounts/{id}/statement](#913-get-accountsidstatement)
    - [9.14 GET /customers/me](#914-get-customersme)
    - [9.15 POST /security/transaction-password](#915-post-securitytransaction-password)
  - [10. Bruno Setup](#10-bruno-setup)
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
- `POST /auth/cpf-check`, `POST /auth/contact-verifications`, `POST /auth/contact-verifications/confirm`, `POST /auth/register`, and `POST /auth/login` require header `X-App-Token: <app_token>`
- `POST /auth/login` also requires header `X-Installation-Id: <canonical_uuid_v4>`
- `POST /auth/refresh` requires `X-Installation-Id` and a valid `refresh_token` in the request body
- `GET /auth/session`, `GET /auth/me`, all `/accounts` and `/accounts/*`, all `/customers/*`, and operational `/security/*` routes require JWT Bearer token plus `X-Installation-Id`
- `POST /security/installations` requires a restricted access token, `X-Installation-Id`, and a valid `X-Step-Up-Token`
- Send JWT in header `Authorization: Bearer <access_token>`

Access control summary:
- Auth entry routes: AppToken only
- Login entry route: AppToken plus installation identifier header
- Auth refresh route: refresh token plus installation identifier header
- Auth session and service routes: JWT plus installation identifier header
- Installation registration route: restricted access token plus installation identifier and step-up token

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

Notes:
- Clients should depend on `error.code`, not on `error.message`.
- `error.details` is optional and may be populated for selected errors, such as `CONTACT_NOT_VERIFIED`.

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

Before registration, the client must check CPF availability and request and
confirm contact verifications for both e-mail and phone.

### 3.0 Check CPF

- Method: POST
- Path: /auth/cpf-check
- Auth required: AppToken (`X-App-Token`)

Request body:

```json
{
  "cpf": "123.456.789-09"
}
```

Success response (200), available CPF:

```json
{
  "data": {
    "cpf": "12345678909",
    "exists": false,
    "available": true
  },
  "error": null
}
```

Success response (200), existing CPF:

```json
{
  "data": {
    "cpf": "12345678909",
    "exists": true,
    "available": false
  },
  "error": null
}
```

The CPF is normalized before lookup. Invalid CPF values return `400 INVALID_DATA`.

Possible errors:
- 401 INVALID_APP_TOKEN: missing or invalid `X-App-Token`
- 400 INVALID_REQUEST: invalid JSON body
- 400 INVALID_DATA: missing or invalid CPF
- 500 INTERNAL_ERROR: unexpected internal error

### 3.0.1 Request Contact Verification

- Method: POST
- Path: /auth/contact-verifications
- Auth required: AppToken (`X-App-Token`)

Request body:

```json
{
  "channel": "email",
  "target": "user@example.com"
}
```

`channel` accepts `email` or `phone`.

The `target` value is normalized before creating a verification challenge:
- `email`: trimmed and lowercased
- `phone`: trimmed

Success response (201):

```json
{
  "data": {
    "verification_id": "8d9ad65f-f837-4f6f-bd20-63f2c7cefab6",
    "channel": "email",
    "target": "user@example.com",
    "debug_token": "123456",
    "expires_at": "2026-05-18T12:10:00Z"
  },
  "error": null
}
```

`debug_token` is returned by the current implementation because the project does
not yet have an e-mail/SMS notification provider to deliver verification codes.
Clients may use it for local/manual testing and debug logs, but application flow
must not treat it as part of the stable verification contract. The stable data
for the request step is `verification_id`, `channel`, `target`, and `expires_at`.

Possible errors:
- 400 INVALID_DATA: invalid channel or empty target
- 409 USER_ALREADY_EXISTS: target e-mail or phone already belongs to a user
- 500 INTERNAL_ERROR: unexpected internal error

### 3.0.2 Confirm Contact Verification

- Method: POST
- Path: /auth/contact-verifications/confirm
- Auth required: AppToken (`X-App-Token`)

Request body:

```json
{
  "verification_id": "8d9ad65f-f837-4f6f-bd20-63f2c7cefab6",
  "token": "123456"
}
```

Success response (200):

```json
{
  "data": {
    "verification_token": "12adf6b7-2c5f-4895-96a3-a8e45db5c1d1",
    "channel": "email",
    "target": "user@example.com",
    "verified_at": "2026-05-18T12:03:00Z"
  },
  "error": null
}
```

The `verification_token` from this response is required by
`POST /auth/register`.

### 3.1 Register User

- Method: POST
- Path: /auth/register
- Auth required: AppToken (`X-App-Token`)

This endpoint creates a User and an associated Customer atomically in a single transaction. The Customer is created automatically — the client never needs to call a separate customer creation endpoint.

Request body:

```json
{
  "email": "user@example.com",
  "phone": "+5511999999999",
  "password": "P@ssword123",
  "name": "Maria Silva",
  "birth_date": "1990-01-15",
  "cpf": "12345678901",
  "email_verification_token": "12adf6b7-2c5f-4895-96a3-a8e45db5c1d1",
  "phone_verification_token": "95d3102d-c58f-4fa3-a1e8-2f0834bb9a39"
}
```

All fields are required.

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
- 409 USER_ALREADY_EXISTS: duplicate e-mail, phone, or CPF document
- 500 INTERNAL_ERROR: unexpected internal error

### 3.2 Login User

- Method: POST
- Path: /auth/login
- Auth required: AppToken (`X-App-Token`) and installation identifier (`X-Installation-Id`)

Request headers:

```http
X-App-Token: <app_token>
X-Installation-Id: <canonical_uuid_v4>
```

`X-Installation-Id` must be a UUID v4 in canonical lowercase form with hyphens,
for example `550e8400-e29b-41d4-a716-446655440000`. The API classifies the
installation after credential validation.

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

When the installation is new and the user still has an available installation
slot, login returns a restricted authorization instead of an operational
session:

```json
{
  "data": {
    "restricted_access_token": "<jwt>",
    "restricted_token_type": "restricted_access",
    "restricted_scope": "installation.register",
    "restricted_expires_at": "2026-06-17T10:05:00Z",
    "user_id": "d3de5f8b-4892-42e8-9680-979cf3f37844",
    "email": "user@example.com",
    "role": "customer",
    "customer_id": "6f3ebf86-bf82-4b75-a2ce-cd261ca47ec3"
  },
  "error": null
}
```

Restricted login responses do not include `refresh_token`. The client must
authorize step-up for `POST /security/installations` and then register the
installation.

`customer_id` is always populated for users with role `customer`. The JWT embeds this value for use in subsequent requests.

Every login issues a new refresh token and persists a corresponding server-side session. The refresh token is required to obtain a new access token via `POST /auth/refresh`.
Access token and refresh session lifetimes are configured through
`JWT_ACCESS_TOKEN_DURATION` and `JWT_REFRESH_TOKEN_DURATION`. Defaults are `15m`
and `168h` when those variables are omitted.

Customer users can complete login only after admin approval has provisioned at
least one account through `POST /admin/users/{id}/approve`. If approval or
account provisioning is still pending, login fails before token/session creation
with `ACCOUNT_APPROVAL_REQUIRED`.

Possible errors:
- 401 INVALID_APP_TOKEN: missing or invalid `X-App-Token`
- 400 INVALID_INSTALLATION_ID: missing or malformed `X-Installation-Id`
- 400 INVALID_REQUEST: invalid JSON body or unknown fields
- 400 INVALID_DATA: invalid email or password input
- 401 INVALID_CREDENTIALS: invalid email/password
- 403 CONTACT_NOT_VERIFIED: e-mail and/or phone not verified
- 403 ACCOUNT_APPROVAL_REQUIRED: customer user still requires admin approval or account provisioning
- 403 INSTALLATION_REVOKED: installation was revoked and cannot login
- 409 INSTALLATION_LIMIT_REACHED: user already has three known installations
- 500 INTERNAL_ERROR: unexpected internal error

### 3.3 Refresh Access Token

- Method: POST
- Path: /auth/refresh
- Auth required: installation identifier (`X-Installation-Id`) and refresh token in request body

Exchanges a valid refresh token for a new access token and a new refresh token (token rotation). Each refresh token is single-use — after a successful refresh the old token is immediately revoked and a new session is created atomically.

Request headers:

```http
X-Installation-Id: <canonical_uuid_v4>
```

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
- The `X-Installation-Id` header must match the installation bound to the refresh session.

Possible errors:
- 400 INVALID_INSTALLATION_ID: missing or malformed `X-Installation-Id`
- 400 INVALID_REQUEST: missing or blank `refresh_token`, invalid JSON, or unknown fields
- 401 INVALID_TOKEN: invalid, revoked, expired, or unknown refresh token
- 403 INSTALLATION_MISMATCH: header installation does not match the session
- 500 INTERNAL_ERROR: unexpected internal error

### 3.4 Get Auth Session

- Method: GET
- Path: /auth/session
- Auth required: JWT Bearer token

Returns the canonical authenticated session snapshot for clients after login.
This endpoint is intended to replace client-side composition of `GET /auth/me`
and `GET /customers/me` during app bootstrap, while keeping those endpoints
available for compatibility and focused use cases.

Success response (200):

```json
{
  "data": {
    "user": {
      "id": "d3de5f8b-4892-42e8-9680-979cf3f37844",
      "email": "user@example.com",
      "phone": "+5527999999999",
      "role": "customer"
    },
    "customer": {
      "id": "6f3ebf86-bf82-4b75-a2ce-cd261ca47ec3",
      "name": "Maria Silva",
      "cpf": "12345678901",
      "birth_date": "1990-01-15",
      "created_at": "2026-05-29T10:00:00Z"
    },
    "readiness": {
      "onboarding_completed": true,
      "approved": true,
      "has_operational_account": true,
      "transaction_password_status": "active"
    }
  },
  "error": null
}
```

Readiness fields:
- `onboarding_completed`: currently returns `true`; future onboarding steps,
  such as address, may change this value.
- `approved`: whether the authenticated customer user is approved/active.
- `has_operational_account`: whether the customer has at least one active
  account.
- `transaction_password_status`: one of `active`, `not_set`, `locked`, or
  `unknown`.

The API returns objective readiness state and does not decide whether a client
should display or navigate to its Home screen. Clients may derive presentation
and navigation decisions from these fields. Authorization remains enforced by
each protected API endpoint independently of client navigation.

The response intentionally does not include `user.customer_id` or
`customer.email`. The customer identifier is represented by `customer.id`, and
the authenticated contact e-mail belongs to `user.email`.

The endpoint never returns transaction password material, password hashes,
pepper values, step-up tokens, or any other sensitive credential material.

Possible errors:
- 401 UNAUTHORIZED: authentication required
- 401 INVALID_TOKEN: token invalid, malformed, or expired
- 409 INVALID_USER_STATE: customer user has an inconsistent customer link
- 404 CUSTOMER_NOT_FOUND or equivalent: linked customer record was not found
- 500 INTERNAL_ERROR: unexpected internal error

### 3.4.1 Get Current User

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

### 3.6 Create Customer Account (Admin Only)

- Method: POST
- Path: /admin/customers/{customer_id}/accounts
- Auth required: JWT (admin role)

Creates an additional account for an existing customer. Account creation is a
provisioning action and is not exposed as customer self-service.

Onboarding approval remains responsible for creating the first account
automatically. Additional accounts for the same `customer_id` are allowed.

Path parameters:

- `customer_id`: UUID of the customer that will receive the new account

Request body:

Empty object or no body required. Extra fields are rejected.

```json
{}
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
- 400 INVALID_DATA: invalid `customer_id`
- 400 INVALID_REQUEST: invalid JSON body or unexpected fields
- 403 FORBIDDEN: authenticated user does not have admin role
- 404 CUSTOMER_NOT_FOUND: customer does not exist
- 500 INTERNAL_ERROR: unexpected internal error

### 3.7 Create Transaction Password

- Method: POST
- Path: /security/transaction-password
- Auth required: JWT

Creates the authenticated user's initial transaction password. This endpoint is
used only for first credential setup. It does not require a previous transaction
password or a step-up token.

The transaction password is a numeric PIN with exactly 6 digits. The API stores
only the hash. The response never returns the PIN or hash.

Request body:

```json
{
  "transaction_password": "123456",
  "transaction_password_confirmation": "123456"
}
```

Success response (201):

```json
{
  "data": {
    "user_id": "d3de5f8b-4892-42e8-9680-979cf3f37844",
    "status": "active",
    "created_at": "2026-05-28T10:00:00Z"
  },
  "error": null
}
```

Response fields:
- `status`: currently always returns `active` on successful creation. The
  transaction password domain also has `blocked` for later validation/step-up
  flows, but a newly created transaction password is always active.

Possible errors:
- 401 UNAUTHORIZED: authentication required
- 401 INVALID_TOKEN: token invalid, malformed, or expired
- 400 INVALID_REQUEST: invalid JSON body or unexpected fields
- 400 INVALID_DATA: PIN is not numeric with 6 digits or confirmation differs
- 403 FORBIDDEN: authenticated user is not active
- 409 TRANSACTION_PASSWORD_ALREADY_SET: transaction password already exists
- 500 INTERNAL_ERROR: unexpected internal error

### 3.8 Authorize Step-Up

- Method: POST
- Path: /security/step-up/authorize
- Auth required: JWT

Authorizes a sensitive logical endpoint with the authenticated user's
transaction password and returns a short-lived step-up token. In the MVP, the
accepted public operations are `POST /accounts/internal-transfers` and
`POST /security/installations`. The installation registration operation may be
authorized with a restricted access token from login; other operations require
an operational access token.

The step-up token is an `HS256` JWT. It lasts 120 seconds, is scoped to the
requested public operation, and is tracked by a persisted `jti` so it can be
consumed once during enforcement. The response never returns the transaction
password, password hash, or operation payload.

The client must send only the public HTTP operation (`method` + `path`).
Internal policy keys (for example `internal_transfer.create`) are resolved by
the backend and are not part of the public request contract.

The canonical client representation uses an uppercase HTTP method, such as
`POST`. For input tolerance, the API trims surrounding whitespace and
normalizes `method` to uppercase before resolving the operation allowlist.
The API only trims surrounding whitespace from `path`; it does not normalize
or rewrite the path itself.

Request body:

```json
{
  "method": "POST",
  "path": "/accounts/internal-transfers",
  "transaction_password": "123456"
}
```

Success response (200):

```json
{
  "data": {
    "step_up_token": "<token>",
    "expires_in": 120
  },
  "error": null
}
```

Possible errors:
- 401 UNAUTHORIZED: authentication required
- 401 INVALID_TOKEN: token invalid, malformed, or expired
- 401 TRANSACTION_PASSWORD_INVALID: transaction password PIN is incorrect
- 400 INVALID_REQUEST: invalid JSON body or unexpected fields
- 400 INVALID_DATA: PIN is not numeric with 6 digits
- 403 FORBIDDEN: authenticated user is not active
- 403 TRANSACTION_PASSWORD_LOCKED: transaction password is temporarily blocked
- 403 STEP_UP_ENDPOINT_NOT_ALLOWED: public operation is not allowed for step-up
- 409 TRANSACTION_PASSWORD_NOT_SET: transaction password does not exist
- 500 INTERNAL_ERROR: unexpected internal error

### 3.9 Register Installation

- Method: POST
- Path: /security/installations
- Auth required: restricted access token, `X-Installation-Id`, and `X-Step-Up-Token`

Request headers:

```http
Authorization: Bearer <restricted_access_token>
X-Installation-Id: <canonical_uuid_v4>
X-Step-Up-Token: <step_up_token>
```

Success response (201):

```json
{
  "data": {
    "access_token": "<jwt>",
    "refresh_token": "<opaque-token>",
    "installation_resource_id": "2e4a8e20-272a-4e7b-b782-bc7f6b1d0442",
    "installation_status": "known"
  },
  "error": null
}
```

Possible errors:
- 400 INVALID_INSTALLATION_ID: missing or malformed `X-Installation-Id`
- 401 INVALID_TOKEN: restricted token is invalid, expired, consumed, or revoked
- 401 STEP_UP_TOKEN_REQUIRED: missing `X-Step-Up-Token`
- 401 STEP_UP_TOKEN_INVALID: invalid step-up token
- 403 INSTALLATION_MISMATCH: header does not match restricted token
- 409 INSTALLATION_LIMIT_REACHED: no available installation slot

### 3.10 List Installations

- Method: GET
- Path: /security/installations
- Auth required: operational JWT and `X-Installation-Id`

Success response (200):

```json
{
  "data": {
    "installations": [
      {
        "resource_id": "2e4a8e20-272a-4e7b-b782-bc7f6b1d0442",
        "status": "known",
        "first_seen_at": "2026-06-17T10:00:00Z",
        "last_seen_at": "2026-06-17T10:00:00Z",
        "created_at": "2026-06-17T10:00:00Z",
        "updated_at": "2026-06-17T10:00:00Z"
      }
    ]
  },
  "error": null
}
```

The response never exposes the raw `installation_id` generated by the client.

### 3.11 Revoke Installation

- Method: DELETE
- Path: /security/installations/{installation_resource_id}
- Auth required: operational JWT and `X-Installation-Id`

Success response (200):

```json
{
  "data": {
    "resource_id": "2e4a8e20-272a-4e7b-b782-bc7f6b1d0442",
    "status": "revoked",
    "revoked_at": "2026-06-17T10:10:00Z"
  },
  "error": null
}
```

The current installation cannot revoke itself. Revoking another installation
invalidates refresh sessions bound to that installation.

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

### 4.2 Customer Account Creation Removed

The customer-facing `POST /accounts` route is not registered. Account creation is
available only through onboarding approval and the admin provisioning route:

```http
POST /admin/customers/{customer_id}/accounts
```

### 4.3 Deposit

- Method: POST
- Path: /terminal/accounts/{id}/deposit
- Auth required: n/a (route disabled)

Operational note:
- This endpoint is not intended for mobile or customer-facing web clients. It
  directly injects balance into the ledger and is positioned as a terminal
  operation.
- The route is intentionally disabled in the API wiring and is not callable.
- A real terminal channel is outside the current project scope.

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
- Path: /terminal/accounts/{id}/withdraw
- Auth required: n/a (route disabled)

Operational note:
- This endpoint is not intended for mobile or customer-facing web clients. It
  directly removes balance from the ledger and is positioned as a terminal
  operation.
- The route is intentionally disabled in the API wiring and is not callable.
- A real terminal channel is outside the current project scope.

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

This endpoint is step-up protected. The client must send a valid
`X-Step-Up-Token` issued for `POST /accounts/internal-transfers`.

Request headers:

```http
Authorization: Bearer <access_token>
X-Step-Up-Token: <step_up_token>
```

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
- `X-Step-Up-Token` is mandatory and must be a valid token issued for
  `POST /accounts/internal-transfers`.
- The step-up token is consumed atomically before the transfer use case is
  executed.
- Because the token is single-use, retrying with the same `X-Step-Up-Token`
  returns `STEP_UP_TOKEN_CONSUMED`.
- `idempotency_key` is required and must be stable across retries of the same
  transfer attempt.
- A new step-up token may be required even when retrying with the same
  `idempotency_key`.
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
- 401 STEP_UP_TOKEN_REQUIRED: step-up token header was not provided
- 401 STEP_UP_TOKEN_INVALID: step-up token is invalid or malformed
- 401 STEP_UP_TOKEN_EXPIRED: step-up token has expired
- 401 STEP_UP_TOKEN_CONSUMED: step-up token was already consumed
- 400 INVALID_REQUEST: invalid JSON body
- 400 INVALID_DATA: missing or malformed account IDs, or missing
  `idempotency_key`
- 400 INVALID_AMOUNT: amount must be greater than zero
- 400 SAME_ACCOUNT_TRANSFER: source and destination are equal
- 403 FORBIDDEN: authenticated user cannot operate the source account
- 403 STEP_UP_ENDPOINT_MISMATCH: step-up token operation does not match
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
        "description": "Aluguel de maio",
        "created_at": "2026-04-02T12:00:00Z"
      }
    ],
    "next_cursor": null
  },
  "error": null
}
```

When there are more results, `next_cursor` is an object — pass both fields as query params for the next page:

Statement item notes:
- `description` is optional and is omitted when the transaction has no description.

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

Notes:
- `cpf` is resolved from the customer's primary `customer_documents` row (`type=cpf`, `country=BR`).

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
- CONTACT_NOT_VERIFIED
- ACCOUNT_APPROVAL_REQUIRED
- UNAUTHORIZED
- INVALID_TOKEN
- FORBIDDEN
- ACCOUNT_NOT_FOUND
- ACCOUNT_INACTIVE
- INSUFFICIENT_FUNDS
- SAME_ACCOUNT_TRANSFER
- TRANSACTION_NOT_FOUND
- TRANSACTION_PASSWORD_ALREADY_SET
- TRANSACTION_PASSWORD_NOT_SET
- TRANSACTION_PASSWORD_INVALID
- TRANSACTION_PASSWORD_LOCKED
- TRANSACTION_PASSWORD_REQUIRED
- STEP_UP_ENDPOINT_NOT_ALLOWED
- STEP_UP_TOKEN_REQUIRED
- STEP_UP_TOKEN_INVALID
- STEP_UP_TOKEN_EXPIRED
- STEP_UP_TOKEN_CONSUMED
- STEP_UP_ENDPOINT_MISMATCH
- INVALID_INSTALLATION_ID
- INSTALLATION_MISMATCH
- INSTALLATION_REVOKED
- INSTALLATION_LIMIT_REACHED
- INTERNAL_ERROR

`INVALID_APP_TOKEN` (HTTP 401) is returned when onboarding routes protected by AppToken (`POST /auth/cpf-check`, `POST /auth/contact-verifications`, `POST /auth/contact-verifications/confirm`, `POST /auth/register`, `POST /auth/login`) are called without `X-App-Token` or with an invalid app token.

`INVALID_INSTALLATION_ID` (HTTP 400) is returned when `X-Installation-Id` is
missing or is not a canonical UUID v4.

`INSTALLATION_MISMATCH` (HTTP 403) is returned when `X-Installation-Id` does
not match the installation bound to the token, refresh session, or restricted
authorization.

`INSTALLATION_REVOKED` (HTTP 403) is returned when a revoked installation tries
to login again.

`INSTALLATION_LIMIT_REACHED` (HTTP 409) is returned when a user attempts to add
a fourth known installation.

`ACCOUNT_APPROVAL_REQUIRED` (HTTP 403) is returned by `POST /auth/login` when a
customer user has valid credentials but cannot enter the app because admin
approval/account provisioning is incomplete. Mobile clients should use this code
to show an approval-pending guidance message.

`CONTACT_NOT_VERIFIED` (HTTP 403) is returned by `POST /auth/login` when one or
both contact channels are not verified. The payload may include
`error.details.email_verified` and `error.details.phone_verified`.

`INVALID_TOKEN` (HTTP 401) is returned for any of the following conditions on the `/auth/refresh` endpoint: token not found, already revoked, expired, or refresh token signature invalid.

`INVALID_USER_STATE` (HTTP 409) indicates the system detected an invariant violation: a user with role `customer` has no linked `customer_id`. This should never occur under normal operation; it signals a data consistency bug.

`TRANSACTION_PASSWORD_ALREADY_SET` (HTTP 409) is returned when the authenticated
user tries to create a transaction password after one already exists.

`TRANSACTION_PASSWORD_NOT_SET` (HTTP 409) is returned by transaction-password
dependent flows when the authenticated user has not created the credential yet.

`TRANSACTION_PASSWORD_INVALID` (HTTP 401) is returned when a transaction password
validation attempt receives an incorrect PIN.

`TRANSACTION_PASSWORD_LOCKED` (HTTP 403) is returned when the transaction
password is temporarily locked after repeated invalid attempts.

`TRANSACTION_PASSWORD_REQUIRED` (HTTP 403) is reserved for challenge/policy
flows that require step-up before a sensitive operation can proceed. It is not
returned by `POST /accounts/internal-transfers`; that endpoint returns
`STEP_UP_TOKEN_REQUIRED` when `X-Step-Up-Token` is missing.

`STEP_UP_TOKEN_REQUIRED` (HTTP 401) is returned by protected endpoints when
header `X-Step-Up-Token` is missing.

`STEP_UP_ENDPOINT_NOT_ALLOWED` (HTTP 403) belongs to step-up authorization and
is returned when the requested public operation (`method` + `path`) is outside
the backend allowlist.

`STEP_UP_TOKEN_INVALID` (HTTP 401) is returned for malformed step-up tokens,
invalid signatures, required claims missing, `scope` different from `step_up`,
or missing persisted `jti`.

`STEP_UP_TOKEN_EXPIRED` (HTTP 401) is returned when the step-up token is
already expired.

`STEP_UP_TOKEN_CONSUMED` (HTTP 401) is returned when a previously consumed
step-up token is reused.

`STEP_UP_ENDPOINT_MISMATCH` (HTTP 403) is returned when the token was issued
for a different operation than the protected endpoint.

## 8. Domain Notes for API Consumers

- Monetary values are represented as integer cents
- UUID is used for all resource identifiers
- Financial operations are synchronous and strongly consistent
- Transfer operation is atomic: debit and credit are committed together

## 9. Error Scenarios by Endpoint (with Payload)

This section lists common error situations and the expected payload shape.

### 9.0 POST /auth/cpf-check

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

Scenario: invalid cpf format or missing cpf
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

Examples of `error.message` for this scenario include `Invalid CPF format` and
`CPF is required`.

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

Scenario: account approval required
- Status: 403
- Code: ACCOUNT_APPROVAL_REQUIRED

```json
{
  "data": null,
  "error": {
    "code": "ACCOUNT_APPROVAL_REQUIRED",
    "message": "Account approval required"
  }
}
```

Scenario: contact not verified
- Status: 403
- Code: CONTACT_NOT_VERIFIED

```json
{
  "data": null,
  "error": {
    "code": "CONTACT_NOT_VERIFIED",
    "message": "Contact not verified",
    "details": {
      "email_verified": true,
      "phone_verified": false
    }
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

### 9.4 GET /auth/session

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

Scenario: invalid customer user state
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

### 9.4.1 GET /auth/me

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

### 9.6 POST /admin/customers/{customer_id}/accounts

Scenario: authenticated user is not admin
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

### 9.7 POST /terminal/accounts/{id}/deposit

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

### 9.8 POST /terminal/accounts/{id}/withdraw

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

Scenario: missing `X-Step-Up-Token`
- Status: 401
- Code: STEP_UP_TOKEN_REQUIRED

```json
{
  "data": null,
  "error": {
    "code": "STEP_UP_TOKEN_REQUIRED",
    "message": "Step-up token required"
  }
}
```

Scenario: invalid or malformed step-up token
- Status: 401
- Code: STEP_UP_TOKEN_INVALID

```json
{
  "data": null,
  "error": {
    "code": "STEP_UP_TOKEN_INVALID",
    "message": "Invalid step-up token"
  }
}
```

Scenario: expired step-up token
- Status: 401
- Code: STEP_UP_TOKEN_EXPIRED

```json
{
  "data": null,
  "error": {
    "code": "STEP_UP_TOKEN_EXPIRED",
    "message": "Step-up token expired"
  }
}
```

Scenario: retry with same consumed step-up token
- Status: 401
- Code: STEP_UP_TOKEN_CONSUMED

```json
{
  "data": null,
  "error": {
    "code": "STEP_UP_TOKEN_CONSUMED",
    "message": "Step-up token already consumed"
  }
}
```

Scenario: token issued for another operation
- Status: 403
- Code: STEP_UP_ENDPOINT_MISMATCH

```json
{
  "data": null,
  "error": {
    "code": "STEP_UP_ENDPOINT_MISMATCH",
    "message": "Step-up endpoint mismatch"
  }
}
```

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

### 9.15 POST /security/transaction-password

Scenario: invalid PIN or confirmation mismatch
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

Scenario: authenticated user is not active
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

Scenario: transaction password already exists
- Status: 409
- Code: TRANSACTION_PASSWORD_ALREADY_SET

```json
{
  "data": null,
  "error": {
    "code": "TRANSACTION_PASSWORD_ALREADY_SET",
    "message": "Transaction password already set"
  }
}
```

## 10. Bruno Setup

The repository uses Bruno for local API exploration under `tools/bruno`.

### 10.1 Files in Repository

- `tools/bruno/README.md`
- Bruno collection files (`*.bru`) when requests are exported/versioned

### 10.2 Environment Variables

Use these variables when configuring the Bruno environment:

- `base_url`: API base URL (default: `http://localhost:8080`)
- `app_token`: application token used by auth entry routes (`/auth/cpf-check`, `/auth/contact-verifications`, `/auth/contact-verifications/confirm`, `/auth/register`, and `/auth/login`)
- `access_token`: JWT used for protected routes
- `refresh_token`: opaque refresh token used by `/auth/refresh`
- `account_id`: account UUID for account operations
- `from_account_id`: source account UUID used by internal transfer requests
- `to_account_id`: destination account UUID used by internal transfer requests
- `recipient_branch`: branch used for recipient lookup examples
- `recipient_account_number`: account number used for recipient lookup examples
- `recipient_document`: CPF used for recipient lookup examples
- `email_verification_id`: id returned by `POST /auth/contact-verifications` (email)
- `phone_verification_id`: id returned by `POST /auth/contact-verifications` (phone)
- `email_verification_token`: token returned by `POST /auth/contact-verifications/confirm` (email)
- `phone_verification_token`: token returned by `POST /auth/contact-verifications/confirm` (phone)
- `transaction_reference`: public reference returned by successful transfer requests
- `id`: user UUID used by admin approval route (`/admin/users/{id}/approve`)

### 10.3 How to Import and Configure

1. Open Bruno.
2. Open the collection directory under `tools/bruno` when collection files are present.
3. Configure the environment variables listed above.
4. Adjust `base_url` if your API is not running on `http://localhost:8080`.
5. Confirm `app_token` matches the value configured in your local API environment.
6. Run auth requests to obtain tokens and update `access_token` / `refresh_token`.

### 10.4 Recommended Execution Flow

Use this flow to bootstrap test data and credentials quickly:

1. `Auth/CPFCheck`
2. `Auth/ContactVerifications/RequestEmail`
3. `Auth/ContactVerifications/ConfirmEmail`
4. `Auth/ContactVerifications/RequestPhone`
5. `Auth/ContactVerifications/ConfirmPhone`
6. `Auth/Register`
7. `Auth/Login` (copy `access_token` and `refresh_token` from response)
8. `Auth/Me` (validate JWT)
9. `Account/User/Approve` (admin only, using `id`)
10. Account endpoints using `account_id` as needed
11. Recipient lookup endpoint using `recipient_branch` + `recipient_account_number`, or `recipient_document`
12. Internal transfer endpoint using `from_account_id` and `to_account_id`
13. Receipt endpoint using `transaction_reference` from a successful transfer response

Notes:
- Keep `X-App-Token` for onboarding requests (`/auth/cpf-check`, `/auth/contact-verifications`, `/auth/contact-verifications/confirm`, `/auth/register`, `/auth/login`) as documented in this file.
- For protected routes, send `Authorization: Bearer <access_token>`.
- If a request returns `401 INVALID_TOKEN`, run `Auth/Refresh` and update tokens.
