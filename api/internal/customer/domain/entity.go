package domain

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID
	Name      string
	CPF       string
	CreatedAt time.Time
}

func NewCustomer(name, cpf string) (*Customer, error) {
	if name == "" {
		return nil, ErrNameRequired
	}

	if cpf == "" {
		return nil, ErrCPFRequired
	}

	return &Customer{
		ID:        uuid.New(),
		Name:      name,
		CPF:       cpf,
		CreatedAt: time.Now().UTC(),
	}, nil
}
