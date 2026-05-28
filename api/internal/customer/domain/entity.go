package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID
	Name      string
	BirthDate time.Time
	CreatedAt time.Time
}

// NewCustomer creates a new Customer entity with the provided name and birth date.
// It validates the input parameters and returns an error if any validation fails.
func NewCustomer(name string, birthDate time.Time) (*Customer, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrNameRequired
	}

	if birthDate.IsZero() {
		return nil, ErrBirthDateRequired
	}

	return &Customer{
		Name:      name,
		BirthDate: birthDate,
		CreatedAt: time.Now().UTC(),
	}, nil
}

type CustomerProfile struct {
	Customer Customer
	Email    string
	CPF      string
}
