package domain

import (
	"context"

	"github.com/google/uuid"
)

type CustomerRepository interface {
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}
