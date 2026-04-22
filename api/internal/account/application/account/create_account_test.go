package application

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
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

func (m *accountRepositoryMock) GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*domain.Transaction, error) {
	return nil, nil
}

func (m *accountRepositoryMock) GetTransactionByReference(ctx context.Context, accountID uuid.UUID, referenceID uuid.UUID, typeName domain.TransactionType) (*domain.Transaction, error) {
	return nil, nil
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

func (m *accountRepositoryMock) IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return 0, nil
}

func (m *accountRepositoryMock) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return 0, nil
}

func (m *accountRepositoryMock) BeginTx(ctx context.Context) (domain.Tx, error) {
	return nil, nil
}

func (m *accountRepositoryMock) WithTransaction(ctx context.Context, fn func(tx domain.Tx) error) error {
	return errors.New("transactions are not used in this test")
}

type customerRepositoryMock struct {
	existsCalls int
	existsValue bool
	existsErr   error
}

type userRepositoryMock struct {
	findByIDCalls int
	findByIDValue *authdomain.User
	findByIDErr   error
}

func (m *userRepositoryMock) Create(ctx context.Context, user *authdomain.User) error {
	return nil
}

func (m *userRepositoryMock) UpdateStatus(ctx context.Context, userID uuid.UUID, status authdomain.UserStatus) error {
	return nil
}

func (m *userRepositoryMock) FindByEmail(ctx context.Context, email string) (*authdomain.User, error) {
	return nil, nil
}

func (m *userRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (*authdomain.User, error) {
	m.findByIDCalls++
	return m.findByIDValue, m.findByIDErr
}

func (m *userRepositoryMock) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return false, nil
}

func (m *userRepositoryMock) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*authdomain.User, error) {
	return nil, nil
}

func (m *customerRepositoryMock) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	m.existsCalls++
	return m.existsValue, m.existsErr
}

func TestCreateAccount_Execute_MissingUserCustomerID(t *testing.T) {
	accountRepo := &accountRepositoryMock{}
	customerRepo := &customerRepositoryMock{}
	userRepo := &userRepositoryMock{}
	useCase := NewCreateAccount(accountRepo, customerRepo, userRepo)

	account, err := useCase.Execute(context.Background(), CreateAccountInput{User: testAdminUser()})

	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected error %v, got %v", domain.ErrForbidden, err)
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

	if userRepo.findByIDCalls != 0 {
		t.Fatalf("expected FindByID not to be called, got %d calls", userRepo.findByIDCalls)
	}
}

func TestCreateAccount_Execute_UserRepositoryNotConfigured(t *testing.T) {
	accountRepo := &accountRepositoryMock{}
	customerRepo := &customerRepositoryMock{}
	customerID := uuid.New()
	useCase := NewCreateAccount(accountRepo, customerRepo, nil)

	account, err := useCase.Execute(context.Background(), CreateAccountInput{User: testCustomerUser(customerID)})

	if err == nil || err.Error() != "user repository not configured" {
		t.Fatalf("expected error %q, got %v", "user repository not configured", err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if customerRepo.existsCalls != 0 {
		t.Fatalf("expected Exists not to be called, got %d calls", customerRepo.existsCalls)
	}

	if accountRepo.nextAccountNumberCalls != 0 {
		t.Fatalf("expected NextAccountNumber not to be called, got %d calls", accountRepo.nextAccountNumberCalls)
	}
}

func TestCreateAccount_Execute_CustomerNotFound(t *testing.T) {
	accountRepo := &accountRepositoryMock{}
	customerRepo := &customerRepositoryMock{}
	customerID := uuid.New()
	userRepo := &userRepositoryMock{findByIDValue: testUserWithStatus(customerID, authdomain.UserStatusActive)}
	useCase := NewCreateAccount(accountRepo, customerRepo, userRepo)

	account, err := useCase.Execute(context.Background(), CreateAccountInput{User: testCustomerUser(customerID)})

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

	if userRepo.findByIDCalls != 1 {
		t.Fatalf("expected FindByID to be called once, got %d calls", userRepo.findByIDCalls)
	}
}

func TestCreateAccount_Execute_CustomerExistsReturnsError(t *testing.T) {
	expectedErr := errors.New("database unavailable")
	accountRepo := &accountRepositoryMock{}
	customerRepo := &customerRepositoryMock{existsErr: expectedErr}
	customerID := uuid.New()
	userRepo := &userRepositoryMock{findByIDValue: testUserWithStatus(customerID, authdomain.UserStatusActive)}
	useCase := NewCreateAccount(accountRepo, customerRepo, userRepo)

	account, err := useCase.Execute(context.Background(), CreateAccountInput{User: testCustomerUser(customerID)})

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

func TestCreateAccount_Execute_UserLookupReturnsError(t *testing.T) {
	expectedErr := errors.New("database unavailable")
	accountRepo := &accountRepositoryMock{}
	customerRepo := &customerRepositoryMock{}
	customerID := uuid.New()
	userRepo := &userRepositoryMock{findByIDErr: expectedErr}
	useCase := NewCreateAccount(accountRepo, customerRepo, userRepo)

	account, err := useCase.Execute(context.Background(), CreateAccountInput{User: testCustomerUser(customerID)})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if customerRepo.existsCalls != 0 {
		t.Fatalf("expected Exists not to be called, got %d calls", customerRepo.existsCalls)
	}

	if accountRepo.nextAccountNumberCalls != 0 {
		t.Fatalf("expected NextAccountNumber not to be called, got %d calls", accountRepo.nextAccountNumberCalls)
	}
}

func TestCreateAccount_Execute_UserNotFound(t *testing.T) {
	accountRepo := &accountRepositoryMock{}
	customerRepo := &customerRepositoryMock{}
	customerID := uuid.New()
	userRepo := &userRepositoryMock{}
	useCase := NewCreateAccount(accountRepo, customerRepo, userRepo)

	account, err := useCase.Execute(context.Background(), CreateAccountInput{User: testCustomerUser(customerID)})

	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected error %v, got %v", domain.ErrForbidden, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if customerRepo.existsCalls != 0 {
		t.Fatalf("expected Exists not to be called, got %d calls", customerRepo.existsCalls)
	}

	if accountRepo.nextAccountNumberCalls != 0 {
		t.Fatalf("expected NextAccountNumber not to be called, got %d calls", accountRepo.nextAccountNumberCalls)
	}
}

func TestCreateAccount_Execute_PendingUserIsForbidden(t *testing.T) {
	accountRepo := &accountRepositoryMock{}
	customerRepo := &customerRepositoryMock{}
	customerID := uuid.New()
	userRepo := &userRepositoryMock{findByIDValue: testUserWithStatus(customerID, authdomain.UserStatusPending)}
	useCase := NewCreateAccount(accountRepo, customerRepo, userRepo)

	account, err := useCase.Execute(context.Background(), CreateAccountInput{User: testCustomerUser(customerID)})

	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected error %v, got %v", domain.ErrForbidden, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if customerRepo.existsCalls != 0 {
		t.Fatalf("expected Exists not to be called, got %d calls", customerRepo.existsCalls)
	}

	if accountRepo.nextAccountNumberCalls != 0 {
		t.Fatalf("expected NextAccountNumber not to be called, got %d calls", accountRepo.nextAccountNumberCalls)
	}
}

func TestCreateAccount_Execute_BlockedUserIsForbidden(t *testing.T) {
	accountRepo := &accountRepositoryMock{}
	customerRepo := &customerRepositoryMock{}
	customerID := uuid.New()
	userRepo := &userRepositoryMock{findByIDValue: testUserWithStatus(customerID, authdomain.UserStatusBlocked)}
	useCase := NewCreateAccount(accountRepo, customerRepo, userRepo)

	account, err := useCase.Execute(context.Background(), CreateAccountInput{User: testCustomerUser(customerID)})

	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected error %v, got %v", domain.ErrForbidden, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if customerRepo.existsCalls != 0 {
		t.Fatalf("expected Exists not to be called, got %d calls", customerRepo.existsCalls)
	}

	if accountRepo.nextAccountNumberCalls != 0 {
		t.Fatalf("expected NextAccountNumber not to be called, got %d calls", accountRepo.nextAccountNumberCalls)
	}
}

func TestCreateAccount_Execute_NextAccountNumberReturnsError(t *testing.T) {
	expectedErr := errors.New("sequence unavailable")
	accountRepo := &accountRepositoryMock{nextAccountNumberErr: expectedErr}
	customerRepo := &customerRepositoryMock{existsValue: true}
	customerID := uuid.New()
	userRepo := &userRepositoryMock{findByIDValue: testUserWithStatus(customerID, authdomain.UserStatusActive)}
	useCase := NewCreateAccount(accountRepo, customerRepo, userRepo)

	account, err := useCase.Execute(context.Background(), CreateAccountInput{User: testCustomerUser(customerID)})

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
	customerID := uuid.New()
	userRepo := &userRepositoryMock{findByIDValue: testUserWithStatus(customerID, authdomain.UserStatusActive)}
	useCase := NewCreateAccount(accountRepo, customerRepo, userRepo)

	account, err := useCase.Execute(context.Background(), CreateAccountInput{User: testCustomerUser(customerID)})

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
	userRepo := &userRepositoryMock{findByIDValue: testUserWithStatus(inputCustomerID, authdomain.UserStatusActive)}
	useCase := NewCreateAccount(accountRepo, customerRepo, userRepo)

	account, err := useCase.Execute(context.Background(), CreateAccountInput{User: testCustomerUser(inputCustomerID)})

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
	customerID := uuid.New()
	userRepo := &userRepositoryMock{findByIDValue: testUserWithStatus(customerID, authdomain.UserStatusActive)}
	useCase := NewCreateAccount(accountRepo, customerRepo, userRepo)

	account, err := useCase.Execute(context.Background(), CreateAccountInput{User: testCustomerUser(customerID)})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if account == nil {
		t.Fatal("expected account to be non-nil")
	}

	if customerRepo.existsCalls != 1 {
		t.Fatalf("expected Exists to be called once, got %d calls", customerRepo.existsCalls)
	}

	if userRepo.findByIDCalls != 1 {
		t.Fatalf("expected FindByID to be called once, got %d calls", userRepo.findByIDCalls)
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
	customerID := uuid.New()
	userRepo := &userRepositoryMock{findByIDValue: testUserWithStatus(customerID, authdomain.UserStatusActive)}
	useCase := NewCreateAccount(accountRepo, customerRepo, userRepo)

	_, _ = useCase.Execute(context.Background(), CreateAccountInput{User: testCustomerUser(customerID)})

	if accountRepo.createCalls != 0 {
		t.Fatalf("Create must not be called when customer does not exist, got %d calls", accountRepo.createCalls)
	}
}

func TestCreateAccount_Execute_InvalidWhenUserIsNil(t *testing.T) {
	accountRepo := &accountRepositoryMock{}
	customerRepo := &customerRepositoryMock{}
	userRepo := &userRepositoryMock{}
	useCase := NewCreateAccount(accountRepo, customerRepo, userRepo)

	account, err := useCase.Execute(context.Background(), CreateAccountInput{
		User: nil,
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

	if userRepo.findByIDCalls != 0 {
		t.Fatalf("expected FindByID not to be called, got %d calls", userRepo.findByIDCalls)
	}
}

func TestCreateAccount_Execute_AdminWithoutCustomerIDIsInvalid(t *testing.T) {
	accountRepo := &accountRepositoryMock{nextAccountNumberValue: "12345678"}
	customerRepo := &customerRepositoryMock{existsValue: true}
	userRepo := &userRepositoryMock{}
	useCase := NewCreateAccount(accountRepo, customerRepo, userRepo)

	account, err := useCase.Execute(context.Background(), CreateAccountInput{
		User: testAdminUser(),
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

	if userRepo.findByIDCalls != 0 {
		t.Fatalf("expected FindByID not to be called, got %d calls", userRepo.findByIDCalls)
	}
}

func TestCreateAccount_Execute_ZeroCustomerIDIsForbidden(t *testing.T) {
	accountRepo := &accountRepositoryMock{}
	customerRepo := &customerRepositoryMock{}
	userRepo := &userRepositoryMock{}
	useCase := NewCreateAccount(accountRepo, customerRepo, userRepo)

	zeroCustomerID := uuid.Nil
	account, err := useCase.Execute(context.Background(), CreateAccountInput{
		User: &authdomain.AuthenticatedUser{
			UserID:     uuid.New(),
			Role:       authdomain.RoleCustomer,
			CustomerID: &zeroCustomerID,
		},
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

func testUserWithStatus(customerID uuid.UUID, status authdomain.UserStatus) *authdomain.User {
	now := time.Now().UTC()
	return &authdomain.User{
		ID:         uuid.New(),
		Email:      fmt.Sprintf("%s@example.com", uuid.NewString()),
		Role:       authdomain.RoleCustomer,
		CustomerID: &customerID,
		Status:     status,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
