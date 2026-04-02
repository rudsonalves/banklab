package application

import (
	"errors"

	accountdomain "github.com/seu-usuario/bank-api/internal/account/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

var ErrForbidden = errors.New("forbidden")

func CanAccessAccount(user *authdomain.AuthenticatedUser, account *accountdomain.Account) bool {
	return authdomain.CanAccessAccount(user, account)
}

func RequireAccountAccess(user *authdomain.AuthenticatedUser, account *accountdomain.Account) error {
	if CanAccessAccount(user, account) {
		return nil
	}

	return ErrForbidden
}
