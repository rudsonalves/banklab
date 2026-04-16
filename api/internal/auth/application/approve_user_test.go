package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	accountdomain "github.com/seu-usuario/bank-api/internal/account/domain"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

// approveUserRepoMock implements domain.UserRepository for ApproveUser tests.
type approveUserRepoMock struct {
	findByIDForUpdateUser *domain.User
	findByIDForUpdateErr  error
	updateStatusErr       error
}

func (m *approveUserRepoMock) Create(ctx context.Context, user *domain.User) error {
	return nil
}

func (m *approveUserRepoMock) UpdateStatus(ctx context.Context, userID uuid.UUID, status domain.UserStatus) error {
	return m.updateStatusErr
}

func (m *approveUserRepoMock) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return m.findByIDForUpdateUser, m.findByIDForUpdateErr
}

func (m *approveUserRepoMock) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}

func (m *approveUserRepoMock) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return nil, nil
}

func (m *approveUserRepoMock) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return false, nil
}

// approveTransactorMock implements domain.Transactor for ApproveUser tests.
type approveTransactorMock struct{}

func (m *approveTransactorMock) RunInTx(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

// approveAccountRepoMock implements accountdomain.AccountRepository for ApproveUser tests.
type approveAccountRepoMock struct {
	nextAccountNumberValue string
	nextAccountNumberErr   error
	createErr              error
}

type approveCustomerRepoMock struct {
	existsValue bool
	existsErr   error
}

func (m *approveCustomerRepoMock) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return m.existsValue, m.existsErr
}

func (m *approveAccountRepoMock) Create(ctx context.Context, account *accountdomain.Account) error {
	return m.createErr
}

func (m *approveAccountRepoMock) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
	return false, nil
}

func (m *approveAccountRepoMock) NextAccountNumber(ctx context.Context) (string, error) {
	return m.nextAccountNumberValue, m.nextAccountNumberErr
}

func (m *approveAccountRepoMock) GetByID(ctx context.Context, id uuid.UUID) (*accountdomain.Account, error) {
	return nil, nil
}

func (m *approveAccountRepoMock) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*accountdomain.Account, error) {
	return nil, nil
}

func (m *approveAccountRepoMock) GetTransactions(
	ctx context.Context,
	accountID uuid.UUID,
	limit int,
	cursorTime *time.Time,
	cursorID *uuid.UUID,
	from *time.Time,
	to *time.Time,
) ([]accountdomain.Transaction, error) {
	return nil, nil
}

func (m *approveAccountRepoMock) IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return 0, nil
}

func (m *approveAccountRepoMock) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return 0, nil
}

func (m *approveAccountRepoMock) CreateTransaction(ctx context.Context, tx *accountdomain.Transaction) error {
	return nil
}

func (m *approveAccountRepoMock) GetOperationByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*accountdomain.Operation, error) {
	return nil, nil
}

func (m *approveAccountRepoMock) CreateOperation(ctx context.Context, op *accountdomain.Operation) error {
	return nil
}

func (m *approveAccountRepoMock) BeginTx(ctx context.Context) (accountdomain.Tx, error) {
	return nil, nil
}

func (m *approveAccountRepoMock) WithTransaction(ctx context.Context, fn func(tx accountdomain.Tx) error) error {
	return nil
}

// helpers

func newPendingUser() *domain.User {
	customerID := uuid.New()
	return &domain.User{
		ID:         uuid.New(),
		Email:      "user@example.com",
		Role:       domain.RoleCustomer,
		CustomerID: &customerID,
		Status:     domain.UserStatusPending,
	}
}

func newApproveUseCase(
	userRepo domain.UserRepository,
	accountRepo accountdomain.AccountRepository,
	customerRepo accountdomain.CustomerRepository,
) *ApproveUserUseCase {
	return NewApproveUserUseCase(userRepo, accountRepo, customerRepo, &approveTransactorMock{})
}

// Tests

func TestApproveUserUseCase_Execute_Success(t *testing.T) {
	user := newPendingUser()
	userRepo := &approveUserRepoMock{findByIDForUpdateUser: user}
	accountRepo := &approveAccountRepoMock{nextAccountNumberValue: "00000001"}
	customerRepo := &approveCustomerRepoMock{existsValue: true}
	uc := newApproveUseCase(userRepo, accountRepo, customerRepo)

	output, err := uc.Execute(context.Background(), ApproveUserInput{UserID: user.ID})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output == nil {
		t.Fatal("expected non-nil output")
	}
	if output.UserID != user.ID {
		t.Errorf("expected UserID %v, got %v", user.ID, output.UserID)
	}
	if output.AccountID == uuid.Nil {
		t.Error("expected non-nil AccountID")
	}
}

func TestApproveUserUseCase_Execute_UserNotFound(t *testing.T) {
	userRepo := &approveUserRepoMock{findByIDForUpdateUser: nil}
	accountRepo := &approveAccountRepoMock{}
	customerRepo := &approveCustomerRepoMock{existsValue: true}
	uc := newApproveUseCase(userRepo, accountRepo, customerRepo)

	_, err := uc.Execute(context.Background(), ApproveUserInput{UserID: uuid.New()})

	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestApproveUserUseCase_Execute_UserAlreadyActive(t *testing.T) {
	user := newPendingUser()
	user.Status = domain.UserStatusActive
	userRepo := &approveUserRepoMock{findByIDForUpdateUser: user}
	accountRepo := &approveAccountRepoMock{}
	customerRepo := &approveCustomerRepoMock{existsValue: true}
	uc := newApproveUseCase(userRepo, accountRepo, customerRepo)

	_, err := uc.Execute(context.Background(), ApproveUserInput{UserID: user.ID})

	if !errors.Is(err, domain.ErrUserAlreadyActive) {
		t.Errorf("expected ErrUserAlreadyActive, got %v", err)
	}
}

func TestApproveUserUseCase_Execute_UserNoCustomerID(t *testing.T) {
	user := newPendingUser()
	user.CustomerID = nil
	userRepo := &approveUserRepoMock{findByIDForUpdateUser: user}
	accountRepo := &approveAccountRepoMock{}
	customerRepo := &approveCustomerRepoMock{existsValue: true}
	uc := newApproveUseCase(userRepo, accountRepo, customerRepo)

	_, err := uc.Execute(context.Background(), ApproveUserInput{UserID: user.ID})

	if !errors.Is(err, domain.ErrInvalidUserState) {
		t.Errorf("expected ErrInvalidUserState, got %v", err)
	}
}

func TestApproveUserUseCase_Execute_AccountCreationFails(t *testing.T) {
	user := newPendingUser()
	accountErr := errors.New("db error")
	userRepo := &approveUserRepoMock{findByIDForUpdateUser: user}
	accountRepo := &approveAccountRepoMock{
		nextAccountNumberValue: "00000001",
		createErr:              accountErr,
	}
	customerRepo := &approveCustomerRepoMock{existsValue: true}
	uc := newApproveUseCase(userRepo, accountRepo, customerRepo)

	_, err := uc.Execute(context.Background(), ApproveUserInput{UserID: user.ID})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestApproveUserUseCase_Execute_CustomerNotFound(t *testing.T) {
	user := newPendingUser()
	userRepo := &approveUserRepoMock{findByIDForUpdateUser: user}
	accountRepo := &approveAccountRepoMock{nextAccountNumberValue: "00000001"}
	customerRepo := &approveCustomerRepoMock{existsValue: false}
	uc := newApproveUseCase(userRepo, accountRepo, customerRepo)

	_, err := uc.Execute(context.Background(), ApproveUserInput{UserID: user.ID})

	if !errors.Is(err, accountdomain.ErrCustomerNotFound) {
		t.Errorf("expected ErrCustomerNotFound, got %v", err)
	}
}
