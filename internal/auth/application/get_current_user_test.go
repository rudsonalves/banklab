package application

import (
	"context"
	"errors"
	"testing"

	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

type currentUserRepositoryMock struct {
	findByIDCalls int
	findByIDValue string
	findByIDUser  *domain.User
	findByIDErr   error
}

func (m *currentUserRepositoryMock) Create(ctx context.Context, user *domain.User) error {
	return nil
}

func (m *currentUserRepositoryMock) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}

func (m *currentUserRepositoryMock) FindByID(ctx context.Context, id string) (*domain.User, error) {
	m.findByIDCalls++
	m.findByIDValue = id
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	return m.findByIDUser, nil
}

func (m *currentUserRepositoryMock) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return false, nil
}

func TestGetCurrentUserUseCase_Execute_Success(t *testing.T) {
	customerID := "customer-1"
	userRepo := &currentUserRepositoryMock{
		findByIDUser: &domain.User{
			ID:         "user-1",
			Email:      "user@example.com",
			Role:       domain.RoleCustomer,
			CustomerID: &customerID,
		},
	}
	useCase := NewGetCurrentUserUseCase(userRepo)
	ctx := WithAuthenticatedUser(context.Background(), AuthenticatedUser{
		UserID:     "user-1",
		Role:       domain.RoleAdmin,
		CustomerID: nil,
	})

	output, err := useCase.Execute(ctx)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output == nil {
		t.Fatal("expected output to be non-nil")
	}

	if output.ID != "user-1" {
		t.Fatalf("expected ID %q, got %q", "user-1", output.ID)
	}

	if output.Email != "user@example.com" {
		t.Fatalf("expected email %q, got %q", "user@example.com", output.Email)
	}

	if output.Role != string(domain.RoleCustomer) {
		t.Fatalf("expected role %q, got %q", domain.RoleCustomer, output.Role)
	}

	if output.CustomerID == nil || *output.CustomerID != customerID {
		t.Fatalf("expected customer ID %q, got %v", customerID, output.CustomerID)
	}

	if userRepo.findByIDCalls != 1 {
		t.Fatalf("expected FindByID to be called once, got %d", userRepo.findByIDCalls)
	}

	if userRepo.findByIDValue != "user-1" {
		t.Fatalf("expected FindByID to be called with %q, got %q", "user-1", userRepo.findByIDValue)
	}
}

func TestGetCurrentUserUseCase_Execute_MissingContext(t *testing.T) {
	userRepo := &currentUserRepositoryMock{}
	useCase := NewGetCurrentUserUseCase(userRepo)

	output, err := useCase.Execute(context.Background())

	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected error %v, got %v", ErrUnauthorized, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if userRepo.findByIDCalls != 0 {
		t.Fatalf("expected FindByID not to be called, got %d calls", userRepo.findByIDCalls)
	}
}

func TestGetCurrentUserUseCase_Execute_UserNotFound(t *testing.T) {
	userRepo := &currentUserRepositoryMock{}
	useCase := NewGetCurrentUserUseCase(userRepo)
	ctx := WithAuthenticatedUser(context.Background(), AuthenticatedUser{
		UserID: "user-1",
		Role:   domain.RoleCustomer,
	})

	output, err := useCase.Execute(ctx)

	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected error %v, got %v", ErrUnauthorized, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}
}

func TestGetCurrentUserUseCase_Execute_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database unavailable")
	userRepo := &currentUserRepositoryMock{findByIDErr: expectedErr}
	useCase := NewGetCurrentUserUseCase(userRepo)
	ctx := WithAuthenticatedUser(context.Background(), AuthenticatedUser{
		UserID: "user-1",
		Role:   domain.RoleCustomer,
	})

	output, err := useCase.Execute(ctx)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}
}