package infrastructure

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/seu-usuario/bank-api/internal/account/bankaccount/domain"
)

type recipientExecutorMock struct {
	query    string
	args     []any
	rows     pgx.Rows
	queryErr error
}

func (m *recipientExecutorMock) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (m *recipientExecutorMock) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	m.query = query
	m.args = args
	if m.queryErr != nil {
		return nil, m.queryErr
	}
	return m.rows, nil
}

func (m *recipientExecutorMock) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return recipientRowMock{}
}

type recipientRowMock struct{}

func (recipientRowMock) Scan(dest ...any) error {
	return nil
}

type recipientRowsMock struct {
	values [][]any
	index  int
	closed bool
}

func (r *recipientRowsMock) Close() {
	r.closed = true
}

func (r *recipientRowsMock) Err() error {
	return nil
}

func (r *recipientRowsMock) CommandTag() pgconn.CommandTag {
	return pgconn.CommandTag{}
}

func (r *recipientRowsMock) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (r *recipientRowsMock) Next() bool {
	if r.index >= len(r.values) {
		r.closed = true
		return false
	}
	r.index++
	return true
}

func (r *recipientRowsMock) Scan(dest ...any) error {
	row := r.values[r.index-1]
	for i := range dest {
		switch target := dest[i].(type) {
		case *uuid.UUID:
			*target = row[i].(uuid.UUID)
		case *string:
			*target = row[i].(string)
		}
	}
	return nil
}

func (r *recipientRowsMock) Values() ([]any, error) {
	return r.values[r.index-1], nil
}

func (r *recipientRowsMock) RawValues() [][]byte {
	return nil
}

func (r *recipientRowsMock) Conn() *pgx.Conn {
	return nil
}

func TestRepository_FindTransferRecipientsByBranchAndNumber_FiltersActiveAccounts(t *testing.T) {
	legacyCPFRef := "customers." + "cpf"

	accountID := uuid.New()
	rows := &recipientRowsMock{
		values: [][]any{{
			accountID,
			"Maria Silva",
			"12345678901",
			"0001",
			"12345678",
		}},
	}
	exec := &recipientExecutorMock{rows: rows}
	repo := &Repository{exec: exec}

	recipients, err := repo.FindTransferRecipientsByBranchAndNumber(context.Background(), "0001", "12345678")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if strings.Contains(exec.query, legacyCPFRef) {
		t.Fatalf("expected query to avoid legacy cpf column reference, got %s", exec.query)
	}
	if !strings.Contains(exec.query, "customer_documents") {
		t.Fatalf("expected customer_documents join in query, got %s", exec.query)
	}
	if !strings.Contains(exec.query, "a.status = $3") {
		t.Fatalf("expected active status filter in query, got %s", exec.query)
	}
	if len(exec.args) != 3 || exec.args[0] != "0001" || exec.args[1] != "12345678" || exec.args[2] != domain.AccountActive {
		t.Fatalf("unexpected query args: %+v", exec.args)
	}
	if len(recipients) != 1 {
		t.Fatalf("expected 1 recipient, got %d", len(recipients))
	}
	if recipients[0].AccountID != accountID {
		t.Fatalf("expected account id %v, got %v", accountID, recipients[0].AccountID)
	}
	if recipients[0].MaskedDocument != "***.456.789-**" {
		t.Fatalf("expected masked document, got %q", recipients[0].MaskedDocument)
	}
}

func TestRepository_FindTransferRecipientsByDocument_FiltersActiveAccounts(t *testing.T) {
	legacyCPFRef := "customers." + "cpf"

	exec := &recipientExecutorMock{rows: &recipientRowsMock{}}
	repo := &Repository{exec: exec}

	recipients, err := repo.FindTransferRecipientsByDocument(context.Background(), "12345678901")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if strings.Contains(exec.query, legacyCPFRef) {
		t.Fatalf("expected query to avoid legacy cpf column reference, got %s", exec.query)
	}
	if !strings.Contains(exec.query, "customer_documents") {
		t.Fatalf("expected customer_documents join in query, got %s", exec.query)
	}
	if !strings.Contains(exec.query, "cd.value = $1") {
		t.Fatalf("expected document filter by customer_documents value, got %s", exec.query)
	}
	if !strings.Contains(exec.query, "a.status = $2") {
		t.Fatalf("expected active status filter in query, got %s", exec.query)
	}
	if len(exec.args) != 2 || exec.args[0] != "12345678901" || exec.args[1] != domain.AccountActive {
		t.Fatalf("unexpected query args: %+v", exec.args)
	}
	if recipients == nil {
		t.Fatal("expected empty recipients slice, got nil")
	}
}
