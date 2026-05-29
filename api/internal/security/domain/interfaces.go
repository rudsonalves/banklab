package domain

import (
	"context"
	"strings"
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

type StepUpTokenVerifier interface {
	Verify(token string) (*VerifiedStepUpTokenClaims, error)
}

type VerifiedStepUpTokenClaims struct {
	UserID      uuid.UUID
	EndpointKey string
	Scope       string
	JTI         string
	ExpiresAt   time.Time
	IssuedAt    time.Time
}

func NewVerifiedStepUpTokenClaims(
	userID uuid.UUID,
	endpointKey string,
	scope string,
	jti string,
	expiresAt time.Time,
	issuedAt time.Time,
) (*VerifiedStepUpTokenClaims, error) {
	claims := &VerifiedStepUpTokenClaims{
		UserID:      userID,
		EndpointKey: strings.TrimSpace(endpointKey),
		Scope:       strings.TrimSpace(scope),
		JTI:         strings.TrimSpace(jti),
		ExpiresAt:   expiresAt.UTC(),
		IssuedAt:    issuedAt.UTC(),
	}

	if err := claims.Validate(); err != nil {
		return nil, err
	}

	return claims, nil
}

func (c *VerifiedStepUpTokenClaims) Validate() error {
	if c == nil {
		return ErrInvalidStepUpToken
	}

	if c.UserID == uuid.Nil ||
		c.EndpointKey == "" ||
		c.Scope != StepUpTokenScope ||
		c.JTI == "" ||
		c.ExpiresAt.IsZero() ||
		c.IssuedAt.IsZero() ||
		!c.ExpiresAt.After(c.IssuedAt) {
		return ErrInvalidStepUpToken
	}

	return nil
}

type StepUpEndpointPolicy interface {
	Validate(endpointKey string) error
}
