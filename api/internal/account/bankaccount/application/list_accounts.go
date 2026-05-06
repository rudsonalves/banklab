package application

import (
	"context"

	"github.com/seu-usuario/bank-api/internal/account/bankaccount/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

type ListAccountsInput struct {
	User *authdomain.AuthenticatedUser
}

type ListAccounts struct {
	repo domain.AccountRepository
}

// NewListAccounts creates a new instance of the ListAccounts use case with the
// provided account repository.
func NewListAccounts(repo domain.AccountRepository) *ListAccounts {
	return &ListAccounts{repo: repo}
}

// Execute retrieves a list of accounts associated with the authenticated user. It
// checks if the user has permission to list their own accounts and returns the
// accounts if access is granted. If the user does not have permission, it returns
// an appropriate error.
func (uc *ListAccounts) Execute(ctx context.Context, input ListAccountsInput) ([]domain.Account, error) {
	if !CanListOwnAccounts(input.User) {
		return nil, domain.ErrForbidden
	}

	return uc.repo.ListByCustomerID(ctx, *input.User.CustomerID)
}
