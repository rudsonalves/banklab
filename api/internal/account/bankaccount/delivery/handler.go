package delivery

import (
	"context"

	accountapp "github.com/seu-usuario/bank-api/internal/account/bankaccount/application"
	"github.com/seu-usuario/bank-api/internal/account/bankaccount/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
)

type createAccountUseCase interface {
	Execute(ctx context.Context, input accountapp.CreateAccountInput) (*domain.Account, error)
}

type listAccountsUseCase interface {
	Execute(ctx context.Context, input accountapp.ListAccountsInput) ([]domain.Account, error)
}

type getBalanceUseCase interface {
	Execute(ctx context.Context, input accountapp.GetAccountBalanceInput) (*accountapp.AccountBalance, error)
}

type lookupInternalTransferRecipientsUseCase interface {
	Execute(ctx context.Context, input accountapp.LookupInternalTransferRecipientsInput) ([]domain.TransferRecipient, error)
}

type Handler struct {
	listAccounts                     listAccountsUseCase
	createAccount                    createAccountUseCase
	balance                          getBalanceUseCase
	lookupInternalTransferRecipients lookupInternalTransferRecipientsUseCase
}

// New creates a new instance of the account Handler with the provided use cases.
func New(
	listAccounts listAccountsUseCase,
	createAccount createAccountUseCase,
	balance getBalanceUseCase,
	lookupInternalTransferRecipients lookupInternalTransferRecipientsUseCase,
) *Handler {
	return &Handler{
		listAccounts:                     listAccounts,
		createAccount:                    createAccount,
		balance:                          balance,
		lookupInternalTransferRecipients: lookupInternalTransferRecipients,
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
