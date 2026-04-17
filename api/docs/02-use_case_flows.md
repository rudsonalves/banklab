# Use Case Flows — Bank API

## 1. Overview

Use cases represent the execution of system operations.

Each flow describes:

* input
* validations
* execution sequence
* system effects
* possible failures

All flows must comply with:

* domain invariants
* atomicity
* balance consistency

---

## 2. Create Customer

### Input

* name
* cpf
* email

### Flow

1. validate input format
2. check if CPF already exists
3. check if email already exists
4. create Customer
5. persist to database

### Output

* created customer data

### Possible Errors

* CPF already registered
* email already registered
* invalid input data

---

## 3. Approve User

### Input

* userId
* administrator credentials

### Flow

1. start transaction
2. load user with lock (`FOR UPDATE`)
3. verify user exists
4. validate current status = `pending`
5. update status → `active`
6. create Account with status `active` (balance = 0)
7. persist changes
8. commit transaction

### Output

* approved user data with status `active`
* account data created atomically

### Behavior

* approval and account creation are **atomic**
* no partial state transitions are allowed
* if any step fails, the entire operation is rolled back

### Possible Errors

* user not found
* user already active
* associated customer does not exist
* account creation failure

---

## 4. Open Account

### Input

* customerId
* authenticated user with `active` status

### Flow

1. verify user is `active`
2. verify customer exists
3. generate account number
4. generate branch
5. create Account with status `active`
6. initial balance = 0
7. persist to database

### Output

* created account data

### Possible Errors

* user not active
* customer not found

---

## 5. Get Balance

### Input

* accountId

### Flow

1. verify account exists
2. return current balance

### Output

* balance

### Possible Errors

* account not found

---

## 6. Deposit

### Input

* accountId
* amount

### Flow

1. validate amount > 0
2. verify account exists
3. verify account is active
4. start transaction
5. update balance (balance + amount)
6. record ledger entry in `account_transactions` (type: deposit)
7. persist changes
8. commit transaction

### Output

* updated balance

### Possible Errors

* invalid amount
* account not found
* inactive account

---

## 7. Withdraw

### Input

* accountId
* amount

### Flow

1. validate amount > 0
2. verify account exists
3. verify account is active
4. start transaction
5. verify sufficient balance
6. update balance (balance - amount)
7. record ledger entry in `account_transactions` (type: withdraw)
8. persist changes
9. commit transaction

### Output

* updated balance

### Possible Errors

* invalid amount
* insufficient balance
* account not found
* inactive account

---

## 8. Transfer

### Input

* fromAccountId
* toAccountId
* amount
* idempotencyKey (optional)

### Flow

1. validate amount > 0
2. validate accounts are different
3. validate idempotency (if applicable)
4. verify accounts exist
5. verify both accounts are active
6. start transaction
7. lock both accounts for update
8. verify sufficient balance in source account
9. debit source account
10. credit destination account
11. record two ledger entries in `account_transactions`:

    * `transfer_out` (source, with idempotencyKey)
    * `transfer_in` (destination)
    * both sharing the same referenceId
12. persist changes
13. commit transaction

### Output

* updated source account balance

### Possible Errors

* invalid amount
* same account transfer
* insufficient balance
* account not found
* inactive account
* duplicate request (idempotency)

---

## 9. Get Statement

### Input

* accountId
* limit (optional)
* cursor (optional)
* cursorId (optional)
* from (optional)
* to (optional)

### Flow

1. verify account exists
2. retrieve account transactions
3. sort by date (descending)
4. apply cursor-based pagination (`created_at` + `id`)

### Output

* list of transactions

### Possible Errors

* account not found

---

## 10. Flow Guarantees

All balance-changing operations must ensure:

### 10.1 Atomicity

* the operation is fully completed or not executed

### 10.2 Consistency

* balance never enters an invalid state

### 10.3 Isolation

* concurrent operations do not corrupt data

### 10.4 Durability

* once committed, the operation persists

---

## 11. Common Patterns

### 11.1 Early Validation

Before starting a transaction:

* validate input
* validate account state

---

### 11.2 Explicit Transactions

All financial operations:

* start a transaction
* execute changes
* complete with commit or rollback

---

### 11.3 Mandatory Recording

Every balance change:

* generates a corresponding Transaction

---

### 11.4 Idempotency

When present:

* prevents duplication
* ensures safe retries

---

## 12. Conclusion

The defined flows ensure:

* predictable behavior
* operational consistency
* alignment with the domain model

This document bridges domain rules and system execution.