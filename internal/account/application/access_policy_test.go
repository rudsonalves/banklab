package application

import (
	"testing"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

func TestCanAccessCustomer_SameCustomer(t *testing.T) {
	customerID := uuid.New()
	user := &authdomain.AuthenticatedUser{
		UserID:     "user-1",
		Role:       authdomain.RoleCustomer,
		CustomerID: &customerID,
	}

	if !CanAccessCustomer(user, customerID) {
		t.Fatal("expected access to be allowed")
	}
}

func TestCanAccessCustomer_DifferentCustomer(t *testing.T) {
	userCustomerID := uuid.New()
	user := &authdomain.AuthenticatedUser{
		UserID:     "user-1",
		Role:       authdomain.RoleCustomer,
		CustomerID: &userCustomerID,
	}

	if CanAccessCustomer(user, uuid.New()) {
		t.Fatal("expected access to be denied")
	}
}

func TestCanAccessCustomer_AdminOverride(t *testing.T) {
	user := &authdomain.AuthenticatedUser{
		UserID: "admin-1",
		Role:   authdomain.RoleAdmin,
	}

	if !CanAccessCustomer(user, uuid.New()) {
		t.Fatal("expected admin access to be allowed")
	}
}

func TestCanAccessCustomer_NilCustomerDenied(t *testing.T) {
	user := &authdomain.AuthenticatedUser{
		UserID: "user-1",
		Role:   authdomain.RoleCustomer,
	}

	if CanAccessCustomer(user, uuid.New()) {
		t.Fatal("expected access to be denied")
	}
}

func TestCanAccessCustomer_UnknownCustomerDenied(t *testing.T) {
	customerID := uuid.New()
	user := &authdomain.AuthenticatedUser{
		UserID:     "user-1",
		Role:       authdomain.RoleCustomer,
		CustomerID: &customerID,
	}

	if CanAccessCustomer(user, uuid.Nil) {
		t.Fatal("expected access to be denied")
	}
}

func TestCanAccessAccount_FromOwnerCustomerID(t *testing.T) {
	customerID := uuid.New()
	user := &authdomain.AuthenticatedUser{
		UserID:     "user-1",
		Role:       authdomain.RoleCustomer,
		CustomerID: &customerID,
	}
	account := &domain.Account{CustomerID: customerID}

	if !CanAccessAccount(user, account) {
		t.Fatal("expected access to be allowed")
	}
}

func TestCanAccessAccount_NilAccountDenied(t *testing.T) {
	user := &authdomain.AuthenticatedUser{UserID: "user-1", Role: authdomain.RoleCustomer}
	if CanAccessAccount(user, nil) {
		t.Fatal("expected access to be denied")
	}
}
