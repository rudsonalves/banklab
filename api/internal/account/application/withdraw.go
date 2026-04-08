package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

type Withdraw struct {
	accountRepo domain.AccountRepository
}

func NewWithdraw(accountRepo domain.AccountRepository) *Withdraw {
	return &Withdraw{accountRepo: accountRepo}
}

type WithdrawInput struct {
	User      *authdomain.AuthenticatedUser
	AccountID uuid.UUID
	Amount    int64
}

func (uc *Withdraw) Execute(ctx context.Context, input WithdrawInput) (_ *domain.Account, err error) {
	if input.AccountID == uuid.Nil {
		return nil, domain.ErrInvalidData
	}

	if input.Amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	var account *domain.Account
	err = uc.accountRepo.WithTransaction(ctx, func(tx domain.Tx) error {
		account, err = tx.GetByIDForUpdate(ctx, input.AccountID)
		if err != nil {
			if errors.Is(err, domain.ErrAccountNotFound) {
				return err
			}
			return fmt.Errorf("get account by id: %w", err)
		}

		if !CanAccessAccount(input.User, account) {
			return domain.ErrForbidden
		}

		if err := account.CanWithdraw(input.Amount); err != nil {
			return err
		}

		updatedBalance, err := tx.DecreaseBalance(ctx, input.AccountID, input.Amount)
		if err != nil {
			if errors.Is(err, domain.ErrAccountNotFound) || errors.Is(err, domain.ErrInsufficientBalance) {
				return err
			}
			return fmt.Errorf("decrease balance: %w", err)
		}

		account.Balance = updatedBalance

		ledgerTx := domain.NewTransaction(
			input.AccountID,
			domain.TransactionWithdraw,
			input.Amount,
			account.Balance,
			nil,
		)
		if err := tx.CreateTransaction(ctx, ledgerTx); err != nil {
			return fmt.Errorf("create withdraw ledger transaction: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return account, nil
}
