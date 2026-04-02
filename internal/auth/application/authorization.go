package application

import (
	"errors"

	"github.com/google/uuid"
	accountdomain "github.com/seu-usuario/bank-api/internal/account/domain"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

var ErrForbidden = errors.New("forbidden")

func CanAccessAccount(user *AuthenticatedUser, account *accountdomain.Account) bool {
	if user == nil || account == nil {
		return false
	}

	if user.Role == domain.RoleAdmin {
		return true
	}

	if user.CustomerID == nil || account.CustomerID == uuid.Nil {
		return false
	}

	customerID, err := uuid.Parse(*user.CustomerID)
	if err != nil {
		return false
	}

	return customerID == account.CustomerID
}

func RequireAccountAccess(user *AuthenticatedUser, account *accountdomain.Account) error {
	if CanAccessAccount(user, account) {
		return nil
	}

	return ErrForbidden
}
