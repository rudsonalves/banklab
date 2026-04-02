# Changelog

## 2026/04/02 — account/withdraw-01

Implements the **withdraw operation** with strong domain validation, transactional safety, and consistent API exposure. Additionally refines balance handling semantics and consolidates domain rules within the entity.

### 1. Application Layer — Withdraw Use Case

* Introduced `Withdraw` use case with input contract (`AccountID`, `Amount`)
* Enforced validations:

  * non-nil UUID
  * positive amount
* Delegated business rules to domain (`Account.CanWithdraw`)
* Implemented transactional flow:

  * `BeginTx → GetByID → DecreaseBalance → Commit`
  * rollback on any failure
* Ensures atomicity and consistency for debit operations, aligned with financial invariants 

### 2. Domain Layer — Business Rule Consolidation

* Introduced domain methods:

  * `Account.CanDeposit`
  * `Account.CanWithdraw`
* Centralized validation logic:

  * account must be active
  * amount must be > 0
  * sufficient balance required for withdraw
* Added new domain error:

  * `ErrInsufficientBalance`
* This is a **notable improvement in design quality**, moving rule enforcement out of the application layer into the domain, increasing cohesion and correctness

### 3. Repository Contract Evolution

* Updated `AccountRepository`:

  * `UpdateBalance` now returns updated balance (`int64`)
  * introduced `DecreaseBalance` for debit operations
* Separation of credit vs debit operations improves:

  * semantic clarity
  * safety (explicit constraint on balance ≥ amount)

### 4. Infrastructure Layer — PostgreSQL Enhancements

* Updated `UpdateBalance` to use `RETURNING balance`

  * eliminates need for manual balance mutation in memory
* Implemented `DecreaseBalance` with safeguard:

  ```sql
  UPDATE accounts
  SET balance = balance - $1
  WHERE id = $2 AND balance >= $1
  ```
* Guarantees:

  * no negative balance at DB level
  * concurrency-safe debit operation
* Mirrors domain invariant enforcement at persistence level (defensive design) 

### 5. Application Refinement — Deposit

* Refactored deposit flow:

  * now uses `Account.CanDeposit`
  * uses returned balance from repository instead of manual increment
* Removes duplication and aligns deposit with new domain-centric validation approach

### 6. Delivery Layer — Withdraw Endpoint

* Added new endpoint:

  * `POST /accounts/{id}/withdraw`
* Implemented request parsing (`WithdrawRequest`)
* Error mapping aligned with API standard:

  * `INVALID_AMOUNT → 400`
  * `INVALID_DATA → 400`
  * `ACCOUNT_NOT_FOUND → 404`
  * `INSUFFICIENT_BALANCE → 422`
  * `ACCOUNT_INACTIVE → 422`
* Maintains response contract consistency (`data` / `error`)

### 7. Handler & Wiring

* Extended handler to include `withdraw` use case via interface
* Updated constructor and dependency injection
* Registered new route in `main.go`:

  * `POST /accounts/{id}/withdraw`
* Maintains modular composition aligned with layered architecture

### 8. Test Coverage

#### 8.1 Application Tests

* Added comprehensive tests for withdraw:

  * invalid amount
  * invalid account ID
  * account not found
  * insufficient balance
  * repository failure
  * success path
* Validates transactional behavior (commit/rollback)

#### 8.2 Domain Tests

* Added unit tests for:

  * `CanDeposit`
  * `CanWithdraw`
* Ensures correctness of core invariants at domain level

#### 8.3 Delivery Tests

* Added handler test:

  * `INSUFFICIENT_BALANCE → 422` mapping

#### 8.4 Test Infrastructure

* Extended mocks:

  * support for `DecreaseBalance`
  * updated `UpdateBalance` signature
* Adjusted integration setup for handler constructor changes

### Conclusion

This commit introduces the **withdraw operation as a first-class financial capability**, with robust safeguards at both domain and database levels.

The most relevant architectural improvement is the **migration of business rules into the domain layer**, combined with **database-level enforcement for balance constraints**, resulting in a highly reliable and consistent implementation.


## 2026/04/02 — account/deposit-01

Implements the **deposit operation** across all architectural layers, introducing transactional consistency, domain validations, HTTP exposure, and full test coverage (unit + integration).

---

### 1. Application Layer — Deposit Use Case

* Introduced `Deposit` use case with explicit input contract (`AccountID`, `Amount`)
* Enforced domain invariants:

  * non-zero/valid UUID
  * amount > 0
  * account must be active
* Implemented **transactional control at the application layer**, ensuring:

  * `BeginTx → GetByID → UpdateBalance → Commit`
  * automatic rollback on failure
* Error wrapping preserves infrastructure context while exposing domain errors
* Aligns with transactional orchestration responsibility defined in the application layer 

---

### 2. Domain Layer Enhancements

* Expanded `AccountRepository` contract:

  * `GetByID`
  * `UpdateBalance`
  * `BeginTx`
* Introduced `Tx` interface to support transactional operations
* Added new domain errors:

  * `ErrInvalidAmount`
  * `ErrAccountNotFound`
  * `ErrAccountInactive`
* Strengthens enforcement of domain invariants for financial operations 

---

### 3. Infrastructure Layer — PostgreSQL Implementation

* Implemented transactional repository (`txRepository`) using `pgx.Tx`
* Added support for:

  * account retrieval (`SELECT`)
  * atomic balance update (`UPDATE balance = balance + $1`)
  * transaction lifecycle (`BeginTx`, `Commit`, `Rollback`)
* Explicit handling for:

  * `ErrNoRows → ErrAccountNotFound`
  * prevention of nested transactions
* Ensures ACID compliance and consistency guarantees as required by financial operations 

---

### 4. Delivery Layer — HTTP Endpoint

* Added new endpoint:

  * `POST /accounts/{id}/deposit`
* Implemented request parsing (`DepositRequest`) and path param validation
* Mapped domain errors to standardized HTTP responses:

  * `INVALID_AMOUNT → 400`
  * `INVALID_DATA → 400`
  * `ACCOUNT_NOT_FOUND → 404`
  * `ACCOUNT_INACTIVE → 422`
* Response structure follows the defined API contract (`data` / `error`) 
* Extended handler with `deposit` use case via interface injection (improves decoupling)

---

### 5. Main Wiring (Composition Root)

* Integrated account module alongside customer module
* Registered new routes:

  * `POST /accounts`
  * `POST /accounts/{id}/deposit`
* Proper dependency wiring:

  * shared repository
  * separate use cases (`CreateAccount`, `Deposit`)
* Reinforces modular monolith structure and clear separation of concerns 

---

### 6. Test Coverage

#### 6.1 Application Tests

* Full coverage for deposit use case:

  * invalid amount
  * account not found
  * inactive account
  * repository failure
  * successful execution
* Validates transaction behavior (commit vs rollback)

#### 6.2 Delivery Tests

* Added handler tests for deposit:

  * error mapping (`ACCOUNT_INACTIVE`)
  * correct HTTP status and response structure

#### 6.3 Integration Test

* End-to-end validation with PostgreSQL:

  * schema setup (customers, accounts, sequence)
  * data seeding
  * HTTP request execution
  * verification of:

    * response payload
    * persisted balance in database
* Confirms real transactional consistency and persistence correctness

---

### 7. Test Infrastructure Adjustments

* Extended mocks to support new repository methods (`GetByID`, `UpdateBalance`, `BeginTx`)
* Introduced transactional mock (`txMock`) to simulate commit/rollback behavior

---

### Conclusion

This commit introduces a **critical financial operation (deposit)** with proper domain validation, strong transactional guarantees, and full-stack test coverage.

The implementation is technically sound and aligned with the system’s architectural principles, particularly regarding **transaction control in the application layer and consistency at the persistence level**.


## 2026/04/02 — tests/general-01

Introduces comprehensive test coverage for the account module across domain, application, and delivery layers, along with improvements to migration structure and handler decoupling.

### 1. Database and Migrations

* Renamed `migrations/main_tables.sql` to `db/schema.sql` to better reflect its role as the base schema definition
* Added versioned migration files for account number sequence:

  * `000001_account_number_sequence.up.sql` with proper `START WITH` and `INCREMENT BY`
  * `000001_account_number_sequence.down.sql` for rollback support
* Removed legacy `001_account_number_sequence.sql` to enforce consistent migration versioning
* Aligns migration strategy with incremental and reversible schema evolution 

### 2. Application Layer Tests (`create_account_test.go`)

* Added full test suite for `CreateAccount` use case, covering:

  * invalid input (nil UUID)
  * customer not found
  * repository error propagation
  * account number generation failure
  * persistence failure
  * successful account creation
* Validates interaction contracts (call counts) with repositories
* Ensures strict adherence to use case flow defined in “Abrir Conta” 

### 3. Delivery Layer Tests (`account_handler_test.go`)

* Introduced HTTP handler tests for account creation endpoint:

  * invalid JSON → `400 INVALID_REQUEST`
  * invalid UUID → `400 INVALID_DATA`
  * customer not found → `404 CUSTOMER_NOT_FOUND`
  * success → `201 Created` with full response validation
* Verifies response structure consistency with defined error/response pattern 
* Ensures proper input validation and use case invocation boundaries

### 4. Domain Layer Tests (`account_test.go`)

* Added tests for `NewAccount`:

  * validation of invalid customer ID
  * successful entity creation with correct defaults:

    * status = active
    * balance = 0
    * non-zero ID and timestamp
* Reinforces domain invariants for Account entity 

### 5. Handler Refactor (`handleer.go`)

* Introduced `createAccountUseCase` interface to decouple handler from concrete implementation
* Updated constructor to depend on interface instead of struct
* Improves testability and aligns with dependency inversion principle described in architecture 

### Conclusion

This commit establishes a solid testing foundation for the account module, ensuring correctness across all architectural layers, improving decoupling, and standardizing database migrations. It significantly increases confidence in core flows while aligning implementation with the system’s architectural and domain principles.


## 2026-04-01 — customer/create-01

### feat(account): implement account creation flow with layered architecture

Introduces the complete **account creation use case**, following the project’s layered architecture (Delivery → Application → Domain → Infrastructure) .

### **Application Layer**

* Added `CreateAccount` use case:

  * Validates `CustomerID`
  * Checks customer existence via `CustomerRepository`
  * Generates sequential account number
  * Applies default branch (`0001`)
  * Creates and persists `Account`
* Optional rule prepared for **one account per customer** (commented for future enforcement)

### **Domain Layer**

* Introduced `Account` entity with:

  * Status management (`active`, `inactive`, `blocked`)
  * Invariants (non-nil customer, initial balance = 0)
* Defined repository contracts:

  * `AccountRepository`
  * `CustomerRepository`
* Added domain errors:

  * `ErrInvalidData`
  * `ErrCustomerNotFound`

### **Infrastructure Layer (PostgreSQL)**

* Implemented `AccountRepository`:

  * `Create`
  * `ExistsByCustomerID`
  * `NextAccountNumber` using DB sequence
* Added migration:

  * `account_number_seq` for deterministic account numbering
  * Unique constraint on `accounts.number`

### **Customer Module Update**

* Extended `CustomerRepository` with `Exists` method
* Implemented existence check in PostgreSQL adapter

### **Delivery Layer (HTTP)**

* Added handler `CreateAccount`:

  * Parses and validates request payload
  * Converts `customer_id` to UUID
  * Maps domain errors to HTTP responses
* Introduced DTOs:

  * `CreateAccountRequest`
  * `AccountData`
* Implemented standardized response structure:

  * `{ data, error }` pattern aligned with API guidelines 

### **Observations (Technical Assessment)**

* The implementation is **structurally sound** and adheres well to the defined architecture and domain boundaries.
* The use case correctly enforces **existence validation before persistence**, aligning with the “Abrir Conta” flow .
* Delegating account number generation to the database is a **strong consistency decision**, avoiding race conditions.
* The optional uniqueness rule (one account per customer) is appropriately deferred—this avoids premature constraint hardening.

### **Minor Issues / Improvements**

* Typo in file names:

  * `accont_request.go` → `account_request.go`
  * `handleer.go` → `handler.go`
* Missing explicit transaction handling (acceptable at this stage, but should be considered for future multi-step operations)

### **Result**

This commit establishes the **foundation for the Account domain**, enabling account creation with proper validation, persistence, and API exposure, while remaining consistent with the system’s transactional and architectural principles.
