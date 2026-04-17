# Database Documentation — Bank API

## 1. Overview

The database is a **first-class component** of the system architecture, alongside:

* API (application layer)
* Mobile client (consumer layer)
* Database (consistency and persistence layer)

It is not treated as a passive storage, but as an **active consistency boundary** responsible for:

* enforcing data integrity
* supporting transactional guarantees
* enabling deterministic financial operations
* preserving auditability

---

## 2. Design Principles

### 2.1 Strong Consistency

All financial operations are executed under:

* ACID transactions
* row-level locking
* deterministic update order

There is **no eventual consistency** in balance-changing operations.

---

### 2.2 Ledger-Based Model

The system follows a **ledger + snapshot model**:

* `account_transactions` → immutable ledger (source of truth)
* `accounts.balance` → current state (derived snapshot)

The ledger is authoritative. The snapshot exists for performance.

---

### 2.3 Immutability

Financial records are **append-only**:

* no UPDATE
* no DELETE

Enforced via database trigger.

---

### 2.4 Idempotency at Data Level

Idempotency is implemented directly in the ledger:

* scoped per account
* enforced via unique index
* guarantees safe retries

---

### 2.5 Minimal Redundancy

Each concept is represented once:

* no duplicated ledgers
* no parallel transaction tables

---

## 3. Schema Overview

### Core Tables

* `customers`
* `users`
* `user_sessions`
* `accounts`
* `account_transactions`

### Support Tables

* `schema_migrations`

---

## 4. Table Definitions

## 4.1 customers

Represents the business entity (account owner).

Fields:

* `id` (UUID, PK)
* `name`
* `cpf` (unique, validated)
* `email` (unique)
* `created_at`

Constraints:

* CPF format validation
* unique CPF and email

---

## 4.2 users

Represents system authentication identity.

Fields:

* `id` (UUID, PK)
* `email` (unique)
* `password_hash`
* `role`
* `customer_id` (nullable)
* `status`
* `created_at`
* `updated_at`

Constraints:

* `customer_id` must exist when role = `customer`

---

## 4.3 user_sessions

Represents authentication sessions (refresh tokens).

Fields:

* `id` (UUID, PK)
* `user_id` (FK → users)
* `token_hash` (unique)
* `expires_at`
* `revoked_at`
* `created_at`

---

## 4.4 accounts

Represents a financial account.

Fields:

* `id` (UUID, PK)
* `customer_id` (FK → customers)
* `number` (unique)
* `branch`
* `balance` (BIGINT, cents)
* `status` (enum: active, inactive, blocked)
* `created_at`

Notes:

* `balance` is a derived snapshot
* must only be modified via transactional operations

---

## 4.5 account_transactions

Represents the **financial ledger**.

This is the most critical table in the system.

Fields:

* `id` (UUID, PK)
* `account_id` (FK → accounts)
* `type` (enum: deposit, withdraw, transfer_in, transfer_out)
* `amount` (BIGINT, cents)
* `balance_after` (BIGINT)
* `reference_id` (UUID, groups related entries)
* `related_account_id` (UUID, used for transfers)
* `idempotency_key` (optional)
* `created_at`

---

## 5. Ledger Semantics

### 5.1 Source of Truth

All balance changes are recorded here.

No balance change exists without a corresponding ledger entry.

---

### 5.2 Transfer Model

A transfer is represented by **two entries**:

* `transfer_out` (source account)
* `transfer_in` (destination account)

Both share:

* same `reference_id`

And contain:

* `related_account_id` pointing to the other side

---

### 5.3 Snapshot Relationship

```text
accounts.balance = last(account_transactions.balance_after)
```

---

### 5.4 Immutability Enforcement

A trigger prevents:

* UPDATE
* DELETE

Ensuring full auditability.

---

## 6. Idempotency Model

### 6.1 Scope

Idempotency is defined as:

```text
(account_id, idempotency_key)
```

---

### 6.2 Enforcement

Implemented via partial unique index:

* only applies when `idempotency_key IS NOT NULL`

---

### 6.3 Behavior

* first execution → success
* retry → conflict → replay result

---

### 6.4 Replay Strategy

Replay does NOT read current state.

Instead, it reconstructs the result using the ledger:

* locate `transfer_out` by `(account_id, idempotency_key)`
* use `reference_id` to find matching `transfer_in`
* return stored `balance_after` values

This guarantees deterministic responses.

---

## 7. Indexing Strategy

### Performance indexes:

* `(account_id, created_at DESC)` → statement queries
* `(reference_id)` → transfer grouping
* `(reference_id, type)` → transfer pairing

### Transfer pair integrity index:

* UNIQUE `(reference_id, type)` for `type IN ('transfer_in', 'transfer_out')` and `reference_id IS NOT NULL`

### Idempotency index:

* `(account_id, idempotency_key)` (partial unique)

---

## 8. Consistency Guarantees

The system guarantees:

### Atomicity

* all operations succeed or fail entirely

### Consistency

* domain rules always enforced

### Isolation

* concurrent operations do not corrupt balance

### Durability

* committed operations are permanent

---

## 9. Concurrency Strategy

* `SELECT ... FOR UPDATE` used for account locking
* deterministic lock ordering prevents deadlocks
* balance updates are atomic

---

## 10. Invariants

The database enforces or supports:

* no negative balance without explicit rule
* no duplicate idempotent operations
* no partial transfers
* no mutation of historical data

---

## 11. Known Limitations

* no FK constraint on `reference_id`
* no explicit double-entry accounting abstraction
* no multi-currency support
* no ledger partitioning (yet)

---

## 12. Evolution Path

Future improvements may include:

* explicit debit/credit model
* partitioned ledger
* audit hashing / cryptographic integrity

---

## 13. Summary

The database acts as:

* the **consistency engine**
* the **audit source**
* the **financial truth layer**

It is intentionally simple in structure but strict in behavior.

The system prioritizes:

* correctness over complexity
* determinism over flexibility
* safety over premature optimization

This makes it a solid foundation for financial applications.
