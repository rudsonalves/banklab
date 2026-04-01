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
