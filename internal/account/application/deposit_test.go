package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/domain"
)

type depositAccountRepositoryMock struct {
	beginTxCalls int
	beginTxErr   error
	tx           domain.Tx
}

func (m *depositAccountRepositoryMock) Create(ctx context.Context, account *domain.Account) error {
	return nil
}

func (m *depositAccountRepositoryMock) CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
	return nil
}

func (m *depositAccountRepositoryMock) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
	return false, nil
}

func (m *depositAccountRepositoryMock) NextAccountNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (m *depositAccountRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return nil, nil
}

func (m *depositAccountRepositoryMock) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return nil, nil
}

func (m *depositAccountRepositoryMock) GetTransactions(
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

func (m *depositAccountRepositoryMock) UpdateBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return 0, nil
}

func (m *depositAccountRepositoryMock) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) error {
	return nil
}

func (m *depositAccountRepositoryMock) BeginTx(ctx context.Context) (domain.Tx, error) {
	m.beginTxCalls++
	if m.beginTxErr != nil {
		return nil, m.beginTxErr
	}
	return m.tx, nil
}

type txMock struct {
	getByIDCalls           int
	getByIDForUpdateCalls  int
	updateBalanceCalls     int
	decreaseBalanceCalls   int
	createTransactionCalls int
	commitCalls            int
	rollbackCalls          int

	account              *domain.Account
	getByIDErr           error
	updateBalanceValue   int64
	updateBalanceErr     error
	decreaseBalanceErr   error
	createTransactionErr error
	createdTransactions  []*domain.Transaction
	commitErr            error
	rollbackErr          error
}

func (m *txMock) Create(ctx context.Context, account *domain.Account) error {
	return nil
}

func (m *txMock) CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
	m.createTransactionCalls++
	m.createdTransactions = append(m.createdTransactions, tx)
	return m.createTransactionErr
}

func (m *txMock) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
	return false, nil
}

func (m *txMock) NextAccountNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (m *txMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	m.getByIDCalls++
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	return m.account, nil
}

func (m *txMock) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	m.getByIDForUpdateCalls++
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	return m.account, nil
}

func (m *txMock) GetTransactions(
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

func (m *txMock) UpdateBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	m.updateBalanceCalls++
	return m.updateBalanceValue, m.updateBalanceErr
}

func (m *txMock) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) error {
	m.decreaseBalanceCalls++
	return m.decreaseBalanceErr
}

func (m *txMock) BeginTx(ctx context.Context) (domain.Tx, error) {
	return nil, nil
}

func (m *txMock) Commit(ctx context.Context) error {
	m.commitCalls++
	return m.commitErr
}

func (m *txMock) Rollback(ctx context.Context) error {
	m.rollbackCalls++
	return m.rollbackErr
}

func TestDeposit_Execute_InvalidAmount(t *testing.T) {
	repo := &depositAccountRepositoryMock{}
	useCase := NewDeposit(repo)

	account, err := useCase.Execute(context.Background(), DepositInput{
		AccountID: uuid.New(),
		Amount:    0,
	})

	if !errors.Is(err, domain.ErrInvalidAmount) {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidAmount, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if repo.beginTxCalls != 0 {
		t.Fatalf("expected BeginTx not to be called, got %d calls", repo.beginTxCalls)
	}
}

func TestDeposit_Execute_AccountNotFound(t *testing.T) {
	tx := &txMock{getByIDErr: domain.ErrAccountNotFound}
	repo := &depositAccountRepositoryMock{tx: tx}
	useCase := NewDeposit(repo)

	account, err := useCase.Execute(context.Background(), DepositInput{
		AccountID: uuid.New(),
		Amount:    100,
	})

	if !errors.Is(err, domain.ErrAccountNotFound) {
		t.Fatalf("expected error %v, got %v", domain.ErrAccountNotFound, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if tx.rollbackCalls != 1 {
		t.Fatalf("expected Rollback to be called once, got %d calls", tx.rollbackCalls)
	}

	if tx.commitCalls != 0 {
		t.Fatalf("expected Commit not to be called, got %d calls", tx.commitCalls)
	}
}

func TestDeposit_Execute_AccountInactive(t *testing.T) {
	tx := &txMock{
		account: &domain.Account{
			ID:     uuid.New(),
			Status: domain.AccountInactive,
		},
	}
	repo := &depositAccountRepositoryMock{tx: tx}
	useCase := NewDeposit(repo)

	account, err := useCase.Execute(context.Background(), DepositInput{
		AccountID: uuid.New(),
		Amount:    100,
	})

	if !errors.Is(err, domain.ErrAccountInactive) {
		t.Fatalf("expected error %v, got %v", domain.ErrAccountInactive, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if tx.updateBalanceCalls != 0 {
		t.Fatalf("expected UpdateBalance not to be called, got %d calls", tx.updateBalanceCalls)
	}

	if tx.createTransactionCalls != 0 {
		t.Fatalf("expected CreateTransaction not to be called, got %d calls", tx.createTransactionCalls)
	}

	if tx.rollbackCalls != 1 {
		t.Fatalf("expected Rollback to be called once, got %d calls", tx.rollbackCalls)
	}

	if tx.commitCalls != 0 {
		t.Fatalf("expected Commit not to be called, got %d calls", tx.commitCalls)
	}
}

func TestDeposit_Execute_Success(t *testing.T) {
	initialBalance := int64(100)
	depositAmount := int64(50)
	accountID := uuid.New()
	dbReturnedBalance := int64(999)

	tx := &txMock{
		account: &domain.Account{
			ID:      accountID,
			Balance: initialBalance,
			Status:  domain.AccountActive,
		},
		updateBalanceValue: dbReturnedBalance,
	}
	repo := &depositAccountRepositoryMock{tx: tx}
	useCase := NewDeposit(repo)

	account, err := useCase.Execute(context.Background(), DepositInput{
		AccountID: accountID,
		Amount:    depositAmount,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if account == nil {
		t.Fatal("expected account to be non-nil")
	}

	if account.Balance != dbReturnedBalance {
		t.Fatalf("expected balance %d, got %d", dbReturnedBalance, account.Balance)
	}

	if tx.createTransactionCalls != 1 {
		t.Fatalf("expected CreateTransaction to be called once, got %d calls", tx.createTransactionCalls)
	}

	created := tx.createdTransactions[0]
	if created.Type != domain.TransactionDeposit {
		t.Fatalf("expected ledger type %s, got %s", domain.TransactionDeposit, created.Type)
	}

	if created.BalanceAfter != dbReturnedBalance {
		t.Fatalf("expected ledger balance_after %d, got %d", dbReturnedBalance, created.BalanceAfter)
	}

	if tx.updateBalanceCalls != 1 {
		t.Fatalf("expected UpdateBalance to be called once, got %d calls", tx.updateBalanceCalls)
	}

	if tx.commitCalls != 1 {
		t.Fatalf("expected Commit to be called once, got %d calls", tx.commitCalls)
	}

	if tx.rollbackCalls != 0 {
		t.Fatalf("expected Rollback not to be called, got %d calls", tx.rollbackCalls)
	}
}

func TestDeposit_Execute_RepositoryFailure(t *testing.T) {
	expectedErr := errors.New("update failed")
	tx := &txMock{
		account:          &domain.Account{ID: uuid.New(), Balance: 200, Status: domain.AccountActive},
		updateBalanceErr: expectedErr,
	}
	repo := &depositAccountRepositoryMock{tx: tx}
	useCase := NewDeposit(repo)

	account, err := useCase.Execute(context.Background(), DepositInput{
		AccountID: uuid.New(),
		Amount:    10,
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if tx.rollbackCalls != 1 {
		t.Fatalf("expected Rollback to be called once, got %d calls", tx.rollbackCalls)
	}

	if tx.createTransactionCalls != 0 {
		t.Fatalf("expected CreateTransaction not to be called, got %d calls", tx.createTransactionCalls)
	}

	if tx.commitCalls != 0 {
		t.Fatalf("expected Commit not to be called, got %d calls", tx.commitCalls)
	}
}

func TestDeposit_Execute_LedgerInsertFailure(t *testing.T) {
	expectedErr := errors.New("ledger insert failed")
	tx := &txMock{
		account:              &domain.Account{ID: uuid.New(), Balance: 100, Status: domain.AccountActive},
		updateBalanceValue:   150,
		createTransactionErr: expectedErr,
	}
	repo := &depositAccountRepositoryMock{tx: tx}
	useCase := NewDeposit(repo)

	account, err := useCase.Execute(context.Background(), DepositInput{
		AccountID: uuid.New(),
		Amount:    50,
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if tx.createTransactionCalls != 1 {
		t.Fatalf("expected CreateTransaction to be called once, got %d calls", tx.createTransactionCalls)
	}

	if tx.rollbackCalls != 1 {
		t.Fatalf("expected Rollback to be called once, got %d calls", tx.rollbackCalls)
	}

	if tx.commitCalls != 0 {
		t.Fatalf("expected Commit not to be called, got %d calls", tx.commitCalls)
	}
}
