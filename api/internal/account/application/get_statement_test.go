package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/domain"
)

type statementRepositoryMock struct {
	getByIDCalls           int
	getByIDErr             error
	account                *domain.Account
	getTransactionsCalls   int
	getTransactionsErr     error
	getTransactionsResult  []domain.Transaction
	lastGetTransactionsArg struct {
		accountID  uuid.UUID
		limit      int
		cursorTime *time.Time
		cursorID   *uuid.UUID
		from       *time.Time
		to         *time.Time
	}
}

func (m *statementRepositoryMock) Create(ctx context.Context, account *domain.Account) error {
	return nil
}

func (m *statementRepositoryMock) CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
	return nil
}

func (m *statementRepositoryMock) GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*domain.Transaction, error) {
	return nil, nil
}

func (m *statementRepositoryMock) GetTransactionByReference(ctx context.Context, accountID uuid.UUID, referenceID uuid.UUID, typeName domain.TransactionType) (*domain.Transaction, error) {
	return nil, nil
}

func (m *statementRepositoryMock) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
	return false, nil
}

func (m *statementRepositoryMock) NextAccountNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (m *statementRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	m.getByIDCalls++
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	if m.account != nil {
		return m.account, nil
	}

	return &domain.Account{ID: id}, nil
}

func (m *statementRepositoryMock) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return nil, nil
}

func (m *statementRepositoryMock) GetTransactions(
	ctx context.Context,
	accountID uuid.UUID,
	limit int,
	cursorTime *time.Time,
	cursorID *uuid.UUID,
	from *time.Time,
	to *time.Time,
) ([]domain.Transaction, error) {
	m.getTransactionsCalls++
	m.lastGetTransactionsArg.accountID = accountID
	m.lastGetTransactionsArg.limit = limit
	m.lastGetTransactionsArg.cursorTime = cursorTime
	m.lastGetTransactionsArg.cursorID = cursorID
	m.lastGetTransactionsArg.from = from
	m.lastGetTransactionsArg.to = to

	if m.getTransactionsErr != nil {
		return nil, m.getTransactionsErr
	}

	return m.getTransactionsResult, nil
}

func (m *statementRepositoryMock) IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return 0, nil
}

func (m *statementRepositoryMock) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return 0, nil
}

func (m *statementRepositoryMock) BeginTx(ctx context.Context) (domain.Tx, error) {
	return nil, nil
}

func (m *statementRepositoryMock) WithTransaction(ctx context.Context, fn func(tx domain.Tx) error) error {
	return errors.New("transactions are not used in this test")
}

func TestGetStatement_Execute_InvalidAccountID(t *testing.T) {
	repo := &statementRepositoryMock{}
	uc := NewGetStatement(repo)

	result, err := uc.Execute(context.Background(), GetStatementInput{})

	if !errors.Is(err, domain.ErrInvalidData) {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidData, err)
	}

	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}

	if repo.getByIDCalls != 0 {
		t.Fatalf("expected GetByID not to be called, got %d", repo.getByIDCalls)
	}
}

func TestGetStatement_Execute_DefaultAndCappedLimit(t *testing.T) {
	repo := &statementRepositoryMock{}
	uc := NewGetStatement(repo)
	accountID := uuid.New()
	customerID := uuid.New()
	repo.account = &domain.Account{ID: accountID, CustomerID: customerID}

	_, err := uc.Execute(context.Background(), GetStatementInput{User: testCustomerUser(customerID), AccountID: accountID})
	if err != nil {
		t.Fatalf("expected no error for default limit, got %v", err)
	}

	if repo.lastGetTransactionsArg.limit != 50 {
		t.Fatalf("expected default limit 50, got %d", repo.lastGetTransactionsArg.limit)
	}

	_, err = uc.Execute(context.Background(), GetStatementInput{User: testCustomerUser(customerID), AccountID: accountID, Limit: 500})
	if err != nil {
		t.Fatalf("expected no error for capped limit, got %v", err)
	}

	if repo.lastGetTransactionsArg.limit != 100 {
		t.Fatalf("expected capped limit 100, got %d", repo.lastGetTransactionsArg.limit)
	}
}

func TestGetStatement_Execute_AccountNotFound(t *testing.T) {
	repo := &statementRepositoryMock{getByIDErr: domain.ErrAccountNotFound}
	uc := NewGetStatement(repo)

	result, err := uc.Execute(context.Background(), GetStatementInput{AccountID: uuid.New()})

	if !errors.Is(err, domain.ErrAccountNotFound) {
		t.Fatalf("expected error %v, got %v", domain.ErrAccountNotFound, err)
	}

	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}

	if repo.getTransactionsCalls != 0 {
		t.Fatalf("expected GetTransactions not to be called, got %d", repo.getTransactionsCalls)
	}
}

func TestGetStatement_Execute_GetTransactionsAccountNotFoundPropagates(t *testing.T) {
	repo := &statementRepositoryMock{getTransactionsErr: domain.ErrAccountNotFound}
	uc := NewGetStatement(repo)
	accountID := uuid.New()
	customerID := uuid.New()
	repo.account = &domain.Account{ID: accountID, CustomerID: customerID}

	result, err := uc.Execute(context.Background(), GetStatementInput{User: testCustomerUser(customerID), AccountID: accountID})

	if !errors.Is(err, domain.ErrAccountNotFound) {
		t.Fatalf("expected error %v, got %v", domain.ErrAccountNotFound, err)
	}

	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
}

func TestGetStatement_Execute_Success(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	accountID := uuid.New()
	customerID := uuid.New()
	txID := uuid.New()
	refID := uuid.New()
	cursorID := uuid.New()
	from := now.Add(-24 * time.Hour)
	to := now

	repo := &statementRepositoryMock{
		account: &domain.Account{ID: accountID, CustomerID: customerID},
		getTransactionsResult: []domain.Transaction{
			{
				ID:           txID,
				AccountID:    accountID,
				Type:         domain.TransactionTransferIn,
				Amount:       250,
				BalanceAfter: 1250,
				ReferenceID:  &refID,
				CreatedAt:    now,
			},
		},
	}
	uc := NewGetStatement(repo)

	result, err := uc.Execute(context.Background(), GetStatementInput{
		User:      testCustomerUser(customerID),
		AccountID: accountID,
		Limit:     1,
		Cursor:    &to,
		CursorID:  &cursorID,
		From:      &from,
		To:        &to,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result to be non-nil")
	}

	if result.AccountID != accountID.String() {
		t.Fatalf("expected account id %q, got %q", accountID.String(), result.AccountID)
	}

	if len(result.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result.Items))
	}

	if result.Items[0].TransactionID != txID.String() {
		t.Fatalf("expected transaction id %q, got %q", txID.String(), result.Items[0].TransactionID)
	}

	if result.Items[0].BalanceAfter != 1250 {
		t.Fatalf("expected balance_after 1250, got %d", result.Items[0].BalanceAfter)
	}

	if result.Items[0].ReferenceID == nil || *result.Items[0].ReferenceID != refID.String() {
		t.Fatalf("expected reference id %q, got %v", refID.String(), result.Items[0].ReferenceID)
	}

	if result.NextCursor == nil {
		t.Fatal("expected next cursor to be non-nil")
	}

	if result.NextCursor.ID != txID.String() {
		t.Fatalf("expected next cursor id %q, got %q", txID.String(), result.NextCursor.ID)
	}
}

func TestGetStatement_Execute_ForbiddenForDifferentCustomer(t *testing.T) {
	accountID := uuid.New()
	repo := &statementRepositoryMock{
		account: &domain.Account{ID: accountID, CustomerID: uuid.New()},
	}
	uc := NewGetStatement(repo)

	result, err := uc.Execute(context.Background(), GetStatementInput{
		User:      testCustomerUser(uuid.New()),
		AccountID: accountID,
	})

	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected error %v, got %v", domain.ErrForbidden, err)
	}

	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}

	if repo.getTransactionsCalls != 0 {
		t.Fatalf("expected GetTransactions not to be called, got %d", repo.getTransactionsCalls)
	}
}
