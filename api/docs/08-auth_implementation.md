# Auth & Authorization — Bank API

## 1. Overview

This document describes the **current authentication and authorization model implemented in the system**.

It reflects the actual runtime behavior of the API and complements the implementation plan defined previously.

The model is intentionally simple and designed to:

* control access to onboarding endpoints
* validate end-to-end authentication flows
* enable functional testing of the system
* avoid premature introduction of advanced security mechanisms

---

## 2. Authentication Model

At the current stage, the system adopts a **multi-stage authentication model**, separating:

* system entry control
* session establishment
* resource access control

This separation allows the system to remain simple while preserving a clear path for future evolution.

---

## 3. Token Types

The system operates with these token types:

### 3.1 App Token

A static token used to restrict access to onboarding endpoints.

**Purpose:**

* prevent unauthorized clients from creating users
* prevent automated abuse of authentication endpoints
* limit onboarding to known applications

**Characteristics:**

* defined via environment variable (`APP_TOKEN`)
* sent via HTTP header: `X-App-Token`
* validated at the HTTP boundary (Delivery layer)
* not associated with any user identity

---

### 3.2 Access Token (JWT)

A user-scoped token issued after successful authentication.

**Purpose:**

* identify the authenticated user
* authorize access to protected resources

**Characteristics:**

* issued during login
* short-lived
* contains claims (`sub`, `role`, `customer_id`)
* validated via JWT middleware
* lifetime is configured by `JWT_ACCESS_TOKEN_DURATION` and defaults to `15m`

### 3.3 Restricted Access Token (JWT)

A short-lived token issued after successful credentials when the presented app
installation is new and cannot be associated automatically.

**Purpose:**

* carry only the context needed to register the presented installation
* allow step-up authorization for `POST /security/installations`
* avoid creating an operational session before the installation is registered

**Characteristics:**

* sent via `Authorization: Bearer <restricted_access_token>`
* contains `token_type = restricted_access`
* contains `scope = installation.register`
* contains the presented `installation_id`
* does not create a refresh token
* cannot access account, customer, statement, transfer, or general security
  operations
* expires after five minutes and is tracked by persisted `jti`

---

## 4. Request Flow

### 4.1 Onboarding Endpoints (AppToken Required)

Endpoints:

* `POST /auth/cpf-check`
* `POST /auth/contact-verifications`
* `POST /auth/contact-verifications/confirm`
* `POST /auth/register`
* `POST /auth/login`

Flow:

```text
Request
  ↓
[AppToken Middleware]
  ↓
Auth Handler
```

**Requirements:**

* `X-App-Token` header is mandatory
* JWT is not required
* `POST /auth/cpf-check` validates CPF format and availability before register
* `POST /auth/contact-verifications` rejects already used e-mail/phone with `USER_ALREADY_EXISTS`
* registration requires confirmed verification tokens for e-mail and phone

On register, CPF is persisted as a primary customer document (`customer_documents`)
instead of a direct `customers.cpf` column.

---

### 4.2 Authenticated Endpoints

Endpoints:

* `GET /auth/me`

Flow:

```text
Request
  ↓
[JWT Middleware]
  ↓
Auth Handler
```

**Requirements:**

* valid `Authorization: Bearer <access_token>`
* App Token is not required

Login may fail with `CONTACT_NOT_VERIFIED` when contact verification was not
completed. In this case, error details include channel verification flags to
drive client onboarding UX.

### 4.3 Refresh Endpoint

Endpoint:

* `POST /auth/refresh`

The refresh endpoint is authenticated by the refresh token in the request body,
not by the expired access token. The handler delegates validation to the refresh
use case, which checks the token signature, persisted session, revocation state,
expiration, and user identity before rotating the session.
Refresh session lifetime is configured by `JWT_REFRESH_TOKEN_DURATION` and
defaults to `168h`.

When a refresh session is bound to an installation, `POST /auth/refresh` also
requires `X-Installation-Id`. The header must match the installation recorded
for that refresh session. Revoking an installation marks matching refresh
sessions as revoked, so later refresh attempts fail even if the refresh token
has not reached its cryptographic expiration.

---

### 4.4 Protected Resource Endpoints (JWT Required)

Examples:

* `/accounts/*`
* `/customers/me`

Flow:

```text
Request
  ↓
[JWT Middleware]
  ↓
Protected Handler
```

**Requirements:**

* valid `Authorization: Bearer <access_token>`
* valid `X-Installation-Id` matching the access token claim and server-side
  session context for operational routes
* App Token is not required

---

### 4.5 Installation Registration Flow

Endpoint:

* `POST /security/installations`

This route accepts a restricted access token, the same `X-Installation-Id`
presented at login, and a valid `X-Step-Up-Token` authorized for the public
operation `POST /security/installations`.

Successful registration consumes the restricted authorization, reserves one of
the user's known installation slots, creates an operational refresh session, and
returns normal access and refresh tokens. The API still treats the
`installation_id` as a weak contextual signal. It is never sufficient by itself
for authentication or authorization.

---

## 5. Audit, Retention, and Logging

Installation identity events that must be observable in operational telemetry:

* login with known installation
* first installation bootstrap
* login requiring restricted installation registration
* login denied for revoked installation
* login denied because the installation limit was reached
* restricted authorization creation, consumption, revocation, and expiration
* installation registration
* installation listing
* installation revocation and session invalidation
* installation mismatch between header, token claim, and session context

Logs and audit payloads must not include access tokens, refresh tokens,
restricted access tokens, step-up tokens, transaction passwords, password hashes,
raw request bodies for sensitive operations, or environment attributes beyond
the minimum fields needed to operate the flow. When correlation is needed, use
event names, user/resource identifiers, status, timestamps, and error codes
instead of secrets.

Retention policy:

* revoked installations are retained as historical records because they preserve
  bootstrap history and support support/audit review;
* active restricted installation authorizations expired for more than 24 hours
  are deleted;
* consumed restricted installation authorizations are deleted 24 hours after
  `consumed_at`;
* revoked restricted installation authorizations are deleted 24 hours after
  `created_at`.

The cleanup function
`cleanup_installation_registration_authorizations()` is scheduled by `pg_cron`
for 03:30 daily.

## 6. Authorization Model

Authorization is based on the authenticated user identity extracted from the JWT.

### Rules:

* users with role `customer` can only access their own resources
* users with role `admin` can access any resource
* ownership is enforced at the **application layer**, not only in handlers

### Source of truth:

* `customer_id` is always derived from JWT
* client input for ownership is ignored

---

## 6.1 Operational Status (UserStatus)

**IMPORTANT DISTINCTION:** Authorization, user lifecycle, and account operability are separate concerns.

### Definition

Beyond **authentication** (who you are) and **authorization** (what role you have), the system enforces a **user lifecycle status** (`users.status`) that controls onboarding progression and eligibility to open accounts.

### Three Layers

1. **Authentication** (JWT)
   - Identity verification
   - Claims: `sub` (user ID), `role`, `customer_id`
   - Validity: short-lived, token-based

2. **Authorization** (Role)
   - Access control: customer vs. admin vs. other roles
   - Enforced at application layer
   - Tied to JWT claims

3. **User Lifecycle Status** (UserStatus)
   - Controls onboarding progression and account-opening eligibility
   - Stored in `users.status` column
   - Values: `pending`, `active`, `blocked`
   - Required: **user must be `active` to complete approval-dependent flows such as account opening**

4. **Account Operational Status** (AccountStatus)
   - Controls whether a specific account can execute financial operations
   - Stored in `accounts.status` column
   - Values: `active`, `inactive`, `blocked`
   - Enforced by account-domain rules such as `CanDeposit`, `CanWithdraw`, and `CanTransfer`

### Invariant

```text
Authentication + Authorization ≠ User Lifecycle Eligibility

Active user status allows onboarding-complete flows such as account creation.
Account status governs whether a specific account is operational.
```

### Examples

| Scenario                                | JWT Valid? | Role OK? | User Status | Account Status | Can Open Account? | Can Move Funds? |
| --------------------------------------- | ---------- | -------- | ----------- | -------------- | ----------------- | --------------- |
| Pending user, valid JWT                 | ✓          | ✓        | pending     | -              | ✗                 | N/A             |
| Active user, valid JWT, active account  | ✓          | ✓        | active      | active         | ✓                 | ✓               |
| Active user, valid JWT, blocked account | ✓          | ✓        | active      | blocked        | ✓                 | ✗               |
| Blocked user, valid JWT, active account | ✓          | ✓        | blocked     | active         | ✗                 | account-defined |
| No JWT                                  | ✗          | -        | -           | -              | ✗                 | ✗               |

### Responsibility

* **Authentication**: verified at JWT middleware (HTTP boundary)
* **Authorization**: enforced at application layer (use cases)
* **User Lifecycle Status**: enforced at application layer where onboarding/account-opening rules apply
* **Account Operational Status**: enforced in account-domain business rules and account use cases

---

## 7. Design Rationale

This model reflects a deliberate decision to:

* restrict **who can initiate authentication**
* separate onboarding from authenticated usage
* delegate **resource access control** to JWT-based identity
* isolate authentication concerns from future security layers

The system assumes:

> a request that successfully obtained a valid access token has passed through a controlled entry point

This assumption is **intentionally limited to this stage of the system**.

---

## 7. Known Limitations

This model does **not** provide:

* continuous validation of client integrity
* request-level contextual verification
* differentiation between client applications after authentication
* protection against token reuse outside the original client
* device or environment validation

Once a valid JWT is issued:

```text
JWT alone is sufficient to access protected resources
```

---

## 8. Security Model Interpretation

The current system follows a **trusted boundary at authentication time**.

This means:

* control is enforced at onboarding (`AppToken`)
* identity is enforced via JWT
* no additional validation is performed per request

This is a **controlled simplification**, not a final security model.

---

## 9. Future Evolution

This model is designed to evolve toward a **Zero Trust Architecture (ZTA)**.

Planned improvements include:

* enforcing App Token on all requests
* introducing client identity into request context
* incorporating device and environment signals
* implementing request-level decision models

Future pipeline:

```text
Request
  ↓
[App Identity]
  ↓
[User Identity]
  ↓
[Context Evaluation]
  ↓
Decision
```

---

## 10. Relationship with Implementation Plan

The authentication system was implemented following a phased approach:

* user persistence
* password hashing
* JWT generation
* session management (refresh token)
* middleware enforcement
* authorization rules

This document describes the **result of those phases**, not the steps themselves.

---

## 11. Conclusion

The current authentication model provides:

* controlled system entry
* clear separation between onboarding and usage
* consistent JWT-based authorization
* sufficient security for functional validation

It intentionally prioritizes:

* simplicity
* clarity
* architectural stability

over premature complexity.

This forms a solid baseline for future evolution toward more advanced security models.
