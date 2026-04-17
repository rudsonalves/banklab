# Domain Model — Bank API

## 1. Overview

The system domain is centered on:

> **balance management through financial transaction records**

Within this model:

* **Transaction** represents the fundamental event
* **Account** represents the current state
* the **balance is derived from valid operations**

This approach ensures consistency, traceability, and auditability.

---

## 2. Domain Entities

### 2.1 Customer

Represents the account holder within the system.

**Attributes**

* id
* name
* cpf
* email
* createdAt

**Constraints**

* CPF must be unique
* email must be unique

---

### 2.2 User

Represents an authenticated system user with lifecycle management.

**Attributes**

* id
* email
* status
* customerId (for users with the customer role)
* createdAt

**Status**

* `pending` — newly registered, awaiting approval
* `active` — authorized to perform operations
* `blocked` — suspended or deactivated

**Constraints**

* email must be unique
* all users start with `pending` status
* only `active` users may access financial operations

---

### 2.3 Account

Represents a financial account.

**Attributes**

* id
* customerId
* number
* branch
* balance (in cents)
* status
* createdAt

**Status**

* `active`
* `inactive`
* `blocked`

---

### 2.4 Transaction

Represents a financial movement within the system.

**Attributes**

* id
* accountId
* type
* amount (in cents)
* balanceAfter
* referenceId (optional)
* relatedAccountId (optional)
* idempotencyKey (optional)
* createdAt

**Types**

* `deposit`
* `withdraw`
* `transfer_in`
* `transfer_out`

---

## 3. Domain Invariants

### 3.1 Monetary Value

* must be represented as an integer (cents)
* must be greater than zero

---

### 3.2 Account

* must exist for any operation
* must be in `active` status to allow transactions

---

### 3.3 Balance

* must never become inconsistent
* cannot be modified without a corresponding transaction

---

### 3.4 Transaction

* must always be associated with an account
* must always represent a financial impact

---

### 3.5 User Lifecycle

* only users with `active` status may perform operations
* status transitions are handled by the application layer
* no automatic transition from `pending` to `active`

---

### 3.6 Transfer

* source account must differ from destination account
* must generate two transactions:

  * debit on the source account
  * credit on the destination account

---

## 4. Business Rules

### 4.1 Deposit

* amount must be positive
* account balance is increased
* a `deposit` transaction is recorded

---

### 4.2 Withdrawal

* amount must be positive
* account must have sufficient balance
* account balance is decreased
* a `withdraw` transaction is recorded

---

### 4.3 Transfer

* amount must be positive
* source account must have sufficient balance
* operation must be atomic
* two transactions must be recorded:

  * `transfer_out` (source)
  * `transfer_in` (destination)

---

### 4.4 Statement

* must reflect all account transactions
* must be ordered by date

---

## 5. Balance Consistency

The account balance must follow:

```text
balance = sum(deposits + transfer_in) - sum(withdraw + transfer_out)
```

---

## 6. Idempotency

Transfer operations may include an `idempotencyKey`.

**Rules**

* the same key (within the scope of the source account) must not produce duplicate operations
* repeated execution must return the same historical result

**Applies to**

* transfer operations

---

## 7. Concurrency Behavior

* balance-modifying operations must be protected
* concurrent operations must not corrupt the balance
* the system must guarantee consistency under concurrent access

---

## 8. Integrity Rules

* no balance change occurs without a corresponding record in `account_transactions`
* no invalid intermediate states are allowed
* operations must be atomic (all-or-nothing)

---

## 9. Design Decisions

### 9.1 Stored Balance

The balance is stored directly in the account for performance reasons.

**Trade-off**

* improved read performance
* increased responsibility for consistency control

---

### 9.2 Transactions as Source of Truth

Ledger entries (`account_transactions`) are:

* immutable
* auditable
* the single source of truth for financial history

---

### 9.3 Simplified Scope

The system intentionally does not support:

* automatic transaction reversal
* complex chargebacks
* multiple currencies

---

## 10. Model Limitations

* no support for overdraft (negative balance)
* no transaction limits
* no distinction between account types (individual vs. business)

---

## 11. Conclusion

The domain model defines a system that is:

* transaction-centric
* consistent by design
* simple, yet robust

This foundation supports future evolution without compromising financial integrity.