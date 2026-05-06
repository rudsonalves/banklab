package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	bankaccountdomain "github.com/seu-usuario/bank-api/internal/account/bankaccount/domain"
	statementdomain "github.com/seu-usuario/bank-api/internal/account/statement/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

const (
	defaultStatementLimit = 50
	maxStatementLimit     = 100
)

type GetStatementInput struct {
	User *authdomain.AuthenticatedUser

	AccountID uuid.UUID

	Limit    int
	Cursor   *time.Time
	CursorID *uuid.UUID

	From *time.Time
	To   *time.Time
}

type StatementCursor struct {
	CreatedAt time.Time
	ID        string
}

type StatementItem struct {
	TransactionID string
	Type          string
	Amount        int64
	BalanceAfter  int64
	ReferenceID   *string
	CreatedAt     time.Time
}

type Statement struct {
	AccountID  string
	Items      []StatementItem
	NextCursor *StatementCursor
}

type GetStatement struct {
	repo statementdomain.Repository
}

// NewGetStatement creates a new instance of the GetStatement use case with
// the provided statement repository.
func NewGetStatement(repo statementdomain.Repository) *GetStatement {
	return &GetStatement{repo: repo}
}

// Execute retrieves the account statement for the specified account ID, applying
// pagination and date filtering based on the input parameters. It checks for
// valid input, verifies access permissions, and returns a structured statement
// containing the transaction history and pagination information.
func (uc *GetStatement) Execute(ctx context.Context, input GetStatementInput) (*Statement, error) {
	if input.AccountID == uuid.Nil {
		return nil, bankaccountdomain.ErrInvalidData
	}

	if input.From != nil && input.To != nil && input.From.After(*input.To) {
		return nil, bankaccountdomain.ErrInvalidData
	}

	if (input.Cursor == nil) != (input.CursorID == nil) {
		return nil, bankaccountdomain.ErrInvalidData
	}

	limit := input.Limit
	if limit == 0 {
		limit = defaultStatementLimit
	}
	if limit < 0 {
		return nil, bankaccountdomain.ErrInvalidData
	}
	if limit > maxStatementLimit {
		limit = maxStatementLimit
	}

	account, err := uc.repo.GetByID(ctx, input.AccountID)
	if err != nil {
		return nil, err
	}

	if !CanAccessAccount(input.User, account) {
		return nil, bankaccountdomain.ErrForbidden
	}

	transactions, err := uc.repo.GetTransactions(
		ctx,
		input.AccountID,
		limit,
		input.Cursor,
		input.CursorID,
		input.From,
		input.To,
	)
	if err != nil {
		if errors.Is(err, bankaccountdomain.ErrAccountNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("get account transactions: %w", err)
	}

	items := make([]StatementItem, 0, len(transactions))
	for _, tx := range transactions {
		item := StatementItem{
			TransactionID: tx.ID.String(),
			Type:          string(tx.Type),
			Amount:        tx.Amount,
			BalanceAfter:  tx.BalanceAfter,
			CreatedAt:     tx.CreatedAt,
		}
		if tx.ReferenceID != nil {
			referenceID := tx.ReferenceID.String()
			item.ReferenceID = &referenceID
		}

		items = append(items, item)
	}

	statement := &Statement{
		AccountID: input.AccountID.String(),
		Items:     items,
	}

	if len(transactions) == limit && len(transactions) > 0 {
		last := transactions[len(transactions)-1]
		statement.NextCursor = &StatementCursor{
			CreatedAt: last.CreatedAt,
			ID:        last.ID.String(),
		}
	}

	return statement, nil
}
