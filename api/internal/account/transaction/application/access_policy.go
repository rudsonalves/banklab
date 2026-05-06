package application

import (
	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/transaction/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

// CanAccessCustomer checks if the authenticated user has access to the
// specified customer.
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
// specified account.
func CanAccessAccount(user *authdomain.AuthenticatedUser, account *domain.Account) bool {
	if account == nil {
		return false
	}

	return CanAccessCustomer(user, account.CustomerID)
}
