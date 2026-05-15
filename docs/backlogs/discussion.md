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
2. Device Registration
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
* device validation
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

## 5. Feature 2 — Device Registration

### Purpose

Bind user sessions to a **known execution environment**.

### Current Limitation

Sessions (`user_sessions`) are not device-aware 

### Proposed Model

New entity:

```text
devices
- id
- user_id
- device_id (client-provided identifier)
- trusted (bool)
- created_at
- last_seen_at
```

### Integration Points

**Login flow:**

```text
login → create session → register/update device
```

**Request flow:**

```text
JWT → extract user
device_id (header) → validate device
```

### Important Constraint

Device identification is **not a strong factor**.

It must be treated as:

> a weak signal, never a standalone trust decision

### Required Changes

* Create `devices` table
* Add endpoints:

  * register device
  * list devices
  * revoke device
* Introduce header:

  * `X-Device-Id`
* Application-level validation:

  * sensitive operations require trusted device

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
 → user confirms (password / device / biometrics)
 → operation proceeds
```

---

### Design Options

| Approach               | Complexity | Strength |
| ---------------------- | ---------- | -------- |
| Password re-check      | Low        | Medium   |
| Device confirmation    | Medium     | Good     |
| Biometric (local)      | Medium     | Good     |
| Real liveness (camera) | High       | Strong   |

### Recommendation

Start with:

* transactional password
* device confirmation

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
| Transfer   | JWT + transactional password + trusted device (+ step-up optional) |

---

### Implementation Direction

Introduce a **policy evaluator** in the application layer:

```text
Evaluate(user, operation, context) → allow / deny
```

Inputs:

* user (JWT context)
* device
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
