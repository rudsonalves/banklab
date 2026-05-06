package domain

import (
	"context"

	"github.com/google/uuid"
)

type TransactionRepository interface {
	// CreateTransaction creates a new transaction in the repository.
	CreateTransaction(ctx context.Context, tx *Transaction) error
}

type AccountRepository interface {
	TransactionRepository

	// GetTransactionByIdempotencyKey returns the transaction with the specified idempotency key for the given account.
	// It returns ErrTransactionNotFound if no transaction exists for the given account and idempotency key.
	GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*Transaction, error)

	// GetTransactionByReference returns the transaction with the specified reference ID and type for the given account.
	// It returns ErrTransactionNotFound if no transaction exists for the given account, reference ID, and type.
	GetTransactionByReference(ctx context.Context, accountID uuid.UUID, referenceID uuid.UUID, typeName TransactionType) (*Transaction, error)

	// Create creates a new account in the repository.
	Create(ctx context.Context, account *Account) error
	// ListByCustomerID returns a list of accounts associated with the specified customer ID.
	ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]Account, error)
	// ExistsByCustomerID checks if there are any accounts associated with the specified customer ID.
	ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error)
	// NextAccountNumber returns the next available account number.
	NextAccountNumber(ctx context.Context) (string, error)

	// GetByID returns ErrAccountNotFound when no account exists for the given id.
	// Implementations must never return (nil, nil).
	GetByID(ctx context.Context, id uuid.UUID) (*Account, error)
	// GetByIDForUpdate returns the account with the specified ID and locks it for update.
	// It returns ErrAccountNotFound when no account exists for the given id.
	// Implementations must never return (nil, nil).
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*Account, error)
	// IncreaseBalance performs an atomic balance increment.
	// It returns ErrAccountNotFound when the account does not exist.
	IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error)
	// DecreaseBalance performs an atomic balance decrement.
	// It returns ErrAccountNotFound when the account does not exist and
	// ErrInsufficientBalance when the account exists but has insufficient funds.
	DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error)

	// BeginTx starts a new transaction and returns a Tx interface for performing operations within that transaction.
	BeginTx(ctx context.Context) (Tx, error)
	// WithTransaction executes the provided function within a transaction. It handles committing or rolling back the transaction based on whether the function returns an error.
	WithTransaction(ctx context.Context, fn func(tx Tx) error) error
}

type Tx interface {
	AccountRepository
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
