package application

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	accountdomain "github.com/seu-usuario/bank-api/internal/account/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

func TestCanAccessAccount_SameCustomer(t *testing.T) {
	customerID := uuid.New()
	userCustomerID := customerID.String()
	user := &AuthenticatedUser{
		UserID:     "user-1",
		Role:       authdomain.RoleCustomer,
		CustomerID: &userCustomerID,
	}
	account := &accountdomain.Account{CustomerID: customerID}

	if !CanAccessAccount(user, account) {
		t.Fatal("expected access to be allowed")
	}
}

func TestCanAccessAccount_DifferentCustomer(t *testing.T) {
	userCustomerID := uuid.NewString()
	user := &AuthenticatedUser{
		UserID:     "user-1",
		Role:       authdomain.RoleCustomer,
		CustomerID: &userCustomerID,
	}
	account := &accountdomain.Account{CustomerID: uuid.New()}

	if CanAccessAccount(user, account) {
		t.Fatal("expected access to be denied")
	}
}

func TestCanAccessAccount_AdminOverride(t *testing.T) {
	user := &AuthenticatedUser{
		UserID: "admin-1",
		Role:   authdomain.RoleAdmin,
	}
	account := &accountdomain.Account{CustomerID: uuid.New()}

	if !CanAccessAccount(user, account) {
		t.Fatal("expected admin access to be allowed")
	}
}

func TestCanAccessAccount_NilCustomerDenied(t *testing.T) {
	user := &AuthenticatedUser{
		UserID: "user-1",
		Role:   authdomain.RoleCustomer,
	}
	account := &accountdomain.Account{CustomerID: uuid.New()}

	if CanAccessAccount(user, account) {
		t.Fatal("expected access to be denied")
	}
}

func TestCanAccessAccount_AccountWithoutCustomerDenied(t *testing.T) {
	userCustomerID := uuid.NewString()
	user := &AuthenticatedUser{
		UserID:     "user-1",
		Role:       authdomain.RoleCustomer,
		CustomerID: &userCustomerID,
	}
	account := &accountdomain.Account{}

	if CanAccessAccount(user, account) {
		t.Fatal("expected access to be denied")
	}
}

func TestRequireAccountAccess_ReturnsForbidden(t *testing.T) {
	userCustomerID := uuid.NewString()
	user := &AuthenticatedUser{
		UserID:     "user-1",
		Role:       authdomain.RoleCustomer,
		CustomerID: &userCustomerID,
	}
	account := &accountdomain.Account{CustomerID: uuid.New()}

	err := RequireAccountAccess(user, account)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected error %v, got %v", ErrForbidden, err)
	}
}
