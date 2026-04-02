package domain

import "time"

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCustomer Role = "customer"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         Role
	CustomerID   *string
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
