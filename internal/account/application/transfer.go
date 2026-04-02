package application

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

type Transfer struct {
	accountRepo domain.AccountRepository
}

func NewTransfer(accountRepo domain.AccountRepository) *Transfer {
	return &Transfer{accountRepo: accountRepo}
}

type TransferInput struct {
	User          *authdomain.AuthenticatedUser
	FromAccountID uuid.UUID
	ToAccountID   uuid.UUID
	Amount        int64
}

type TransferResult struct {
	FromAccountID uuid.UUID
	ToAccountID   uuid.UUID
	Amount        int64
	FromBalance   int64
	ToBalance     int64
}

func (uc *Transfer) Execute(ctx context.Context, input TransferInput) (_ *TransferResult, err error) {
	if input.FromAccountID == uuid.Nil || input.ToAccountID == uuid.Nil {
		return nil, domain.ErrInvalidData
	}

	if input.Amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	if input.FromAccountID == input.ToAccountID {
		return nil, domain.ErrSameAccountTransfer
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

	// Lock both accounts in deterministic UUID order to reduce deadlock risk
	// when concurrent transfers touch the same rows in opposite directions.
	firstID, secondID := orderedUUIDs(input.FromAccountID, input.ToAccountID)
	firstAccount, err := tx.GetByIDForUpdate(ctx, firstID)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("get first account for update: %w", err)
	}

	secondAccount, err := tx.GetByIDForUpdate(ctx, secondID)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("get second account for update: %w", err)
	}

	fromAccount, toAccount := mapTransferAccounts(input.FromAccountID, firstAccount, secondAccount)

	if !authdomain.CanAccessAccount(input.User, fromAccount) {
		return nil, domain.ErrForbidden
	}

	if err := fromAccount.CanTransfer(input.Amount, input.ToAccountID); err != nil {
		return nil, err
	}

	if err := toAccount.CanDeposit(input.Amount); err != nil {
		return nil, err
	}

	if err := tx.DecreaseBalance(ctx, input.FromAccountID, input.Amount); err != nil {
		return nil, fmt.Errorf("decrease source balance: %w", err)
	}
	fromAccount.Balance -= input.Amount

	updatedToBalance, err := tx.UpdateBalance(ctx, input.ToAccountID, input.Amount)
	if err != nil {
		return nil, fmt.Errorf("increase destination balance: %w", err)
	}
	toAccount.Balance = updatedToBalance

	referenceID := uuid.New()

	outgoing := domain.NewTransaction(
		input.FromAccountID,
		domain.TransactionTransferOut,
		input.Amount,
		fromAccount.Balance,
		&referenceID,
	)
	if err := tx.CreateTransaction(ctx, outgoing); err != nil {
		return nil, fmt.Errorf("create transfer out ledger transaction: %w", err)
	}

	incoming := domain.NewTransaction(
		input.ToAccountID,
		domain.TransactionTransferIn,
		input.Amount,
		toAccount.Balance,
		&referenceID,
	)
	if err := tx.CreateTransaction(ctx, incoming); err != nil {
		return nil, fmt.Errorf("create transfer in ledger transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	committed = true

	return &TransferResult{
		FromAccountID: input.FromAccountID,
		ToAccountID:   input.ToAccountID,
		Amount:        input.Amount,
		FromBalance:   fromAccount.Balance,
		ToBalance:     toAccount.Balance,
	}, nil
}

func orderedUUIDs(left, right uuid.UUID) (uuid.UUID, uuid.UUID) {
	if bytes.Compare(left[:], right[:]) <= 0 {
		return left, right
	}

	return right, left
}

func mapTransferAccounts(fromID uuid.UUID, first, second *domain.Account) (*domain.Account, *domain.Account) {
	if first.ID == fromID {
		return first, second
	}

	return second, first
}
