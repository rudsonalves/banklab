package delivery

import (
	"context"

	accountapp "github.com/seu-usuario/bank-api/internal/account/application/account"
	statementapp "github.com/seu-usuario/bank-api/internal/account/application/statement"
	"github.com/seu-usuario/bank-api/internal/account/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
)

type createAccountUseCase interface {
	Execute(ctx context.Context, input accountapp.CreateAccountInput) (*domain.Account, error)
}

type listAccountsUseCase interface {
	Execute(ctx context.Context, input accountapp.ListAccountsInput) ([]domain.Account, error)
}

type statementUseCase interface {
	Execute(ctx context.Context, input statementapp.GetStatementInput) (*statementapp.Statement, error)
}

type getBalanceUseCase interface {
	Execute(ctx context.Context, input accountapp.GetAccountBalanceInput) (*accountapp.AccountBalance, error)
}

type Handler struct {
	listAccounts  listAccountsUseCase
	createAccount createAccountUseCase
	statement     statementUseCase
	balance       getBalanceUseCase
}

// New creates a new instance of the account Handler with the provided use cases.
func New(
	listAccounts listAccountsUseCase,
	createAccount createAccountUseCase,
	statement statementUseCase,
	balance getBalanceUseCase,
) *Handler {
	return &Handler{
		listAccounts:  listAccounts,
		createAccount: createAccount,
		statement:     statement,
		balance:       balance,
	}
}

// RequireUser retrieves the authenticated user from the context.
// If no user is found, it returns an unauthorized error.
func RequireUser(ctx context.Context) (*authdomain.AuthenticatedUser, error) {
	user, ok := sharedauthctx.GetAuthenticatedUser(ctx)
	if !ok || user == nil {
		return nil, authdomain.ErrUnauthorized
	}

	return user, nil
}
