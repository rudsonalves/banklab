package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	accountapplication "github.com/seu-usuario/bank-api/internal/account/application/account"
	accountdomain "github.com/seu-usuario/bank-api/internal/account/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	customerdomain "github.com/seu-usuario/bank-api/internal/customer/domain"
)

type ApproveUserUseCase struct {
	userRepo     authdomain.UserRepository
	accountRepo  accountdomain.AccountRepository
	customerRepo customerdomain.CustomerRepository
	transactor   authdomain.Transactor
	branchPolicy accountapplication.BranchPolicy
}

// NewApproveUserUseCase creates a new instance of the ApproveUserUseCase with the
// provided dependencies.
func NewApproveUserUseCase(
	userRepo authdomain.UserRepository,
	accountRepo accountdomain.AccountRepository,
	customerRepo customerdomain.CustomerRepository,
	transactor authdomain.Transactor,
	branchPolicy accountapplication.BranchPolicy,
) *ApproveUserUseCase {
	return &ApproveUserUseCase{
		userRepo:     userRepo,
		accountRepo:  accountRepo,
		customerRepo: customerRepo,
		transactor:   transactor,
		branchPolicy: branchPolicy,
	}
}

type ApproveUserInput struct {
	UserID uuid.UUID
}

type ApproveUserOutput struct {
	UserID    uuid.UUID
	Status    string
	AccountID uuid.UUID
}

// Execute approves a pending user by updating their status to active and creating
// a new bank account for them. It performs necessary validations, checks the
// existence of the associated customer, generates an account number, and
// determines the branch using the branch policy. The entire operation is
// executed within a transaction to ensure data consistency.
func (uc *ApproveUserUseCase) Execute(ctx context.Context, input ApproveUserInput) (*ApproveUserOutput, error) {
	var output *ApproveUserOutput

	err := uc.transactor.RunInTx(ctx, func(txCtx context.Context) error {
		user, err := uc.userRepo.FindByIDForUpdate(txCtx, input.UserID)
		if err != nil {
			return fmt.Errorf("load user: %w", err)
		}
		if user == nil {
			return authdomain.ErrUserNotFound
		}
		if user.Status != authdomain.UserStatusPending {
			return authdomain.ErrUserAlreadyActive
		}

		if err := uc.userRepo.UpdateStatus(txCtx, user.ID, authdomain.UserStatusActive); err != nil {
			return fmt.Errorf("update user status: %w", err)
		}
		user.Status = authdomain.UserStatusActive

		if user.CustomerID == nil {
			return authdomain.ErrInvalidUserState
		}

		exists, err := uc.customerRepo.Exists(txCtx, *user.CustomerID)
		if err != nil {
			return fmt.Errorf("check customer existence: %w", err)
		}
		if !exists {
			return accountdomain.ErrCustomerNotFound
		}

		number, err := uc.accountRepo.NextAccountNumber(txCtx)
		if err != nil {
			return fmt.Errorf("generate account number: %w", err)
		}
		if uc.branchPolicy == nil {
			return fmt.Errorf("branch policy not configured")
		}

		account, err := accountdomain.NewAccount(*user.CustomerID, number, uc.branchPolicy.Branch())
		if err != nil {
			return fmt.Errorf("create account: %w", err)
		}

		if err := uc.accountRepo.Create(txCtx, account); err != nil {
			return fmt.Errorf("persist account: %w", err)
		}

		output = &ApproveUserOutput{
			UserID:    user.ID,
			Status:    string(user.Status),
			AccountID: account.ID,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return output, nil
}
