package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	// Create persists the full User entity, including optional CustomerID.
	Create(ctx context.Context, user *User) error
	// FindByEmail returns the full User entity, including optional CustomerID.
	FindByEmail(ctx context.Context, email string) (*User, error)
	// FindByID returns the full User entity, including optional CustomerID.
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
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
	ParseAccessToken(token string) (*TokenClaims, error)
}
