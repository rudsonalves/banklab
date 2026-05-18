package application

import (
	"context"
	"time"

	"github.com/seu-usuario/bank-api/internal/customer/domain"
)

type CreateCustomer struct {
	repo domain.CustomerRepository
}

// NewCreateCustomer creates a new instance of CreateCustomer with the provided
// customer repository.
func NewCreateCustomer(repo domain.CustomerRepository) *CreateCustomer {
	return &CreateCustomer{repo: repo}
}

type Input struct {
	Name      string
	BirthDate time.Time
}

// Execute creates a new customer using the provided input data. It validates the
// input, creates a new customer entity, and persists it using the customer
// repository.
func (uc *CreateCustomer) Execute(ctx context.Context, input Input) (*domain.Customer, error) {
	customer, err := domain.NewCustomer(
		input.Name,
		input.BirthDate,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}
