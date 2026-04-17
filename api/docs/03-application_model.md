# Application Model

## Role of the Application Model

The **Application Model** describes how the domain is represented and manipulated inside the application (Go).

In practical terms, it is the bridge between:

* **Domain** (concepts and invariants)
* **Database** (persistence)
* **API** (input and output)

In one line:

> **The Application Model is the system's operational model in code.**

---

## What the Application Model Must Contain

### 1. Application Structures (Go structs)

Models used to operate use cases in code.

Includes:

* representations of domain entities adapted for application use
* value objects adapted for practical use

Examples:

* `Account`
* `Transaction`
* `User`

Concrete examples from the current implementation:

```go
type Account struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	Number     string
	Branch     string
	Balance    int64
	Status     AccountStatus
	CreatedAt  time.Time
}

type Transaction struct {
	ID               uuid.UUID
	AccountID        uuid.UUID
	Type             TransactionType
	Amount           int64
	BalanceAfter     int64
	ReferenceID      *uuid.UUID
	RelatedAccountID *uuid.UUID
	IdempotencyKey   *string
	CreatedAt        time.Time
}

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Role         Role
	CustomerID   *uuid.UUID
	Status       UserStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Customer struct {
	ID        uuid.UUID
	Name      string
	CPF       string
	Email     string
	CreatedAt time.Time
}
```

### 2. DTOs (Input and Output)

Transport models used between application layers.

Includes:

* use case inputs
* intermediate outputs

Examples:

* `CreateCustomerInput`
* `TransferInput`

Concrete use case DTOs:

```go
type TransferInput struct {
	User           *authdomain.AuthenticatedUser
	FromAccountID  uuid.UUID
	ToAccountID    uuid.UUID
	Amount         int64
	IdempotencyKey string
}

type LoginUserInput struct {
	Email    string
	Password string
}

type LoginUserOutput struct {
	AccessToken  string
	RefreshToken string
	UserID       uuid.UUID
	Email        string
	Role         string
	CustomerID   *uuid.UUID
}

// Customer create use case currently uses a generic name.
type Input struct {
	Name  string
	CPF   string
	Email string
}
```

### 3. Contracts (Interfaces)

Abstractions that define application dependencies.

Includes:

* repositories
* external services

Examples:

* `AccountRepository`
* `AuthService`

Concrete contracts used by application use cases:

```go
type AccountRepository interface {
	CreateTransaction(ctx context.Context, tx *Transaction) error
	GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*Transaction, error)
	GetTransactionByReference(ctx context.Context, accountID uuid.UUID, referenceID uuid.UUID, typeName TransactionType) (*Transaction, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Account, error)
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*Account, error)
	IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error)
	DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error)
	WithTransaction(ctx context.Context, fn func(tx Tx) error) error
}

type CustomerRepository interface {
	Create(ctx context.Context, c *Customer) error
	GetByID(ctx context.Context, id uuid.UUID) (*Customer, error)
}

// External service contract used by auth application use cases.
type TokenService interface {
	GenerateAccessToken(claims TokenClaims) (string, error)
	GenerateRefreshToken(userID uuid.UUID) (string, error)
	ParseAccessToken(token string) (*TokenClaims, error)
	ParseRefreshToken(token string) (uuid.UUID, error)
}
```

### 4. Rules Applied in the Application

Technical and operational rules that support use case execution, without redefining domain invariants.

Includes:

* additional validations (beyond the domain)
* normalizations
* data adaptations

Typical patterns in current use cases:

```go
// transfer: input-level checks before domain calls
if input.FromAccountID == uuid.Nil || input.ToAccountID == uuid.Nil {
	return nil, domain.ErrInvalidData
}
if input.Amount <= 0 {
	return nil, domain.ErrInvalidAmount
}
if input.FromAccountID == input.ToAccountID {
	return nil, domain.ErrSameAccountTransfer
}

// login/register: normalization and basic validation
email := normalizeEmail(input.Email) // strings.ToLower(strings.TrimSpace(...))
if email == "" {
	return nil, domain.ErrInvalidEmail
}

// create account: enforce active user before account creation
if user.Status != authdomain.UserStatusActive {
	return nil, domain.ErrForbidden
}
```

### 5. Mappings

Explicit transformations between domain, persistence, and transport representations.

Includes:

* DB -> Entity
* Entity -> DTO
* DTO -> Response

Concrete mapping examples:

```go
// DB -> Entity (infrastructure)
var t domain.Transaction
err := row.Scan(
	&t.ID,
	&t.AccountID,
	&t.Type,
	&t.Amount,
	&t.BalanceAfter,
	&t.ReferenceID,
	&t.RelatedAccountID,
	&t.IdempotencyKey,
	&t.CreatedAt,
)

// DTO -> use case input (delivery -> application)
result, err := h.transfer.Execute(ctx, application.TransferInput{
	User:           user,
	FromAccountID:  fromAccountID,
	ToAccountID:    toAccountID,
	Amount:         req.Amount,
	IdempotencyKey: req.IdempotencyKey,
})

// Entity/output -> HTTP response DTO (delivery)
sharedhttp.WriteJSON(w, http.StatusOK, TransferData{
	FromAccountID: result.FromAccountID.String(),
	ToAccountID:   result.ToAccountID.String(),
	Amount:        result.Amount,
	FromBalance:   result.FromBalance,
	ToBalance:     result.ToBalance,
})
```

---

## What Is NOT the Responsibility of the Application Model

The following does not belong to the Application Model:

* core business rules (**Domain**)
* use case execution flow (**Use Cases**)
* SQL, schema, and tables (**Database**)
* HTTP protocols and JSON serialization (**API/Delivery**)

---

## Precise Definition

> The Application Model defines **the concrete data representations used by the application layer to execute use cases**, bridging domain concepts, persistence structures, and transport formats.

---

## Expected Outcome

With this document well defined, it should be possible to:

* quickly understand how the code is structured
* map domain <-> implementation
* reduce ambiguity between layers

