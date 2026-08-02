package authctx

import (
	"context"
	"testing"

	"github.com/google/uuid"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

func TestGetAuthenticatedUser_ContextContainsUser(t *testing.T) {
	customerID := uuid.New()
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	ctx := WithAuthenticatedUser(context.Background(), AuthenticatedUser{
		UserID:     userID,
		Role:       authdomain.RoleCustomer,
		CustomerID: &customerID,
	})

	user, ok := GetAuthenticatedUser(ctx)
	if !ok {
		t.Fatal("expected authenticated user in context")
	}

	if user == nil || user.UserID != userID {
		t.Fatalf("expected user id %q, got %#v", userID, user)
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

func TestOperationalSessionContext(t *testing.T) {
	userID := uuid.New()
	customerID := uuid.New()
	installationID := uuid.New()
	ctx := WithOperationalSession(context.Background(), OperationalSession{
		UserID:         userID,
		Role:           authdomain.RoleCustomer,
		CustomerID:     &customerID,
		InstallationID: &installationID,
	})

	session, ok := GetOperationalSession(ctx)
	if !ok {
		t.Fatal("expected operational session in context")
	}
	if session.UserID != userID {
		t.Fatalf("expected user id %q, got %q", userID, session.UserID)
	}
	if session.InstallationID == nil || *session.InstallationID != installationID {
		t.Fatalf("expected installation id %q, got %#v", installationID, session.InstallationID)
	}
}

func TestRequireOperationalSession_ReturnsErrorWhenMissing(t *testing.T) {
	session, err := RequireOperationalSession(context.Background())
	if err != ErrOperationalSessionNotFound {
		t.Fatalf("expected error %v, got %v", ErrOperationalSessionNotFound, err)
	}
	if session != nil {
		t.Fatalf("expected nil session, got %#v", session)
	}
}

func TestRestrictedSessionContext(t *testing.T) {
	userID := uuid.New()
	installationID := uuid.New()
	ctx := WithRestrictedSession(context.Background(), RestrictedSession{
		UserID:         userID,
		InstallationID: installationID,
		JTI:            "jti-1",
		Scope:          "installation.register",
	})

	session, ok := GetRestrictedSession(ctx)
	if !ok {
		t.Fatal("expected restricted session in context")
	}
	if session.UserID != userID || session.InstallationID != installationID {
		t.Fatalf("unexpected restricted session: %#v", session)
	}
	if session.JTI != "jti-1" || session.Scope != "installation.register" {
		t.Fatalf("unexpected restricted token metadata: %#v", session)
	}
}

func TestRequireRestrictedSession_ReturnsErrorWhenMissing(t *testing.T) {
	session, err := RequireRestrictedSession(context.Background())
	if err != ErrRestrictedSessionNotFound {
		t.Fatalf("expected error %v, got %v", ErrRestrictedSessionNotFound, err)
	}
	if session != nil {
		t.Fatalf("expected nil session, got %#v", session)
	}
}
