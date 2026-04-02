package domain

import (
	"context"

	"github.com/google/uuid"
)

type AccountRepository interface {
	Create(ctx context.Context, account *Account) error
	ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error)
	NextAccountNumber(ctx context.Context) (string, error)

	GetByID(ctx context.Context, id uuid.UUID) (*Account, error)
	UpdateBalance(ctx context.Context, id uuid.UUID, amount int64) error

	BeginTx(ctx context.Context) (Tx, error)
}

type Tx interface {
	AccountRepository
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
