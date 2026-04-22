package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

type GetAccountBalanceInput struct {
	User *authdomain.AuthenticatedUser

	AccountID uuid.UUID
}

type AccountBalance struct {
	AccountID uuid.UUID
	Balance   int64
}

type GetAccountBalance struct {
	repo domain.AccountRepository
}

func NewGetAccountBalance(repo domain.AccountRepository) *GetAccountBalance {
	return &GetAccountBalance{repo: repo}
}

func (uc *GetAccountBalance) Execute(ctx context.Context, input GetAccountBalanceInput) (*AccountBalance, error) {
	if input.AccountID == uuid.Nil {
		return nil, domain.ErrInvalidData
	}

	account, err := uc.repo.GetByID(ctx, input.AccountID)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, domain.ErrAccountNotFound
	}

	if !CanAccessAccount(input.User, account) {
		return nil, domain.ErrForbidden
	}

	return &AccountBalance{
		AccountID: account.ID,
		Balance:   account.Balance,
	}, nil
}
