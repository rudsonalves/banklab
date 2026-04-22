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

var errTransferDuplicateConflict = errors.New("transfer duplicate conflict")

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
		// Idempotency check: if a ledger entry already exists for this key, replay result.
		if input.IdempotencyKey != "" {
			existing, err := tx.GetTransactionByIdempotencyKey(ctx, input.FromAccountID, input.IdempotencyKey)
			if err != nil {
				return fmt.Errorf("get transaction by idempotency key: %w", err)
			}

			if existing != nil {
				result, err = transferResultFromLedger(ctx, tx, existing)
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

		// Origin side carries idempotency_key and related_account_id.
		var outgoing *domain.Transaction
		if input.IdempotencyKey != "" {
			outgoing = domain.NewTransactionWithIdempotency(
				input.FromAccountID,
				domain.TransactionTransferOut,
				input.Amount,
				fromAccount.Balance,
				&referenceID,
				&input.ToAccountID,
				input.IdempotencyKey,
			)
		} else {
			outgoing = domain.NewTransaction(
				input.FromAccountID,
				domain.TransactionTransferOut,
				input.Amount,
				fromAccount.Balance,
				&referenceID,
			)
			outgoing.RelatedAccountID = &input.ToAccountID
		}
		if err := tx.CreateTransaction(ctx, outgoing); err != nil {
			if errors.Is(err, domain.ErrTransferDuplicate) {
				existing, getErr := tx.GetTransactionByIdempotencyKey(ctx, input.FromAccountID, input.IdempotencyKey)
				if getErr != nil {
					return fmt.Errorf("reload transaction by idempotency key: %w", getErr)
				}
				if existing == nil {
					return fmt.Errorf("reload transaction by idempotency key: not found")
				}
				result, getErr = transferResultFromLedger(ctx, tx, existing)
				if getErr != nil {
					return getErr
				}
				// Rollback the duplicate execution while preserving the replay result.
				return errTransferDuplicateConflict
			}
			return fmt.Errorf("create transfer out ledger transaction: %w", err)
		}

		incoming := domain.NewTransaction(
			input.ToAccountID,
			domain.TransactionTransferIn,
			input.Amount,
			toAccount.Balance,
			&referenceID,
		)
		incoming.RelatedAccountID = &input.FromAccountID
		if err := tx.CreateTransaction(ctx, incoming); err != nil {
			return fmt.Errorf("create transfer in ledger transaction: %w", err)
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
		if errors.Is(err, errTransferDuplicateConflict) && result != nil {
			return result, nil
		}
		return nil, err
	}

	return result, nil
}

func transferResultFromLedger(ctx context.Context, tx domain.Tx, outgoing *domain.Transaction) (*TransferResult, error) {
	if outgoing == nil {
		return nil, fmt.Errorf("ledger inconsistency: outgoing transaction is nil")
	}

	if outgoing.Type != domain.TransactionTransferOut {
		return nil, fmt.Errorf("ledger inconsistency: expected transfer_out, got %s", outgoing.Type)
	}

	if outgoing.RelatedAccountID == nil {
		return nil, fmt.Errorf("ledger inconsistency: missing related_account_id on transfer_out")
	}

	if outgoing.ReferenceID == nil {
		return nil, fmt.Errorf("ledger inconsistency: missing reference_id on transfer_out")
	}

	incoming, err := tx.GetTransactionByReference(
		ctx,
		*outgoing.RelatedAccountID,
		*outgoing.ReferenceID,
		domain.TransactionTransferIn,
	)
	if err != nil {
		return nil, fmt.Errorf("get transfer_in by reference: %w", err)
	}

	if incoming == nil {
		return nil, fmt.Errorf("ledger inconsistency: transfer_in not found for reference %s", outgoing.ReferenceID.String())
	}

	if incoming.RelatedAccountID == nil || *incoming.RelatedAccountID != outgoing.AccountID {
		return nil, fmt.Errorf("ledger inconsistency: transfer_in related_account_id mismatch")
	}

	return &TransferResult{
		FromAccountID: outgoing.AccountID,
		ToAccountID:   *outgoing.RelatedAccountID,
		Amount:        outgoing.Amount,
		FromBalance:   outgoing.BalanceAfter,
		ToBalance:     incoming.BalanceAfter,
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
