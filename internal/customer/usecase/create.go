package usecase

import (
	"context"

	"github.com/seu-usuario/bank-api/internal/customer/domain"
	"github.com/seu-usuario/bank-api/internal/customer/repository"
)

type CreateCustomer struct {
	repo *repository.Repository
}

func NewCreateCustomer(repo *repository.Repository) *CreateCustomer {
	return &CreateCustomer{repo: repo}
}

type Input struct {
	Name  string
	CPF   string
	Email string
}

func (uc *CreateCustomer) Execute(ctx context.Context, input Input) error {
	customer, err := domain.NewCustomer(
		input.Name,
		input.CPF,
		input.Email,
	)
	if err != nil {
		return err
	}

	return uc.repo.Create(ctx, customer)
}
