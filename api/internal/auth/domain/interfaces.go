package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	// Create persists the full User entity, including optional CustomerID.
	Create(ctx context.Context, user *User) error
	// UpdateStatus updates only the user's status. This is used
	UpdateStatus(ctx context.Context, userID uuid.UUID, status UserStatus) error
	// FindByEmail returns the full User entity, including optional CustomerID.
	FindByEmail(ctx context.Context, email string) (*User, error)
	// FindByID returns the full User entity, including optional CustomerID.
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	// ExistsByEmail checks if a user with the given email already exists.
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash string, password string) error
}

type TokenClaims struct {
	UserID     uuid.UUID
	Role       Role
	CustomerID *uuid.UUID
}

type TokenService interface {
	GenerateAccessToken(claims TokenClaims) (string, error)
	GenerateRefreshToken(userID uuid.UUID) (string, error)

	ParseAccessToken(token string) (*TokenClaims, error)
	ParseRefreshToken(token string) (uuid.UUID, error)
}

type SessionRepository interface {
	Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	FindByTokenHash(ctx context.Context, tokenHash string) (userID uuid.UUID, expiresAt time.Time, revoked bool, err error)
	Revoke(ctx context.Context, tokenHash string) error
}

type Transactor interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}
