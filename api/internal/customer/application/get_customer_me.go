package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/customer/domain"
)

type GetCustomerMe struct {
	repo domain.CustomerRepository
}

// NewGetCustomerMe creates a new instance of GetCustomerMe with the provided
// customer repository.
func NewGetCustomerMe(repo domain.CustomerRepository) *GetCustomerMe {
	return &GetCustomerMe{repo: repo}
}

type GetCustomerMeInput struct {
	CustomerID uuid.UUID
}

// Execute retrieves the customer information for the authenticated user based on
// the provided input. It validates the input, fetches the customer data from
// the repository, and returns the customer information along with the associated
// email. If the customer is not found or if there is an error during retrieval,
// it returns an appropriate error.
func (uc *GetCustomerMe) Execute(ctx context.Context, input GetCustomerMeInput) (*domain.Customer, string, error) {
	if input.CustomerID == uuid.Nil {
		return nil, "", domain.ErrInvalidData
	}

	customer, email, err := uc.repo.GetByID(ctx, input.CustomerID)
	if err != nil {
		return nil, "", err
	}
	if customer == nil {
		return nil, "", domain.ErrNotFound
	}

	return customer, email, nil
}
