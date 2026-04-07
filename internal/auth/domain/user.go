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
