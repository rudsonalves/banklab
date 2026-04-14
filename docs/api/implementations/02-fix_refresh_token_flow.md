# AI Prompt — Fix Refresh Token Flow (Consistency & Safety)

You are working on a Go-based banking API with layered architecture:

```text
Delivery → Application → Domain
                  ↓
             Infrastructure
```

## Critical rules

* Do NOT violate dependency direction
* Do NOT move business logic out of Application/Domain
* Do NOT introduce new frameworks
* All changes must be incremental and keep the system working
* Each task must be independently testable

---

# 🎯 Goal

Fix the refresh token implementation to ensure:

* **Consistency of session lifecycle**
* **Transactional safety during refresh**
* **Correct test behavior**
* **Strict contract integrity**

---

# Task 1 — Make SessionRepository mandatory in Login

## Problem

Login currently allows issuing refresh tokens without persisting sessions.

This breaks the contract: a refresh token may be returned but unusable.

## Objective

Ensure every issued refresh token has a persisted session.

## Scope

Update constructor:

```go
func NewLoginUserUseCase(
	userRepo domain.UserRepository,
	hasher domain.PasswordHasher,
	tokenService domain.TokenService,
	sessionRepo domain.SessionRepository,
) *LoginUserUseCase
```

Remove:

* variadic dependency
* nil checks

Remove:

```go
if uc.sessionRepo != nil {
	...
}
```

Always persist session.

## Rules

* Do NOT change business logic
* Do NOT change token generation

## Done when

* Login always creates a session
* No optional session behavior remains

---

# Task 2 — Make refresh token rotation atomic

## Problem

Refresh flow:

* revokes old token
* creates new session

If Create fails → session is lost

## Objective

Ensure refresh token rotation is **atomic**

## Scope

Wrap this block in a transaction:

```go
Revoke(oldToken)
Create(newToken)
```

## Implementation constraint

* Use existing transaction mechanism (context-based if available)
* DO NOT move logic to infrastructure
* Transaction must be controlled in Application layer

## Expected behavior

* If any step fails → rollback everything
* Old token must remain valid if new session fails

## Done when

* Rotation is atomic
* No partial state possible

---

# Task 3 — Fix broken test assertion

## Problem

This condition is invalid:

```go
errors.Is(err, err)
```

Always true → test is meaningless

## Objective

Fix test validation logic

## Scope

Replace:

```go
if !errors.Is(err, err) && !strings.Contains(err.Error(), tc.errText)
```

With:

```go
if !strings.Contains(err.Error(), tc.errText)
```

## Rules

* Do NOT weaken test assertions
* Prefer explicit error validation

## Done when

* Tests correctly fail when expected
* Assertions are meaningful

---

# Task 4 — Validate revoke operation

## Problem

`Revoke()` does not confirm if a session was actually updated

## Objective

Ensure revoke operation is effective

## Scope

In repository:

* Check affected rows
* If zero rows → return error

## Example

```go
result, err := Exec(...)
rows := result.RowsAffected()
if rows == 0 {
	return error
}
```

## Rules

* Do NOT change method signature
* Keep behavior simple

## Done when

* Revoke fails when token does not exist

---

# Task 5 — Align handler construction

## Problem

Handler uses setter injection for refresh use case

This weakens construction safety

## Objective

Make handler fully initialized at construction

## Scope

Update constructor:

```go
func New(
	registerUser registerUserUseCase,
	loginUser loginUserUseCase,
	getCurrentUser getCurrentUserUseCase,
	refreshAccessToken refreshAccessTokenUseCase,
) *Handler
```

Remove:

```go
SetRefreshAccessTokenUseCase
```

## Rules

* Do NOT change handler behavior
* Only improve initialization safety

## Done when

* Handler is immutable after creation
* No nil dependencies at runtime

---

# Task 6 — Add missing nil check consistency in Login handler

## Problem

Login handler does not check `output == nil`, unlike others

## Objective

Ensure consistent defensive behavior

## Scope

Add:

```go
if output == nil {
	sharedhttp.WriteError(...)
	return
}
```

## Done when

* Login handler matches other handlers

---

# Task 7 — Add Application-level tests for refresh flow

## Objective

Validate real system behavior (not just token service)

## Scope

Add tests for:

### Case 1 — Successful refresh

* valid token
* session exists
* returns new access + refresh token

### Case 2 — Revoked token

* must fail

### Case 3 — Expired token

* must fail

### Case 4 — Token not found

* must fail

### Case 5 — Rotation integrity

* old token no longer usable
* new token works

## Rules

* Mock repositories
* Test only Application layer

## Done when

* Refresh flow is fully covered

---

# Global constraints

* Do NOT refactor unrelated code
* Do NOT introduce new abstractions
* Maintain current naming conventions
* Preserve error mapping behavior

---

# Final note

This is not a feature change.

This is a **consistency and correctness fix** in authentication flow.

Focus on:

* eliminating invalid states
* guaranteeing atomic operations
* ensuring contracts are never violated
