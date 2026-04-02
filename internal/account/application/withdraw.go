package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/domain"
)

type Withdraw struct {
	accountRepo domain.AccountRepository
}

func NewWithdraw(accountRepo domain.AccountRepository) *Withdraw {
	return &Withdraw{accountRepo: accountRepo}
}

type WithdrawInput struct {
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

	if err := account.CanWithdraw(input.Amount); err != nil {
		return nil, err
	}

	if err := tx.DecreaseBalance(ctx, input.AccountID, input.Amount); err != nil {
		return nil, fmt.Errorf("decrease balance: %w", err)
	}

	account.Balance -= input.Amount

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	committed = true

	return account, nil
}
