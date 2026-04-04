package delivery

import (
	"context"

	"github.com/seu-usuario/bank-api/internal/account/application"
	"github.com/seu-usuario/bank-api/internal/account/domain"
	authdelivery "github.com/seu-usuario/bank-api/internal/auth/delivery"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
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

type transferUseCase interface {
	Execute(ctx context.Context, input application.TransferInput) (*application.TransferResult, error)
}

type statementUseCase interface {
	Execute(ctx context.Context, input application.GetStatementInput) (*application.Statement, error)
}

type Handler struct {
	createAccount createAccountUseCase
	deposit       depositUseCase
	withdraw      withdrawUseCase
	transfer      transferUseCase
	statement     statementUseCase
}

func New(
	createAccount createAccountUseCase,
	deposit depositUseCase,
	withdraw withdrawUseCase,
	transfer transferUseCase,
	statement statementUseCase,
) *Handler {
	return &Handler{
		createAccount: createAccount,
		deposit:       deposit,
		withdraw:      withdraw,
		transfer:      transfer,
		statement:     statement,
	}
}

func RequireUser(ctx context.Context) (*authdelivery.AuthenticatedUser, *sharederrors.AppError) {
	user, ok := authdelivery.GetAuthenticatedUser(ctx)
	if !ok || user == nil {
		return nil, sharederrors.ErrUnauthorized
	}

	return user, nil
}
