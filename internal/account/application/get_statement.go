package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/domain"
)

const (
	defaultStatementLimit = 50
	maxStatementLimit     = 100
)

type GetStatementInput struct {
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
	repo domain.AccountRepository
}

func NewGetStatement(repo domain.AccountRepository) *GetStatement {
	return &GetStatement{repo: repo}
}

func (uc *GetStatement) Execute(ctx context.Context, input GetStatementInput) (*Statement, error) {
	if input.AccountID == uuid.Nil {
		return nil, domain.ErrInvalidData
	}

	if input.From != nil && input.To != nil && input.From.After(*input.To) {
		return nil, domain.ErrInvalidData
	}

	if (input.Cursor == nil) != (input.CursorID == nil) {
		return nil, domain.ErrInvalidData
	}

	limit := input.Limit
	if limit == 0 {
		limit = defaultStatementLimit
	}
	if limit < 0 {
		return nil, domain.ErrInvalidData
	}
	if limit > maxStatementLimit {
		limit = maxStatementLimit
	}

	if _, err := uc.repo.GetByID(ctx, input.AccountID); err != nil {
		return nil, err
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
		if errors.Is(err, domain.ErrAccountNotFound) {
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
