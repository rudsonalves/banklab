package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID
	Name      string
	CPF       string
	CreatedAt time.Time
}

// NewCustomer creates a new Customer entity with the provided name and CPF. It
// validates the input parameters and returns an error if any required field is
// missing or invalid. If the input is valid, it generates a new UUID for the
// customer and sets the CreatedAt timestamp to the current time in UTC.
func NewCustomer(name, cpf string) (*Customer, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrNameRequired
	}

	normalizedCPF := normalizeCPF(cpf)
	if normalizedCPF == "" {
		return nil, ErrCPFRequired
	}

	if !ValidateCPF(normalizedCPF) {
		return nil, ErrCPFInvalid
	}

	return &Customer{
		ID:        uuid.New(),
		Name:      name,
		CPF:       normalizedCPF,
		CreatedAt: time.Now().UTC(),
	}, nil
}
