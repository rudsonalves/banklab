package application

import (
	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

// CanAccessCustomer checks if the authenticated user has access to the
// specified customer. It returns true if the user is an admin or if the
// user's associated customer ID matches the provided customer ID.
func CanAccessCustomer(user *authdomain.AuthenticatedUser, customerID uuid.UUID) bool {
	if user == nil {
		return false
	}

	if user.Role == authdomain.RoleAdmin {
		return true
	}

	if user.CustomerID == nil || customerID == uuid.Nil {
		return false
	}

	return *user.CustomerID == customerID
}

// CanAccessAccount checks if the authenticated user has access to the
// specified account. It returns true if the user has access to the customer
// associated with the account.
func CanAccessAccount(user *authdomain.AuthenticatedUser, account *domain.Account) bool {
	if account == nil {
		return false
	}

	return CanAccessCustomer(user, account.CustomerID)
}
