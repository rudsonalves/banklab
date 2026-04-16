# Changelog

## 2026/04/16 — api/user_status-04

Introduces **user status enforcement across account creation and admin approval flow**, strengthening authorization guarantees and aligning onboarding with an explicit lifecycle (pending → active).

### 1. Application Layer — CreateAccount Hardening

* Extended `CreateAccount` use case to depend on `UserRepository`
* Added validation pipeline before account creation:

  * user must exist
  * user must have valid `UserID`
  * user must be in `active` status
* Enforced strict access control:

  * any non-active user (pending, blocked) is rejected with `ErrForbidden`
* Prevents invalid states where accounts could be created for non-approved users
* This change aligns account creation with the authentication model where identity alone is insufficient without valid state 

### 2. New Admin Capability — Approve User

* Integrated `ApproveUserUseCase` into application wiring
* Added new protected endpoint:

  * `POST /admin/users/{id}/approve`
* Approval flow:

  * transitions user status to `active`
  * creates associated account atomically
* Extended output contract:

  * now includes `status` field alongside `user_id` and `account_id`
* Establishes explicit onboarding lifecycle:

  * register → pending → approved → active → operational
* This is a **structural improvement**, not just a feature addition

### 3. Delivery Layer — Authorization Enforcement

* Implemented `ApproveUser` handler with strict guards:

  * requires authenticated user
  * enforces `admin` role
  * validates UUID path parameter
* Maps domain errors consistently:

  * `FORBIDDEN`, `INVALID_DATA`, `USER_NOT_FOUND`, etc.
* Response includes:

  * `user_id`
  * `status`
  * `account_id`
* Maintains API contract consistency with existing envelope pattern 

### 4. Error Handling Standardization

* Added missing domain error mappings in auth layer:

  * `ErrForbidden`
  * `ErrInvalidData`
* Unified error message:

  * `"Access denied"` replaces inconsistent variants
* Removed duplicate error registration guard in shared mapper:

  * simplifies registry behavior
  * shifts responsibility to developer discipline
* Keeps alignment with global error strategy and response model 

### 5. Dependency Wiring (Composition Root)

* Updated `main.go`:

  * `CreateAccount` now receives `userRepo`
  * `ApproveUserUseCase` added to auth handler
* Ensures correct dependency flow:

  * Delivery → Application → Domain
* Reinforces modular monolith structure and explicit wiring rules 

### 6. Test Coverage Expansion

#### 6.1 CreateAccount Tests

* Added scenarios for:

  * user repository not configured
  * user not found
  * user lookup failure
  * pending user rejection
  * blocked user rejection
* Verified:

  * no repository side-effects on failure paths
  * correct interaction counts

#### 6.2 ApproveUser Tests

* Added full coverage:

  * success case (admin)
  * unauthorized request
  * non-admin rejection
  * invalid UUID
  * domain error mapping (not found, conflict, forbidden, internal)
* Validates both:

  * authorization boundary
  * HTTP contract correctness

#### 6.3 Integration Adjustments

* Updated handler constructors to include new dependency
* Ensured compatibility with existing integration tests

### 7. Domain Alignment

* Reinforces the concept that:

  * **user status is part of authorization context**, not just authentication
* Prevents illegal transitions such as:

  * financial operations executed by pending users
* Aligns with domain invariants where operations depend on valid state, not only identity 

### Conclusion

This commit introduces a **critical correction in the authorization model** by incorporating user status into the decision process.

The system now eliminates an important invalid state:
users could previously authenticate and operate without being formally approved.

From an architectural perspective, this is a **consistency and correctness fix**, ensuring that onboarding, authorization, and financial operations are coherently integrated.


## 2026/04/16 — api/user_status-03

Implements the **user approval flow with automatic account creation**, introducing transactional consistency across auth and account modules, strengthening user lifecycle control, and expanding domain and infrastructure support for status transitions.

### 1. Application Layer — ApproveUser Use Case

* Added `ApproveUserUseCase` to handle transition from `pending` to `active`
* Full transactional orchestration via `Transactor.RunInTx`
* Execution flow:

  * load user with `FindByIDForUpdate` (row-level lock)
  * validate user existence and current status
  * update user status to `active`
  * validate `customer_id` presence and existence
  * generate account number
  * create and persist new account
* Ensures **atomic activation + account creation**, preventing partial states
* Reuses `accountapplication.GenerateBranch()` for consistency with account module

### 2. Domain Layer Enhancements

* Introduced new domain error:

  * `ErrUserAlreadyActive`
* Reinforces user lifecycle invariants:

  * only `pending` users can be approved
  * active users cannot be reprocessed
* Aligns with invariant enforcement strategy described in the domain model 

### 3. Repository Contract Evolution

* Extended `UserRepository`:

  * added `FindByIDForUpdate` for pessimistic locking
* Enables safe concurrent approval handling, consistent with system-wide locking strategy 

### 4. Infrastructure Layer — PostgreSQL

* Implemented `FindByIDForUpdate` using:

  * `SELECT ... FOR UPDATE`
* Ensures:

  * row-level locking during approval
  * protection against concurrent status transitions
* Behavior:

  * returns `nil` when user not found (mapped at application level)

### 5. Error Handling Standardization

* Registered new domain errors in error registry:

  * `USER_NOT_FOUND → 404`
  * `USER_ALREADY_ACTIVE → 409`
* Added corresponding error codes in shared layer:

  * `ErrCodeUserNotFound`
  * `ErrCodeUserAlreadyActive`
* Maintains consistency with global API error contract 

### 6. Account Module Adjustment

* Refactored branch generation:

  * `generateBranch` → `GenerateBranch`
* Promotes reuse across modules and avoids duplication in account creation logic

### 7. Test Coverage

#### 7.1 Application Tests — ApproveUser

* Added comprehensive test suite:

  * success scenario (user activation + account creation)
  * user not found
  * user already active
  * missing customer_id
  * customer not found
  * account creation failure
* Validates transactional integrity and invariant enforcement

#### 7.2 Test Infrastructure Updates

* Updated mocks across auth tests to support:

  * `FindByIDForUpdate`
* Ensures compatibility with new repository contract without breaking existing tests

### 8. Architectural Impact

* Introduces a **controlled onboarding progression**:

  * register → pending user
  * approve → active user + account creation
* Eliminates invalid intermediate states:

  * active user without account
  * account created for unapproved user
* Strengthens consistency guarantees across modules, aligned with system transaction model 

### Conclusion

This commit introduces a **critical lifecycle transition for users**, integrating authentication and financial domains under a single transactional boundary.

The implementation is technically robust, particularly due to:

* explicit use of pessimistic locking
* strict invariant enforcement
* elimination of partial states

From an architectural perspective, this significantly improves the correctness and consistency of the onboarding flow, aligning it with the system’s financial integrity requirements.


## 2026/04/16 — api/user_status-02

Refactors user registration transaction handling, introduces **user status as a first-class domain attribute**, and standardizes persistence behavior across repository and application layers.

### 1. Application Layer — Transaction Handling Refactor

* Replaced repository-specific transaction coupling (`WithTransaction`) with a **generic `Transactor` abstraction**
* Updated `RegisterUserUseCase` to depend on `domain.Transactor` instead of casting the repository
* Execution now uses:

  * `transactor.RunInTx(ctx, fn)`
* Removes implicit assumptions about repository capabilities and enforces **explicit transaction orchestration at the application layer**
* Improves architectural consistency with the layered design already used in account operations 

### 2. Domain Layer — User Status Introduction

* Added `status` attribute to `User` entity lifecycle
* Introduced new domain error:

  * `ErrUserNotFound`
* Extended `UserRepository` contract:

  * added `UpdateStatus(userID, status)`
* Clarified repository responsibilities with explicit method semantics (e.g. `ExistsByEmail`)
* Aligns user lifecycle with a **state-driven model**, enabling future approval/activation flows

### 3. Persistence Layer — PostgreSQL Updates

* Added `status` column to `users` table:

  * `VARCHAR(20) NOT NULL DEFAULT 'pending'`
* Updated repository behavior:

  * `Create` now persists `status`
  * `FindByEmail` and `FindByID` now map `status`
  * implemented `UpdateStatus` with:

    * `UPDATE users SET status = $1, updated_at = NOW()`
    * returns `ErrUserNotFound` when no rows affected
* Removed legacy `WithTransaction` implementation from repository
* Consolidates responsibility: **repository handles persistence, application handles transactions**

### 4. Migration Layer

* Added migration:

  * `000006_user_status.up.sql`
  * `000006_user_status.down.sql`
* Ensures schema evolution is:

  * incremental
  * reversible
  * aligned with existing migration strategy

### 5. Test Infrastructure Refactor

* Updated all mocks to support new contract:

  * added `UpdateStatus` to `UserRepository` mocks
* Replaced transaction mocking:

  * removed `WithTransaction`
  * introduced `registerTransactorMock` with `RunInTx`
* Adjusted assertions:

  * validate `RunInTx` invocation instead of repository transaction calls

### 6. Integration Tests

* Updated integration setup:

  * ensures `users.status` column exists via `ALTER TABLE IF NOT EXISTS`
* Added validation for:

  * default status (`pending`) on creation
  * status transition via `UpdateStatus`
* Strengthens alignment between:

  * schema
  * repository
  * domain behavior

### 7. Wiring Adjustments

* Updated `main.go`:

  * `RegisterUserUseCase` now receives `transactor`
* Keeps dependency graph explicit and consistent with other use cases

### 8. Design Considerations

This change is particularly relevant from an architectural standpoint:

* eliminates hidden coupling between repository and transaction control
* enforces **application-layer ownership of transactional boundaries**
* introduces a **stateful user lifecycle**, which is essential for:

  * approval flows
  * onboarding pipelines
  * access control evolution

The decision to move away from `WithTransaction` is correct and aligns the auth module with the same rigor already present in financial operations.

### Conclusion

This commit is a structural improvement rather than a feature addition.

It establishes:

* a **clean transaction boundary model**
* a **state-driven user lifecycle**
* a **more consistent repository contract**

These changes prepare the system for more advanced flows such as user approval, activation, and policy-based authorization without requiring further architectural refactoring.


## 2026/04/16 — api/user_status-01

Introduces **user status management** into the authentication domain and restructures the HTTP layer to support a clearer multi-stage authentication model, including AppToken-based onboarding and JWT-protected routes. 

### 1. Domain Layer — User Status

* Added `UserStatus` type with explicit states:

  * `pending`
  * `active`
  * `blocked`
* Extended `User` entity to include `Status` field
* Updated `NewUser` factory:

  * initializes all new users with `UserStatusPending`
* This change establishes a **foundation for lifecycle control** (approval, activation, blocking), which was previously absent in the model

### 2. Bootstrap — Environment Configuration

* Introduced automatic `.env` loading using `github.com/joho/godotenv`
* Implemented flexible resolution strategy:

  * local `.env`
  * `api/.env`
  * executable-relative paths
* Ensures configuration is available regardless of execution context
* Aligns with **fail-fast configuration validation** already present in `main.go`

### 3. Main Wiring Refactor (cmd/api/main.go)

* Reorganized startup into explicit sections:

  * Config
  * Repositories
  * Services
  * Use Cases
  * Handlers
  * Middlewares
  * Routers
* Enforced validation of critical environment variables:

  * `APP_TOKEN`
  * `JWT_SECRET`
* Improved readability and maintainability of composition root
* This is a **structural improvement**, not just cosmetic; it clarifies dependency boundaries

### 4. Routing Architecture — Separation of Concerns

* Split routing into three layers:

  * `authRouter` (authentication endpoints)
  * `apiRouter` (business endpoints)
  * `mainRouter` (composition)
* Introduced explicit middleware application per route group:

  * AppToken for onboarding
  * JWT for authenticated access
* Eliminated global middleware wrapping, replacing it with **route-level control**, which is more precise and safer

### 5. Authentication Model — AppToken + JWT

* Applied AppToken middleware to:

  * `POST /auth/register`
  * `POST /auth/login`
* Applied JWT middleware to:

  * `POST /auth/refresh`
  * `GET /auth/me`
  * all `/accounts/*`
  * `/customers/me`
* This formalizes a **two-phase authentication model**:

  * controlled entry (AppToken)
  * authenticated session (JWT)
* Matches the intended design described in the authentication documentation 

### 6. API Contract Documentation Updates

* Updated REST documentation to reflect:

  * AppToken requirement for onboarding endpoints
  * JWT requirement for all protected endpoints
* Added new error code:

  * `INVALID_APP_TOKEN` (HTTP 401)
* Expanded error scenarios with concrete payload examples
* Clarified access control rules and authentication flows
* Improves alignment between implementation and public contract 

### 7. Dependency Updates

* Added `godotenv` dependency to `go.mod` and `go.sum`
* Enables environment-based configuration without external tooling

### 8. Architectural Impact

* Introduces the first step toward **user lifecycle governance** via status
* Establishes a clearer boundary between:

  * onboarding security
  * session-based authentication
  * resource authorization
* Prepares the system for future features such as:

  * user approval workflows
  * account activation
  * access blocking

### Conclusion

This commit is a **strategic evolution of the authentication layer**, not merely a feature addition.

It introduces user lifecycle semantics and formalizes a multi-stage authentication model, improving both **security posture and architectural clarity**, while keeping the system aligned with its current simplicity goals and ready for future extensions.


## 2026/04/16 — api/app_token-01

Introduces **application-level request validation via App Token middleware**, enforces stricter environment configuration, and refactors HTTP server initialization to support middleware composition and improved security boundaries.

### 1. Security — App Token Middleware

* Implemented `AppToken` middleware to enforce presence and validity of `X-App-Token` header
* Uses `crypto/subtle.ConstantTimeCompare` to prevent timing attacks during token comparison
* Rejects unauthorized requests early in the pipeline with standardized error response
* Integrates with existing error contract via `ErrInvalidAppToken`
* Establishes a clear separation between:

  * client identity (JWT)
  * client application validation (App Token)

This aligns with a layered security model, where multiple independent signals are validated before request execution, consistent with the system’s architectural direction 

### 2. Error Standardization

* Added `ErrInvalidAppToken` to shared sentinel errors
* Ensures consistency with API response envelope (`data` / `error`) 
* Avoids inline error construction, reinforcing centralized error definitions

### 3. HTTP Pipeline Refactor

* Replaced global `http.Handle` usage with explicit `http.ServeMux`
* Introduced handler composition:

  * `router → auth middleware → app token middleware`
* Final pipeline:

  * `handler := AppToken(...) (router)`
* Improves:

  * composability
  * testability
  * visibility of request flow

### 4. Environment Hardening

* Enforced mandatory environment variables:

  * `APP_TOKEN` now required (fail-fast with `log.Fatal`)
  * `JWT_SECRET` no longer has fallback value
* Eliminates insecure default configurations
* Guarantees that authentication and application validation cannot run in an invalid state

### 5. Routing Organization

* Centralized all routes into `ServeMux`
* Maintains explicit registration of:

  * auth endpoints
  * customer endpoints
  * account endpoints (including transfer)
* Preserves existing authorization behavior via JWT middleware

### 6. Test Coverage — Middleware

* Added comprehensive tests for `AppToken`:

  * missing header
  * invalid token
  * valid token (happy path)
* Validates:

  * correct HTTP status (`401`)
  * response envelope structure
  * prevention of downstream handler execution on failure

### 7. Developer Experience

* Added `api-run` command to `Makefile` for local server execution
* Introduced `.gitignore` for API build artifacts

### 8. Documentation Reorganization

* Moved API documentation into `api/docs`
* Improves cohesion between code and documentation
* Maintains consistency with modular project structure

### Conclusion

This commit introduces a **critical security boundary at the application level**, ensuring that requests are validated not only by user identity (JWT) but also by client context (App Token).

From an architectural standpoint, this is a meaningful step toward a **multi-signal validation model**, where authentication alone is no longer treated as sufficient.


## 2026/04/16 — api/refresh_token-02

Refines authentication flow and project documentation, introducing **refresh token persistence on the client side** and restructuring the repository documentation to reflect the current architecture and usage model.

### 1. Mobile — Refresh Token Support

* Extended `LoggedUser` model:

  * added `refreshToken` field
  * updated `fromMap` to deserialize `refresh_token`
* Updated authentication repository:

  * persist `refreshToken` alongside `accessToken` on login
  * ensure removal of both tokens on logout
* This aligns the mobile client with the backend session model, enabling **token rotation and session continuity**

### 2. Authentication Flow Consistency

* Ensures that client-side storage reflects server expectations:

  * access + refresh token pair becomes the canonical session representation
* Prepares the mobile layer for:

  * automatic token refresh via interceptor
  * retry logic on `401` responses
* Reinforces contract implied by auth endpoints and JWT usage 

### 3. Monorepo Documentation Simplification

* Rewrote root `README.md`:

  * removed narrative-heavy content
  * introduced concise structure and quick-start flow
  * clarified dual-app nature (API + mobile)
* Focus shifted from conceptual description to **operational clarity**

### 4. API Documentation Restructuring

* Simplified `api/README.md`:

  * clearer separation of stack, architecture, and features
  * explicit route listing
  * streamlined setup instructions
* Removed legacy architecture document:

  * replaced with new `ARCHITECTURE.md`
* New architecture document:

  * formalizes modular monolith structure
  * clarifies layer responsibilities and dependency direction
  * documents authentication and refresh flow behavior
* Maintains alignment with layered architecture principles 

### 5. Documentation Standardization (Mobile)

* Rewrote `mobile/README.md`:

  * emphasizes role as integration client
  * adds environment configuration guidance
  * documents test and build workflows
* Introduced `docs/mobile/ARCHITECTURE.md`:

  * defines layered structure (UI, Data, Domain, Core)
  * formalizes request flow and interceptor behavior

### 6. Licensing

* Added MIT license to API module:

  * clarifies usage and distribution terms
  * aligns repository with open-source conventions

### 7. Structural Improvements

* Normalized directory descriptions across README files
* Improved onboarding flow:

  * Docker → migrations → API → mobile
* Reduced redundancy across documentation layers

### Conclusion

This commit is primarily a **consistency and alignment step** between client, API, and documentation.

It establishes:

* a **complete token lifecycle on the client (access + refresh)**
* a **clearer and more operational documentation structure**
* a **more explicit architectural baseline for future evolution**

From a design standpoint, this is a necessary consolidation step before advancing into more complex authentication concerns such as concurrent refresh handling and session control.


## 2026/04/15 — api/refresh_token-02

Refactors and hardens the refresh token flow to guarantee **atomic rotation, consistency, and correctness of session management**, while simplifying dependency contracts and eliminating invalid execution paths.

---

### 1. Application Layer — Refresh Token Atomicity

* Introduced `Transactor` as a first-class dependency in `RefreshAccessTokenUseCase`
* Implemented atomic rotation using `RunInTx`:

  * `Revoke(old_token)` + `Create(new_token)` executed in a single transaction
* Removed non-transactional revoke operation
* Guarantees:

  * no partial state (no “revoked without replacement” or “duplicated sessions”)
  * rollback preserves original token validity on failure
* Aligns transaction control with Application layer responsibilities 

---

### 2. Infrastructure Layer — Transaction Support

* Added `PostgresTransactor`:

  * wraps `pgx` transaction lifecycle (`Begin → Commit / Rollback`)
  * injects transaction into context (`ContextWithTx`)
* Enables multiple repositories to share the same transaction transparently
* Strengthens infrastructure compliance with domain contracts (`Transactor` interface)

---

### 3. Domain Layer — Contract Expansion

* Introduced `Transactor` interface:

  * explicit control of transactional boundaries at use case level
* Added `ErrSessionNotFound`:

  * ensures revoke failures are explicit and not silently ignored
* Reinforces domain-driven consistency for session lifecycle

---

### 4. Session Repository — Correctness Enforcement

* Updated `Revoke` implementation:

  * now validates `RowsAffected()`
  * returns `ErrSessionNotFound` when token does not exist
* Eliminates silent inconsistencies in session state

---

### 5. Login Use Case — Contract Tightening

* Removed optional `SessionRepository` dependency (variadic → required)
* Enforced invariant:

  * **no refresh token is issued without persisted session**
* Simplified logic:

  * always hashes and persists refresh token
  * failure in session creation aborts login
* This removes previously possible invalid states

---

### 6. Delivery Layer — Dependency Integrity

* Refactored `Handler` constructor:

  * now requires all use cases upfront (including refresh)
  * removed setter injection (`SetRefreshAccessTokenUseCase`)
* Ensures:

  * no partially initialized handlers
  * no runtime mutation of dependencies
* Added defensive check for `nil` output in login flow

---

### 7. Wiring (main.go & integration)

* Registered `PostgresTransactor` in composition root
* Injected into `RefreshAccessTokenUseCase`
* Updated handler initialization to reflect new constructor contract
* Ensures consistent dependency graph across application and tests

---

### 8. Test Suite — Coverage Expansion

#### 8.1 Refresh Flow Tests

* Updated all tests to include `Transactor` dependency
* Added `transactorMock` for transactional execution

#### 8.2 Rotation Integrity Test (Stateful)

* Introduced `statefulSessionMock` to simulate real session lifecycle
* Validates full rotation behavior:

  * old token becomes unusable after refresh
  * new token is immediately valid
  * reuse of revoked token fails correctly
* This is a critical validation of **system invariants**

#### 8.3 Login Tests

* Updated to reflect mandatory `SessionRepository`
* Ensures session persistence is always exercised

#### 8.4 Infrastructure Test Fix

* Simplified JWT error assertion (removed invalid `errors.Is` usage)

---

### 9. API & Documentation Updates

* Login now explicitly returns `refresh_token`
* Introduced `/auth/refresh` endpoint with:

  * token rotation semantics
  * single-use refresh tokens
  * atomic revoke + create behavior
* Documented error scenarios:

  * invalid, expired, revoked, or missing sessions
* Clarifies contract for clients and aligns behavior with implementation 

---

### Conclusion

This commit is a **consistency and correctness milestone** for authentication flow.

It eliminates entire classes of invalid states by enforcing:

* **atomic token rotation**
* **mandatory session persistence**
* **explicit failure handling**
* **constructor-level dependency integrity**

From an architectural standpoint, this is a decisive improvement:
transaction boundaries are now correctly owned by the Application layer, and the authentication model becomes **predictable, verifiable, and resilient under failure conditions**.


## 2026/04/11 — api/refresh_token-01

Implements a **complete refresh token flow with session management**, evolving the authentication model from stateless JWT-only to a **stateful, revocable, and rotating session-based approach**. This change aligns the system with a more robust security posture while preserving the layered architecture principles 

---

### 1. Domain Layer — Contracts Expansion

* Extended `TokenService`:

  * added `GenerateRefreshToken(userID)`
  * added `ParseRefreshToken(token)`
* Introduced `SessionRepository`:

  * `Create`
  * `FindByTokenHash`
  * `Revoke`
* Establishes **explicit session lifecycle control** at the domain boundary

---

### 2. Application Layer — Login Flow Evolution

* `LoginUserUseCase` updated to:

  * generate **access token + refresh token**
  * hash refresh token using `SHA-256`
  * persist session with expiration (`30 days TTL`)
* Output now includes:

  * `AccessToken`
  * `RefreshToken`
* Optional session dependency supported (backward-safe injection)

**Key observation (architectural):**
This is the first point where authentication becomes **state-aware**, breaking the purely stateless JWT model intentionally.

---

### 3. Application Layer — Refresh Token Use Case

* Introduced `RefreshAccessTokenUseCase`:

  * validates refresh token integrity
  * validates session (existence, expiration, revocation, ownership)
  * loads user from repository
  * generates new access token
  * performs **refresh token rotation**:

    * revoke old token
    * generate new refresh token
    * persist new session

**Security properties introduced:**

* replay protection (rotation)
* server-side invalidation
* binding between token and stored session

---

### 4. Infrastructure Layer — Token Service

* Extended `JWTTokenService`:

  * added **opaque refresh token generation**

    * payload: `userID + nonce`
    * signature: `HMAC-SHA256`
    * encoding: `base64url`
  * implemented `ParseRefreshToken`

    * signature validation using constant-time comparison
    * strict payload validation

* Access token improvements:

  * ensured `exp` claim correctness with TTL enforcement

**Technical decision:**
Refresh token is **not JWT**, which is a correct choice to:

* reduce attack surface
* simplify validation
* avoid overloading JWT semantics

---

### 5. Infrastructure Layer — Session Persistence

* Added `PostgresSessionRepository`

  * `Create`: inserts hashed token
  * `FindByTokenHash`: retrieves session state
  * `Revoke`: soft-revokes via `revoked_at`

* Supports transaction-aware execution via context

* Migration `000005_user_sessions`:

  * new table `user_sessions`
  * indexed by `user_id` and `expires_at`
  * unique constraint on `token_hash`

**Critical design choice:**

* only **hashed tokens are stored**
* prevents token leakage from DB compromise

---

### 6. Delivery Layer — HTTP Contract Updates

* `/auth/login`:

  * now returns `refresh_token`

* New endpoint:

  * `POST /auth/refresh`

* Handler additions:

  * request validation (`refresh_token`)
  * consistent error mapping
  * response envelope preserved

* Introduced DTOs:

  * `refreshAccessTokenRequest`
  * `refreshAccessTokenData`

**Important:**
This extends the API contract beyond what is currently documented  and requires documentation update.

---

### 7. Dependency Wiring (main.go)

* Registered:

  * `SessionRepository`
  * `RefreshAccessTokenUseCase`
* Injected into:

  * `LoginUserUseCase`
  * handler via setter
* Exposed route:

  * `POST /auth/refresh`

---

### 8. Test Coverage

#### 8.1 Application Tests

* Login:

  * validates access + refresh generation
  * verifies session persistence (hash + expiration)
  * covers failure scenarios:

    * access token failure
    * refresh token failure
    * session persistence failure

* Refresh flow:

  * success path
  * invalid token
  * session not found
  * revoked session
  * expired session
  * user mismatch
  * repository failures
  * rotation integrity

#### 8.2 Infrastructure Tests

* Refresh token:

  * generation + parsing
  * entropy validation
  * tampering detection
  * malformed token handling
* Access token:

  * expiration correctness

#### 8.3 Delivery Tests

* Login:

  * validates response now includes `refresh_token`
* Refresh:

  * success case
  * invalid token → `401 INVALID_TOKEN`

#### 8.4 Integration Tests

* End-to-end validation:

  * login returns both tokens
  * refresh endpoint wired correctly
* Test DB isolation:

  * switched to `bank_test`
* CPF constraint repair added for test consistency

---

### 9. Test & Environment Adjustments

* Updated default test database:

  * `bank → bank_test`
* Added defensive SQL for constraint repair:

  * prevents flaky test runs due to regex mismatch

---

### 10. Behavioral Changes Summary

* Access tokens are now **short-lived**
* Refresh tokens are:

  * generated securely
  * persisted as hashed values
  * validated against DB
  * rotated on use
* Authentication becomes:

  * **stateful**
  * **revocable**
  * **traceable**

---

### Conclusion

This commit represents a **major security and architectural milestone**:

* transitions authentication from **stateless JWT** to **session-backed model**
* introduces **refresh token rotation**, a critical protection against replay attacks
* enforces **server-side control over sessions**, enabling future features such as:

  * logout
  * device/session listing
  * anomaly detection

From a technical standpoint, the implementation is **well-aligned with Clean Architecture principles**, keeping:

* domain contracts pure
* application responsible for orchestration
* infrastructure isolated

The only architectural caveat is the **optional session dependency in login**, which introduces a potential inconsistency. In a production-grade system, this should be mandatory to avoid issuing unusable refresh tokens.

Overall, this is a **production-grade foundation for authentication**, suitable for fintech-level requirements.


## 2026/04/10 — infra/layout-01

Introduces a **UI layout standardization layer** for the Flutter application, centralizing structural concerns and improving consistency across authentication screens, while also refining routing behavior and state handling patterns.

### 1. Routing Adjustment

* Updated initial route:

  * from `HomeRoutes.home` to `AuthRoutes.login`
* Aligns application startup with authentication flow, enforcing a more realistic entry point for protected systems
* This change is consistent with the backend contract where authentication precedes access to account resources 

---

### 2. Introduction of SafeScaffold

* Added new base component: `SafeScaffold`
* Encapsulates:

  * `SafeArea` handling for body and bottom navigation
  * consistent horizontal constraints (`maxWidth: 460`)
  * standardized padding for bottom actions
* Provides a **reusable layout abstraction**, reducing duplication and enforcing UI consistency
* Conceptually aligns with separation of responsibilities seen in the backend architecture, isolating structural concerns from business/UI logic 

---

### 3. Login Page Refactor

* Migrated from `Scaffold` to `SafeScaffold`
* Introduced `AppBar` for clearer navigation structure
* Refactored state handling:

  * replaced `setState` with `ValueNotifier<bool>` for password visibility
* Improved layout:

  * consistent spacing using `Column.spacing`
  * moved primary action to `bottomNavigationBar`
  * added `GestureDetector` to dismiss keyboard
* Decoupled navigation logic into dedicated methods (`_navToRegister`)
* Replaced direct widget access with local `_viewModel` reference for better readability and lifecycle control

---

### 4. Register Page Refactor

* Applied same structural pattern as Login:

  * `SafeScaffold`
  * `AppBar`
  * bottom action bar for primary CTA
* Introduced local `_viewmodel` reference
* Improved layout consistency:

  * removed redundant spacing widgets
  * standardized vertical rhythm using `spacing`
* Added explicit navigation method (`_navToLogin`)
* Ensures both auth screens follow the same **visual and interaction contract**

---

### 5. UI Behavior Improvements

* Centralized primary actions (Entrar / Cadastrar) in bottom area:

  * improves ergonomics on mobile devices
  * creates a consistent interaction pattern
* Added loading state handling directly in action buttons
* Improved keyboard UX with tap-to-dismiss behavior

---

### 6. Architectural Considerations

This change is subtle but important from a design perspective:

* Introduces a **UI composition layer**, analogous to how backend layers isolate responsibilities
* Reduces duplication while preserving flexibility
* Moves toward a **design system mindset**, even without formalizing one yet

A critical observation:
this abstraction is well-scoped. It does not attempt to generalize business logic or navigation, only layout concerns. This is a good boundary and avoids premature over-engineering.

---

### Conclusion

This commit establishes a **foundation for consistent UI composition**, improving maintainability, readability, and user experience.

The introduction of `SafeScaffold` combined with the refactoring of authentication screens represents a **clear step toward a scalable UI architecture**, mirroring the layered discipline already present in the backend.


## 2026/04/10 — infra/routing-01

Introduces a **structured routing architecture using GoRouter**, along with UI composition, dependency injection integration, and initial authentication flows. This commit establishes a clear separation of routing concerns aligned with a modular layered approach 

### 1. Routing Architecture Refactor

* Replaced monolithic route definition with **modular route groups**:

  * `authRoutes()`
  * `homeRoutes()`
* Router now composes routes using spread operators, improving scalability and readability
* Updated `initialLocation` to use `HomeRoutes.home.path`, removing reliance on generic enums

### 2. Route Definition Strategy

* Replaced generic `Routes` enum with **domain-oriented route enums**:

  * `AuthRoutes`
  * `HomeRoutes`
* Each enum encapsulates its own path, improving cohesion and reducing accidental coupling
* Introduced dedicated route files:

  * `routes/auth_routes.dart`
  * `routes/home_routes.dart`

Opinion: This is a strong architectural move. It prevents the typical “god enum” anti-pattern and aligns routing with feature boundaries.

### 3. GoRouter Integration

* Migrated from `MaterialApp` to `MaterialApp.router`
* Centralized router creation via `router()` factory
* Added `ExtraCodec` support for serialization:

  * now explicitly supports `null` values
  * prevents runtime failures when passing optional navigation data

### 4. Dependency Injection Integration

* Introduced `Uis.add(injector)` into dependency setup
* ViewModels are now resolved directly in route builders via injector:

  * `LoginViewModel`
  * `RegisterViewmodel`
  * `HomeViewmodel`
* Removed redundant LocalSecureStorage registration from `Data` layer, keeping DI responsibilities better distributed

Opinion: Injecting ViewModels at the routing boundary is a pragmatic choice. It keeps UI decoupled while avoiding premature abstraction layers.

### 5. Application Entry Point Refactor

* Renamed `MainApp` to `AppWidget`
* Moved it into `/uis`, reinforcing UI ownership
* Introduced internal router instance (`GoRouter`) inside the widget
* Replaced `home:` with `routerConfig`, aligning app initialization with navigation system

### 6. Authentication UI Implementation

#### Login Flow

* Implemented full `LoginPage`:

  * form validation (email/password)
  * loading state via `Command`
  * success/failure feedback using `SnackBar`
* Navigation:

  * success → `HomeRoutes.home`
  * register link → `AuthRoutes.register`

#### Register Flow

* Replaced placeholder with full implementation:

  * fields: name, email, cpf, password
  * validation rules for each field
  * command-based execution
* Navigation:

  * success → `AuthRoutes.login`

### 7. ViewModel Layer Introduction

* Added ViewModels:

  * `LoginViewModel`
  * `RegisterViewmodel`
  * `HomeViewmodel`
* Standardized usage of `Command1` for async actions
* Established consistent interaction pattern:

  * UI observes command state
  * ViewModel delegates to repository

### 8. UI Composition Adjustments

* `HomePage` now receives `HomeViewmodel` via constructor
* Ensures consistency with DI-driven UI pattern
* Created centralized `uis.dart` for ViewModel registration

### 9. Codebase Cleanup and Direction

* Removed unused imports and redundant DI registrations
* Added note to relocate `getProfile` from `AuthApi` to a future profile service
* Introduced (commented) navigation extension for future evaluation

### Conclusion

This commit represents a **foundational shift in navigation and UI architecture**, achieving:

* modular routing aligned with feature boundaries
* clean integration between routing and dependency injection
* consistent ViewModel-driven UI pattern
* scalable structure for future expansion (auth, home, and beyond)

From an architectural standpoint, this is a well-directed evolution. The system moves closer to a **feature-oriented modular design**, reducing global coupling and improving long-term maintainability.


## 2026/04/10 — infra/http-client-setup-01

Establishes a **centralized and environment-driven HTTP client configuration**, removing runtime mutation patterns and aligning the mobile client with a more deterministic and infrastructure-oriented design.

### 1. Environment Configuration Refactor

* Introduced `AppEnv` as the single source of truth for runtime configuration:

  * `baseUrl` with strict validation (non-empty and valid URI)
  * `connectTimeout` and `receiveTimeout` via compile-time environment variables
  * `AppMode` enum with explicit parsing and validation
* Removed legacy `EnviromentKey`, eliminating loosely validated configuration access
* This change enforces **fail-fast behavior**, which is a critical improvement for reliability in distributed systems

### 2. HTTP Client Design Simplification

* Removed `setBaseUrl` from `RestClient` interface and its implementation
* Eliminated runtime base URL mutation across the application layer
* All configuration is now resolved at instantiation time via `DioFactory`
* This is a **significant architectural improvement**, as it:

  * removes hidden side effects
  * avoids per-request configuration inconsistencies
  * enforces immutability of infrastructure concerns

### 3. DioFactory Redesign

* Refactored `DioFactory` to return a configured `Dio` instance instead of `RestClient`
* Integrated `AppEnv` directly into `BaseOptions`:

  * `baseUrl`
  * timeouts
  * default headers
* Added support for optional `defaultHeaders`
* Improved interceptor registration:

  * avoids duplicate interceptor instances using type comparison
* This aligns the HTTP client with an **infrastructure-first responsibility model**, consistent with layered architecture principles 

### 4. Dependency Injection Restructuring

* Reorganized `CoreServices` with explicit layering:

  1. `FlutterSecureStorage`
  2. `LocalSecureStorage` abstraction
  3. base `Dio` instance
  4. `AuthInterceptor` with isolated configuration
  5. `RestClient` composed from `Dio`
* Notable design decision:

  * `AuthInterceptor` uses a dedicated `Dio` instance to avoid recursive interception
* This setup improves:

  * testability
  * separation of concerns
  * predictability of request flow

### 5. API Layer Cleanup

* Removed manual base URL overrides from `AuthApi`
* All endpoints now rely on centralized configuration
* This eliminates duplication and prevents divergence across API calls
* Aligns the client with a **contract-driven API consumption model** 

### 6. Interceptor Behavior Clarification

* Updated `AuthInterceptor` comment to explicitly document behavior:

  * skips token injection when `Authorization` header is already present
* Improves readability and reduces ambiguity in request handling

### 7. Test Adjustments

* Updated `DioRestClient` tests:

  * removed dependency on `setBaseUrl`
  * now validate behavior based on `Dio` configuration
* Ensures tests reflect the new immutable configuration model

### Conclusion

This commit represents a **structural upgrade of the HTTP client layer**, shifting from mutable, scattered configuration to a **centralized, deterministic, and environment-driven approach**.

From an architectural standpoint, the most relevant gain is the clear separation between **application logic and infrastructure concerns**, reinforcing the principles of layered architecture and significantly reducing the risk of inconsistent network behavior across the application.


## 2026/04/09 — infra/di-and-env-setup-01

Establishes the **foundational infrastructure layer for dependency injection and environment configuration** in the Flutter client, aligning the mobile architecture with a modular, scalable structure and enabling controlled environment-based execution.

### 1. Development Environment Configuration

* Added `.vscode/launch.json` with predefined run configurations:

  * Dev, Staging, Prod
  * Integration test profile (Dev)
* Each configuration uses `--dart-define-from-file`, enabling **externalized environment configuration**
* Introduced `.env` file strategy (`dev.env`, `staging.env`, `prod.env`) and ensured they are ignored via `.gitignore`
* This approach is technically sound and aligns with production-grade practices for **environment isolation and reproducibility**

### 2. Dependency Injection Setup

* Introduced centralized DI configuration via `dependencies.dart`
* Adopted `AutoInjector` as DI container
* Implemented idempotent initialization (`_initialized` guard)
* Structured registration into modular layers:

  * `CoreServices`
  * `Services`
  * `Data`
* This is a **critical architectural improvement**, bringing the mobile project closer to the same separation principles already present in the backend 

### 3. Core Services Layer

* Added `CoreServices` module:

  * Registers `FlutterSecureStorage`
  * Configures `RestClient` via `DioFactory`
* Environment-driven configuration:

  * `baseUrl` via `EnviromentKey`
  * timeouts defined explicitly
* This enforces **centralized HTTP client configuration**, avoiding scattered setup across the codebase

### 4. Environment Abstraction

* Introduced `EnviromentKey`:

  * Maps compile-time variables using `String.fromEnvironment` and `int.fromEnvironment`
* Supports:

  * base URL
  * timeouts
  * app mode
  * access token (for internal usage)
* This design is particularly robust, as it avoids runtime parsing and ensures **compile-time guarantees**

### 5. Data Layer Composition

* Introduced `Data` module for DI registration:

  * `LocalSecureStorage` abstraction
  * `AuthRepository` implementation
* Proper dependency chaining:

  * Repository depends on API + storage
* This reinforces the **Repository as SSOT pattern**, consistent with your architectural direction

### 6. Services Layer Refactor

* Introduced `Services` module:

  * Registers `AuthApi` with injected `RestClient`
* Removed legacy empty `services.dart`
* Clean separation between:

  * core infrastructure (HTTP, storage)
  * feature services (API layer)

### 7. Authentication Repository Implementation

* Added `AuthRepository` contract and `AuthRepositoryImpl`
* Responsibilities:

  * manage authentication state (`currentUser`, `isLoggedIn`)
  * persist access token
  * handle login, logout, register, and profile
* Introduced explicit unauthenticated handling:

  * new `AppErrorCode.unauthenticated`
* This is a **well-structured implementation**, with clear boundaries between:

  * API (remote)
  * storage (local)
  * state (in-memory)

### 8. Storage and Auth Adjustments

* Renamed `authToken` → `accessToken` for semantic clarity
* Updated `AuthInterceptor` to use new key consistently
* Improved session lifecycle:

  * proper token write on login
  * cleanup on logout and refresh failure
* These changes reduce ambiguity and improve long-term maintainability

### 9. Application Bootstrap

* Updated `main.dart`:

  * introduced `setupDependencies()` before `runApp`
* Ensures all dependencies are resolved prior to UI initialization
* Aligns with proper application lifecycle control

### 10. Minor Improvements

* Adjusted imports in `AuthApi`
* Improved test launch configuration for integration tests
* Small consistency fixes across modules

### Conclusion

This commit introduces a **structural turning point in the mobile application architecture**.

Key gains:

* centralized dependency management
* environment-driven configuration
* clear separation of layers (core, services, data)
* improved authentication flow consistency

From an architectural perspective, this is a **necessary and well-executed foundation**, enabling the project to scale without accumulating coupling or configuration debt.


## 2026/04/09 — theme/composition-01

Introduces a structured **theme composition system** for the Flutter application, including dynamic theme resolution, Material 3 integration, custom typography, and improvements in developer tooling via Makefile refinements.

### 1. Theme Composition Architecture

* Refactored `MainApp` from `StatelessWidget` to `StatefulWidget` to support context-dependent initialization
* Introduced controlled theme composition flow:

  * resolve system brightness (`platformBrightness`)
  * select base theme (`light` / `dark`)
  * apply app-level overrides via `_buildAppTheme`
* Encapsulates theme creation logic, improving cohesion and avoiding scattered configuration across widgets
* This approach is conceptually aligned with layered responsibility principles, where configuration is centralized and isolated 

### 2. Material Theme Abstraction

* Added `MaterialTheme` class:

  * centralizes all `ColorScheme` definitions
  * supports multiple variants:

    * light / dark
    * medium contrast
    * high contrast
* Provides factory methods:

  * `light()`, `dark()`, and contrast variations
* Uses Material 3 (`useMaterial3: true`)
* Ensures consistency and scalability of design tokens across the application
* This is a **notable improvement in design maturity**, replacing ad-hoc theming with a reusable and extensible system

### 3. Typography System with Google Fonts

* Introduced `createTextTheme` helper:

  * composes two font families:

    * body font (Quicksand)
    * display font (EB Garamond)
* Uses `google_fonts` package for runtime font resolution
* Merges text styles to preserve semantic roles (`body`, `label`, etc.)
* Enables consistent typography without coupling UI components to font configuration

### 4. Dynamic Theme Initialization

* Theme is initialized in `didChangeDependencies`:

  * ensures access to `BuildContext`
  * avoids unnecessary recomputation
* Separation between:

  * theme construction (`MaterialTheme`)
  * runtime selection (`brightness`)
  * UI overrides (`AppBarTheme`)
* Improves maintainability and testability of UI configuration

### 5. UI Adjustments

* Updated `AppBar` styling:

  * uses `primaryContainer` and `onPrimaryContainer`
  * enforces semi-bold title (`FontWeight.w600`)
* Minor text change in HomePage:

  * "Home Page" → "Type Home Page"

### 6. Dependency Updates

* Added `google_fonts` dependency for typography support
* Introduced transitive dependency `http` (via ecosystem resolution)

### 7. Makefile Improvements

* Added `tests` target:

  * aggregates `api-test` and `mobile-test`
* Renamed Flutter commands for consistency and ergonomics:

  * `flutter-clean` → `fclean`
  * `flutter-build` → `fbuild`
* Added new utility:

  * `fadd pkg=<name>` to simplify dependency installation
* Improves developer experience and standardizes command usage across environments

### Conclusion

This commit establishes a **robust and scalable theming foundation**, transitioning from a basic configuration to a **composable design system** with clear separation of concerns.

From a technical standpoint, the introduction of a dedicated theme layer combined with dynamic resolution and Material 3 alignment significantly improves maintainability, consistency, and long-term extensibility of the UI layer.


## 2026/04/08 — main

Restructures the repository into a cohesive **monorepo architecture**, consolidating backend, mobile, infrastructure, and documentation while improving developer experience, build orchestration, and project clarity.

### 1. Monorepo Consolidation

* Introduced unified repository structure:

  * `api/` (Go backend)
  * `mobile/` (Flutter client)
  * `infra/` (Docker/infrastructure)
  * `docs/` (centralized documentation)
* Promoted project to a **full-stack system workspace**, aligning backend and mobile under a single lifecycle
* Reinforces the modular monolith approach described in the architecture documentation 

### 2. Documentation Reorganization

* Moved all API documentation from `api/docs/` → `docs/api/`
* Updated all internal references to reflect new structure
* Centralized architectural and API design artifacts:

  * architecture
  * domain model
  * use cases
  * API contract
* Improves discoverability and enforces documentation as a **first-class artifact of the system design**

### 3. Root-Level README Overhaul

* Replaced minimal README with comprehensive project documentation:

  * system purpose and engineering goals
  * architectural overview (layered modular monolith)
  * API capabilities and guarantees
  * mobile role as integration validator
  * local development workflow
* Explicitly documents:

  * transactional consistency strategy
  * concurrency handling (row-level locking)
  * API contract conventions
* Aligns with the REST contract and system behavior expectations 

### 4. Build and Tooling Unification

* Introduced root-level `Makefile` as a **monorepo task runner**
* Added commands:

  * Docker lifecycle (`docker-up`, `docker-down`, `docker-logs`)
  * Flutter utilities (`flutter-clean`, `flutter-build`)
* Removed duplicated Makefiles from:

  * `api/`
  * `mobile/`
* Establishes a **single entry point for all development workflows**, reducing operational fragmentation

### 5. Infrastructure Standardization

* Moved `docker-compose.yml` to repository root
* Simplifies environment setup and aligns with monorepo conventions
* Enables consistent orchestration across backend and mobile dependencies

### 6. Dependency Management Improvements (Go)

* Promoted key dependencies from indirect to direct:

  * `jwt`
  * `uuid`
  * `pgx`
  * `crypto`
* Updated `go.sum` with explicit versions and additional test dependencies (`testify`, `difflib`)
* Improves dependency clarity and reproducibility of builds

### 7. Repository Hygiene

* Added `.gitignore` covering:

  * Go build artifacts
  * Flutter build/cache directories
  * environment files and OS artifacts
* Introduced MIT `LICENSE`, formalizing project usage and distribution rights

### 8. API Project Adjustments

* Updated `api/README.md`:

  * aligned commands with new root Makefile
  * corrected build paths (`api/build/`)
  * updated documentation links to `docs/api/`
* Ensures consistency between documentation and actual project structure

### Conclusion

This commit represents a **structural milestone** rather than a feature addition.

Key impacts:

* Establishes a **clean monorepo foundation**
* Improves **developer ergonomics and workflow consistency**
* Elevates documentation to a **core part of the system design**
* Aligns project organization with its architectural principles

From an engineering perspective, this is a highly valuable refactor that reduces cognitive load, eliminates duplication, and prepares the codebase for scalable evolution across both backend and mobile layers.
