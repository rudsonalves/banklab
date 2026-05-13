package infrastructure

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	transactiondomain "github.com/seu-usuario/bank-api/internal/account/transaction/domain"
)

type statementExecutorMock struct {
	query string
	args  []any
	rows  pgx.Rows
}

func (m *statementExecutorMock) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	m.query = query
	m.args = args
	return m.rows, nil
}

func (m *statementExecutorMock) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return statementRowMock{}
}

func (m *statementExecutorMock) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

type statementRowMock struct{}

func (statementRowMock) Scan(dest ...any) error {
	return nil
}

type statementRowsMock struct {
	values [][]any
	index  int
	closed bool
}

func (r *statementRowsMock) Close() {
	r.closed = true
}

func (r *statementRowsMock) Err() error {
	return nil
}

func (r *statementRowsMock) CommandTag() pgconn.CommandTag {
	return pgconn.CommandTag{}
}

func (r *statementRowsMock) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (r *statementRowsMock) Next() bool {
	if r.index >= len(r.values) {
		r.closed = true
		return false
	}
	r.index++
	return true
}

func (r *statementRowsMock) Scan(dest ...any) error {
	row := r.values[r.index-1]
	for i := range dest {
		switch target := dest[i].(type) {
		case *uuid.UUID:
			*target = row[i].(uuid.UUID)
		case *transactiondomain.TransactionType:
			*target = row[i].(transactiondomain.TransactionType)
		case *int64:
			*target = row[i].(int64)
		case **uuid.UUID:
			if row[i] == nil {
				*target = nil
			} else {
				value := row[i].(uuid.UUID)
				*target = &value
			}
		case **string:
			if row[i] == nil {
				*target = nil
			} else {
				value := row[i].(string)
				*target = &value
			}
		case *time.Time:
			*target = row[i].(time.Time)
		}
	}
	return nil
}

func (r *statementRowsMock) Values() ([]any, error) {
	return r.values[r.index-1], nil
}

func (r *statementRowsMock) RawValues() [][]byte {
	return nil
}

func (r *statementRowsMock) Conn() *pgx.Conn {
	return nil
}

func TestRepository_GetTransactions_ReturnsDescription(t *testing.T) {
	accountID := uuid.New()
	transactionID := uuid.New()
	referenceID := uuid.New()
	relatedAccountID := uuid.New()
	idempotencyKey := "statement-description-key"
	description := "Aluguel de maio"
	createdAt := time.Now().UTC().Truncate(time.Second)

	rows := &statementRowsMock{
		values: [][]any{{
			transactionID,
			accountID,
			transactiondomain.TransactionTransferOut,
			int64(-100),
			int64(900),
			referenceID,
			relatedAccountID,
			idempotencyKey,
			description,
			createdAt,
		}},
	}
	exec := &statementExecutorMock{rows: rows}
	repo := &Repository{exec: exec}

	transactions, err := repo.GetTransactions(
		context.Background(),
		accountID,
		10,
		nil,
		nil,
		nil,
		nil,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(exec.query, "description") {
		t.Fatalf("expected query to select description, got %s", exec.query)
	}
	if len(transactions) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(transactions))
	}
	if transactions[0].Description == nil || *transactions[0].Description != description {
		t.Fatalf("expected description %q, got %+v", description, transactions[0].Description)
	}
	if transactions[0].CreatedAt != createdAt {
		t.Fatalf("expected created_at %v, got %v", createdAt, transactions[0].CreatedAt)
	}
}
