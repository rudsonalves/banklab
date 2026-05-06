package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/bankaccount/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	customerdomain "github.com/seu-usuario/bank-api/internal/customer/domain"
)

type CreateAccount struct {
	accountRepo  domain.AccountRepository
	customerRepo customerdomain.CustomerRepository
	userRepo     authdomain.UserRepository
	branchPolicy BranchPolicy
}

// NewCreateAccount creates a new instance of the CreateAccount use case with the
// provided dependencies.
func NewCreateAccount(
	accountRepo domain.AccountRepository,
	customerRepo customerdomain.CustomerRepository,
	userRepo authdomain.UserRepository,
	branchPolicy BranchPolicy,
) *CreateAccount {
	return &CreateAccount{
		accountRepo:  accountRepo,
		customerRepo: customerRepo,
		userRepo:     userRepo,
		branchPolicy: branchPolicy,
	}
}

type CreateAccountInput struct {
	User *authdomain.AuthenticatedUser
}

// Execute creates a new bank account for the authenticated user. It performs
// necessary validations, checks access permissions, and ensures that the associated customer exists before creating the account. The account number is
// generated using the account repository, and the branch is determined by the
// branch policy.
func (uc *CreateAccount) Execute(ctx context.Context, input CreateAccountInput) (*domain.Account, error) {
	if input.User == nil || input.User.CustomerID == nil {
		return nil, domain.ErrForbidden
	}
	if uc.userRepo == nil {
		return nil, fmt.Errorf("user repository not configured")
	}
	if input.User.UserID == uuid.Nil {
		return nil, domain.ErrForbidden
	}

	user, err := uc.userRepo.FindByID(ctx, input.User.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	if user == nil {
		return nil, domain.ErrForbidden
	}
	if user.Status != authdomain.UserStatusActive {
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

	if uc.branchPolicy == nil {
		return nil, fmt.Errorf("branch policy not configured")
	}
	branch := uc.branchPolicy.Branch()

	account, err := domain.NewAccount(customerID, number, branch)
	if err != nil {
		return nil, err
	}

	if err := uc.accountRepo.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}

	return account, nil
}
