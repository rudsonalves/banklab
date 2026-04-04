package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

type CreateAccount struct {
	accountRepo  domain.AccountRepository
	customerRepo domain.CustomerRepository
}

func NewCreateAccount(
	accountRepo domain.AccountRepository,
	customerRepo domain.CustomerRepository,
) *CreateAccount {
	return &CreateAccount{
		accountRepo:  accountRepo,
		customerRepo: customerRepo,
	}
}

type CreateAccountInput struct {
	User       *authdomain.AuthenticatedUser
	CustomerID uuid.UUID
}

func (uc *CreateAccount) Execute(ctx context.Context, input CreateAccountInput) (*domain.Account, error) {
	if input.CustomerID == uuid.Nil {
		return nil, domain.ErrInvalidData
	}

	if input.User == nil {
		return nil, domain.ErrForbidden
	}

	if input.User.Role != authdomain.RoleAdmin {
		if input.User.CustomerID == nil || *input.User.CustomerID != input.CustomerID {
			return nil, domain.ErrForbidden
		}
	}

	exists, err := uc.customerRepo.Exists(ctx, input.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("check customer existence: %w", err)
	}
	if !exists {
		return nil, domain.ErrCustomerNotFound
	}

	// Optional business rule: one account per customer
	// Uncomment if needed
	/*
		exists, err = uc.accountRepo.ExistsByCustomerID(ctx, input.CustomerID)
		if err != nil {
			return nil, fmt.Errorf("check account existence: %w", err)
		}
		if exists {
			return nil, domain.ErrAccountAlreadyExists
		}
	*/

	number, err := uc.accountRepo.NextAccountNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("generate account number: %w", err)
	}
	branch := generateBranch()

	account, err := domain.NewAccount(input.CustomerID, number, branch)
	if err != nil {
		return nil, err
	}

	if err := uc.accountRepo.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}

	return account, nil
}

func generateBranch() string {
	return "0001"
}
