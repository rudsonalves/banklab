package delivery

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/application"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

func TestGetAuthenticatedUser_ContextContainsUser(t *testing.T) {
	customerID := uuid.New()
	ctx := application.WithAuthenticatedUser(context.Background(), application.AuthenticatedUser{
		UserID:     "user-1",
		Role:       domain.RoleCustomer,
		CustomerID: &customerID,
	})

	user, ok := GetAuthenticatedUser(ctx)
	if !ok {
		t.Fatal("expected authenticated user in context")
	}

	if user == nil || user.UserID != "user-1" {
		t.Fatalf("expected user id %q, got %#v", "user-1", user)
	}
}

func TestRequireAuthenticatedUser_ReturnsErrorWhenMissing(t *testing.T) {
	user, err := RequireAuthenticatedUser(context.Background())
	if err == nil {
		t.Fatal("expected error when authenticated user is missing")
	}

	if user != nil {
		t.Fatalf("expected nil user, got %#v", user)
	}

	if err != ErrAuthenticatedUserNotFound {
		t.Fatalf("expected error %v, got %v", ErrAuthenticatedUserNotFound, err)
	}
}
