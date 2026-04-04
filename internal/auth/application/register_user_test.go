package application

import (
	"context"
	"errors"
	"testing"

	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

type userRepositoryMock struct {
	existsByEmailCalls int
	existsByEmailValue bool
	existsByEmailErr   error
	createCalls        int
	createErr          error
	createdUser        *domain.User
}

func (m *userRepositoryMock) Create(ctx context.Context, user *domain.User) error {
	m.createCalls++
	m.createdUser = user
	return m.createErr
}

func (m *userRepositoryMock) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}

func (m *userRepositoryMock) FindByID(ctx context.Context, id string) (*domain.User, error) {
	return nil, nil
}

func (m *userRepositoryMock) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	m.existsByEmailCalls++
	return m.existsByEmailValue, m.existsByEmailErr
}

type passwordHasherMock struct {
	hashCalls int
	hashValue string
	hashErr   error
}

func (m *passwordHasherMock) Hash(password string) (string, error) {
	m.hashCalls++
	if m.hashErr != nil {
		return "", m.hashErr
	}
	return m.hashValue, nil
}

func (m *passwordHasherMock) Compare(hash string, password string) error {
	return nil
}

func TestRegisterUserUseCase_Execute_Success(t *testing.T) {
	userRepo := &userRepositoryMock{}
	hasher := &passwordHasherMock{hashValue: "hashed-password"}
	useCase := NewRegisterUserUseCase(userRepo, hasher)

	output, err := useCase.Execute(context.Background(), RegisterUserInput{
		Email:    "  USER@Example.com ",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output == nil {
		t.Fatal("expected output to be non-nil")
	}

	if output.ID == "" {
		t.Fatal("expected output ID to be set")
	}

	if output.Email != "user@example.com" {
		t.Fatalf("expected normalized email, got %q", output.Email)
	}

	if output.Role != string(domain.RoleCustomer) {
		t.Fatalf("expected role %q, got %q", domain.RoleCustomer, output.Role)
	}

	if output.CustomerID != nil {
		t.Fatalf("expected customer ID to be nil, got %v", *output.CustomerID)
	}

	if userRepo.existsByEmailCalls != 1 {
		t.Fatalf("expected ExistsByEmail to be called once, got %d", userRepo.existsByEmailCalls)
	}

	if hasher.hashCalls != 1 {
		t.Fatalf("expected Hash to be called once, got %d", hasher.hashCalls)
	}

	if userRepo.createCalls != 1 {
		t.Fatalf("expected Create to be called once, got %d", userRepo.createCalls)
	}

	if userRepo.createdUser == nil {
		t.Fatal("expected created user to be captured")
	}

	if userRepo.createdUser.PasswordHash != "hashed-password" {
		t.Fatalf("expected hashed password to be persisted, got %q", userRepo.createdUser.PasswordHash)
	}

	if userRepo.createdUser.Role != domain.RoleCustomer {
		t.Fatalf("expected persisted role %q, got %q", domain.RoleCustomer, userRepo.createdUser.Role)
	}

	if userRepo.createdUser.CustomerID != nil {
		t.Fatalf("expected persisted customer ID to be nil, got %v", *userRepo.createdUser.CustomerID)
	}

	if userRepo.createdUser.CreatedAt.IsZero() {
		t.Fatal("expected created_at to be set")
	}

	if userRepo.createdUser.UpdatedAt.IsZero() {
		t.Fatal("expected updated_at to be set")
	}

	if !userRepo.createdUser.CreatedAt.Equal(userRepo.createdUser.UpdatedAt) {
		t.Fatal("expected created_at and updated_at to match on creation")
	}
}

func TestRegisterUserUseCase_Execute_DuplicateEmail(t *testing.T) {
	userRepo := &userRepositoryMock{existsByEmailValue: true}
	hasher := &passwordHasherMock{hashValue: "hashed-password"}
	useCase := NewRegisterUserUseCase(userRepo, hasher)

	output, err := useCase.Execute(context.Background(), RegisterUserInput{
		Email:    "user@example.com",
		Password: "password123",
	})

	if !errors.Is(err, ErrEmailAlreadyExists) {
		t.Fatalf("expected error %v, got %v", ErrEmailAlreadyExists, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if hasher.hashCalls != 0 {
		t.Fatalf("expected Hash not to be called, got %d calls", hasher.hashCalls)
	}

	if userRepo.createCalls != 0 {
		t.Fatalf("expected Create not to be called, got %d calls", userRepo.createCalls)
	}
}

func TestRegisterUserUseCase_Execute_InvalidEmail(t *testing.T) {
	userRepo := &userRepositoryMock{}
	hasher := &passwordHasherMock{}
	useCase := NewRegisterUserUseCase(userRepo, hasher)

	output, err := useCase.Execute(context.Background(), RegisterUserInput{
		Email:    "invalid-email",
		Password: "password123",
	})

	if !errors.Is(err, ErrInvalidEmail) {
		t.Fatalf("expected error %v, got %v", ErrInvalidEmail, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if userRepo.existsByEmailCalls != 0 {
		t.Fatalf("expected ExistsByEmail not to be called, got %d calls", userRepo.existsByEmailCalls)
	}

	if hasher.hashCalls != 0 {
		t.Fatalf("expected Hash not to be called, got %d calls", hasher.hashCalls)
	}

	if userRepo.createCalls != 0 {
		t.Fatalf("expected Create not to be called, got %d calls", userRepo.createCalls)
	}
}

func TestRegisterUserUseCase_Execute_InvalidPassword(t *testing.T) {
	userRepo := &userRepositoryMock{}
	hasher := &passwordHasherMock{}
	useCase := NewRegisterUserUseCase(userRepo, hasher)

	output, err := useCase.Execute(context.Background(), RegisterUserInput{
		Email:    "user@example.com",
		Password: "short",
	})

	if !errors.Is(err, ErrInvalidPassword) {
		t.Fatalf("expected error %v, got %v", ErrInvalidPassword, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if userRepo.existsByEmailCalls != 0 {
		t.Fatalf("expected ExistsByEmail not to be called, got %d calls", userRepo.existsByEmailCalls)
	}

	if hasher.hashCalls != 0 {
		t.Fatalf("expected Hash not to be called, got %d calls", hasher.hashCalls)
	}

	if userRepo.createCalls != 0 {
		t.Fatalf("expected Create not to be called, got %d calls", userRepo.createCalls)
	}
}

func TestRegisterUserUseCase_Execute_HashingFailure(t *testing.T) {
	expectedErr := errors.New("hash unavailable")
	userRepo := &userRepositoryMock{}
	hasher := &passwordHasherMock{hashErr: expectedErr}
	useCase := NewRegisterUserUseCase(userRepo, hasher)

	output, err := useCase.Execute(context.Background(), RegisterUserInput{
		Email:    "user@example.com",
		Password: "password123",
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if userRepo.createCalls != 0 {
		t.Fatalf("expected Create not to be called, got %d calls", userRepo.createCalls)
	}
}
