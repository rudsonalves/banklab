package delivery

import (
	"context"

	"github.com/seu-usuario/bank-api/internal/account/application"
	"github.com/seu-usuario/bank-api/internal/account/domain"
)

type createAccountUseCase interface {
	Execute(ctx context.Context, input application.CreateAccountInput) (*domain.Account, error)
}

type Handler struct {
	createAccount createAccountUseCase
}

func New(createAccount createAccountUseCase) *Handler {
	return &Handler{
		createAccount: createAccount,
	}
}
