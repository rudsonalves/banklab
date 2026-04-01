package domain

import (
	"context"

	"github.com/google/uuid"
)

type AccountRepository interface {
	Create(ctx context.Context, account *Account) error
	ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error)
	NextAccountNumber(ctx context.Context) (string, error)
}
