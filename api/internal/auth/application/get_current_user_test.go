package application

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

type currentUserRepositoryMock struct {
	findByIDCalls int
	findByIDValue uuid.UUID
	findByIDUser  *domain.User
	findByIDErr   error
}

func (m *currentUserRepositoryMock) Create(ctx context.Context, user *domain.User) error {
	return nil
}

func (m *currentUserRepositoryMock) UpdateStatus(ctx context.Context, userID uuid.UUID, status domain.UserStatus) error {
	return nil
}

func (m *currentUserRepositoryMock) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}

func (m *currentUserRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
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
	customerID := uuid.New()
	testUserID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	userRepo := &currentUserRepositoryMock{
		findByIDUser: &domain.User{
			ID:         testUserID,
			Email:      "user@example.com",
			Role:       domain.RoleCustomer,
			CustomerID: &customerID,
		},
	}
	useCase := NewGetCurrentUserUseCase(userRepo)
	ctx := WithAuthenticatedUser(context.Background(), AuthenticatedUser{
		UserID:     testUserID,
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

	if output.ID != testUserID {
		t.Fatalf("expected ID %q, got %q", testUserID, output.ID)
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

	if userRepo.findByIDValue != testUserID {
		t.Fatalf("expected FindByID to be called with %q, got %q", testUserID, userRepo.findByIDValue)
	}
}

func TestGetCurrentUserUseCase_Execute_MissingContext(t *testing.T) {
	userRepo := &currentUserRepositoryMock{}
	useCase := NewGetCurrentUserUseCase(userRepo)

	output, err := useCase.Execute(context.Background())

	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("expected error %v, got %v", domain.ErrUnauthorized, err)
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
		UserID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Role:   domain.RoleCustomer,
	})

	output, err := useCase.Execute(ctx)

	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("expected error %v, got %v", domain.ErrUnauthorized, err)
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
		UserID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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
