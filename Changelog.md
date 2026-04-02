# Changelog

## 2026/04/02 — auth/phase-00

Introduces the **account statement (ledger retrieval) capability** as the first step toward read-oriented financial visibility, along with infrastructure improvements, repository extensions, and enhanced tooling support.

### 1. Application Layer — Get Statement Use Case

* Added `GetStatement` use case with support for:

  * pagination (`limit`, `cursor`, `cursor_id`)
  * time filtering (`from`, `to`)
* Implemented validation rules:

  * non-nil account ID
  * `from <= to`
  * cursor and cursor_id must be provided together
  * limit normalization (default = 50, max = 100)
* Execution flow:

  * validate input
  * ensure account existence (`GetByID`)
  * retrieve transactions via repository
  * map to response DTO
  * build cursor for pagination
* Returns structured result (`Statement`) with:

  * items
  * next cursor (for forward pagination)

### 2. Domain Layer — Repository Evolution

* Extended `AccountRepository` with:

  * `GetTransactions(...)`
* Enables:

  * cursor-based pagination
  * time-range filtering
  * ordered retrieval of ledger entries
* Maintains separation between:

  * write operations (balance changes)
  * read operations (ledger queries)

### 3. Infrastructure Layer — PostgreSQL Implementation

* Implemented `GetTransactions` in:

  * base repository
  * transactional repository (`txRepository`)
* Query characteristics:

  * ordered by `(created_at DESC, id DESC)`
  * cursor-based pagination using tuple comparison
  * optional filters (`from`, `to`)
* Ensures:

  * stable ordering
  * efficient pagination without offset
  * consistency with append-only ledger model

### 4. Delivery Layer — Statement Endpoint

* Added new endpoint:

  * `GET /accounts/{id}/statement`
* Implemented:

  * query parsing (`limit`, `cursor`, `cursor_id`, `from`, `to`)
  * validation helpers:

    * `parseOptionalInt`
    * `parseOptionalTime`
    * `parseOptionalUUID`
* Error handling:

  * `INVALID_DATA → 400`
  * `ACCOUNT_NOT_FOUND → 404`
* Response includes:

  * transaction list
  * pagination cursor (`next_cursor`)
* Maintains consistent API contract (`data` / `error`)

### 5. Handler Refactor & Wiring

* Added `statementUseCase` interface to handler
* Updated constructor to include new dependency
* Registered route in `main.go`:

  * `GET /accounts/{id}/statement`
* Renamed `handleer.go` → `handler.go` (naming correction)

### 6. Data Structures — Statement DTOs

* Introduced:

  * `StatementItemData`
  * `StatementData`
  * `StatementCursorData`
* Provides clear separation between:

  * domain models
  * API response representation

### 7. Test Coverage

#### 7.1 Application Tests

* Full coverage for `GetStatement`:

  * invalid input scenarios
  * default and capped limits
  * account not found
  * repository error propagation
  * successful mapping and cursor generation

#### 7.2 Delivery Tests

* Added handler tests for:

  * invalid query params
  * missing cursor pair
  * account not found
  * success scenario with full query validation
* Ensures:

  * correct HTTP status mapping
  * proper request parsing behavior

#### 7.3 Test Infrastructure

* Extended mocks across use cases to support:

  * `GetTransactions`
* Maintains consistency across all existing tests

### 8. Tooling — Makefile Introduction

* Added `Makefile` with commands:

  * `migration` → run DB migrations
  * `commit` → standardized commit workflow
  * `diff` → staged diff + line count
  * `push` / `pull` → branch-aware Git operations
* Improves developer experience and workflow consistency

### 9. Documentation

* Added `docs/06-implementation.md`:

  * comprehensive description of current system implementation
  * covers architecture, domain, use cases, persistence, and testing
* Serves as a **baseline reference for future phases**, including authentication

### Conclusion

This commit introduces the **read side of the financial ledger (account statement)**, completing the core CRUD + transactional flow with observability over historical operations.

From an architectural standpoint, this is a **crucial milestone**, as it:

* separates read concerns from write flows
* introduces cursor-based pagination (scalable pattern)
* reinforces the ledger model as the source of truth

Additionally, the inclusion of documentation and tooling indicates a transition toward a more **structured and maintainable development process**, which is essential before introducing authentication and authorization layers in subsequent phases.


## 2026/04/02 — account/statement-01

Implements **account statement retrieval (ledger visualization)** with cursor-based pagination, date filtering, and full-stack integration (application, delivery, and persistence). Also introduces developer tooling improvements and formalizes implementation documentation.

### 1. Application Layer — GetStatement Use Case

* Introduced `GetStatement` use case with input:

  * `AccountID`
  * `Limit`
  * `Cursor` + `CursorID` (pagination)
  * `From` / `To` (date filters)
* Implemented validation rules:

  * non-nil account ID
  * `from <= to`
  * cursor consistency (cursor + cursor_id must coexist)
  * limit normalization:

    * default = 50
    * max cap = 100
* Flow:

  * validate input
  * ensure account exists
  * fetch transactions via repository
  * map to statement DTO
  * generate `NextCursor` when applicable
* Returns structured result (`Statement`, `StatementItem`, `StatementCursor`)

### 2. Domain Layer — Repository Expansion

* Extended `AccountRepository` with:

  * `GetTransactions(...)`
* Supports:

  * cursor-based pagination (created_at + id)
  * optional date filtering
* Enables read-side access to the **ledger (account_transactions)**
* Maintains domain-driven contract for data retrieval

### 3. Infrastructure Layer — PostgreSQL Implementation

* Implemented `GetTransactions` in both:

  * base repository
  * transactional repository (`txRepository`)
* Query characteristics:

  * ordered by `(created_at DESC, id DESC)`
  * cursor pagination using tuple comparison:

    ```sql
    (created_at, id) < ($cursor_time, $cursor_id)
    ```
  * optional filters:

    * `from` (>= created_at)
    * `to` (<= created_at)
* Ensures:

  * stable pagination
  * deterministic ordering
  * efficient range queries via indexes
* Defensive handling:

  * limit normalization at repository level

### 4. Delivery Layer — Statement Endpoint

* Added endpoint:

  * `GET /accounts/{id}/statement`
* Implemented query parsing:

  * `limit` (int)
  * `cursor` (RFC3339 datetime)
  * `cursor_id` (UUID)
  * `from` / `to` (RFC3339 datetime)
* Introduced helpers:

  * `parseOptionalInt`
  * `parseOptionalTime`
  * `parseOptionalUUID`
* Validation rules enforced at handler boundary:

  * invalid formats → `400 INVALID_DATA`
  * cursor mismatch → `400 INVALID_DATA`
* Response structure:

  * `StatementData`
  * `StatementItemData`
  * `StatementCursorData`
* Maintains API contract consistency (`data` / `error`)

### 5. Handler & Wiring

* Extended handler to include `statement` use case via interface
* Updated constructor signature (now includes 5 use cases)
* Registered route in `main.go`:

  * `GET /accounts/{id}/statement`
* Renamed file:

  * `handleer.go` → `handler.go` (corrects naming inconsistency)

### 6. Test Coverage

#### 6.1 Application Tests

* Added `get_statement_test.go` covering:

  * invalid account ID
  * default and capped limits
  * account not found
  * repository error propagation
  * successful pagination + cursor generation
* Verifies:

  * correct argument propagation to repository
  * cursor construction logic

#### 6.2 Delivery Tests

* Added handler tests for:

  * invalid query parameters (`from`, cursor mismatch)
  * account not found → `404`
  * successful response with:

    * items mapping
    * cursor serialization

#### 6.3 Test Infrastructure Updates

* Updated mocks across use cases to support:

  * `GetTransactions`
* Ensures compatibility with expanded repository interface

### 7. Documentation

* Added `docs/06-implementation.md`:

  * comprehensive description of:

    * architecture
    * domain model
    * use cases
    * persistence strategy
    * concurrency model
    * test coverage
* This is a **significant improvement in project maturity**, providing a clear implementation reference

### 8. Developer Experience

* Added `Makefile` with utilities:

  * `migration` (run DB migrations)
  * `commit` (standardized commit flow)
  * `diff` (staged diff + line count)
  * `push` / `pull` (branch-aware git operations)
* Improves local workflow consistency and productivity

### Conclusion

This commit introduces **read-side capabilities for the account ledger**, completing the core financial lifecycle:

* write operations (deposit, withdraw, transfer)
* read operations (statement with pagination and filtering)

Key strengths:

* **robust cursor-based pagination design**
* **consistent domain-driven validation**
* **clean separation between layers**
* **strong test coverage across boundaries**

From an architectural perspective, this marks the transition from a purely transactional system to a **fully observable financial ledger**, which is essential for real-world banking systems.


## 2026/04/02 — account/ledger-02

Refines the **deposit operation** to align with ledger-level consistency guarantees by introducing row-level locking and improving error handling semantics.

---

### 1. Application Layer — Deposit Consistency Improvement

* Replaced `GetByID` with `GetByIDForUpdate` inside the deposit transaction flow:

  * ensures **row-level locking (`SELECT ... FOR UPDATE`)**
  * prevents race conditions during concurrent balance updates
* This change aligns deposit behavior with other financial operations (withdraw and transfer), establishing a **uniform concurrency model**

---

### 2. Error Handling Enhancement

* Added explicit handling for `ErrAccountNotFound`:

  * now returned directly without wrapping
* Improves:

  * error transparency
  * correct propagation to delivery layer (HTTP mapping)
* Other errors remain wrapped with contextual information for observability

---

### 3. Concurrency and Ledger Integrity

* Deposit now guarantees:

  * **read-after-lock semantics**
  * no stale balance reads
  * safe concurrent updates under high contention
* This is a critical improvement for financial correctness, especially in scenarios with simultaneous deposits and withdrawals

---

### 4. Architectural Consistency

* Aligns deposit with previously implemented operations:

  * withdraw → already transactional + safe
  * transfer → deterministic locking + multi-entity safety
* Establishes a **cohesive ledger model**, where all balance mutations:

  * occur within transactions
  * operate on locked rows
  * respect ACID guarantees

---

### Conclusion

This commit addresses a subtle but important gap in the deposit flow by introducing **proper locking semantics and precise error handling**.

From a technical standpoint, this is a **high-value correction**, ensuring that even simple credit operations adhere to the same robustness standards required for a reliable financial ledger.


## 2026/04/02 — account/ledger-01

Introduces the **transfer operation as part of the account ledger capabilities**, including transactional coordination, row-level locking, API exposure, and extensive test coverage. Also refines repository contracts and strengthens domain invariants. 

### 1. Application Layer — Transfer Use Case

* Added `Transfer` use case with input:

  * `FromAccountID`
  * `ToAccountID`
  * `Amount`
* Implemented full transactional workflow:

  * `BeginTx → Lock दोनों accounts → Validate → Debit → Credit → Commit`
  * rollback guaranteed via deferred handler
* Introduced **deterministic locking strategy**:

  * accounts locked via `GetByIDForUpdate`
  * ordered by UUID (`orderedUUIDs`) to reduce deadlock risk
* Explicit separation between:

  * domain validation (`CanTransfer`, `CanDeposit`)
  * persistence operations (`DecreaseBalance`, `UpdateBalance`)
* Returns enriched result (`TransferResult`) with both balances after operation

### 2. Domain Layer — Ledger Rules Consolidation

* Added `Account.CanTransfer(amount, destinationID)`:

  * prevents self-transfer
  * reuses withdraw validation chain
* Introduced new domain error:

  * `ErrSameAccountTransfer`
* Reinforces domain as the authoritative layer for financial rules

### 3. Repository Contract Evolution

* Added `GetByIDForUpdate`:

  * enables row-level locking (`SELECT ... FOR UPDATE`)
* Clarified `DecreaseBalance` contract:

  * assumes account existence
  * enforces balance constraint at DB level
* These changes are critical for **multi-entity consistency in ledger operations**

### 4. Infrastructure Layer — PostgreSQL

* Implemented `GetByIDForUpdate` in both:

  * base repository
  * transactional repository (`txRepository`)
* Uses:

  ```sql
  SELECT ... FOR UPDATE
  ```
* Guarantees:

  * row-level locking
  * concurrency safety
  * compatibility with ordered locking strategy
* Strengthens ACID guarantees for transfer operations

### 5. Delivery Layer — Transfer Endpoint

* Added endpoint:

  * `POST /accounts/transfer`
* Introduced DTOs:

  * `TransferRequest`
  * `TransferData`
* Implemented validation and error mapping:

  * `INVALID_DATA → 400`
  * `INVALID_AMOUNT → 400`
  * `SAME_ACCOUNT_TRANSFER → 400`
  * `ACCOUNT_NOT_FOUND → 404`
  * `INSUFFICIENT_BALANCE → 422`
  * `ACCOUNT_INACTIVE → 422`
* Maintains API response standard (`data` / `error`)

### 6. Handler & Wiring

* Extended handler to support `transfer` use case via interface
* Updated constructor signature to include new dependency
* Registered route in `main.go`:

  * `POST /accounts/transfer`
* Updated integration setup to reflect new handler signature

### 7. Test Coverage

#### 7.1 Application Tests

* Comprehensive test suite covering:

  * invalid inputs (UUID, amount, same account)
  * account not found (source/destination)
  * insufficient balance
  * inactive destination
  * debit/credit failures
  * commit failure
  * success scenario
* Validates:

  * rollback behavior
  * commit execution
  * deterministic locking order

#### 7.2 Domain Tests

* Added tests for `CanTransfer`:

  * self-transfer
  * invalid amount
  * inactive account
  * insufficient balance
  * success

#### 7.3 Delivery Tests

* Added handler test for:

  * `SAME_ACCOUNT_TRANSFER → 400`

#### 7.4 Test Infrastructure

* Extended mocks:

  * support for `GetByIDForUpdate`
  * lock order tracking
  * transactional behavior validation

### 8. Test Refinements

* Updated deposit tests:

  * now assert returned DB balance instead of computed value
* Improved withdraw tests:

  * stricter assertions on repository interaction (no unintended calls)

### Conclusion

This commit elevates the system to support **ledger-level operations via account transfers**, with strong guarantees around:

* **atomicity (single transaction)**
* **consistency (domain + DB constraints)**
* **concurrency safety (row-level locking + deterministic ordering)**

From an architectural perspective, this is a **robust and production-ready implementation**, particularly due to its explicit handling of deadlocks and strict separation of responsibilities across layers.


## 2026/04/02 — account/transfer-01

Implements the **transfer operation** as a fully transactional, concurrency-safe use case, including domain validation, deterministic locking strategy, HTTP exposure, and comprehensive test coverage.

### 1. Application Layer — Transfer Use Case

* Introduced `Transfer` use case with input:

  * `FromAccountID`
  * `ToAccountID`
  * `Amount`
* Enforced validations:

  * non-nil UUIDs
  * amount > 0
  * source and destination must be different
* Implemented **single transaction orchestration**:

  * `BeginTx → Lock Accounts → Validate → Debit → Credit → Commit`
  * automatic rollback on failure
* Introduced **deterministic locking strategy**:

  * accounts are locked using `GetByIDForUpdate`
  * ordered by UUID (`orderedUUIDs`) to **prevent deadlocks** in concurrent transfers
* Clear separation between:

  * domain validation (`CanTransfer`, `CanDeposit`)
  * persistence execution (`DecreaseBalance`, `UpdateBalance`)
* This is a **high-quality implementation**, particularly due to explicit deadlock mitigation and strict transactional boundaries 

### 2. Domain Layer — Transfer Rules

* Added `Account.CanTransfer(amount, destinationID)`:

  * prevents same-account transfers
  * reuses withdraw validation (`CanWithdraw`)
* Introduced new domain error:

  * `ErrSameAccountTransfer`
* Reinforces domain as the **single source of truth for business rules**, eliminating duplication across use cases

### 3. Repository Contract Evolution

* Added `GetByIDForUpdate` to `AccountRepository`:

  * enables row-level locking (`SELECT ... FOR UPDATE`)
* Clarified `DecreaseBalance` contract:

  * does not validate existence
  * enforces balance constraint at DB level
* These changes are essential to support **safe concurrent transfers**

### 4. Infrastructure Layer — PostgreSQL

* Implemented `GetByIDForUpdate` using:

  ```sql
  SELECT ... FOR UPDATE
  ```
* Available both in base repository and transactional (`txRepository`)
* Guarantees:

  * row-level locking
  * prevention of race conditions
  * compatibility with ordered locking strategy
* Strengthens ACID guarantees for multi-entity operations

### 5. Delivery Layer — Transfer Endpoint

* Added new endpoint:

  * `POST /accounts/transfer`
* Introduced request/response DTOs:

  * `TransferRequest`
  * `TransferData`
* Implemented full validation and error mapping:

  * `INVALID_DATA → 400`
  * `INVALID_AMOUNT → 400`
  * `SAME_ACCOUNT_TRANSFER → 400`
  * `ACCOUNT_NOT_FOUND → 404`
  * `INSUFFICIENT_BALANCE → 422`
  * `ACCOUNT_INACTIVE → 422`
* Maintains API response contract consistency (`data` / `error`)

### 6. Handler Refactor

* Extended handler to include `transfer` use case via interface
* Updated constructor and dependency injection
* Registered route in `main.go`:

  * `POST /accounts/transfer`
* Preserves decoupling and adherence to layered architecture

### 7. Test Coverage

#### 7.1 Application Tests

* Extensive test suite covering:

  * invalid inputs (UUID, amount, same account)
  * account not found (source/destination)
  * insufficient balance
  * inactive destination account
  * debit failure
  * credit failure
  * commit failure
  * success scenario
* Validates:

  * rollback behavior
  * commit correctness
  * deterministic locking order

#### 7.2 Domain Tests

* Added tests for `CanTransfer`:

  * same account
  * invalid amount
  * inactive account
  * insufficient balance
  * success

#### 7.3 Delivery Tests

* Added handler test:

  * `SAME_ACCOUNT_TRANSFER → 400` mapping

#### 7.4 Test Infrastructure

* Extended mocks:

  * support for `GetByIDForUpdate`
  * tracking of lock order
  * transactional behavior validation

### 8. Additional Improvements

* Minor refinement in deposit test:

  * now validates returned DB balance instead of computed value
* Improved withdraw test assertions:

  * ensures no unintended repository calls occur

### Conclusion

This commit introduces the **most complex financial operation (transfer)** with a robust and production-grade design.

Key highlights:

* **Deterministic locking to prevent deadlocks**
* **Strict transactional integrity**
* **Domain-driven validation**
* **Database-level safety guarantees**

From an architectural standpoint, this is a **mature and well-executed implementation**, significantly elevating the reliability of the system’s financial core.


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
