package domain

import (
	"context"

	"github.com/google/uuid"
)

// CustomerRepository defines the canonical customer persistence contract.
type CustomerRepository interface {
	Create(ctx context.Context, c *Customer) error
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Customer, string, error)
}
