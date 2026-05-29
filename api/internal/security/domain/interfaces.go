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

type StepUpTokenRepository interface {
	Create(ctx context.Context, token *StepUpToken) error
	FindByJTI(ctx context.Context, jti string) (*StepUpToken, error)
	ConsumeByJTI(ctx context.Context, jti string, now time.Time) (*StepUpToken, error)
}

type StepUpTokenSigner interface {
	Sign(token *StepUpToken) (string, error)
}

type StepUpEndpointPolicy interface {
	Validate(endpointKey string) error
}
