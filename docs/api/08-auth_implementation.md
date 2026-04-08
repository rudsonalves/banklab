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
  * `ParseAccessToken`

**Claims:**

* `sub`
* `role`
* `cid`
* `exp`, `iat`

**Done when:**

* Token is generated and parsed correctly
* Invalid/expired tokens are rejected

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
* generate JWT

**Done when:**

* valid login returns token
* invalid password fails

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
* `GET /auth/me`

**Requirements:**

* DTO separation
* no domain exposure

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
* JWT parsing
* register use case
* login use case
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

---

### **Task 8.2 — Token validation edge cases**

* expired token
* malformed token
* missing header

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
