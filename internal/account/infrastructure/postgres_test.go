package infrastructure

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/seu-usuario/bank-api/internal/account/domain"
)

type fakeExecutor struct {
	rows []pgx.Row
}

func (f *fakeExecutor) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, errors.New("not implemented")
}

func (f *fakeExecutor) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeExecutor) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	if len(f.rows) == 0 {
		return fakeRow{err: errors.New("no more fake rows")}
	}

	row := f.rows[0]
	f.rows = f.rows[1:]
	return row
}

type fakeRow struct {
	err      error
	intValue int
	i64Value int64
	setInt   bool
	setI64   bool
}

func (r fakeRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) == 0 {
		return nil
	}

	if r.setInt {
		if v, ok := dest[0].(*int); ok {
			*v = r.intValue
			return nil
		}
	}

	if r.setI64 {
		if v, ok := dest[0].(*int64); ok {
			*v = r.i64Value
			return nil
		}
	}

	return nil
}

type fakeTx struct {
	commitCalls   int
	rollbackCalls int
	commitErr     error
	rollbackErr   error
}

func (f *fakeTx) Create(ctx context.Context, account *domain.Account) error { return nil }

func (f *fakeTx) CreateTransaction(ctx context.Context, tx *domain.Transaction) error { return nil }

func (f *fakeTx) GetOperationByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*domain.Operation, error) {
	return nil, nil
}

func (f *fakeTx) CreateOperation(ctx context.Context, op *domain.Operation) error { return nil }

func (f *fakeTx) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
	return false, nil
}

func (f *fakeTx) NextAccountNumber(ctx context.Context) (string, error) { return "", nil }

func (f *fakeTx) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) { return nil, nil }

func (f *fakeTx) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return nil, nil
}

func (f *fakeTx) GetTransactions(
	ctx context.Context,
	accountID uuid.UUID,
	limit int,
	cursorTime *time.Time,
	cursorID *uuid.UUID,
	from *time.Time,
	to *time.Time,
) ([]domain.Transaction, error) {
	return nil, nil
}

func (f *fakeTx) IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return 0, nil
}

func (f *fakeTx) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return 0, nil
}

func (f *fakeTx) BeginTx(ctx context.Context) (domain.Tx, error) {
	return nil, errors.New("nested transactions are not supported")
}

func (f *fakeTx) WithTransaction(ctx context.Context, fn func(tx domain.Tx) error) error {
	return errors.New("nested transactions are not supported")
}

func (f *fakeTx) Commit(ctx context.Context) error {
	f.commitCalls++
	return f.commitErr
}

func (f *fakeTx) Rollback(ctx context.Context) error {
	f.rollbackCalls++
	return f.rollbackErr
}

func TestDecreaseBalance_AccountNotFound(t *testing.T) {
	repo := baseRepository{exec: &fakeExecutor{rows: []pgx.Row{
		fakeRow{err: pgx.ErrNoRows},
		fakeRow{err: pgx.ErrNoRows},
	}}}

	_, err := repo.DecreaseBalance(context.Background(), uuid.New(), 10)
	if !errors.Is(err, domain.ErrAccountNotFound) {
		t.Fatalf("expected ErrAccountNotFound, got %v", err)
	}
}

func TestDecreaseBalance_InsufficientBalance(t *testing.T) {
	repo := baseRepository{exec: &fakeExecutor{rows: []pgx.Row{
		fakeRow{err: pgx.ErrNoRows},
		fakeRow{setInt: true, intValue: 1},
	}}}

	_, err := repo.DecreaseBalance(context.Background(), uuid.New(), 10)
	if !errors.Is(err, domain.ErrInsufficientBalance) {
		t.Fatalf("expected ErrInsufficientBalance, got %v", err)
	}
}

func TestDecreaseBalance_Success(t *testing.T) {
	repo := baseRepository{exec: &fakeExecutor{rows: []pgx.Row{
		fakeRow{setI64: true, i64Value: 90},
	}}}

	balance, err := repo.DecreaseBalance(context.Background(), uuid.New(), 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if balance != 90 {
		t.Fatalf("expected balance 90, got %d", balance)
	}
}

func TestDecreaseBalance_RepositoryAndTxRepositoryParity(t *testing.T) {
	repo := &Repository{base: baseRepository{exec: &fakeExecutor{rows: []pgx.Row{
		fakeRow{err: pgx.ErrNoRows},
		fakeRow{setInt: true, intValue: 1},
	}}}}
	txRepo := &txRepository{base: baseRepository{exec: &fakeExecutor{rows: []pgx.Row{
		fakeRow{err: pgx.ErrNoRows},
		fakeRow{setInt: true, intValue: 1},
	}}}}

	_, errRepo := repo.DecreaseBalance(context.Background(), uuid.New(), 10)
	_, errTx := txRepo.DecreaseBalance(context.Background(), uuid.New(), 10)
	if !errors.Is(errRepo, domain.ErrInsufficientBalance) {
		t.Fatalf("expected repository error ErrInsufficientBalance, got %v", errRepo)
	}
	if !errors.Is(errTx, domain.ErrInsufficientBalance) {
		t.Fatalf("expected tx repository error ErrInsufficientBalance, got %v", errTx)
	}
}

func TestRunInTransaction_RollsBackOnCallbackError(t *testing.T) {
	tx := &fakeTx{}
	expectedErr := errors.New("callback failed")

	err := runInTransaction(context.Background(), tx, func(tx domain.Tx) error {
		return expectedErr
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected callback error, got %v", err)
	}
	if tx.rollbackCalls != 1 {
		t.Fatalf("expected one rollback call, got %d", tx.rollbackCalls)
	}
	if tx.commitCalls != 0 {
		t.Fatalf("expected zero commit calls, got %d", tx.commitCalls)
	}
}

func TestRunInTransaction_CommitsOnSuccess(t *testing.T) {
	tx := &fakeTx{}

	err := runInTransaction(context.Background(), tx, func(tx domain.Tx) error {
		return nil
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tx.commitCalls != 1 {
		t.Fatalf("expected one commit call, got %d", tx.commitCalls)
	}
	if tx.rollbackCalls != 0 {
		t.Fatalf("expected zero rollback calls, got %d", tx.rollbackCalls)
	}
}

func TestRunInTransaction_CommitFailureTriggersRollback(t *testing.T) {
	tx := &fakeTx{commitErr: errors.New("commit failed")}

	err := runInTransaction(context.Background(), tx, func(tx domain.Tx) error {
		return nil
	})

	if err == nil || !strings.Contains(err.Error(), "commit transaction") {
		t.Fatalf("expected commit transaction error, got %v", err)
	}
	if tx.commitCalls != 1 {
		t.Fatalf("expected one commit call, got %d", tx.commitCalls)
	}
	if tx.rollbackCalls != 1 {
		t.Fatalf("expected one rollback call, got %d", tx.rollbackCalls)
	}
}
