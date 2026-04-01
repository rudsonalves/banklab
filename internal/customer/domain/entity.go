package domain

import "github.com/google/uuid"

type Customer struct {
	ID    uuid.UUID
	Name  string
	CPF   string
	Email string
}

func NewCustomer(name, cpf, email string) (*Customer, error) {
	if name == "" {
		return nil, ErrNameRequired
	}

	if cpf == "" {
		return nil, ErrCPFRequired
	}

	if email == "" {
		return nil, ErrEmailRequired
	}

	return &Customer{
		ID:    uuid.New(),
		Name:  name,
		CPF:   cpf,
		Email: email,
	}, nil
}
