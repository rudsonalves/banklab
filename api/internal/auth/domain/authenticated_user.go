package domain

import "github.com/google/uuid"

type AuthenticatedUser struct {
	UserID     uuid.UUID
	Role       Role
	CustomerID *uuid.UUID
}
