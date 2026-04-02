package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/domain"
)

type accountRepositoryMock struct {
	createCalls             int
	createErr               error
	existsByCustomerIDCalls int
	nextAccountNumberCalls  int
	nextAccountNumberValue  string
	nextAccountNumberErr    error
}

func (m *accountRepositoryMock) Create(ctx context.Context, account *domain.Account) error {
	m.createCalls++
	return m.createErr
}

func (m *accountRepositoryMock) CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
	return nil
}

func (m *accountRepositoryMock) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
	m.existsByCustomerIDCalls++
	return false, nil
}

func (m *accountRepositoryMock) NextAccountNumber(ctx context.Context) (string, error) {
	m.nextAccountNumberCalls++
	return m.nextAccountNumberValue, m.nextAccountNumberErr
}

func (m *accountRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return nil, nil
}

func (m *accountRepositoryMock) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return nil, nil
}

func (m *accountRepositoryMock) GetTransactions(
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

func (m *accountRepositoryMock) UpdateBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return 0, nil
}

func (m *accountRepositoryMock) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) error {
	return nil
}

func (m *accountRepositoryMock) BeginTx(ctx context.Context) (domain.Tx, error) {
	return nil, nil
}

type customerRepositoryMock struct {
	existsCalls int
	existsValue bool
	existsErr   error
}

func (m *customerRepositoryMock) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	m.existsCalls++
	return m.existsValue, m.existsErr
}

func TestCreateAccount_Execute_InvalidCustomerID(t *testing.T) {
	accountRepo := &accountRepositoryMock{}
	customerRepo := &customerRepositoryMock{}
	useCase := NewCreateAccount(accountRepo, customerRepo)

	account, err := useCase.Execute(context.Background(), CreateAccountInput{CustomerID: uuid.Nil})

	if !errors.Is(err, domain.ErrInvalidData) {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidData, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if accountRepo.createCalls != 0 {
		t.Fatalf("expected Create not to be called, got %d calls", accountRepo.createCalls)
	}

	if accountRepo.existsByCustomerIDCalls != 0 {
		t.Fatalf("expected ExistsByCustomerID not to be called, got %d calls", accountRepo.existsByCustomerIDCalls)
	}

	if accountRepo.nextAccountNumberCalls != 0 {
		t.Fatalf("expected NextAccountNumber not to be called, got %d calls", accountRepo.nextAccountNumberCalls)
	}

	if customerRepo.existsCalls != 0 {
		t.Fatalf("expected Exists not to be called, got %d calls", customerRepo.existsCalls)
	}
}

func TestCreateAccount_Execute_CustomerNotFound(t *testing.T) {
	accountRepo := &accountRepositoryMock{}
	customerRepo := &customerRepositoryMock{}
	useCase := NewCreateAccount(accountRepo, customerRepo)
	customerID := uuid.New()

	account, err := useCase.Execute(context.Background(), CreateAccountInput{User: testCustomerUser(customerID), CustomerID: customerID})

	if !errors.Is(err, domain.ErrCustomerNotFound) {
		t.Fatalf("expected error %v, got %v", domain.ErrCustomerNotFound, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if accountRepo.createCalls != 0 {
		t.Fatalf("expected Create not to be called, got %d calls", accountRepo.createCalls)
	}

	if accountRepo.existsByCustomerIDCalls != 0 {
		t.Fatalf("expected ExistsByCustomerID not to be called, got %d calls", accountRepo.existsByCustomerIDCalls)
	}

	if accountRepo.nextAccountNumberCalls != 0 {
		t.Fatalf("expected NextAccountNumber not to be called, got %d calls", accountRepo.nextAccountNumberCalls)
	}
}

func TestCreateAccount_Execute_CustomerExistsReturnsError(t *testing.T) {
	expectedErr := errors.New("database unavailable")
	accountRepo := &accountRepositoryMock{}
	customerRepo := &customerRepositoryMock{existsErr: expectedErr}
	useCase := NewCreateAccount(accountRepo, customerRepo)
	customerID := uuid.New()

	account, err := useCase.Execute(context.Background(), CreateAccountInput{User: testCustomerUser(customerID), CustomerID: customerID})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if accountRepo.createCalls != 0 {
		t.Fatalf("expected Create not to be called, got %d calls", accountRepo.createCalls)
	}

	if accountRepo.existsByCustomerIDCalls != 0 {
		t.Fatalf("expected ExistsByCustomerID not to be called, got %d calls", accountRepo.existsByCustomerIDCalls)
	}

	if accountRepo.nextAccountNumberCalls != 0 {
		t.Fatalf("expected NextAccountNumber not to be called, got %d calls", accountRepo.nextAccountNumberCalls)
	}
}

func TestCreateAccount_Execute_NextAccountNumberReturnsError(t *testing.T) {
	expectedErr := errors.New("sequence unavailable")
	accountRepo := &accountRepositoryMock{nextAccountNumberErr: expectedErr}
	customerRepo := &customerRepositoryMock{existsValue: true}
	useCase := NewCreateAccount(accountRepo, customerRepo)
	customerID := uuid.New()

	account, err := useCase.Execute(context.Background(), CreateAccountInput{User: testCustomerUser(customerID), CustomerID: customerID})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if accountRepo.createCalls != 0 {
		t.Fatalf("expected Create not to be called, got %d calls", accountRepo.createCalls)
	}
}

func TestCreateAccount_Execute_CreateReturnsError(t *testing.T) {
	expectedErr := errors.New("insert failed")
	accountRepo := &accountRepositoryMock{
		nextAccountNumberValue: "123456",
		createErr:              expectedErr,
	}
	customerRepo := &customerRepositoryMock{existsValue: true}
	useCase := NewCreateAccount(accountRepo, customerRepo)
	customerID := uuid.New()

	account, err := useCase.Execute(context.Background(), CreateAccountInput{User: testCustomerUser(customerID), CustomerID: customerID})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if customerRepo.existsCalls != 1 {
		t.Fatalf("expected Exists to be called once, got %d calls", customerRepo.existsCalls)
	}

	if accountRepo.nextAccountNumberCalls != 1 {
		t.Fatalf("expected NextAccountNumber to be called once, got %d calls", accountRepo.nextAccountNumberCalls)
	}

	if accountRepo.createCalls != 1 {
		t.Fatalf("expected Create to be called once, got %d calls", accountRepo.createCalls)
	}
}

func TestCreateAccount_Execute_Success(t *testing.T) {
	inputCustomerID := uuid.New()
	accountRepo := &accountRepositoryMock{nextAccountNumberValue: "12345678"}
	customerRepo := &customerRepositoryMock{existsValue: true}
	useCase := NewCreateAccount(accountRepo, customerRepo)

	account, err := useCase.Execute(context.Background(), CreateAccountInput{User: testCustomerUser(inputCustomerID), CustomerID: inputCustomerID})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if account == nil {
		t.Fatal("expected account to be non-nil")
	}

	if account.CustomerID != inputCustomerID {
		t.Fatalf("expected CustomerID %v, got %v", inputCustomerID, account.CustomerID)
	}

	if account.Number != "12345678" {
		t.Fatalf("expected Number %q, got %q", "12345678", account.Number)
	}

	if account.Balance != 0 {
		t.Fatalf("expected Balance %d, got %d", 0, account.Balance)
	}

	if account.Status != "active" {
		t.Fatalf("expected Status %q, got %q", "active", account.Status)
	}
}

func TestCreateAccount_Execute_InteractionCountsOnSuccess(t *testing.T) {
	accountRepo := &accountRepositoryMock{nextAccountNumberValue: "12345678"}
	customerRepo := &customerRepositoryMock{existsValue: true}
	useCase := NewCreateAccount(accountRepo, customerRepo)
	customerID := uuid.New()

	account, err := useCase.Execute(context.Background(), CreateAccountInput{User: testCustomerUser(customerID), CustomerID: customerID})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if account == nil {
		t.Fatal("expected account to be non-nil")
	}

	if customerRepo.existsCalls != 1 {
		t.Fatalf("expected Exists to be called once, got %d calls", customerRepo.existsCalls)
	}

	if accountRepo.nextAccountNumberCalls != 1 {
		t.Fatalf("expected NextAccountNumber to be called once, got %d calls", accountRepo.nextAccountNumberCalls)
	}

	if accountRepo.createCalls != 1 {
		t.Fatalf("expected Create to be called once, got %d calls", accountRepo.createCalls)
	}
}

func TestCreateAccount_Execute_DoesNotCallCreateWhenCustomerNotFound(t *testing.T) {
	accountRepo := &accountRepositoryMock{}
	customerRepo := &customerRepositoryMock{existsValue: false}
	useCase := NewCreateAccount(accountRepo, customerRepo)
	customerID := uuid.New()

	_, _ = useCase.Execute(context.Background(), CreateAccountInput{User: testCustomerUser(customerID), CustomerID: customerID})

	if accountRepo.createCalls != 0 {
		t.Fatalf("Create must not be called when customer does not exist, got %d calls", accountRepo.createCalls)
	}
}

func TestCreateAccount_Execute_ForbiddenForDifferentCustomer(t *testing.T) {
	accountRepo := &accountRepositoryMock{}
	customerRepo := &customerRepositoryMock{}
	useCase := NewCreateAccount(accountRepo, customerRepo)

	account, err := useCase.Execute(context.Background(), CreateAccountInput{
		User:       testCustomerUser(uuid.New()),
		CustomerID: uuid.New(),
	})

	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected error %v, got %v", domain.ErrForbidden, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if customerRepo.existsCalls != 0 {
		t.Fatalf("expected Exists not to be called, got %d calls", customerRepo.existsCalls)
	}

	if accountRepo.createCalls != 0 {
		t.Fatalf("expected Create not to be called, got %d calls", accountRepo.createCalls)
	}
}

func TestCreateAccount_Execute_AdminCanCreateForAnyCustomer(t *testing.T) {
	customerID := uuid.New()
	accountRepo := &accountRepositoryMock{nextAccountNumberValue: "12345678"}
	customerRepo := &customerRepositoryMock{existsValue: true}
	useCase := NewCreateAccount(accountRepo, customerRepo)

	account, err := useCase.Execute(context.Background(), CreateAccountInput{
		User:       testAdminUser(),
		CustomerID: customerID,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if account == nil || account.CustomerID != customerID {
		t.Fatalf("expected account for customer %v, got %+v", customerID, account)
	}

	if accountRepo.createCalls != 1 {
		t.Fatalf("expected Create to be called once, got %d calls", accountRepo.createCalls)
	}

	if customerRepo.existsCalls != 1 {
		t.Fatalf("expected Exists to be called once, got %d calls", customerRepo.existsCalls)
	}
}
