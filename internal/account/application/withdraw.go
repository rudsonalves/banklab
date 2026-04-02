package application

import (
	"context"
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

	account, err := tx.GetByIDForUpdate(ctx, input.AccountID)
	if err != nil {
		return nil, fmt.Errorf("get account by id: %w", err)
	}

	if !authdomain.CanAccessAccount(input.User, account) {
		return nil, domain.ErrForbidden
	}

	if err := account.CanWithdraw(input.Amount); err != nil {
		return nil, err
	}

	if err := tx.DecreaseBalance(ctx, input.AccountID, input.Amount); err != nil {
		return nil, fmt.Errorf("decrease balance: %w", err)
	}

	account.Balance -= input.Amount

	ledgerTx := domain.NewTransaction(
		input.AccountID,
		domain.TransactionWithdraw,
		input.Amount,
		account.Balance,
		nil,
	)
	if err := tx.CreateTransaction(ctx, ledgerTx); err != nil {
		return nil, fmt.Errorf("create withdraw ledger transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	committed = true

	return account, nil
}
