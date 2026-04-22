package delivery

import (
	"context"

	accountapp "github.com/seu-usuario/bank-api/internal/account/application/account"
	statementapp "github.com/seu-usuario/bank-api/internal/account/application/statement"
	transactionapp "github.com/seu-usuario/bank-api/internal/account/application/transaction"
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

type depositUseCase interface {
	Execute(ctx context.Context, input transactionapp.DepositInput) (*domain.Account, error)
}

type withdrawUseCase interface {
	Execute(ctx context.Context, input transactionapp.WithdrawInput) (*domain.Account, error)
}

type transferUseCase interface {
	Execute(ctx context.Context, input transactionapp.TransferInput) (*transactionapp.TransferResult, error)
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
	deposit       depositUseCase
	withdraw      withdrawUseCase
	transfer      transferUseCase
	statement     statementUseCase
	balance       getBalanceUseCase
}

func New(
	listAccounts listAccountsUseCase,
	createAccount createAccountUseCase,
	deposit depositUseCase,
	withdraw withdrawUseCase,
	transfer transferUseCase,
	statement statementUseCase,
	balance getBalanceUseCase,
) *Handler {
	return &Handler{
		listAccounts:  listAccounts,
		createAccount: createAccount,
		deposit:       deposit,
		withdraw:      withdraw,
		transfer:      transfer,
		statement:     statement,
		balance:       balance,
	}
}

func RequireUser(ctx context.Context) (*authdomain.AuthenticatedUser, error) {
	user, ok := sharedauthctx.GetAuthenticatedUser(ctx)
	if !ok || user == nil {
		return nil, authdomain.ErrUnauthorized
	}

	return user, nil
}
