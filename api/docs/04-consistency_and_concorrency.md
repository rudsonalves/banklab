# Consistency and Concurrency Strategy — Bank API

## 1. Overview

System consistency is ensured through:

* ACID transactions in PostgreSQL
* explicit concurrency control
* application-level validations

The system adopts a model of:

> **strong and immediate consistency**

There is no tolerance for eventual inconsistency in financial operations.

---

## 2. Objectives

Ensure that:

* balance is never corrupted
* concurrent operations do not introduce inconsistencies
* no invalid intermediate states exist
* duplicate operations are prevented

---

## 3. Concurrency Model

### 3.1 Expected Concurrency

The system must support:

* multiple simultaneous deposits
* concurrent withdrawals on the same account
* simultaneous transfers involving the same accounts

---

### 3.2 Main Risks

Without proper control:

* race conditions
* unintended negative balance
* duplicate operations

---

## 4. Adopted Strategy

### 4.1 Explicit Transactions

All balance-changing operations must:

* start a transaction
* execute all steps within it
* complete with commit or rollback

---

### 4.2 Pessimistic Locking

Using:

```sql id="jhlhnr"
SELECT ... FOR UPDATE
```

Purpose:

* lock involved records
* prevent concurrent reads/writes

---

## 5. Critical Operations

---

### 5.1 Deposit

#### Strategy

1. start transaction
2. lock account:

   ```sql id="jz1f6y"
   SELECT * FROM accounts WHERE id = ? FOR UPDATE;
   ```
3. update balance
4. insert ledger entry in `account_transactions`
5. commit

---

### 5.2 Withdrawal

#### Strategy

1. start transaction
2. lock account
3. validate balance
4. update balance
5. insert ledger entry in `account_transactions`
6. commit

---

### 5.3 Transfer

#### Strategy

1. start transaction
2. lock source account
3. lock destination account
4. validate source balance
5. update source balance
6. update destination balance
7. insert two entries in `account_transactions` (`transfer_out` and `transfer_in`, sharing the same `reference_id`)
8. commit

---

### 5.4 User Approval (Lifecycle)

#### Objective

Transition user from `pending` → `active` with atomic account creation.

#### Strategy

1. start transaction
2. load user with lock (`FOR UPDATE`):

   ```sql id="d1c0og"
   SELECT * FROM users WHERE id = ? FOR UPDATE;
   ```
3. validate status = `pending`
4. update status → `active`
5. generate account number
6. insert account with status `active` (balance = 0)
7. update `customer_id` if necessary
8. commit

#### Guarantees

* **atomicity**: status update and account creation occur in a single transaction
* **isolation**: user is locked during approval
* **consistency**: invariants ensure valid state
* **durability**: once committed, the change persists

#### Invariants

* no approval without account creation
* no account without an active user
* transition `pending` → `active` is irreversible at this stage

#### Avoided Risks

* approval without account creation
* duplicate approval
* concurrent approval of the same user

---

## 6. Lock Ordering (CRITICAL)

To prevent deadlocks:

> Accounts must always be locked in a consistent order.

### Rule

* sort account IDs
* lock the lower ID first
* then lock the higher ID

---

## 7. Transaction Isolation

### Recommended Level

```text id="bnylrd"
READ COMMITTED (PostgreSQL default)
```

Reason:

* combined with `FOR UPDATE`, it provides sufficient safety
* avoids unnecessary overhead

---

### Alternative (if needed)

```text id="xg4l5z"
REPEATABLE READ
```

Used only if more complex consistency issues arise.

---

## 8. Idempotency

### Objective

Prevent duplication in:

* network retries
* repeated requests

---

### Strategy

* use of `idempotency_key`
* partial unique index in the ledger

```sql id="y0nm0g"
(account_id, idempotency_key)
```

Current scope:

* applied to `transfer_out` (source account)

---

### Behavior

* repeated operation → not executed again
* replay returns the original historical result
* replay is reconstructed from the ledger (`account_transactions`), independent of current `accounts.balance`

---

## 9. System Guarantees

### 9.1 Atomicity

* operations are fully completed or not executed

---

### 9.2 Consistency

* domain rules are always enforced

---

### 9.3 Isolation

* concurrent operations do not interfere with active ones

---

### 9.4 Durability

* committed operations are permanent

---

## 10. Failure and Recovery

### 10.1 Transaction Failure

* automatic rollback
* no partial state is persisted

---

### 10.2 Connection Timeout

* operation must be considered incomplete
* client may retry using idempotency

---

## 11. Intentional Decisions

### 11.1 Pessimistic Locking

Chosen instead of optimistic locking.

Reason:

* higher predictability
* lower risk in financial systems

---

### 11.2 Immediate Consistency

No use of:

* queues
* asynchronous processing
* eventual consistency

---

### 11.3 Database as Coordinator

PostgreSQL is responsible for:

* effective operation serialization
* concurrency control

---

## 12. Limitations

* no distributed database support
* no distributed concurrency control
* scalability is constrained by the database

---

## 13. Future Evolution

Possible improvements:

* selective use of optimistic locking
* queues for non-critical operations
* account sharding
* read replication

---

## 14. Conclusion

The adopted strategy ensures:

* strong consistency
* predictability
* safety in financial operations

This approach prioritizes correctness over premature performance optimization, which is essential for this type of system.
