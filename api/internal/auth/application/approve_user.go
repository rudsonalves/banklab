package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	accountapplication "github.com/seu-usuario/bank-api/internal/account/application"
	accountdomain "github.com/seu-usuario/bank-api/internal/account/domain"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

type ApproveUserUseCase struct {
	userRepo     domain.UserRepository
	accountRepo  accountdomain.AccountRepository
	customerRepo accountdomain.CustomerRepository
	transactor   domain.Transactor
}

func NewApproveUserUseCase(
	userRepo domain.UserRepository,
	accountRepo accountdomain.AccountRepository,
	customerRepo accountdomain.CustomerRepository,
	transactor domain.Transactor,
) *ApproveUserUseCase {
	return &ApproveUserUseCase{
		userRepo:     userRepo,
		accountRepo:  accountRepo,
		customerRepo: customerRepo,
		transactor:   transactor,
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

func (uc *ApproveUserUseCase) Execute(ctx context.Context, input ApproveUserInput) (*ApproveUserOutput, error) {
	var output *ApproveUserOutput

	err := uc.transactor.RunInTx(ctx, func(txCtx context.Context) error {
		user, err := uc.userRepo.FindByIDForUpdate(txCtx, input.UserID)
		if err != nil {
			return fmt.Errorf("load user: %w", err)
		}
		if user == nil {
			return domain.ErrUserNotFound
		}
		if user.Status != domain.UserStatusPending {
			return domain.ErrUserAlreadyActive
		}

		if err := uc.userRepo.UpdateStatus(txCtx, user.ID, domain.UserStatusActive); err != nil {
			return fmt.Errorf("update user status: %w", err)
		}
		user.Status = domain.UserStatusActive

		if user.CustomerID == nil {
			return domain.ErrInvalidUserState
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

		account, err := accountdomain.NewAccount(*user.CustomerID, number, accountapplication.GenerateBranch())
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
