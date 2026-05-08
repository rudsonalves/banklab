package application

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/transaction/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

var errTransferDuplicateConflict = errors.New("transfer duplicate conflict")

type Transfer struct {
	accountRepo domain.Repository
}

// NewTransfer creates a new instance of the Transfer use case with the provided account repository.
func NewTransfer(accountRepo domain.Repository) *Transfer {
	return &Transfer{accountRepo: accountRepo}
}

type TransferInput struct {
	User           *authdomain.AuthenticatedUser
	FromAccountID  uuid.UUID
	ToAccountID    uuid.UUID
	Amount         int64
	IdempotencyKey string
	Description    *string
}

type TransferResult struct {
	FromAccountID        uuid.UUID
	ToAccountID          uuid.UUID
	TransactionReference uuid.UUID
	Amount               int64
	FromBalance          int64
	ToBalance            int64
}

// Execute performs a transfer transaction from one account to another with the specified amount.
func (uc *Transfer) Execute(ctx context.Context, input TransferInput) (_ *TransferResult, err error) {
	if input.Amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	if input.FromAccountID == uuid.Nil || input.ToAccountID == uuid.Nil {
		return nil, domain.ErrInvalidData
	}

	idempotencyKey := strings.TrimSpace(input.IdempotencyKey)
	if idempotencyKey == "" {
		return nil, domain.ErrInvalidData
	}

	if input.FromAccountID == input.ToAccountID {
		return nil, domain.ErrSameAccountTransfer
	}

	description := normalizeTransferDescription(input.Description)
	var result *TransferResult

	err = uc.accountRepo.WithTransaction(ctx, func(tx domain.Tx) error {
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

		// Idempotency check: if a ledger entry already exists for this key,
		// replay the result after source-account authorization succeeds.
		existing, err := tx.GetTransactionByIdempotencyKey(ctx, input.FromAccountID, idempotencyKey)
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
		outgoing := domain.NewTransactionWithIdempotency(
			input.FromAccountID,
			domain.TransactionTransferOut,
			input.Amount,
			fromAccount.Balance,
			&referenceID,
			&input.ToAccountID,
			idempotencyKey,
		)
		outgoing.Description = description
		if err := tx.CreateTransaction(ctx, outgoing); err != nil {
			if errors.Is(err, domain.ErrTransferDuplicate) {
				existing, getErr := tx.GetTransactionByIdempotencyKey(ctx, input.FromAccountID, idempotencyKey)
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
			FromAccountID:        input.FromAccountID,
			ToAccountID:          input.ToAccountID,
			TransactionReference: referenceID,
			Amount:               input.Amount,
			FromBalance:          fromAccount.Balance,
			ToBalance:            toAccount.Balance,
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

// transferResultFromLedger reconstructs a TransferResult from the ledger
// transactions of a completed transfer. It validates the consistency of the
// ledger entries to ensure the integrity of the transfer data.
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
		FromAccountID:        outgoing.AccountID,
		ToAccountID:          *outgoing.RelatedAccountID,
		TransactionReference: *outgoing.ReferenceID,
		Amount:               outgoing.Amount,
		FromBalance:          outgoing.BalanceAfter,
		ToBalance:            incoming.BalanceAfter,
	}, nil
}

// orderedUUIDs returns the two UUIDs in a consistent order (lowest first). This is used to ensure that when locking two accounts for a transfer, we always
// lock them in the same order to reduce the risk of deadlocks.
func orderedUUIDs(left, right uuid.UUID) (uuid.UUID, uuid.UUID) {
	if bytes.Compare(left[:], right[:]) <= 0 {
		return left, right
	}

	return right, left
}

func normalizeTransferDescription(description *string) *string {
	if description == nil {
		return nil
	}

	normalized := strings.TrimSpace(*description)
	if normalized == "" {
		return nil
	}

	return &normalized
}

// mapTransferAccounts takes the fromAccountID and the two accounts locked
// in deterministic order, and returns them mapped to fromAccount and
// toAccount based on the IDs.
func mapTransferAccounts(fromID uuid.UUID, first, second *domain.Account) (*domain.Account, *domain.Account) {
	if first.ID == fromID {
		return first, second
	}

	return second, first
}
