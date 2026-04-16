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
	User *authdomain.AuthenticatedUser
}

func (uc *CreateAccount) Execute(ctx context.Context, input CreateAccountInput) (*domain.Account, error) {
	if input.User == nil || input.User.CustomerID == nil {
		return nil, domain.ErrForbidden
	}

	customerID := *input.User.CustomerID
	if customerID == uuid.Nil {
		return nil, domain.ErrForbidden
	}

	if !CanAccessCustomer(input.User, customerID) {
		return nil, domain.ErrForbidden
	}

	exists, err := uc.customerRepo.Exists(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("check customer existence: %w", err)
	}
	if !exists {
		return nil, domain.ErrCustomerNotFound
	}

	// Optional business rule: one account per customer
	// Uncomment if needed
	/*
			exists, err = uc.accountRepo.ExistsByCustomerID(ctx, customerID)
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
	branch := GenerateBranch()

	account, err := domain.NewAccount(customerID, number, branch)
	if err != nil {
		return nil, err
	}

	if err := uc.accountRepo.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}

	return account, nil
}

func GenerateBranch() string {
	return "0001"
}
