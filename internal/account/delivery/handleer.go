package delivery

import (
	"context"

	"github.com/seu-usuario/bank-api/internal/account/application"
	"github.com/seu-usuario/bank-api/internal/account/domain"
)

type createAccountUseCase interface {
	Execute(ctx context.Context, input application.CreateAccountInput) (*domain.Account, error)
}

type depositUseCase interface {
	Execute(ctx context.Context, input application.DepositInput) (*domain.Account, error)
}

type withdrawUseCase interface {
	Execute(ctx context.Context, input application.WithdrawInput) (*domain.Account, error)
}

type Handler struct {
	createAccount createAccountUseCase
	deposit       depositUseCase
	withdraw      withdrawUseCase
}

func New(createAccount createAccountUseCase, deposit depositUseCase, withdraw withdrawUseCase) *Handler {
	return &Handler{
		createAccount: createAccount,
		deposit:       deposit,
		withdraw:      withdraw,
	}
}
