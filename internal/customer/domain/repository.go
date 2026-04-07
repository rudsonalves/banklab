package domain

import (
	"context"

	"github.com/google/uuid"
)

// CustomerRepository defines the canonical customer persistence contract.
type CustomerRepository interface {
	Create(ctx context.Context, c *Customer) error
	GetByID(ctx context.Context, id uuid.UUID) (*Customer, error)
}

// Repository is kept for backward compatibility with existing flows that
// only need creation and existence checks.
type Repository interface {
	Create(ctx context.Context, c *Customer) error
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}
