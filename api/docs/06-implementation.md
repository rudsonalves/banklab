# Implementation Documentation - Bank API

## 1. Scope

This document describes the implementation currently present in the codebase.

It focuses on:
- package structure and dependency flow
- runtime startup and wiring
- domain entities and business rules
- application use cases
- HTTP delivery behavior
- PostgreSQL persistence model
- consistency and concurrency strategy
- test coverage and execution

This is an implementation-oriented guide, not only a target architecture description.

## 2. High-Level Architecture

The project follows a layered modular monolith style:

- Delivery layer: HTTP handlers and DTO mapping
- Application layer: use case orchestration
- Domain layer: entities, invariants, contracts, business errors
- Infrastructure layer: PostgreSQL repository implementations

Dependency direction:

- Delivery -> Application -> Domain
- Infrastructure -> Domain

Main entrypoint:

- cmd/api/main.go

## 3. Folder Map (Implemented)

- cmd/api: process bootstrap and route registration
- internal/database: PostgreSQL pool creation
- internal/customer:
  - domain: customer entity, validation errors, repository contract
  - application: create customer use case, get customer me use case
  - delivery: HTTP handler for GET /customers/me
  - infrastructure: PostgreSQL repository implementation
- internal/account:
  - domain: account and transaction entities, repository contracts, account rules
  - application: create account, deposit, withdraw, transfer, statement
  - delivery: HTTP handlers, request parsing, response mapping
  - infrastructure: PostgreSQL repository with transactional and locking support
- db/schema.sql: baseline schema including customers, accounts, account_transactions, users, user_sessions
- migrations: additive files for account number sequence, ledger consolidation, and transfer pair integrity

## 4. Bootstrap and Wiring

Process startup sequence (cmd/api/main.go):

1. Create PostgreSQL pool via internal/database.NewPool
2. Build auth repository and use cases (register, login, get current user)
3. Build customer repository and use cases (get customer me)
4. Build account repository and all account use cases
5. Build auth, customer, and account handlers
6. Register HTTP routes with net/http ServeMux style patterns
7. Start server on port 8080

Current route registration:

- POST /auth/register (public)
- POST /auth/login (public)
- GET /auth/me (JWT required)
- GET /customers/me (JWT required)
- POST /accounts (JWT required)
- POST /accounts/{id}/deposit (JWT required)
- POST /accounts/{id}/withdraw (JWT required)
- GET /accounts/{id}/statement (JWT required)
- POST /accounts/transfer (JWT required)

## 5. Domain Model

### 5.1 Customer Domain

Entity fields:
- id (UUID)
- name
- cpf
- email
- created_at

Factory validation (NewCustomer):
- name required
- cpf required
- email required

Domain errors:
- ErrNameRequired
- ErrCPFRequired
- ErrEmailRequired
- ErrCPFAlreadyExists
- ErrEmailAlreadyExists
- ErrInvalidData
- ErrNotFound

### 5.2 Account Domain

Entity fields:
- id (UUID)
- customer_id (UUID)
- number
- branch
- balance (int64 cents)
- status (active, inactive, blocked)
- created_at

Factory behavior (NewAccount):
- rejects nil customer ID
- initializes balance = 0
- initializes status = active
- sets ID and created_at

Business rule methods:
- CanDeposit(amount)
  - amount must be > 0
  - account must be active
- CanWithdraw(amount)
  - amount must be > 0
  - account must be active
  - balance must be sufficient
- CanTransfer(amount, destinationID)
  - destination must be different from source
  - reuses withdraw validation chain

Account domain errors:
- ErrInvalidData
- ErrInvalidAmount
- ErrAccountNotFound
- ErrInsufficientBalance
- ErrSameAccountTransfer
- ErrCustomerNotFound
- ErrAccountInactive
- ErrForbidden

### 5.3 Transaction Domain

Transaction types:
- deposit
- withdraw
- transfer_out
- transfer_in

Entity fields:
- id
- account_id
- type
- amount
- balance_after
- reference_id (nullable UUID)
- related_account_id (nullable UUID)
- idempotency_key (nullable string)
- created_at

Factory behavior:
- NewTransaction creates an immutable transaction object with generated UUID and UTC timestamp.

## 6. Application Use Cases

### 6.1 Register User (Auth)

Input:
- email
- password
- name
- cpf

Flow:
1. Start transaction
2. Validate email format and password strength
3. Check email uniqueness
4. Create Customer entity and persist
5. Create User entity via `domain.NewUser` (enforces: RoleCustomer requires non-nil customer_id)
6. Persist user
7. Commit transaction
8. Post-transaction invariant check: user.CustomerID must not be nil

The customer is always created before the user. If any step fails, the transaction rolls back and no partial state is persisted.

### 6.2 Approve User

Input:
- user_id
- authenticated user context (admin role required)

Flow:
1. Start transaction
2. Load user with lock (FOR UPDATE)
3. Validate user exists (returns ErrNotFound if not)
4. Validate user status == pending (returns appropriate error if already active or blocked)
5. Update user.status → active
6. Generate account number using sequence
7. Create Account entity (branch fixed as "0001", balance = 0, status = active)
8. Persist account
9. Update user.customer_id with newly created account's customer_id (if needed)
10. Persist user updates
11. Commit transaction

Atomicity guarantee:
- Status transition and account creation happen within same transaction
- No partial state possible: either both succeed or both rollback

Dependencies:
- UserRepository for user lookup and update
- AccountRepository for account creation
- Sequence generator for account number

### 6.3 Create Account

Input:
- authenticated user context (customer_id derived from token principal)

Flow:
1. Validate authenticated user has non-nil customer_id (returns ErrForbidden otherwise)
2. Validate user status == active (returns ErrForbidden if pending or blocked)
3. Validate customer_id is not zero UUID (returns ErrForbidden)
4. Validate user owns the customer_id via CanAccessCustomer
5. Ensure customer exists in the database
6. Generate account number using sequence
7. Build account entity (branch currently fixed as "0001")
8. Persist and return account

Notes:
- The client MUST NOT and CANNOT provide customer_id — it is ignored if sent.
- Optional one-account-per-customer rule exists but is currently commented out.
- User must be active (enforced at application layer) to create accounts

### 6.4 Deposit

Input:
- account_id
- amount

Flow:
1. Validate input
2. Begin DB transaction
3. Load account with row lock (GetByIDForUpdate)
4. Validate with CanDeposit
5. Increment balance (UpdateBalance)
6. Insert ledger row in account_transactions
7. Commit

Rollback strategy:
- deferred rollback executes unless commit succeeds.

### 6.5 Withdraw

Input:
- account_id
- amount

Flow:
1. Validate input
2. Begin DB transaction
3. Load account with row lock
4. Validate with CanWithdraw
5. Decrease balance atomically (DecreaseBalance)
6. Insert withdraw ledger row
7. Commit

### 6.6 Transfer

Input:
- from_account_id
- to_account_id
- amount

Flow:
1. Validate IDs, amount, and distinct accounts
2. Begin DB transaction
3. Lock both accounts in deterministic UUID order
4. Validate source with CanTransfer
5. Validate destination with CanDeposit
6. Decrease source balance
7. Increase destination balance
8. Insert transfer_out and transfer_in ledger rows sharing same reference_id
9. If idempotency key is present, enforce deduplication on transfer_out by `(account_id, idempotency_key)`
10. On duplicate, rollback duplicate mutation and replay deterministic result from ledger pair
11. Commit

Deadlock mitigation:
- lock ordering by UUID bytes using orderedUUIDs.

### 6.7 Get Statement

Input:
- account_id
- limit
- cursor + cursor_id pair
- from / to (optional date filters)

Flow:
1. Validate account ID and query consistency
2. Normalize limit (default 50, cap 100)
3. Ensure account exists
4. Query account_transactions ordered by created_at desc, id desc
5. Map rows to API statement items
6. Build next cursor if full page returned

### 6.8 Get Customer Me

Input:
- authenticated user context (customer_id derived from token principal)

Flow:
1. Validate user.CustomerID is not nil (returns ErrInvalidData otherwise)
2. Query customer by ID via repository
3. Return nil result as ErrNotFound
4. Return customer data

## 7. HTTP Delivery Layer

## 7.1 Account Handler Endpoints

- CreateAccount
- Deposit
- Withdraw
- Transfer
- Statement

Implemented concerns:
- JSON request decoding
- path/query parsing and UUID/time validation
- call application use cases
- map domain errors to HTTP status and stable error codes
- return JSON response with data/error envelope

## 7.2 Customer Handler Endpoints

- Me (GET /customers/me)

Implemented concerns:
- Reads authenticated user from request context
- Returns 401 if no user in context
- Returns 400 if user has no customer_id
- Calls GetCustomerMe use case
- Maps ErrNotFound to 404

## 7.3 Auth Handler Endpoints

- Register (POST /auth/register)
- Login (POST /auth/login)
- Me (GET /auth/me)

Implemented concerns:
- Minimal delivery validation: rejects blank required fields
- Domain/format validation delegated to application layer
- Returns customer_id in register and login responses

## 7.4 Response Contract

Response envelope format:

Success:
- data: object
- error: null

Error:
- data: null
- error.code
- error.message
- error.details (optional)

Common status mapping examples:
- 400 for invalid input
- 404 for not found
- 409 for conflicts (customer duplicates)
- 422 for business rule violations (insufficient balance, inactive account)
- 500 for internal failures

## 8. Persistence Implementation

## 8.1 DB Pool and Connectivity

internal/database.NewPool creates a pgxpool connection using a hard-coded URL:
- postgres://postgres:postgres@localhost:5432/bank?sslmode=disable

## 8.2 Account Repository (PostgreSQL)

Implemented repository behaviors:
- account creation
- transaction row insertion in account_transactions
- account lookup and row-lock lookup
- statement query with cursor pagination
- balance increment and conditional decrement
- transaction begin/commit/rollback

Transaction abstraction:
- BeginTx returns a txRepository implementing domain.Tx
- txRepository implements same account repository methods using pgx.Tx
- nested BeginTx in txRepository is blocked and returns an error

## 8.3 Customer Repository (PostgreSQL)

Implemented behaviors:
- create customer
- check customer existence by ID

Error conversion:
- unique violation on customers_cpf_key -> ErrCPFAlreadyExists
- unique violation on customers_email_key -> ErrEmailAlreadyExists
- check violation -> ErrInvalidData
- unknown infra failures are wrapped with context

## 8.4 Schema and Migrations

Primary relational objects in db/schema.sql:
- customers
- accounts
- account_transactions
- users
- user_sessions

Applied migration files include:
- account_number_seq sequence
- account_transactions table + immutability trigger + indexes
- ledger consolidation from `transactions` to `account_transactions`
- transfer pair integrity indexes (`reference_id`, `reference_id+type` unique per transfer leg)

Important implementation detail:
- account operations persist ledger entries exclusively in `account_transactions`.
- idempotent transfer replay is reconstructed from ledger rows (`transfer_out` + paired `transfer_in`) using `reference_id`.

## 9. Consistency and Concurrency Strategy (Implemented)

Current implementation enforces consistency by:
- explicit database transactions in balance-changing use cases
- row-level locking with SELECT ... FOR UPDATE for deposit, withdraw, transfer
- deterministic dual-row lock ordering in transfer
- atomic conditional update for decrement to prevent overdraft races
- rollback on any intermediate failure

This provides strong immediate consistency for critical account operations.

## 10. Test Coverage

Implemented tests include:

Domain tests:
- account invariants and rule methods
- `domain.NewUser` invariant (RoleCustomer requires non-nil customer_id)

Application tests:
- register user (transactional creation, invariant enforcement)
- create account
- deposit (including ownership enforcement)
- withdraw (including ownership enforcement)
- transfer (including source ownership enforcement)
- get statement (including ownership enforcement)
- access policy helpers (CanAccessCustomer, CanAccessAccount)
- get customer me

Delivery tests:
- account handler unit tests for success and error mappings
- auth handler unit tests
- customer handler unit tests (GET /customers/me)
- deposit integration test with real PostgreSQL

## 11. Local Run and Validation

### 11.1 Dependencies

- Go 1.26.1 (go.mod)
- PostgreSQL 16 (docker-compose)

### 11.2 Start database

- docker compose up -d

### 11.3 Run tests

- go test ./...

### 11.4 Run API

- go run ./cmd/api

Server listens on:
- :8080

## 12. Known Implementation Notes

- README.md is currently empty; docs folder is the main source of project documentation.
- DB connection string is hard-coded in code and not yet externalized via environment variables.
- Branch generation for accounts is currently fixed to "0001".
- A one-account-per-customer rule is scaffolded but not active.
- Ledger is single-source (`account_transactions`) and append-only.

## 13. Summary

The current implementation already contains a complete vertical slice for:
- customer creation
- account lifecycle operations (create, deposit, withdraw, transfer)
- statement retrieval with pagination and date filters

It is implemented with transactional integrity, row-level locking, domain-level invariants, and practical test coverage appropriate for a financial core service baseline.
