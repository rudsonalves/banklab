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
