# Backlog Section — Transactional Security & Context Validation

## 1. Objective

Introduce **context-aware decision making** into the system, extending the current authentication model.

The current system relies on:

* App Token (entry control)
* JWT (identity)
* Role + ownership (authorization)
* User status (operational capability)

This evolution introduces:

> **operation-level trust evaluation**

Instead of trusting all requests equally once a JWT is valid.

---

## 2. Scope

This backlog covers three core capabilities:

1. Transactional Password
2. Installation Registration
3. Liveness / Step-up Authentication

These capabilities are not isolated features. They must be treated as **inputs to a unified decision model**.

---

## 3. Architectural Positioning

The new logic must be implemented at the **application layer**, not:

* not in HTTP handlers (delivery)
* not inside domain entities
* not inside repositories

It acts as a **policy/guard layer** executed before critical use cases.

---

### Updated Request Flow

```text
Request
 → JWT Middleware
 → Context Evaluation (NEW)
 → Use Case
```

Where **Context Evaluation** may include:

* transactional password validation
* installation validation
* liveness validation

---

## 4. Feature 1 — Transactional Password

### Purpose

Provide **explicit user intent validation** for sensitive operations.

### Applies to

* withdraw
* transfer

### Design Decisions

* Stored as a separate credential (hashed)
* Must not reuse login password
* Validated **before transaction begins**

### Important Constraint

Validation must occur **outside the DB transaction** to avoid:

* increased lock duration
* degraded concurrency behavior

### Required Changes

* Extend `users` (or separate table) with:

  * `transaction_password_hash`
* Add endpoints:

  * set transactional password
  * update transactional password
* Update use cases:

  * `Withdraw`
  * `Transfer`

### New Errors

* `INVALID_TRANSACTION_PASSWORD`

---

## 5. Feature 2 — Installation Registration

Esta frente está separada nos backlogs ativos:

* API 010: validação, associação, estados, sessão e revogação;
* Mobile 013: geração, persistência local e envio de `X-Installation-Id`.

### Purpose

Bind user sessions to a **known app installation**.

### Current Limitation

Sessions (`user_sessions`) are not installation-aware.

### Proposed Model

New entity:

```text
app_installations
- id
- user_id
- installation_id (client-provided UUID)
- status
- platform
- app_version
- created_at
- last_seen_at
```

### Integration Points

**Login flow:**

```text
mobile creates installation_id before login
  → login sends credentials + X-Installation-Id
  → known installation: create bound session
  → new account: silently register first installation and create bound session
  → existing account on new installation: issue restricted access
  → transaction password authorizes POST /security/installations through step-up
  → registration endpoint creates the association and operational session
  → new installation at limit: reject until one known installation is revoked
```

**Request flow:**

```text
JWT → extract user
session + X-Installation-Id → validate bound installation
```

### Important Constraint

Installation identification is **not device identification or a strong
factor**.

It must be treated as:

> a weak signal, never a standalone trust decision

### Required Changes

* Create `app_installations` table
* Bootstrap the first installation atomically during login
* Issue restricted access for later installations
* Add `POST /security/installations` to the step-up operation allowlist
* Limit each user to three `known` installations
* Preserve revoked installation associations as audit history
* Validate the session-installation binding on refresh and authenticated calls
* Add installation endpoints:

  * register installation
  * list installations
  * revoke installation
* Introduce header:

  * `X-Installation-Id`
* Application-level validation:

  * installation state may become one contextual input for sensitive operations

Physical device identification, fingerprinting, attestation, and device trust
are deferred to a separate future discussion.

---

## 6. Feature 3 — Liveness / Step-up Authentication

### Purpose

Ensure the user is **actively present at the time of operation**, not just authenticated.

### Problem Addressed

JWT validity does not guarantee:

* user presence
* user awareness of the operation

---

### Recommended Approach (MVP)

Start with a **step-up challenge model**:

* short-lived challenge (TTL)
* linked to operation
* validated before execution

---

### Flow

```text
Request (transfer)
 → requires step-up
 → client requests challenge
 → user confirms (password / another approved factor / biometrics)
 → operation proceeds
```

---

### Design Options

| Approach               | Complexity | Strength |
| ---------------------- | ---------- | -------- |
| Password re-check      | Low        | Medium   |
| Installation context   | Low        | Weak     |
| Biometric (local)      | Medium     | Good     |
| Real liveness (camera) | High       | Strong   |

### Recommendation

Start with:

* transactional password
* installation context as a weak signal

Defer real liveness detection.

---

### Required Changes

* New entity:

  * `auth_challenges`
* Fields:

  * id
  * user_id
  * type
  * expires_at
  * consumed_at
* Endpoints:

  * create challenge
  * validate challenge

---

## 7. Feature 4 — Operation Policy Layer (Critical)

### Purpose

Centralize **decision logic per operation**.

Without this, the system will degrade into scattered validations.

---

### Concept

Each use case defines its **security requirements**.

Example:

| Operation  | Requirements                                                       |
| ---------- | ------------------------------------------------------------------ |
| GetBalance | JWT                                                                |
| Deposit    | JWT                                                                |
| Withdraw   | JWT + transactional password                                       |
| Transfer   | JWT + transactional password + installation context (+ step-up optional) |

---

### Implementation Direction

Introduce a **policy evaluator** in the application layer:

```text
Evaluate(user, operation, context) → allow / deny
```

Inputs:

* user (JWT context)
* installation
* operation type
* optional challenge state

---

### Responsibility

* decide if operation can proceed
* return structured failure (error codes)

---

## 8. Impact Summary

This backlog introduces:

* contextual validation before use cases
* separation between identity and trust
* foundation for Zero Trust evolution

It does **not** change:

* domain invariants
* transaction model
* ledger consistency guarantees

---

## 9. Final Assessment

This is not a “security feature set”.

It is the introduction of a:

> **decision engine for operation trust**

If implemented correctly, it becomes the foundation for:

* adaptive security
* risk-based authentication
* Zero Trust Architecture evolution

If implemented poorly, it becomes:

* duplicated validation logic
* scattered checks
* inconsistent behavior

---

If you want, the next step is to define:

* the **policy interface**
* and how it plugs directly into your existing use cases (especially `Transfer`)
