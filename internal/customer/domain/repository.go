package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, c *Customer) error
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}
