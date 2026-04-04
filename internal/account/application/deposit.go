package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

type Deposit struct {
	accountRepo domain.AccountRepository
}

func NewDeposit(accountRepo domain.AccountRepository) *Deposit {
	return &Deposit{accountRepo: accountRepo}
}

type DepositInput struct {
	User      *authdomain.AuthenticatedUser
	AccountID uuid.UUID
	Amount    int64
}

func (uc *Deposit) Execute(ctx context.Context, input DepositInput) (_ *domain.Account, err error) {
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

		if !authdomain.CanAccessAccount(input.User, account) {
			return domain.ErrForbidden
		}

		if err := account.CanDeposit(input.Amount); err != nil {
			return err
		}

		updatedBalance, err := tx.UpdateBalance(ctx, input.AccountID, input.Amount)
		if err != nil {
			if errors.Is(err, domain.ErrAccountNotFound) {
				return err
			}
			return fmt.Errorf("update balance: %w", err)
		}

		ledgerTx := domain.NewTransaction(
			input.AccountID,
			domain.TransactionDeposit,
			input.Amount,
			updatedBalance,
			nil,
		)
		if err := tx.CreateTransaction(ctx, ledgerTx); err != nil {
			return fmt.Errorf("create deposit ledger transaction: %w", err)
		}

		account.Balance = updatedBalance
		return nil
	})
	if err != nil {
		return nil, err
	}

	return account, nil
}
