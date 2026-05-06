package application

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/bankaccount/domain"
)

type balanceRepositoryMock struct {
	getByIDCalls int
	getByIDErr   error
	account      *domain.Account
}

func (m *balanceRepositoryMock) Create(ctx context.Context, account *domain.Account) error {
	return nil
}

func (m *balanceRepositoryMock) ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]domain.Account, error) {
	return nil, nil
}

func (m *balanceRepositoryMock) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
	return false, nil
}

func (m *balanceRepositoryMock) NextAccountNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (m *balanceRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	m.getByIDCalls++
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	if m.account != nil {
		return m.account, nil
	}

	return &domain.Account{ID: id}, nil
}

func TestGetAccountBalance_Execute_InvalidAccountID(t *testing.T) {
	repo := &balanceRepositoryMock{}
	uc := NewGetAccountBalance(repo)

	result, err := uc.Execute(context.Background(), GetAccountBalanceInput{})

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

func TestGetAccountBalance_Execute_AccountNotFound(t *testing.T) {
	repo := &balanceRepositoryMock{getByIDErr: domain.ErrAccountNotFound}
	uc := NewGetAccountBalance(repo)

	result, err := uc.Execute(context.Background(), GetAccountBalanceInput{AccountID: uuid.New()})

	if !errors.Is(err, domain.ErrAccountNotFound) {
		t.Fatalf("expected error %v, got %v", domain.ErrAccountNotFound, err)
	}

	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
}

func TestGetAccountBalance_Execute_ForbiddenForDifferentCustomer(t *testing.T) {
	accountID := uuid.New()
	repo := &balanceRepositoryMock{
		account: &domain.Account{ID: accountID, CustomerID: uuid.New()},
	}
	uc := NewGetAccountBalance(repo)

	result, err := uc.Execute(context.Background(), GetAccountBalanceInput{
		User:      testCustomerUser(uuid.New()),
		AccountID: accountID,
	})

	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected error %v, got %v", domain.ErrForbidden, err)
	}

	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
}

func TestGetAccountBalance_Execute_Success(t *testing.T) {
	accountID := uuid.New()
	customerID := uuid.New()
	repo := &balanceRepositoryMock{
		account: &domain.Account{ID: accountID, CustomerID: customerID, Balance: 1250},
	}
	uc := NewGetAccountBalance(repo)

	result, err := uc.Execute(context.Background(), GetAccountBalanceInput{
		User:      testCustomerUser(customerID),
		AccountID: accountID,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result to be non-nil")
	}

	if result.AccountID != accountID {
		t.Fatalf("expected account id %q, got %q", accountID.String(), result.AccountID.String())
	}

	if result.Balance != 1250 {
		t.Fatalf("expected balance 1250, got %d", result.Balance)
	}
}
