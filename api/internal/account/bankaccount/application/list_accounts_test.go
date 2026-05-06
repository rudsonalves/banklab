package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/bankaccount/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

type listAccountsRepositoryMock struct {
	listByCustomerIDCalls int
	listByCustomerIDValue []domain.Account
	listByCustomerIDErr   error
}

func (m *listAccountsRepositoryMock) Create(ctx context.Context, account *domain.Account) error {
	return nil
}

func (m *listAccountsRepositoryMock) ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]domain.Account, error) {
	m.listByCustomerIDCalls++
	if m.listByCustomerIDErr != nil {
		return nil, m.listByCustomerIDErr
	}

	return m.listByCustomerIDValue, nil
}

func (m *listAccountsRepositoryMock) CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
	return nil
}

func (m *listAccountsRepositoryMock) GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*domain.Transaction, error) {
	return nil, nil
}

func (m *listAccountsRepositoryMock) GetTransactionByReference(ctx context.Context, accountID uuid.UUID, referenceID uuid.UUID, typeName domain.TransactionType) (*domain.Transaction, error) {
	return nil, nil
}

func (m *listAccountsRepositoryMock) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
	return false, nil
}

func (m *listAccountsRepositoryMock) NextAccountNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (m *listAccountsRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return nil, nil
}

func (m *listAccountsRepositoryMock) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return nil, nil
}

func (m *listAccountsRepositoryMock) GetTransactions(
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

func (m *listAccountsRepositoryMock) IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return 0, nil
}

func (m *listAccountsRepositoryMock) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return 0, nil
}

func (m *listAccountsRepositoryMock) BeginTx(ctx context.Context) (domain.Tx, error) {
	return nil, nil
}

func (m *listAccountsRepositoryMock) WithTransaction(ctx context.Context, fn func(tx domain.Tx) error) error {
	return errors.New("transactions are not used in this test")
}

func TestListAccounts_Execute_ForbiddenWhenUserMissing(t *testing.T) {
	repo := &listAccountsRepositoryMock{}
	uc := NewListAccounts(repo)

	accounts, err := uc.Execute(context.Background(), ListAccountsInput{})

	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected error %v, got %v", domain.ErrForbidden, err)
	}

	if accounts != nil {
		t.Fatalf("expected nil accounts, got %+v", accounts)
	}

	if repo.listByCustomerIDCalls != 0 {
		t.Fatalf("expected ListByCustomerID not to be called, got %d calls", repo.listByCustomerIDCalls)
	}
}

func TestListAccounts_Execute_ForbiddenWhenCustomerIDMissing(t *testing.T) {
	repo := &listAccountsRepositoryMock{}
	uc := NewListAccounts(repo)

	accounts, err := uc.Execute(context.Background(), ListAccountsInput{
		User: &authdomain.AuthenticatedUser{UserID: uuid.New(), Role: authdomain.RoleAdmin},
	})

	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected error %v, got %v", domain.ErrForbidden, err)
	}

	if accounts != nil {
		t.Fatalf("expected nil accounts, got %+v", accounts)
	}

	if repo.listByCustomerIDCalls != 0 {
		t.Fatalf("expected ListByCustomerID not to be called, got %d calls", repo.listByCustomerIDCalls)
	}
}

func TestListAccounts_Execute_Success(t *testing.T) {
	customerID := uuid.New()
	expected := []domain.Account{
		{
			ID:         uuid.New(),
			CustomerID: customerID,
			Number:     "10000001",
			Branch:     "0001",
			Status:     domain.AccountActive,
		},
	}
	repo := &listAccountsRepositoryMock{listByCustomerIDValue: expected}
	uc := NewListAccounts(repo)

	accounts, err := uc.Execute(context.Background(), ListAccountsInput{
		User: &authdomain.AuthenticatedUser{
			UserID:     uuid.New(),
			Role:       authdomain.RoleCustomer,
			CustomerID: &customerID,
		},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.listByCustomerIDCalls != 1 {
		t.Fatalf("expected ListByCustomerID to be called once, got %d", repo.listByCustomerIDCalls)
	}

	if len(accounts) != 1 {
		t.Fatalf("expected 1 account, got %d", len(accounts))
	}

	if accounts[0].ID != expected[0].ID {
		t.Fatalf("expected account id %v, got %v", expected[0].ID, accounts[0].ID)
	}
}
