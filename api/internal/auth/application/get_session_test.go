package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	accountdomain "github.com/seu-usuario/bank-api/internal/account/bankaccount/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	customerdomain "github.com/seu-usuario/bank-api/internal/customer/domain"
	securitydomain "github.com/seu-usuario/bank-api/internal/security/domain"
)

type sessionUserRepositoryMock struct {
	findByIDCalls int
	findByIDValue uuid.UUID
	findByIDUser  *authdomain.User
	findByIDErr   error
}

func (m *sessionUserRepositoryMock) Create(ctx context.Context, user *authdomain.User) error {
	return nil
}

func (m *sessionUserRepositoryMock) UpdateStatus(ctx context.Context, userID uuid.UUID, status authdomain.UserStatus) error {
	return nil
}

func (m *sessionUserRepositoryMock) FindByEmail(ctx context.Context, email string) (*authdomain.User, error) {
	return nil, nil
}

func (m *sessionUserRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (*authdomain.User, error) {
	m.findByIDCalls++
	m.findByIDValue = id
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	return m.findByIDUser, nil
}

func (m *sessionUserRepositoryMock) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return false, nil
}

func (m *sessionUserRepositoryMock) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	return false, nil
}

func (m *sessionUserRepositoryMock) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*authdomain.User, error) {
	return nil, nil
}

type sessionCustomerRepositoryMock struct {
	getByIDCalls int
	getByIDValue uuid.UUID
	profile      *customerdomain.CustomerProfile
	err          error
}

func (m *sessionCustomerRepositoryMock) Create(ctx context.Context, c *customerdomain.Customer) error {
	return nil
}

func (m *sessionCustomerRepositoryMock) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return false, nil
}

func (m *sessionCustomerRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*customerdomain.CustomerProfile, error) {
	m.getByIDCalls++
	m.getByIDValue = id
	if m.err != nil {
		return nil, m.err
	}
	return m.profile, nil
}

type sessionAccountRepositoryMock struct {
	listByCustomerIDCalls int
	listByCustomerIDValue uuid.UUID
	accounts              []accountdomain.Account
	err                   error
}

func (m *sessionAccountRepositoryMock) Create(ctx context.Context, account *accountdomain.Account) error {
	return nil
}

func (m *sessionAccountRepositoryMock) ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]accountdomain.Account, error) {
	m.listByCustomerIDCalls++
	m.listByCustomerIDValue = customerID
	if m.err != nil {
		return nil, m.err
	}
	return m.accounts, nil
}

func (m *sessionAccountRepositoryMock) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
	return false, nil
}

func (m *sessionAccountRepositoryMock) NextAccountNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (m *sessionAccountRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*accountdomain.Account, error) {
	return nil, nil
}

type sessionTransactionPasswordRepositoryMock struct {
	findByUserIDCalls int
	findByUserIDValue uuid.UUID
	password          *securitydomain.TransactionPassword
	err               error
}

func (m *sessionTransactionPasswordRepositoryMock) Create(ctx context.Context, password *securitydomain.TransactionPassword) error {
	return nil
}

func (m *sessionTransactionPasswordRepositoryMock) FindByUserID(ctx context.Context, userID uuid.UUID) (*securitydomain.TransactionPassword, error) {
	m.findByUserIDCalls++
	m.findByUserIDValue = userID
	if m.err != nil {
		return nil, m.err
	}
	return m.password, nil
}

func (m *sessionTransactionPasswordRepositoryMock) SaveValidationState(ctx context.Context, password *securitydomain.TransactionPassword) error {
	return nil
}

func (m *sessionTransactionPasswordRepositoryMock) UpdatePasswordHash(
	ctx context.Context,
	id uuid.UUID,
	passwordHash string,
	changedAt time.Time,
	updatedAt time.Time,
) error {
	return nil
}

func TestGetSessionUseCase_Execute_ActiveReadySession(t *testing.T) {
	userID := uuid.New()
	customerID := uuid.New()
	now := time.Date(2026, 6, 2, 10, 0, 0, 0, time.UTC)
	userRepo := &sessionUserRepositoryMock{
		findByIDUser: &authdomain.User{
			ID:         userID,
			Email:      "user@example.com",
			Phone:      "+5527999999999",
			Role:       authdomain.RoleCustomer,
			CustomerID: &customerID,
			Status:     authdomain.UserStatusActive,
		},
	}
	customerRepo := &sessionCustomerRepositoryMock{
		profile: &customerdomain.CustomerProfile{
			Customer: customerdomain.Customer{
				ID:        customerID,
				Name:      "Maria Silva",
				BirthDate: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
				CreatedAt: now,
			},
			CPF: "12345678901",
		},
	}
	accountRepo := &sessionAccountRepositoryMock{
		accounts: []accountdomain.Account{{
			ID:         uuid.New(),
			CustomerID: customerID,
			Status:     accountdomain.AccountActive,
		}},
	}
	passwordRepo := &sessionTransactionPasswordRepositoryMock{
		password: &securitydomain.TransactionPassword{
			ID:     uuid.New(),
			UserID: userID,
			Status: securitydomain.TransactionPasswordActive,
		},
	}
	useCase := NewGetSessionUseCase(userRepo, customerRepo, accountRepo, passwordRepo)
	ctx := WithAuthenticatedUser(context.Background(), AuthenticatedUser{
		UserID:     userID,
		Role:       authdomain.RoleCustomer,
		CustomerID: &customerID,
	})

	output, err := useCase.Execute(ctx)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.User.ID != userID {
		t.Fatalf("expected user id %q, got %q", userID, output.User.ID)
	}
	if output.User.Phone != "+5527999999999" {
		t.Fatalf("expected phone %q, got %q", "+5527999999999", output.User.Phone)
	}
	if output.Customer.ID != customerID {
		t.Fatalf("expected customer id %q, got %q", customerID, output.Customer.ID)
	}
	if output.Readiness.TransactionPasswordStatus != TransactionPasswordSessionStatusActive {
		t.Fatalf("expected transaction password status active, got %q", output.Readiness.TransactionPasswordStatus)
	}
}

func TestGetSessionUseCase_Execute_WithoutTransactionPassword(t *testing.T) {
	userID := uuid.New()
	customerID := uuid.New()
	userRepo := &sessionUserRepositoryMock{
		findByIDUser: &authdomain.User{
			ID:         userID,
			Email:      "user@example.com",
			Phone:      "+5527999999999",
			Role:       authdomain.RoleCustomer,
			CustomerID: &customerID,
			Status:     authdomain.UserStatusActive,
		},
	}
	customerRepo := &sessionCustomerRepositoryMock{
		profile: &customerdomain.CustomerProfile{
			Customer: customerdomain.Customer{ID: customerID},
			CPF:      "12345678901",
		},
	}
	accountRepo := &sessionAccountRepositoryMock{
		accounts: []accountdomain.Account{{Status: accountdomain.AccountActive}},
	}
	passwordRepo := &sessionTransactionPasswordRepositoryMock{}
	useCase := NewGetSessionUseCase(userRepo, customerRepo, accountRepo, passwordRepo)
	ctx := WithAuthenticatedUser(context.Background(), AuthenticatedUser{
		UserID:     userID,
		Role:       authdomain.RoleCustomer,
		CustomerID: &customerID,
	})

	output, err := useCase.Execute(ctx)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Readiness.TransactionPasswordStatus != TransactionPasswordSessionStatusNotSet {
		t.Fatalf("expected transaction password status not_set, got %q", output.Readiness.TransactionPasswordStatus)
	}
}

func TestGetSessionUseCase_Execute_MissingContext(t *testing.T) {
	useCase := NewGetSessionUseCase(
		&sessionUserRepositoryMock{},
		&sessionCustomerRepositoryMock{},
		&sessionAccountRepositoryMock{},
		&sessionTransactionPasswordRepositoryMock{},
	)

	output, err := useCase.Execute(context.Background())

	if !errors.Is(err, authdomain.ErrUnauthorized) {
		t.Fatalf("expected error %v, got %v", authdomain.ErrUnauthorized, err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
}

func TestGetSessionUseCase_Execute_CustomerUserWithoutCustomerID(t *testing.T) {
	userID := uuid.New()
	userRepo := &sessionUserRepositoryMock{
		findByIDUser: &authdomain.User{
			ID:     userID,
			Email:  "user@example.com",
			Role:   authdomain.RoleCustomer,
			Status: authdomain.UserStatusActive,
		},
	}
	customerRepo := &sessionCustomerRepositoryMock{}
	useCase := NewGetSessionUseCase(
		userRepo,
		customerRepo,
		&sessionAccountRepositoryMock{},
		&sessionTransactionPasswordRepositoryMock{},
	)
	ctx := WithAuthenticatedUser(context.Background(), AuthenticatedUser{
		UserID: userID,
		Role:   authdomain.RoleCustomer,
	})

	output, err := useCase.Execute(ctx)

	if !errors.Is(err, authdomain.ErrInvalidUserState) {
		t.Fatalf("expected error %v, got %v", authdomain.ErrInvalidUserState, err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
	if customerRepo.getByIDCalls != 0 {
		t.Fatalf("expected customer repository not to be called, got %d", customerRepo.getByIDCalls)
	}
}

func TestGetSessionUseCase_Execute_CustomerNotFound(t *testing.T) {
	userID := uuid.New()
	customerID := uuid.New()
	userRepo := &sessionUserRepositoryMock{
		findByIDUser: &authdomain.User{
			ID:         userID,
			Email:      "user@example.com",
			Role:       authdomain.RoleCustomer,
			CustomerID: &customerID,
			Status:     authdomain.UserStatusActive,
		},
	}
	useCase := NewGetSessionUseCase(
		userRepo,
		&sessionCustomerRepositoryMock{},
		&sessionAccountRepositoryMock{},
		&sessionTransactionPasswordRepositoryMock{},
	)
	ctx := WithAuthenticatedUser(context.Background(), AuthenticatedUser{
		UserID:     userID,
		Role:       authdomain.RoleCustomer,
		CustomerID: &customerID,
	})

	output, err := useCase.Execute(ctx)

	if !errors.Is(err, customerdomain.ErrNotFound) {
		t.Fatalf("expected error %v, got %v", customerdomain.ErrNotFound, err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
}
