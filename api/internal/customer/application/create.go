package application

import (
	"context"

	"github.com/seu-usuario/bank-api/internal/customer/domain"
)

type CreateCustomer struct {
	repo domain.Repository
}

func NewCreateCustomer(repo domain.Repository) *CreateCustomer {
	return &CreateCustomer{repo: repo}
}

type Input struct {
	Name string
	CPF  string
}

func (uc *CreateCustomer) Execute(ctx context.Context, input Input) (*domain.Customer, error) {
	customer, err := domain.NewCustomer(
		input.Name,
		input.CPF,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}
