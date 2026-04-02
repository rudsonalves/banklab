package delivery

import (
	"context"
	"testing"

	"github.com/seu-usuario/bank-api/internal/auth/application"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

func TestGetAuthenticatedUser_ContextContainsUser(t *testing.T) {
	customerID := "customer-1"
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

func TestMustGetAuthenticatedUser_PanicsWhenMissing(t *testing.T) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected panic when authenticated user is missing")
		}

		if recovered != "authenticated user not found in context" {
			t.Fatalf("expected panic %q, got %#v", "authenticated user not found in context", recovered)
		}
	}()

	_ = MustGetAuthenticatedUser(context.Background())
}
