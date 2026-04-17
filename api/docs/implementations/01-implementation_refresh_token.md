# Implement Refresh Token with Clean Architecture

You are working on a Go-based banking API with a layered architecture:

```text
Delivery → Application → Domain
                  ↓
             Infrastructure
```

## Mandatory rules

* Do NOT violate dependency direction
* Do NOT introduce business logic outside Domain/Application
* Do NOT couple infrastructure into Domain
* All changes must be incremental and keep the system working
* Avoid unnecessary refactors outside the scope of each task

---

# Goal

Evolve the current authentication system (which uses a non-expiring access token) into a proper model with:

* Expiring Access Token
* Refresh Token
* Future support for session control (implemented incrementally)

---

# Task 1 — Expand TokenService contract

## Goal

Support refresh token generation and parsing

## Scope

Update the domain interface:

```go
type TokenService interface {
	GenerateAccessToken(claims TokenClaims) (string, error)
	GenerateRefreshToken(userID uuid.UUID) (string, error)

	ParseAccessToken(token string) (*TokenClaims, error)
	ParseRefreshToken(token string) (uuid.UUID, error)
}
```

## Rules

* Do NOT implement logic yet (interface only)
* Do NOT modify other layers

## Done when

* Interface compiles
* No concrete implementation changed yet

---

# Task 2 — Update LoginUserUseCase

## Goal

Return refresh token on login

## Scope

Update output:

```go
type LoginUserOutput struct {
	AccessToken  string
	RefreshToken string
	UserID       uuid.UUID
	Email        string
	Role         string
	CustomerID   *uuid.UUID
}
```

Update execution:

```go
accessToken, err := uc.tokenService.GenerateAccessToken(...)
refreshToken, err := uc.tokenService.GenerateRefreshToken(user.ID)
```

## Rules

* Do NOT change existing validations
* Do NOT add persistence yet
* Do NOT change authentication flow

## Done when

* Login returns both tokens
* Code compiles

---

# Task 3 — Update Delivery layer (HTTP)

## Goal

Expose refresh token in REST contract

## Scope

Update `/auth/login` response:

```json
{
  "access_token": "...",
  "refresh_token": "...",
  ...
}
```

## Rules

* Preserve `{ data, error }` envelope
* Do NOT modify global API structure
* Do NOT touch other endpoints

## Done when

* Endpoint returns refresh token

---

# Task 4 — Add Access Token expiration

## Goal

Introduce TTL for access tokens

## Scope

In TokenService implementation:

* Add `exp` claim to JWT
* Suggested TTL: 15 minutes

## Rules

* Do NOT change existing claims
* Do NOT change token format
* Do NOT break ParseAccessToken

## Done when

* Token expires correctly
* Expired tokens are rejected

---

# Task 5 — Implement Refresh Token generation

## Goal

Generate secure refresh tokens

## Scope

* Use either JWT or opaque token (prefer opaque)
* Must identify the user (directly or indirectly)

## Rules

* Do NOT persist yet
* Do NOT implement revocation
* Ensure high entropy (use crypto/rand)

## Done when

* Refresh token is generated and returned

---

# Task 6 — Create /auth/refresh endpoint

## Goal

Allow issuing a new access token

## Scope

Create endpoint:

```
POST /auth/refresh
```

Input:

```json
{
  "refresh_token": "..."
}
```

Output:

```json
{
  "access_token": "..."
}
```

## Flow

1. Parse refresh token
2. Extract userID
3. Load user
4. Generate new access token

## Rules

* Do NOT persist tokens yet
* Do NOT implement rotation yet
* Do NOT modify login

## Done when

* Endpoint works correctly
* New access token is returned

---

# Task 7 — Introduce SessionRepository (preparation)

## Goal

Prepare for session control

## Scope

Create domain interface:

```go
type SessionRepository interface {
	Create(...)
	FindByTokenHash(...)
	Revoke(...)
}
```

## Rules

* Do NOT implement yet
* Do NOT integrate into flow

## Done when

* Interface exists
* No system impact

---

# Task 8 — Persist refresh tokens

## Goal

Store refresh tokens securely

## Scope

Infrastructure:

* Create `user_sessions` table
* Store hashed refresh tokens

## Rules

* NEVER store raw tokens
* Use hashing (e.g., SHA-256)

## Done when

* Session is stored during login

---

# Task 9 — Validate refresh token against storage

## Goal

Ensure real token validity

## Scope

In `/auth/refresh`:

* Validate token hash in DB
* Check expiration
* Check revocation

## Done when

* Invalid tokens are rejected
* Valid tokens succeed

---

# Task 10 — Refresh token rotation (recommended)

## Goal

Prevent replay attacks

## Scope

On refresh:

* Revoke old token
* Generate new refresh token
* Persist new token

## Done when

* Old refresh token no longer works

---

# Global rules for all tasks

* Do NOT move code across layers without clear reason
* Do NOT introduce external frameworks
* Keep consistency with existing error handling
* Each task must be independently testable
* Prefer small, verifiable changes
