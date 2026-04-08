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

var errTransferIdempotencyConflict = errors.New("transfer idempotency conflict")

type Transfer struct {
	accountRepo domain.AccountRepository
}

func NewTransfer(accountRepo domain.AccountRepository) *Transfer {
	return &Transfer{accountRepo: accountRepo}
}

type TransferInput struct {
	User           *authdomain.AuthenticatedUser
	FromAccountID  uuid.UUID
	ToAccountID    uuid.UUID
	Amount         int64
	IdempotencyKey string
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

	var result *TransferResult
	err = uc.accountRepo.WithTransaction(ctx, func(tx domain.Tx) error {
		if input.IdempotencyKey != "" {
			existing, err := tx.GetOperationByIdempotencyKey(ctx, input.FromAccountID, input.IdempotencyKey)
			if err != nil {
				return fmt.Errorf("get operation by idempotency key: %w", err)
			}

			if existing != nil {
				result, err = transferResultFromOperation(ctx, tx, input, existing)
				if err != nil {
					return err
				}
				return nil
			}
		}

		// Lock both accounts in deterministic UUID order to reduce deadlock risk
		// when concurrent transfers touch the same rows in opposite directions.
		firstID, secondID := orderedUUIDs(input.FromAccountID, input.ToAccountID)
		firstAccount, err := tx.GetByIDForUpdate(ctx, firstID)
		if err != nil {
			if errors.Is(err, domain.ErrAccountNotFound) {
				return err
			}
			return fmt.Errorf("get first account for update: %w", err)
		}

		secondAccount, err := tx.GetByIDForUpdate(ctx, secondID)
		if err != nil {
			if errors.Is(err, domain.ErrAccountNotFound) {
				return err
			}
			return fmt.Errorf("get second account for update: %w", err)
		}

		fromAccount, toAccount := mapTransferAccounts(input.FromAccountID, firstAccount, secondAccount)

		if !CanAccessAccount(input.User, fromAccount) {
			return domain.ErrForbidden
		}

		if err := fromAccount.CanTransfer(input.Amount, input.ToAccountID); err != nil {
			return err
		}

		if err := toAccount.CanDeposit(input.Amount); err != nil {
			return err
		}

		updatedFromBalance, err := tx.DecreaseBalance(ctx, input.FromAccountID, input.Amount)
		if err != nil {
			if errors.Is(err, domain.ErrAccountNotFound) || errors.Is(err, domain.ErrInsufficientBalance) {
				return err
			}
			return fmt.Errorf("decrease source balance: %w", err)
		}
		fromAccount.Balance = updatedFromBalance

		updatedToBalance, err := tx.IncreaseBalance(ctx, input.ToAccountID, input.Amount)
		if err != nil {
			if errors.Is(err, domain.ErrAccountNotFound) {
				return err
			}
			return fmt.Errorf("increase destination balance: %w", err)
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
			return fmt.Errorf("create transfer out ledger transaction: %w", err)
		}

		incoming := domain.NewTransaction(
			input.ToAccountID,
			domain.TransactionTransferIn,
			input.Amount,
			toAccount.Balance,
			&referenceID,
		)
		if err := tx.CreateTransaction(ctx, incoming); err != nil {
			return fmt.Errorf("create transfer in ledger transaction: %w", err)
		}

		if input.IdempotencyKey != "" {
			op := domain.NewOperation(
				input.FromAccountID,
				domain.TransactionTransferOut,
				input.Amount,
				&input.ToAccountID,
				&referenceID,
				input.IdempotencyKey,
			)

			if err := tx.CreateOperation(ctx, op); err != nil {
				if errors.Is(err, domain.ErrOperationAlreadyProcessed) {
					existing, getErr := tx.GetOperationByIdempotencyKey(ctx, input.FromAccountID, input.IdempotencyKey)
					if getErr != nil {
						return fmt.Errorf("reload operation by idempotency key: %w", getErr)
					}
					if existing == nil {
						return fmt.Errorf("reload operation by idempotency key: not found")
					}

					result, getErr = transferResultFromOperation(ctx, tx, input, existing)
					if getErr != nil {
						return getErr
					}

					// Force rollback of this duplicate execution while preserving the replay result.
					return errTransferIdempotencyConflict
				}

				return fmt.Errorf("create transfer operation: %w", err)
			}
		}

		result = &TransferResult{
			FromAccountID: input.FromAccountID,
			ToAccountID:   input.ToAccountID,
			Amount:        input.Amount,
			FromBalance:   fromAccount.Balance,
			ToBalance:     toAccount.Balance,
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, errTransferIdempotencyConflict) && result != nil {
			return result, nil
		}
		return nil, err
	}

	return result, nil
}

func transferResultFromOperation(ctx context.Context, tx domain.Tx, input TransferInput, op *domain.Operation) (*TransferResult, error) {
	toAccountID := input.ToAccountID
	if op.RelatedAccountID != nil {
		toAccountID = *op.RelatedAccountID
	}

	fromAccount, err := tx.GetByID(ctx, input.FromAccountID)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("get source account for idempotency replay: %w", err)
	}

	toAccount, err := tx.GetByID(ctx, toAccountID)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("get destination account for idempotency replay: %w", err)
	}

	amount := op.Amount
	if amount <= 0 {
		amount = input.Amount
	}

	return &TransferResult{
		FromAccountID: input.FromAccountID,
		ToAccountID:   toAccountID,
		Amount:        amount,
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
