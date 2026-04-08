package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCustomer Role = "customer"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Role         Role
	CustomerID   *uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleCustomer:
		return true
	default:
		return false
	}
}

// NewUser constructs a User and enforces domain invariants.
// For RoleCustomer, customerID must be non-nil.
func NewUser(id uuid.UUID, email, passwordHash string, role Role, customerID *uuid.UUID, now time.Time) (*User, error) {
	if role == RoleCustomer && customerID == nil {
		return nil, ErrInvalidUserState
	}
	return &User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		CustomerID:   customerID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}
