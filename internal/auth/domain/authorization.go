package domain

import (
	"github.com/google/uuid"
	accountdomain "github.com/seu-usuario/bank-api/internal/account/domain"
)

func CanAccessAccount(user *AuthenticatedUser, account *accountdomain.Account) bool {
	if user == nil || account == nil {
		return false
	}

	if user.Role == RoleAdmin {
		return true
	}

	if user.CustomerID == nil || account.CustomerID == uuid.Nil {
		return false
	}

	return *user.CustomerID == account.CustomerID
}
