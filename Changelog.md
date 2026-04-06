# Changelog

## 2026/04/05 — check/adjustments-03

Introduces a centralized and extensible error registration mechanism, replacing hardcoded mappings with a registry-based approach. Also standardizes error initialization across modules and ensures consistent behavior in both runtime and test environments.

### 1. Bootstrap Initialization

* Added `bootstrap.Init()` in `main.go` to ensure all application errors are registered at startup
* Introduced `internal/bootstrap` package to centralize cross-module initialization
* Created `bootstrap.RegisterErrors()` to orchestrate error registration across:

  * account
  * customer
  * auth

### 2. Modular Error Registration per Domain

* Added `errors_registry.go` in:

  * `account/application`
  * `auth/application`
  * `customer/application`
* Each module now:

  * defines its own domain-to-AppError mappings
  * registers errors using `RegisterDomainError`
* Improves:

  * modularity
  * separation of concerns
  * scalability for new domains

### 3. Shared Error System Refactor

* Replaced large `switch` in `MapError` with a dynamic registry-based system
* Introduced:

  * `entry` struct for matching logic
  * global `registry` slice
  * `Register(match, AppError)` function
* Implemented validation:

  * prevents duplicate error code registration with conflicting definitions
* Added `RegisterDomainError` helper:

  * simplifies domain error mapping using `errors.Is`

### 4. Default Error Handling Improvements

* Introduced `internalError()` helper for fallback responses
* Standardized behavior:

  * `nil` error returns internal error
  * unmatched errors default to internal error
* Added default registration for `ErrInvalidRequest` via `init()`

### 5. Test Environment Consistency

* Added `TestMain` setup in:

  * `account/delivery`
  * `auth/delivery`
* Ensures error registry is initialized before test execution
* Aligns test behavior with application runtime

### 6. Architectural Impact

* Eliminates coupling between shared error mapper and domain packages
* Moves system from:

  * static, tightly coupled error mapping
  * to dynamic, pluggable registration model
* Enables:

  * independent evolution of domains
  * safer extension of error catalog
  * cleaner application boundaries

### Conclusion

This refactor significantly improves the robustness and maintainability of the error handling strategy. By introducing a registry-based approach with centralized bootstrap, the system becomes more modular, predictable, and aligned with clean architecture principles, especially regarding dependency direction and domain isolation.


## 2026/04/05 — check/adjustments-02

Refactors and consolidates the error handling strategy across the application, eliminating duplication, enforcing domain-driven error definitions, and standardizing HTTP responses through shared infrastructure.

### 1. Centralized Error Handling (Shared Layer)

* Introduced a unified error system in `internal/shared/errors`:

  * Added `AppError` struct with `Code`, `Message`, and `Status`
  * Created `codes.go` to define all standardized error codes
  * Implemented `mapper.go` with `MapError(err)` to translate domain errors into HTTP responses
* Removed legacy error definitions and factories (`errors.go`)
* Introduced `ErrInvalidRequest` as a shared sentinel error

This change establishes a single source of truth for error mapping and eliminates scattered error logic across layers

### 2. Domain-Driven Error Definitions

* Moved all auth-related errors into `internal/auth/domain/errors.go`:

  * `ErrUnauthorized`, `ErrInvalidToken`, `ErrInvalidCredentials`, etc.
* Updated all application use cases to rely on domain errors instead of local definitions:

  * `login_user.go`
  * `register_user.go`
  * `get_current_user.go`
* Removed duplicated error declarations from application layer

This enforces a clean architecture boundary where domain owns business semantics

### 3. Removal of Application-Level Error Categorization

* Deleted:

  * `internal/account/application/errors.go`
  * `internal/customer/application/errors.go`
* Removed all categorization helpers (`CategorizeError`, `ValidationField`)
* Eliminated `mapAccountError` and similar mapping logic from delivery layer

Error interpretation is now fully delegated to the shared mapper

### 4. HTTP Layer Standardization

* Refactored `internal/shared/http/response.go`:

  * Replaced `WriteSuccess` with `WriteJSON`
  * Simplified `WriteError` to accept `AppError`
  * Introduced `Response` and `ErrorBody` structs
* Updated all handlers to use:

  * `sharedhttp.WriteJSON`
  * `sharedhttp.WriteError(sharederrors.MapError(err))`

This ensures consistent response structure across all endpoints

### 5. Delivery Layer Simplification

* Removed custom response wrappers:

  * `internal/account/delivery/response.go`
  * custom response structs in customer delivery
* Replaced manual error construction with centralized mapping:

  * Removed inline `NewErrorWithDetails`
  * Removed duplicated HTTP status handling
* Updated handlers in:

  * account
  * auth
  * customer
    to rely exclusively on shared error and response utilities

This significantly reduces boilerplate and improves maintainability

### 6. Authentication Flow Adjustments

* Updated `RequireUser` to return `error` instead of `AppError`
* Standardized unauthorized handling using `domain.ErrUnauthorized`
* Updated JWT middleware to use:

  * `sharederrors.MapError(authdomain.ErrUnauthorized)`
  * `sharederrors.MapError(authdomain.ErrInvalidToken)`

This aligns authentication errors with the global error strategy

### 7. Test Adjustments

* Updated tests to use domain errors instead of application-level errors:

  * auth use cases and handlers
* Adjusted expected error codes:

  * `INSUFFICIENT_BALANCE` → `INSUFFICIENT_FUNDS`
* Updated HTTP response tests to reflect new API:

  * `WriteSuccess` → `WriteJSON`
  * `WriteError` signature change

### 8. Minor Cleanups

* Removed unused imports (`errors` in multiple files)
* Adjusted handler dependencies and imports for consistency
* Updated response writing logic across all modules

### Conclusion

This commit delivers a substantial architectural improvement by **centralizing error handling, enforcing domain ownership of business errors, and standardizing HTTP responses**. The result is a cleaner, more maintainable codebase with reduced duplication and clearer separation of concerns, particularly between domain, application, and delivery layers.


## 2026/04/04 — check/small_adjustments-01

Applies a set of **targeted architectural refinements** focused on access control centralization, error standardization, and safer authentication handling. These changes improve consistency, reduce duplication, and strengthen boundary responsibilities across modules.

### 1. Changelog Update

* Updated reference date for `auth/phase-06` to reflect the correct delivery timeline
* Ensures historical accuracy and alignment with recent changes

### 2. Access Control Centralization (Account Module)

* Introduced `access_policy.go` with:

  * `CanAccessCustomer`
  * `CanAccessAccount`
* Consolidates authorization logic within the **account application layer**, removing cross-module coupling with auth
* Replaced previous scattered checks with unified policy usage across:

  * `CreateAccount`
  * `Deposit`
  * `Withdraw`
  * `Transfer`
  * `GetStatement`
* Improves cohesion and enforces a **single source of truth for access rules**

### 3. Removal of Legacy Authorization Logic (Auth Module)

* Removed:

  * `internal/auth/application/authorization.go`
  * `internal/auth/domain/authorization.go`
* Eliminates duplicated responsibility between modules
* Clarifies architectural boundary:

  * **Auth → identity**
  * **Account → authorization rules**

### 4. Error Categorization Standardization

* Introduced `application/errors.go` (account module):

  * defines `ErrorCategory`
  * implements `CategorizeError`

* Refactored handler error mapping to rely on categories instead of `errors.Is` chains

* Benefits:

  * consistent error translation
  * reduced coupling to domain error types
  * easier extensibility

* Applied same pattern to **customer module**:

  * centralized `CategorizeError`
  * extracted `ValidationField` helper

### 5. Delivery Layer Simplification

* Refactored account handler:

  * removed direct dependency on domain errors
  * now maps errors via application layer categorization
* Refactored customer handler:

  * removed local validation logic duplication
  * delegates to application helpers

### 6. Authentication Context Handling (Safety Improvement)

* Replaced panic-based helper:

  * `MustGetAuthenticatedUser` → `RequireAuthenticatedUser`

* Now returns explicit error (`ErrAuthenticatedUserNotFound`) instead of panicking

* This is a **significant robustness improvement**, preventing uncontrolled crashes in production paths

* Updated JWT middleware:

  * now uses delivery-level context helpers instead of application-level functions

* Reinforces correct layer responsibility

### 7. Test Coverage

* Added tests for:

  * access policy (`CanAccessCustomer`, `CanAccessAccount`)
  * new authentication helper behavior (error instead of panic)
* Removed obsolete authorization tests from auth module
* Adjusted existing tests to align with new abstractions

### 8. Minor Consistency Improvements

* Standardized access checks across all account use cases
* Reduced duplicated logic in validation and error handling
* Improved readability and intent clarity in application layer

### Conclusion

This commit delivers **architectural hygiene and consistency improvements** rather than new features.

Key highlights:

* **Authorization logic correctly placed in the account domain boundary**
* **Error handling standardized via categorization**
* **Safer authentication context management (no panics)**

Overall, this is a **high-value refactor**, reducing technical debt and reinforcing clean architecture principles without increasing system complexity.


## 2026/04/04 — auth/phase-06

Implements authentication and authorization as a first-class concern in the API, introducing JWT-based identity, route protection, and ownership validation across account operations. This phase consolidates security boundaries while maintaining clear separation of concerns across layers.

### 1. API Bootstrap and Wiring

* Integrated full auth stack into `main.go`:

  * user repository (PostgreSQL with pgx)
  * bcrypt password hasher
  * JWT token service with configurable secret and expiration
* Introduced auth use cases:

  * register user
  * login user
  * get current user
* Added JWT middleware and applied it to protected routes
* Protected all `/accounts/*` endpoints and `/auth/me`
* Kept `/auth/register` and `/auth/login` as public endpoints

### 2. Route Protection Layer

* Introduced `RequireAuth` middleware:

  * validates JWT token
  * injects authenticated context into request
* Enforced authentication boundary at HTTP layer instead of use case layer
* Ensures:

  * unauthenticated requests fail early
  * authenticated requests proceed with identity context

### 3. Authorization Enforcement

* Enforced account ownership rules through middleware and use case integration:

  * user must own the account (`customer_id`)
  * admin role bypasses ownership restriction
* Standardized forbidden responses:

  * `FORBIDDEN` with message "access denied to account"
* Covered critical flows:

  * own account access succeeds
  * cross-account access is denied
  * admin override is allowed

### 4. Infrastructure Migration to pgx

* Refactored `PostgresUserRepository`:

  * replaced `database/sql` with `pgxpool`
  * updated query execution (`Exec`, `QueryRow`)
  * replaced `sql.ErrNoRows` with `pgx.ErrNoRows`
* Aligns auth module with existing database infrastructure
* Improves consistency and performance characteristics

### 5. Authentication Contracts and Documentation

* Expanded README:

  * added auth module and endpoints
  * documented protected routes
  * introduced JWT configuration via `JWT_SECRET`
* Extended REST documentation:

  * added authentication section
  * documented `/auth/register`, `/auth/login`, `/auth/me`
  * updated all account endpoints with auth requirements
  * included new error codes: `UNAUTHORIZED`, `INVALID_TOKEN`, `FORBIDDEN`, `USER_ALREADY_EXISTS`, `INVALID_CREDENTIALS`

### 6. Error Standardization

* Refined shared error messages:

  * `UNAUTHORIZED` → "authentication required"
  * `FORBIDDEN` → "access denied to account"
* Updated response tests to reflect new semantics
* Aligns API responses with security expectations and clarity

### 7. Test Coverage

#### 7.1 Integration Tests

* Added end-to-end test suite for auth and authorization:

  * register → login → access `/auth/me`
  * authenticated access to own account
  * forbidden access to чужой account
  * admin override scenarios
* Validates:

  * JWT issuance and parsing
  * middleware enforcement
  * ownership rules

#### 7.2 Infrastructure Tests

* Added bcrypt tests:

  * cost fallback behavior
  * hash and compare validation
  * failure on wrong password
* Added JWT tests:

  * token generation and parsing
  * invalid signature handling
  * expired token rejection
  * malformed token detection
  * invalid signing method protection
  * missing claims validation

#### 7.3 Repository Tests

* Added integration tests for `PostgresUserRepository`:

  * create and query user
  * duplicate email constraint
  * not found behavior
  * existence checks

### 8. Developer Experience

* Added `make tests` target:

  * runs all tests with coverage
* Improves feedback loop and encourages test-driven workflow

### Conclusion

This phase establishes a solid security foundation by introducing authentication, enforcing authorization boundaries, and protecting critical financial routes. The implementation is cohesive, with proper layering, strong test coverage, and clear API contracts, enabling the system to evolve toward production-grade security standards.


## 2026/04/04 — auth/phase-05

Introduces **authorization enforcement and transactional abstraction improvements** across account operations, consolidating transaction management, refining repository design, and strengthening consistency guarantees between application and infrastructure layers.

### 1. Application Layer — Transaction Handling Refactor

* Replaced manual transaction management (`BeginTx`, `Commit`, `Rollback`) with `WithTransaction`
* Eliminated boilerplate patterns:

  * removed `committed` flags
  * removed deferred rollback logic
* Centralized transaction lifecycle:

  * execution, commit, and rollback now handled by repository
* Improved readability and reduced error-prone patterns

### 2. Authorization Enforcement

* Added authorization checks using:

  * `authdomain.CanAccessAccount`
* Applied consistently across use cases:

  * `Deposit`
  * `Withdraw`
  * `Transfer`
* Ensures:

  * only account owners (or authorized roles) can operate
  * unauthorized access returns `domain.ErrForbidden`
* This marks the first concrete enforcement of **Phase 5 — Authorization**

### 3. Transfer Use Case Adjustments

* Migrated full transfer flow to `WithTransaction`
* Preserved deterministic locking strategy using ordered UUIDs
* Maintained atomic behavior:

  * debit → credit → ledger entries
* Simplified error propagation by removing manual transaction boundaries

### 4. Withdraw and Deposit Use Cases

* Refactored both operations to:

  * use `WithTransaction`
  * enforce authorization before business logic
* Improved error handling:

  * explicit propagation of domain errors (`ErrAccountNotFound`, `ErrInsufficientBalance`)
* Ledger creation now fully encapsulated within transactional closure

### 5. Repository Contract Evolution

* Added new method:

  * `WithTransaction(ctx, fn)`
* Updated contract of `DecreaseBalance`:

  * now distinguishes:

    * `ErrAccountNotFound`
    * `ErrInsufficientBalance`
* Improves semantic clarity and enables better error mapping at upper layers

### 6. Infrastructure Layer — PostgreSQL Refactor

* Introduced `baseRepository` abstraction:

  * decouples execution (`exec`) from concrete type (`db` or `tx`)
  * enables reuse across repository and transactional repository
* Added `executor` interface:

  * unifies `pgxpool.Pool` and `pgx.Tx`
* Refactored all data access methods to use `baseRepository`
* Eliminated duplicated logic between:

  * `Repository`
  * `txRepository`

### 7. Transaction Orchestration

* Implemented `runInTransaction` helper:

  * executes callback
  * handles rollback on error
  * ensures rollback on commit failure
* Added `WithTransaction` implementation:

  * wraps `BeginTx` + `runInTransaction`
* Prevented nested transactions:

  * explicit errors returned in `txRepository`

### 8. Balance Consistency Improvements

* Enhanced `DecreaseBalance`:

  * now uses `RETURNING` clause
  * distinguishes between:

    * account not found
    * insufficient balance
  * fallback existence check implemented
* Ensures stronger correctness under concurrent conditions

### 9. Test Updates and Additions

* Updated mocks to support `WithTransaction`
* Added safeguards:

  * prevent nested transactions in tests
* Introduced new infrastructure tests:

  * `DecreaseBalance` scenarios:

    * account not found
    * insufficient balance
    * success
  * parity between repository and transactional repository
  * transaction lifecycle:

    * rollback on callback error
    * commit on success
    * rollback on commit failure
* Improved coverage of transactional behavior and edge cases

### 10. Architectural Impact

* Moves transaction management from application layer to infrastructure boundary
* Enforces clear separation:

  * application defines intent
  * repository controls execution and consistency
* Reduces duplication and aligns with unit-of-work pattern

### Conclusion

This commit represents a significant step toward a **robust authorization and transaction model**, combining:

* centralized transaction orchestration
* consistent authorization enforcement
* reduced duplication via repository abstraction
* improved correctness in concurrent financial operations

The system is now better aligned with production-grade patterns, particularly in terms of **security, consistency, and maintainability**.


### 2026/04/02 auth/phase-04

Introduces authentication context propagation and enforces authorization rules across account use cases, along with a consistent error handling strategy at the delivery layer.

1. internal/account/application/auth_test.go

   * Added helper functions to create test users (customer and admin) for authentication scenarios

2. internal/account/application/create_account.go

   * Added AuthenticatedUser to input
   * Enforced authorization: only admins or matching customers can create accounts
   * Introduced forbidden validation before repository calls

3. internal/account/application/create_account_test.go

   * Updated all tests to include authenticated user context
   * Added tests for forbidden access and admin privileges

4. internal/account/application/deposit.go

   * Added AuthenticatedUser to input
   * Enforced access control using CanAccessAccount

5. internal/account/application/deposit_test.go

   * Updated tests to include user context
   * Added forbidden scenario validation

6. internal/account/application/get_statement.go

   * Added AuthenticatedUser to input
   * Enforced access validation before retrieving transactions

7. internal/account/application/get_statement_test.go

   * Updated tests to include user context
   * Added forbidden access validation

8. internal/account/application/transfer.go

   * Added AuthenticatedUser to input
   * Enforced authorization on source account

9. internal/account/application/transfer_test.go

   * Updated tests to include user context
   * Added forbidden access scenario

10. internal/account/application/withdraw.go

    * Added AuthenticatedUser to input
    * Enforced access validation

11. internal/account/application/withdraw_test.go

    * Updated tests to include user context
    * Added admin access validation

12. internal/account/domain/errors.go

    * Introduced ErrForbidden for authorization failures

13. internal/account/delivery/account_handler.go

    * Added RequireUser to enforce authentication at handler level
    * Injected user context into all use cases
    * Replaced manual error handling with centralized mapAccountError
    * Standardized success responses using writeSuccess

14. internal/account/delivery/account_handler_test.go

    * Updated tests to include authenticated requests
    * Added tests for unauthorized and forbidden scenarios

15. internal/account/delivery/auth_test.go

    * Added helpers to inject authenticated user into request context

16. internal/account/delivery/deposit_integration_test.go

    * Injected authentication context into integration test flow

17. internal/account/delivery/handler.go

    * Added RequireUser helper to extract authenticated user from context

18. internal/account/delivery/response.go

    * Refactored response handling to use shared HTTP utilities
    * Replaced custom response structure with shared abstractions

19. internal/auth/application/authorization.go

    * Delegated authorization logic to auth domain layer
    * Simplified CanAccessAccount and RequireAccountAccess

20. internal/auth/application/authorization_test.go

    * Updated tests to use UUID instead of string for CustomerID

This commit establishes a consistent security boundary across application and delivery layers, ensuring that all account operations are protected by explicit authentication and authorization rules while improving error handling cohesion.


## 2026/04/02 — auth/phase-03

Implements a complete **authentication and authorization layer**, including JWT-based security, access control rules, standardized error handling, and HTTP response normalization. This phase establishes the foundation for secure interaction across the system.

### 1. Application Layer — Authorization Rules

* Introduced `CanAccessAccount` and `RequireAccountAccess`
* Implements access control based on:

  * admin override (full access)
  * customer ownership (account.CustomerID == user.CustomerID)
* Defensive validation:

  * nil checks
  * invalid UUID parsing
* Added `ErrForbidden` to represent authorization failures
* This design is clean and pragmatic, correctly centralizing authorization logic at the application boundary

### 2. Application Layer — Auth Context Refinement

* Replaced struct-based context key with typed key (`contextKey`)
* Eliminates collision risks and improves type safety
* Updated:

  * `WithAuthenticatedUser`
  * `GetAuthenticatedUser`
* Supports both value and pointer retrieval patterns
* This is a subtle but important improvement for robustness in middleware-driven systems

### 3. Delivery Layer — Auth Handlers

* Introduced HTTP handlers:

  * `POST /auth/register`
  * `POST /auth/login`
  * `GET /auth/me`
* Implemented request parsing and response mapping for:

  * user registration
  * authentication (token issuance)
  * current user retrieval
* Centralized error mapping via `MapError`
* Integrated structured logging for failure scenarios
* Handlers are well-structured and follow clear separation between transport and application concerns

### 4. JWT Middleware

* Added `JWTMiddleware` with two modes:

  * `RequireAuth`: enforces authentication
  * `OptionalAuth`: enriches context if token is present
* Implements:

  * Bearer token extraction
  * token parsing via `TokenService`
  * injection of `AuthenticatedUser` into request context
* Handles:

  * missing header
  * malformed token
  * invalid/expired token
* This middleware is technically solid and aligns with common production patterns

### 5. Delivery Layer — Context Helpers

* Added helper functions:

  * `GetAuthenticatedUser`
  * `MustGetAuthenticatedUser`
* Simplifies access to authenticated principal in handlers
* `MustGetAuthenticatedUser` enforces strict contract via panic (appropriate for internal invariants)

### 6. Shared Layer — Error Standardization

* Introduced `AppError` structure:

  * `code`
  * `message`
  * optional `details`
* Added predefined errors:

  * `INVALID_REQUEST`, `INVALID_DATA`
  * `UNAUTHORIZED`, `INVALID_TOKEN`
  * `USER_ALREADY_EXISTS`, etc.
* Provides consistent error contract across the entire API
* This is a strong architectural improvement, reducing duplication and ambiguity

### 7. Shared Layer — HTTP Response Abstraction

* Added standardized response format:

  * `{ data, error }`
* Introduced helpers:

  * `WriteSuccess`
  * `WriteError`
* Ensures:

  * consistent headers (`application/json`)
  * predictable payload structure
* Removes response formatting duplication across handlers

### 8. Test Coverage

* Comprehensive tests across all layers:

  * authorization rules (customer vs admin vs invalid cases)
  * context handling (presence and panic scenarios)
  * handlers (success and error mappings)
  * JWT middleware (valid, invalid, missing, expired tokens)
  * shared response and error structures
* Validates both behavior and contract consistency

### Conclusion

This commit represents a **major architectural milestone**, introducing a cohesive and production-ready authentication system.

Key strengths:

* Clear separation of concerns (application, delivery, shared)
* Robust JWT-based authentication flow
* Centralized and consistent error/response handling
* Well-defined authorization rules at the domain boundary

From a technical perspective, this is a **high-quality and scalable foundation** for securing all future endpoints in the system.


## 2026/04/02 — auth/phase-02

Implements the **core authentication module (phase 02)**, including user registration, login, token generation, and authenticated context handling. Establishes a clean separation between authentication concerns and domain logic, with strong validation and full test coverage.

### 1. Application Layer — User Registration

* Added `RegisterUserUseCase` with input (`Email`, `Password`)
* Implemented validation rules:

  * normalized email (trim + lowercase)
  * structural email validation
  * minimum password length (≥ 8)
* Enforced uniqueness via `ExistsByEmail`
* Integrated password hashing through `PasswordHasher`
* Created user entity with:

  * UUID identifier
  * default role = `customer`
  * timestamps (`CreatedAt`, `UpdatedAt`)
* Returns structured output without exposing sensitive data (no password)

### 2. Application Layer — Login

* Introduced `LoginUserUseCase`
* Responsibilities:

  * normalize email input
  * validate credentials
  * compare password hash
  * generate access token via `TokenService`
* Returns:

  * JWT (or equivalent token)
  * user identity and role
  * optional `CustomerID`
* Proper error handling:

  * `ErrInvalidCredentials` for authentication failures
  * avoids leaking whether email or password is incorrect

### 3. Application Layer — Current User Resolution

* Added `GetCurrentUserUseCase`
* Retrieves authenticated user from context (`context.Context`)
* Supports two modes:

  * **context-only** (no repository) → lightweight resolution
  * **repository-backed** → ensures user still exists
* Introduced:

  * `AuthenticatedUser` (principal abstraction)
  * `WithAuthenticatedUser` / `GetAuthenticatedUser` helpers
* Returns:

  * user ID, email, role, and optional customer binding
* Handles unauthorized scenarios via `ErrUnauthorized`

### 4. Context-Based Authentication Model

* Establishes a **request-scoped identity propagation mechanism**
* Decouples authentication middleware from business logic
* Enables:

  * future integration with JWT middleware
  * role-based authorization at use case level
* This is a solid architectural decision, aligning with Go idioms and clean layering

### 5. Domain Integration

* Leverages existing domain contracts:

  * `UserRepository`
  * `PasswordHasher`
  * `TokenService`
* Introduces no leakage of infrastructure concerns into application layer
* Maintains clear dependency inversion

### 6. Validation & Normalization Strategy

* Centralized helpers:

  * `normalizeEmail`
  * `isValidEmail`
  * `isValidPassword`
* Ensures:

  * consistent input handling
  * predictable authentication behavior
* Notably avoids over-engineering while covering essential edge cases

### 7. Test Coverage

#### 7.1 RegisterUser

* success flow (including normalization and persistence)
* duplicate email
* invalid email
* invalid password
* hashing failure
* verifies:

  * repository interaction
  * correct entity construction
  * timestamp integrity

#### 7.2 LoginUser

* success scenario (full flow: lookup → compare → token)
* user not found
* wrong password
* token generation failure
* validates:

  * normalization
  * secure flow (no unnecessary calls on failure)
  * correct token claims

#### 7.3 GetCurrentUser

* success with repository
* missing context (unauthorized)
* user not found
* repository error propagation
* ensures strict control over authentication state

### 8. Error Handling Strategy

* Clear domain-aligned errors:

  * `ErrInvalidEmail`
  * `ErrInvalidPassword`
  * `ErrEmailAlreadyExists`
  * `ErrInvalidCredentials`
  * `ErrUnauthorized`
* Consistent wrapping for infrastructure errors
* Prevents information leakage in authentication flows

### Conclusion

This commit establishes a **robust and extensible authentication foundation**, covering:

* user lifecycle (register + login)
* secure credential handling
* token-based authentication
* request-scoped identity propagation

From an architectural perspective, the design is **clean, idiomatic, and production-ready**, particularly due to:

* strict separation of concerns
* context-driven authentication model
* defensive error handling
* comprehensive test coverage

It provides a solid base for future extensions such as middleware, authorization policies, and multi-tenant support.


## 2026/04/02 — auth/phase-01

Introduces the **account statement (ledger query) capability**, expands repository contracts to support transaction history retrieval, and improves project operability with tooling and documentation. This marks a transition from pure command operations to **read-side financial visibility**.

### 1. Application Layer — Get Statement Use Case

* Added `GetStatement` use case with support for:

  * pagination (`limit`, `cursor`, `cursor_id`)
  * date filtering (`from`, `to`)
* Implemented validations:

  * non-nil account ID
  * cursor consistency (cursor + cursor_id must coexist)
  * valid date range (`from <= to`)
  * limit normalization (default: 50, max: 100)
* Flow:

  * validate input
  * ensure account exists
  * query transactions via repository
  * map to output DTO
  * build next cursor when applicable
* Returns structured response with:

  * transaction list
  * pagination cursor
* Clean separation between **query logic and persistence concerns**

### 2. Domain & Repository Evolution

* Extended `AccountRepository` with:

  * `GetTransactions(...)`
* Enables:

  * cursor-based pagination
  * time-range filtering
* Maintains domain purity by exposing **query intent without leaking SQL concerns**


### 3. Infrastructure Layer — PostgreSQL Statement Query

* Implemented `GetTransactions` for:

  * base repository
  * transactional repository (`txRepository`)
* SQL characteristics:

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
  * efficient index usage (assuming proper indexing)

### 4. Delivery Layer — Statement Endpoint

* Added endpoint:

  * `GET /accounts/{id}/statement`
* Implemented:

  * path param parsing (account ID)
  * query param parsing:

    * `limit`
    * `cursor`
    * `cursor_id`
    * `from`
    * `to`
* Introduced helper parsers:

  * `parseOptionalInt`
  * `parseOptionalTime`
  * `parseOptionalUUID`
* Error handling:

  * `INVALID_DATA → 400`
  * `ACCOUNT_NOT_FOUND → 404`
* Response includes:

  * list of transactions
  * optional `next_cursor` for pagination

### 5. DTO Additions

* Added:

  * `StatementData`
  * `StatementItemData`
  * `StatementCursorData`
* Clearly separates:

  * internal domain representation
  * external API contract

### 6. Handler & Wiring Updates

* Extended handler to include `statement` use case via interface
* Updated constructor signature accordingly
* Registered new route in `main.go`
* Refactored handler file naming (`handleer.go → handler.go`)
* Maintains consistency with dependency injection and layered design

### 7. Test Coverage

#### 7.1 Application Tests

* Added full suite for `GetStatement`:

  * invalid input scenarios
  * default and capped limits
  * account not found
  * transaction retrieval failure
  * successful pagination flow
* Verifies:

  * correct repository invocation
  * cursor generation logic

#### 7.2 Delivery Tests

* Added handler tests for:

  * invalid query params
  * cursor consistency validation
  * account not found
  * successful response mapping

#### 7.3 Test Infrastructure

* Updated all mocks to support:

  * `GetTransactions`
* Ensures compatibility across all existing use cases

### 8. Tooling — Makefile

* Introduced `Makefile` to standardize developer workflow:

  * `migration` → run DB migrations
  * `commit` → commit using predefined message file
  * `diff` → generate staged diff and line count
  * `push` / `pull` → simplified git operations
* Improves productivity and consistency in development operations

### 9. Documentation

* Added `docs/06-implementation.md`:

  * comprehensive description of:

    * architecture
    * domain model
    * use cases
    * persistence strategy
    * concurrency model
    * test coverage
* Serves as **authoritative reference of the current implementation state**

### Conclusion

This commit introduces the **read-side of the financial ledger (account statement)**, completing a critical capability for any banking system.

Key highlights:

* **Cursor-based pagination with deterministic ordering**
* **Time-range filtering**
* **Clean separation between command (write) and query (read) concerns**
* **Strong alignment with transactional consistency already established**

From an architectural standpoint, this is a **natural and necessary evolution**, transforming the system from a purely operational core into a **queryable financial platform** with proper observability of account activity.


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
