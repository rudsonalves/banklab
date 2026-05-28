package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TransactionPasswordRepository interface {
	Create(ctx context.Context, password *TransactionPassword) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*TransactionPassword, error)
	SaveValidationState(ctx context.Context, password *TransactionPassword) error

	UpdatePasswordHash(
		ctx context.Context,
		id uuid.UUID,
		passwordHash string,
		changedAt time.Time,
		updatedAt time.Time,
	) error
}

type TransactionPasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}
