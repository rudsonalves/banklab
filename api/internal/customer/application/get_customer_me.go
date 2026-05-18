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

// Execute retrieves the customer profile for the authenticated user by ID.
// It fetches the customer data from the repository and returns the customer
// profile.
// Returns an error if the customer is not found or if there is an error during
// retrieval.
func (uc *GetCustomerMe) Execute(
	ctx context.Context,
	input GetCustomerMeInput,
) (*domain.CustomerProfile, error) {
	if input.CustomerID == uuid.Nil {
		return nil, domain.ErrInvalidData
	}

	profile, err := uc.repo.GetByID(ctx, input.CustomerID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, domain.ErrNotFound
	}

	return profile, nil
}
