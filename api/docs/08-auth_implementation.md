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

The system operates with two distinct tokens:

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

---

## 4. Request Flow

### 4.1 Onboarding Endpoints (AppToken Required)

Endpoints:

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

---

### 4.2 Authenticated Endpoints (JWT Required)

Endpoints:

* `POST /auth/refresh`
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

---

### 4.3 Protected Resource Endpoints (JWT Required)

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
* App Token is not required

---

## 5. Authorization Model

Authorization is based on the authenticated user identity extracted from the JWT.

### Rules:

* users with role `customer` can only access their own resources
* users with role `admin` can access any resource
* ownership is enforced at the **application layer**, not only in handlers

### Source of truth:

* `customer_id` is always derived from JWT
* client input for ownership is ignored

---

## 5.1 Operational Status (UserStatus)

**IMPORTANT DISTINCTION:** Authorization and operational capability are separated.

### Definition

Beyond **authentication** (who you are) and **authorization** (what role you have), the system enforces **operational status** (whether you can act).

### Three Layers

1. **Authentication** (JWT)
   - Identity verification
   - Claims: `sub` (user ID), `role`, `customer_id`
   - Validity: short-lived, token-based

2. **Authorization** (Role)
   - Access control: customer vs. admin vs. other roles
   - Enforced at application layer
   - Tied to JWT claims

3. **Operational Status** (UserStatus)
   - Capability to perform financial operations
   - Stored in `users.status` column
   - Values: `pending`, `active`, `blocked`
   - Required: **user must be `active` to operate in the system**

### Invariant

```text
Authentication + Authorization ≠ Operational Capability

Only active users can perform financial operations.
```

### Examples

| Scenario                    | JWT Valid? | Role OK? | Status? | Can Operate? |
| --------------------------- | ---------- | -------- | ------- | ------------ |
| Pending user, has valid JWT | ✓          | ✓        | pending | ✗            |
| Active user, customer role  | ✓          | ✓        | active  | ✓            |
| Blocked user, valid JWT     | ✓          | ✓        | blocked | ✗            |
| No JWT                      | ✗          | -        | -       | ✗            |

### Responsibility

* **Authentication**: verified at JWT middleware (HTTP boundary)
* **Authorization**: enforced at application layer (use cases)
* **Operational Status**: enforced at application layer (business rules)

---

## 6. Design Rationale

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
