package domain

import "github.com/google/uuid"

type AuthenticatedUser struct {
	UserID     string
	Role       Role
	CustomerID *uuid.UUID
}
