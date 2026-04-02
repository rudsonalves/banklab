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
