package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	"github.com/seu-usuario/bank-api/internal/security/domain"
)

type transactionPasswordRepositoryMock struct {
	createCalls      int
	createErr        error
	createdPassword  *domain.TransactionPassword
	findByUserID     *domain.TransactionPassword
	findByUserIDErr  error
	findByUserIDUser uuid.UUID
}

func (m *transactionPasswordRepositoryMock) Create(ctx context.Context, password *domain.TransactionPassword) error {
	m.createCalls++
	m.createdPassword = password
	if m.createErr == nil && password.ID == uuid.Nil {
		password.ID = uuid.New()
	}
	return m.createErr
}

func (m *transactionPasswordRepositoryMock) FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.TransactionPassword, error) {
	m.findByUserIDUser = userID
	return m.findByUserID, m.findByUserIDErr
}

func (m *transactionPasswordRepositoryMock) SaveValidationState(ctx context.Context, password *domain.TransactionPassword) error {
	return nil
}

func (m *transactionPasswordRepositoryMock) UpdatePasswordHash(
	ctx context.Context,
	id uuid.UUID,
	passwordHash string,
	changedAt time.Time,
	updatedAt time.Time,
) error {
	return nil
}

type transactionPasswordUserRepositoryMock struct {
	findByIDCalls int
	findByIDValue *authdomain.User
	findByIDErr   error
	findByIDID    uuid.UUID
}

func (m *transactionPasswordUserRepositoryMock) Create(ctx context.Context, user *authdomain.User) error {
	return nil
}

func (m *transactionPasswordUserRepositoryMock) UpdateStatus(ctx context.Context, userID uuid.UUID, status authdomain.UserStatus) error {
	return nil
}

func (m *transactionPasswordUserRepositoryMock) FindByEmail(ctx context.Context, email string) (*authdomain.User, error) {
	return nil, nil
}

func (m *transactionPasswordUserRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (*authdomain.User, error) {
	m.findByIDCalls++
	m.findByIDID = id
	return m.findByIDValue, m.findByIDErr
}

func (m *transactionPasswordUserRepositoryMock) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return false, nil
}

func (m *transactionPasswordUserRepositoryMock) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	return false, nil
}

func (m *transactionPasswordUserRepositoryMock) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*authdomain.User, error) {
	return nil, nil
}

type transactionPasswordHasherMock struct {
	hashCalls    int
	hashPassword string
	hashValue    string
	hashErr      error
}

func (m *transactionPasswordHasherMock) Hash(password string) (string, error) {
	m.hashCalls++
	m.hashPassword = password
	return m.hashValue, m.hashErr
}

func (m *transactionPasswordHasherMock) Compare(hash, password string) bool {
	return hash == password
}

func TestCreateTransactionPasswordUseCase_Execute_Success(t *testing.T) {
	userID := uuid.New()
	now := time.Date(2026, 5, 28, 10, 0, 0, 0, time.UTC)
	passwordRepo := &transactionPasswordRepositoryMock{}
	userRepo := &transactionPasswordUserRepositoryMock{
		findByIDValue: activeUser(userID),
	}
	hasher := &transactionPasswordHasherMock{hashValue: "hashed-pin"}
	uc := NewCreateTransactionPasswordUseCase(passwordRepo, userRepo, hasher)
	uc.now = func() time.Time { return now }

	output, err := uc.Execute(context.Background(), CreateTransactionPasswordInput{
		User:                            authenticatedUser(userID),
		TransactionPassword:             "123456",
		TransactionPasswordConfirmation: "123456",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output == nil {
		t.Fatal("expected output, got nil")
	}
	if output.UserID != userID.String() {
		t.Fatalf("expected user id %q, got %q", userID.String(), output.UserID)
	}
	if output.Status != string(domain.TransactionPasswordActive) {
		t.Fatalf("expected status %q, got %q", domain.TransactionPasswordActive, output.Status)
	}
	if !output.CreatedAt.Equal(now) {
		t.Fatalf("expected created_at %v, got %v", now, output.CreatedAt)
	}
	if userRepo.findByIDCalls != 1 || userRepo.findByIDID != userID {
		t.Fatalf("expected FindByID once with %q, got calls=%d id=%q", userID, userRepo.findByIDCalls, userRepo.findByIDID)
	}
	if passwordRepo.findByUserIDUser != userID {
		t.Fatalf("expected FindByUserID with %q, got %q", userID, passwordRepo.findByUserIDUser)
	}
	if hasher.hashCalls != 1 || hasher.hashPassword != "123456" {
		t.Fatalf("expected Hash once with PIN, got calls=%d password=%q", hasher.hashCalls, hasher.hashPassword)
	}
	if passwordRepo.createCalls != 1 {
		t.Fatalf("expected Create once, got %d", passwordRepo.createCalls)
	}
	if passwordRepo.createdPassword == nil {
		t.Fatal("expected created password to be captured")
	}
	if passwordRepo.createdPassword.ID == uuid.Nil {
		t.Fatal("expected repository mock to populate generated id")
	}
	if passwordRepo.createdPassword.PasswordHash != "hashed-pin" {
		t.Fatalf("expected hash %q, got %q", "hashed-pin", passwordRepo.createdPassword.PasswordHash)
	}
}

func TestCreateTransactionPasswordUseCase_Execute_MissingAuthenticatedUser(t *testing.T) {
	uc := NewCreateTransactionPasswordUseCase(
		&transactionPasswordRepositoryMock{},
		&transactionPasswordUserRepositoryMock{},
		&transactionPasswordHasherMock{},
	)

	output, err := uc.Execute(context.Background(), CreateTransactionPasswordInput{
		TransactionPassword:             "123456",
		TransactionPasswordConfirmation: "123456",
	})

	if !errors.Is(err, authdomain.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
}

func TestCreateTransactionPasswordUseCase_Execute_InvalidPIN(t *testing.T) {
	passwordRepo := &transactionPasswordRepositoryMock{}
	userRepo := &transactionPasswordUserRepositoryMock{}
	hasher := &transactionPasswordHasherMock{}
	uc := NewCreateTransactionPasswordUseCase(passwordRepo, userRepo, hasher)

	output, err := uc.Execute(context.Background(), CreateTransactionPasswordInput{
		User:                            authenticatedUser(uuid.New()),
		TransactionPassword:             "12345a",
		TransactionPasswordConfirmation: "12345a",
	})

	if !errors.Is(err, domain.ErrInvalidTransactionPasswordPIN) {
		t.Fatalf("expected ErrInvalidTransactionPasswordPIN, got %v", err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
	if userRepo.findByIDCalls != 0 {
		t.Fatalf("expected FindByID not to be called, got %d", userRepo.findByIDCalls)
	}
	if hasher.hashCalls != 0 {
		t.Fatalf("expected Hash not to be called, got %d", hasher.hashCalls)
	}
	if passwordRepo.createCalls != 0 {
		t.Fatalf("expected Create not to be called, got %d", passwordRepo.createCalls)
	}
}

func TestCreateTransactionPasswordUseCase_Execute_ConfirmationMismatch(t *testing.T) {
	passwordRepo := &transactionPasswordRepositoryMock{}
	userRepo := &transactionPasswordUserRepositoryMock{}
	hasher := &transactionPasswordHasherMock{}
	uc := NewCreateTransactionPasswordUseCase(passwordRepo, userRepo, hasher)

	output, err := uc.Execute(context.Background(), CreateTransactionPasswordInput{
		User:                            authenticatedUser(uuid.New()),
		TransactionPassword:             "123456",
		TransactionPasswordConfirmation: "654321",
	})

	if !errors.Is(err, domain.ErrInvalidTransactionPasswordPIN) {
		t.Fatalf("expected ErrInvalidTransactionPasswordPIN, got %v", err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
	if userRepo.findByIDCalls != 0 {
		t.Fatalf("expected FindByID not to be called, got %d", userRepo.findByIDCalls)
	}
}

func TestCreateTransactionPasswordUseCase_Execute_UserNotActive(t *testing.T) {
	userID := uuid.New()
	passwordRepo := &transactionPasswordRepositoryMock{}
	userRepo := &transactionPasswordUserRepositoryMock{
		findByIDValue: &authdomain.User{
			ID:     userID,
			Email:  "pending@example.com",
			Role:   authdomain.RoleCustomer,
			Status: authdomain.UserStatusPending,
		},
	}
	hasher := &transactionPasswordHasherMock{}
	uc := NewCreateTransactionPasswordUseCase(passwordRepo, userRepo, hasher)

	output, err := uc.Execute(context.Background(), CreateTransactionPasswordInput{
		User:                            authenticatedUser(userID),
		TransactionPassword:             "123456",
		TransactionPasswordConfirmation: "123456",
	})

	if !errors.Is(err, authdomain.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
	if hasher.hashCalls != 0 {
		t.Fatalf("expected Hash not to be called, got %d", hasher.hashCalls)
	}
	if passwordRepo.createCalls != 0 {
		t.Fatalf("expected Create not to be called, got %d", passwordRepo.createCalls)
	}
}

func TestCreateTransactionPasswordUseCase_Execute_UserNotFound(t *testing.T) {
	userID := uuid.New()
	uc := NewCreateTransactionPasswordUseCase(
		&transactionPasswordRepositoryMock{},
		&transactionPasswordUserRepositoryMock{},
		&transactionPasswordHasherMock{},
	)

	output, err := uc.Execute(context.Background(), CreateTransactionPasswordInput{
		User:                            authenticatedUser(userID),
		TransactionPassword:             "123456",
		TransactionPasswordConfirmation: "123456",
	})

	if !errors.Is(err, authdomain.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
}

func TestCreateTransactionPasswordUseCase_Execute_AlreadySet(t *testing.T) {
	userID := uuid.New()
	now := time.Date(2026, 5, 28, 10, 0, 0, 0, time.UTC)
	existing, err := domain.NewTransactionPassword(userID, "stored-hash", now)
	if err != nil {
		t.Fatalf("expected no error creating existing password, got %v", err)
	}
	passwordRepo := &transactionPasswordRepositoryMock{findByUserID: existing}
	userRepo := &transactionPasswordUserRepositoryMock{findByIDValue: activeUser(userID)}
	hasher := &transactionPasswordHasherMock{}
	uc := NewCreateTransactionPasswordUseCase(passwordRepo, userRepo, hasher)

	output, err := uc.Execute(context.Background(), CreateTransactionPasswordInput{
		User:                            authenticatedUser(userID),
		TransactionPassword:             "123456",
		TransactionPasswordConfirmation: "123456",
	})

	if !errors.Is(err, domain.ErrTransactionPasswordAlreadySet) {
		t.Fatalf("expected ErrTransactionPasswordAlreadySet, got %v", err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
	if hasher.hashCalls != 0 {
		t.Fatalf("expected Hash not to be called, got %d", hasher.hashCalls)
	}
	if passwordRepo.createCalls != 0 {
		t.Fatalf("expected Create not to be called, got %d", passwordRepo.createCalls)
	}
}

func TestCreateTransactionPasswordUseCase_Execute_HashFailure(t *testing.T) {
	userID := uuid.New()
	expectedErr := errors.New("hash failed")
	passwordRepo := &transactionPasswordRepositoryMock{}
	userRepo := &transactionPasswordUserRepositoryMock{findByIDValue: activeUser(userID)}
	hasher := &transactionPasswordHasherMock{hashErr: expectedErr}
	uc := NewCreateTransactionPasswordUseCase(passwordRepo, userRepo, hasher)

	output, err := uc.Execute(context.Background(), CreateTransactionPasswordInput{
		User:                            authenticatedUser(userID),
		TransactionPassword:             "123456",
		TransactionPasswordConfirmation: "123456",
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
	if passwordRepo.createCalls != 0 {
		t.Fatalf("expected Create not to be called, got %d", passwordRepo.createCalls)
	}
}

func TestCreateTransactionPasswordUseCase_Execute_CreateFailure(t *testing.T) {
	userID := uuid.New()
	expectedErr := errors.New("insert failed")
	passwordRepo := &transactionPasswordRepositoryMock{createErr: expectedErr}
	userRepo := &transactionPasswordUserRepositoryMock{findByIDValue: activeUser(userID)}
	hasher := &transactionPasswordHasherMock{hashValue: "hashed-pin"}
	uc := NewCreateTransactionPasswordUseCase(passwordRepo, userRepo, hasher)

	output, err := uc.Execute(context.Background(), CreateTransactionPasswordInput{
		User:                            authenticatedUser(userID),
		TransactionPassword:             "123456",
		TransactionPasswordConfirmation: "123456",
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
	if passwordRepo.createCalls != 1 {
		t.Fatalf("expected Create to be called once, got %d", passwordRepo.createCalls)
	}
}

func activeUser(userID uuid.UUID) *authdomain.User {
	return &authdomain.User{
		ID:     userID,
		Email:  "user@example.com",
		Role:   authdomain.RoleCustomer,
		Status: authdomain.UserStatusActive,
	}
}

func authenticatedUser(userID uuid.UUID) *authdomain.AuthenticatedUser {
	return &authdomain.AuthenticatedUser{
		UserID: userID,
		Role:   authdomain.RoleCustomer,
	}
}
