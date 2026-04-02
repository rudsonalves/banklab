package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/domain"
)

type Deposit struct {
	accountRepo domain.AccountRepository
}

func NewDeposit(accountRepo domain.AccountRepository) *Deposit {
	return &Deposit{accountRepo: accountRepo}
}

type DepositInput struct {
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

	tx, err := uc.accountRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	account, err := tx.GetByID(ctx, input.AccountID)
	if err != nil {
		return nil, fmt.Errorf("get account by id: %w", err)
	}

	if err := account.CanDeposit(input.Amount); err != nil {
		return nil, err
	}

	updatedBalance, err := tx.UpdateBalance(ctx, input.AccountID, input.Amount)
	if err != nil {
		return nil, fmt.Errorf("update balance: %w", err)
	}

	account.Balance = updatedBalance

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	committed = true

	return account, nil
}
