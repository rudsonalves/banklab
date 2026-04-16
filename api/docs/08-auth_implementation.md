# **Auth & Authorization — Progressive Task Plan**

## **Phase 0 — Foundation (Database + Domain Contracts)**

### **Task 0.1 — Create users table migration**

* Create migration files:

  * `xxxxxx_create_users_table_up.sql`
  * `xxxxxx_create_users_table_down.sql`
* Include:

  * `id (UUID)`
  * `email (UNIQUE)`
  * `password_hash`
  * `role`
  * `customer_id (nullable, unique, FK)`
  * timestamps
* Validate migration runs successfully

**Done when:**

* `migrate up` creates table
* `migrate down` rolls back cleanly

---

### **Task 0.2 — Define User domain entity**

* Create `auth/domain/user.go`
* Define:

  * struct `User`
  * role type (typed, not raw string)
* No external dependencies

**Done when:**

* Entity compiles
* No infra leakage

---

### **Task 0.3 — Define domain interfaces**

Create interfaces in `auth/domain`:

* `UserRepository`
* `PasswordHasher`
* `TokenService`
* `SessionRepository` — `Create`, `FindByTokenHash`, `Revoke`
* `Transactor` — `RunInTx(ctx, fn)` — controls database transactions from the Application layer

**Done when:**

* Interfaces are minimal and clear
* No implementation yet

---

## **Phase 1 — Infrastructure Adapters**

### **Task 1.1 — PostgreSQL UserRepository**

* Implement repository in `auth/infrastructure`
* Methods:

  * `Create`
  * `FindByEmail`
  * `FindByID`
  * `ExistsByEmail`

**Edge cases:**

* duplicate email
* null customer_id

**Done when:**

* Repository works with real DB
* Basic integration test passes

---

### **Task 1.2 — Bcrypt PasswordHasher**

* Implement:

  * `Hash(password)`
  * `Compare(hash, password)`

**Constraints:**

* never expose password
* proper error handling

**Done when:**

* Hash + compare works
* Wrong password fails correctly

---

### **Task 1.3 — JWT TokenService**

* Implement:

  * `GenerateAccessToken`
  * `GenerateRefreshToken` — returns a random opaque token (signed, high-entropy)
  * `ParseAccessToken`
  * `ParseRefreshToken`

**Claims:**

* `sub`
* `role`
* `cid`
* `exp`, `iat`

**Done when:**

* Token is generated and parsed correctly
* Invalid/expired tokens are rejected

---

### **Task 1.4 — PostgreSQL SessionRepository**

* Implement in `auth/infrastructure`
* Methods:

  * `Create` — inserts a new session row
  * `FindByTokenHash` — returns `userID`, `expiresAt`, `revoked`
  * `Revoke` — sets `revoked_at = NOW()`; returns `ErrSessionNotFound` if zero rows affected

* Uses context-based transaction propagation via `database.TxFromContext`

**Done when:**

* Session lifecycle works end-to-end with real DB
* `Revoke` fails cleanly when token hash is unknown

---

### **Task 1.5 — PostgresTransactor**

* Implement `domain.Transactor` in `auth/infrastructure`
* `RunInTx` — begins a pgx transaction, injects it into context via `database.ContextWithTx`, defers rollback, commits on success
* Satisfies `var _ domain.Transactor = (*PostgresTransactor)(nil)`

**Done when:**

* Any two repository calls inside `RunInTx` share the same transaction
* Rollback occurs automatically on error

---

## **Phase 2 — Application Layer (Use Cases)**

### **Task 2.1 — RegisterUser use case**

Flow:

* validate email
* validate password
* check email uniqueness
* hash password
* create user
* persist

**Decisions:**

* do NOT require `customer_id` initially

**Done when:**

* user is created
* duplicate email fails

---

### **Task 2.2 — LoginUser use case**

Flow:

* find user by email
* compare password
* generate access token
* generate refresh token
* hash refresh token with SHA-256
* persist session via `SessionRepository.Create` (TTL: 30 days)
* return both tokens

**Constraints:**

* `SessionRepository` is a required dependency (not optional)
* If session persistence fails, login fails — no token is returned
* No code path issues a refresh token without a persisted session

**Done when:**

* valid login returns both `access_token` and `refresh_token`
* invalid password fails
* session creation failure propagates as an error

---

### **Task 2.4 — RefreshAccessToken use case**

Flow:

* validate token is not blank
* parse JWT signature to extract `userID`
* hash token with SHA-256
* look up session by hash — reject if not found, revoked, expired, or user mismatch
* fetch user from `UserRepository`
* generate new access token
* generate new refresh token
* atomically via `Transactor.RunInTx`:
  * revoke old session (`SessionRepository.Revoke`)
  * create new session (`SessionRepository.Create`)
* return both new tokens

**Atomicity guarantee:**

* Revoke + Create execute in a single database transaction
* If either step fails the transaction is rolled back
* The old token remains valid when rollback occurs — no partial state is possible

**Done when:**

* valid token returns new `access_token` + `refresh_token`
* old token is revoked after rotation
* new token is immediately usable
* rollback prevents partial state

---

### **Task 2.3 — GetCurrentUser use case**

Flow:

* read authenticated principal
* optionally fetch user

**Done when:**

* returns correct identity
* works with middleware context

---

## **Phase 3 — Delivery Layer (HTTP)**

### **Task 3.1 — Auth handlers**

Implement:

* `POST /auth/register`
* `POST /auth/login`
* `POST /auth/refresh`
* `GET /auth/me`

**Requirements:**

* DTO separation
* no domain exposure
* all handlers guard against nil use cases and nil output
* `Handler.New` takes all four use cases as required constructor parameters — no setter injection

**Done when:**

* endpoints respond correctly
* JSON format follows standard

---

### **Task 3.2 — Standardized error handling**

* Integrate auth errors into existing error pattern

Add codes:

* `USER_ALREADY_EXISTS`
* `INVALID_CREDENTIALS`
* `UNAUTHORIZED`
* `INVALID_TOKEN`

**Done when:**

* no raw `http.Error`
* all responses follow `{data, error}` format

---

## **Phase 4 — Authentication Middleware**

### **Task 4.1 — JWT middleware**

* Read `Authorization` header
* Validate token
* Extract claims
* Inject principal into context

**Principal:**

```text
userID
role
customerID
```

**Done when:**

* valid token populates context
* invalid token returns 401

---

### **Task 4.2 — Context helpers**

* Helper functions:

  * `GetAuthenticatedUser(ctx)`
  * `MustGetAuthenticatedUser(ctx)`

**Done when:**

* handlers/use cases can easily access identity

---

## **Phase 5 — Authorization (Critical Layer)**

### **Task 5.1 — Ownership validation logic**

Implement a reusable function/service:

```text
CanAccessAccount(user, account)
```

Rules:

* user.customer_id == account.customer_id
* OR user.role == admin

**Done when:**

* ownership logic is centralized
* no duplication across handlers

---

### **Task 5.2 — Integrate authorization into use cases**

Update:

* GetBalance
* GetStatement
* Deposit
* Withdraw
* Transfer

**Rules:**

* must be authenticated
* must own account (or admin)

**IMPORTANT:**

* enforce in application layer, not only in handler

**Done when:**

* unauthorized access is blocked
* correct access is allowed

---

## **Phase 6 — Route Protection**

### **Task 6.1 — Protect account routes**

Apply middleware to:

```text
/accounts/*
```

**Done when:**

* unauthenticated requests fail
* authenticated requests proceed

---

### **Task 6.2 — Transfer-specific rule**

* Validate:

  * user owns **source account**

**Done when:**

* cannot transfer from чужой account
* admin can override

---

## **Phase 7 — Testing**

### **Task 7.1 — Unit tests**

Cover:

* password hashing
* JWT parsing (including `ParseRefreshToken`)
* register use case
* login use case — including `SessionPersistenceFailure`
* refresh use case:

  * success
  * invalid / malformed token
  * session not found
  * revoked token
  * expired token
  * user-ID mismatch
  * revoke failure (Create not called)
  * create failure
  * rotation integrity (old token unusable, new token works) — uses stateful session mock
* ownership logic

---

### **Task 7.2 — Integration tests**

Test flows:

* register → login → access `/auth/me`
* access own account → success
* access another account → forbidden
* admin access → allowed

---

## **Phase 8 — Hardening (Minimal)**

### **Task 8.1 — Input validation**

* email format
* password length
* refresh token: blank check before any DB access

---

### **Task 8.2 — Token validation edge cases**

* expired token
* malformed token
* missing header
* revoked refresh token
* expired refresh token

---

### **Task 8.3 — Revoke integrity**

* `SessionRepository.Revoke` checks `RowsAffected()` after the UPDATE
* Returns `domain.ErrSessionNotFound` when zero rows are affected
* This prevents silent failures when the hash never existed or was already deleted

---

### **Task 8.4 — Atomic token rotation**

* `Revoke` and `Create` are wrapped in a single `Transactor.RunInTx` call in the Application layer
* Infrastructure repositories detect the transaction in context via `database.TxFromContext` and use it automatically
* If either operation fails the transaction rolls back — the client retains the original token

---

## **Suggested Branching Strategy**

Use incremental branches:

```text
auth/users-table-01
auth/domain-01
auth/infrastructure-01
auth/usecases-01
auth/http-01
auth/middleware-01
auth/authorization-01
auth/tests-01
```

---

# **Execution Strategy (Important)**

Do NOT jump across phases.

Recommended order:

```text
DB → Domain → Infrastructure → UseCases → HTTP → Middleware → Authorization → Tests
```

---

# **Critical Opinion**

The most common failure here would be:

> implementing login before ownership

If you reach a point where users can authenticate but can still access any account by ID, the system is **functionally insecure**.

So treat this as the real milestone:

> Authentication is necessary
> Authorization is what makes it correct
