package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	transactiondomain "github.com/seu-usuario/bank-api/internal/account/transaction/domain"
)

type baseRepository struct {
	exec executor
}

// GetByIDForUpdate retrieves an account by its ID and locks it for update.
func (r *baseRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*transactiondomain.Account, error) {
	var account transactiondomain.Account

	query := `
		SELECT id, customer_id, balance, status
		FROM accounts
		WHERE id = $1
		FOR UPDATE
	`

	err := r.exec.QueryRow(ctx, query, id).Scan(
		&account.ID,
		&account.CustomerID,
		&account.Balance,
		&account.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, transactiondomain.ErrAccountNotFound
		}
		return nil, fmt.Errorf("get account by id for update: %w", err)
	}

	return &account, nil
}

func (r *baseRepository) GetByBranchAndNumber(ctx context.Context, branch, number string) (*transactiondomain.Account, error) {
	var account transactiondomain.Account

	query := `
		SELECT id, customer_id, balance, status
		FROM accounts
		WHERE branch = $1 AND number = $2
	`

	err := r.exec.QueryRow(ctx, query, branch, number).Scan(
		&account.ID,
		&account.CustomerID,
		&account.Balance,
		&account.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, transactiondomain.ErrAccountNotFound
		}
		return nil, fmt.Errorf("get account by branch and number: %w", err)
	}

	return &account, nil
}

// IncreaseBalance increases the balance of the specified account by the
// given amount and returns the new balance.
func (r *baseRepository) IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	var balance int64

	query := `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
		RETURNING balance
	`

	err := r.exec.QueryRow(ctx, query, amount, id).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, transactiondomain.ErrInvalidAmount
		}
		return 0, fmt.Errorf("update balance: %w", err)
	}

	return balance, nil
}

// DecreaseBalance decreases the balance of the specified account by the
// given amount and returns the new balance.
func (r *baseRepository) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	var balance int64

	query := `
		UPDATE accounts
		SET balance = balance - $1
		WHERE id = $2
		  AND balance >= $1
		RETURNING balance
	`

	err := r.exec.QueryRow(ctx, query, amount, id).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			exists, existsErr := r.accountExists(ctx, id)
			if existsErr != nil {
				return 0, fmt.Errorf("decrease balance: %w", existsErr)
			}
			if !exists {
				return 0, transactiondomain.ErrAccountNotFound
			}
			return 0, transactiondomain.ErrInsufficientBalance
		}
		return 0, fmt.Errorf("decrease balance: %w", err)
	}

	return balance, nil
}

// accountExists checks if an account with the given ID exists in the
// database.
func (r *baseRepository) accountExists(ctx context.Context, id uuid.UUID) (bool, error) {
	var dummy int
	query := `
		SELECT 1
		FROM accounts
		WHERE id = $1
	`

	err := r.exec.QueryRow(ctx, query, id).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("check account exists: %w", err)
	}

	return true, nil
}

// CreateTransaction inserts a new transaction record into the database.
// If a transaction with the same idempotency key already exists for the
// account, it returns a specific error.
func (r *baseRepository) CreateTransaction(ctx context.Context, tx *transactiondomain.Transaction) error {
	query := `
		INSERT INTO transactions (
			id, account_id, type, amount, balance_after, reference_id,
			related_account_id, idempotency_key, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (account_id, idempotency_key)
		WHERE idempotency_key IS NOT NULL
		DO NOTHING
	`

	cmd, err := r.exec.Exec(ctx, query,
		tx.ID,
		tx.AccountID,
		tx.Type,
		tx.Amount,
		tx.BalanceAfter,
		tx.ReferenceID,
		tx.RelatedAccountID,
		tx.IdempotencyKey,
		tx.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create account transaction: %w", err)
	}

	if tx.IdempotencyKey != nil && cmd.RowsAffected() == 0 {
		return transactiondomain.ErrTransferDuplicate
	}

	return nil
}

// GetTransactionByIdempotencyKey retrieves a transaction by its
// idempotency key for a specific account. It returns nil if no such
// transaction exists.
func (r *baseRepository) GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*transactiondomain.Transaction, error) {
	var t transactiondomain.Transaction

	query := `
		SELECT id, account_id, type, amount, balance_after, reference_id,
		       related_account_id, idempotency_key, created_at
		FROM transactions
		WHERE account_id = $1 AND idempotency_key = $2
		LIMIT 1
	`

	err := r.exec.QueryRow(ctx, query, accountID, key).Scan(
		&t.ID,
		&t.AccountID,
		&t.Type,
		&t.Amount,
		&t.BalanceAfter,
		&t.ReferenceID,
		&t.RelatedAccountID,
		&t.IdempotencyKey,
		&t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get transaction by idempotency key: %w", err)
	}

	return &t, nil
}

// GetTransactionByReference retrieves a transaction by its reference ID
// for a specific account and type. It returns nil if no such transaction
// exists.
func (r *baseRepository) GetTransactionByReference(ctx context.Context, accountID uuid.UUID, referenceID uuid.UUID, typeName transactiondomain.TransactionType) (*transactiondomain.Transaction, error) {
	var t transactiondomain.Transaction

	query := `
		SELECT id, account_id, type, amount, balance_after, reference_id,
		       related_account_id, idempotency_key, created_at
		FROM transactions
		WHERE account_id = $1 AND reference_id = $2 AND type = $3
		LIMIT 1
	`

	err := r.exec.QueryRow(ctx, query, accountID, referenceID, typeName).Scan(
		&t.ID,
		&t.AccountID,
		&t.Type,
		&t.Amount,
		&t.BalanceAfter,
		&t.ReferenceID,
		&t.RelatedAccountID,
		&t.IdempotencyKey,
		&t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get transaction by reference: %w", err)
	}

	return &t, nil
}

// GetTransferReceiptByReference retrieves the transfer receipt associated with
// the given reference ID. It returns a detailed receipt of the transfer
// operation, including information about the source and destination accounts,
// the amount transferred, and the status of the transfer.
func (r *baseRepository) GetTransferReceiptByReference(ctx context.Context, referenceID uuid.UUID) (*transactiondomain.TransferReceipt, error) {
	var receipt transactiondomain.TransferReceipt

	query := `
		SELECT
			transfer_out.type,
			transfer_out.amount,
			transfer_out.reference_id,
			transfer_out.created_at,
			source_account.id,
			source_account.customer_id,
			source_account.branch,
			source_account.number,
			destination_account.id,
			destination_account.customer_id,
			destination_account.branch,
			destination_account.number,
			recipient_customer.name
		FROM transactions transfer_out
		JOIN transactions transfer_in
		  ON transfer_in.reference_id = transfer_out.reference_id
		 AND transfer_in.type = 'transfer_in'
		JOIN accounts source_account
		  ON source_account.id = transfer_out.account_id
		JOIN accounts destination_account
		  ON destination_account.id = transfer_in.account_id
		JOIN customers recipient_customer
		  ON recipient_customer.id = destination_account.customer_id
		WHERE transfer_out.reference_id = $1
		  AND transfer_out.type = 'transfer_out'
	`

	err := r.exec.QueryRow(ctx, query, referenceID).Scan(
		&receipt.OperationType,
		&receipt.Amount,
		&receipt.TransactionReference,
		&receipt.OperationDate,
		&receipt.SourceAccountID,
		&receipt.SourceCustomerID,
		&receipt.SourceBranch,
		&receipt.SourceAccountNumber,
		&receipt.DestinationAccountID,
		&receipt.DestinationCustomerID,
		&receipt.DestinationBranch,
		&receipt.DestinationAccountNumber,
		&receipt.RecipientName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, transactiondomain.ErrTransactionNotFound
		}
		return nil, fmt.Errorf("get transfer receipt by reference: %w", err)
	}

	receipt.Status = "completed"
	return &receipt, nil
}
