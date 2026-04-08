package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/customer/domain"
)

type GetCustomerMe struct {
	repo domain.CustomerRepository
}

func NewGetCustomerMe(repo domain.CustomerRepository) *GetCustomerMe {
	return &GetCustomerMe{repo: repo}
}

type GetCustomerMeInput struct {
	CustomerID uuid.UUID
}

func (uc *GetCustomerMe) Execute(ctx context.Context, input GetCustomerMeInput) (*domain.Customer, error) {
	if input.CustomerID == uuid.Nil {
		return nil, domain.ErrInvalidData
	}

	customer, err := uc.repo.GetByID(ctx, input.CustomerID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, domain.ErrNotFound
	}

	return customer, nil
}
