package application

import (
	"context"

	"github.com/seu-usuario/bank-api/internal/account/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

type ListAccountsInput struct {
	User *authdomain.AuthenticatedUser
}

type ListAccounts struct {
	repo domain.AccountRepository
}

func NewListAccounts(repo domain.AccountRepository) *ListAccounts {
	return &ListAccounts{repo: repo}
}

func (uc *ListAccounts) Execute(ctx context.Context, input ListAccountsInput) ([]domain.Account, error) {
	if !CanListOwnAccounts(input.User) {
		return nil, domain.ErrForbidden
	}

	return uc.repo.ListByCustomerID(ctx, *input.User.CustomerID)
}
