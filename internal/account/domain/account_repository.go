package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *Transaction) error
}

type AccountRepository interface {
	TransactionRepository

	Create(ctx context.Context, account *Account) error
	ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error)
	NextAccountNumber(ctx context.Context) (string, error)

	GetByID(ctx context.Context, id uuid.UUID) (*Account, error)
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*Account, error)
	GetTransactions(
		ctx context.Context,
		accountID uuid.UUID,
		limit int,
		cursorTime *time.Time,
		cursorID *uuid.UUID,
		from *time.Time,
		to *time.Time,
	) ([]Transaction, error)
	UpdateBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error)
	// DecreaseBalance performs an atomic balance decrement.
	// It returns ErrAccountNotFound when the account does not exist and
	// ErrInsufficientBalance when the account exists but has insufficient funds.
	DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) error

	BeginTx(ctx context.Context) (Tx, error)
	WithTransaction(ctx context.Context, fn func(tx Tx) error) error
}

type Tx interface {
	AccountRepository
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
