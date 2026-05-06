package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, account *Account) error
	ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]Account, error)
	ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error)
	NextAccountNumber(ctx context.Context) (string, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Account, error)
}

type AccountRepository = Repository
