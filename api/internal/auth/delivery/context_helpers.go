package delivery

import (
	"context"

	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
)

type AuthenticatedUser = authdomain.AuthenticatedUser

var ErrAuthenticatedUserNotFound = sharedauthctx.ErrAuthenticatedUserNotFound

func GetAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, bool) {
	return sharedauthctx.GetAuthenticatedUser(ctx)
}

func WithAuthenticatedUser(ctx context.Context, user AuthenticatedUser) context.Context {
	return sharedauthctx.WithAuthenticatedUser(ctx, user)
}

func RequireAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, error) {
	return sharedauthctx.RequireAuthenticatedUser(ctx)
}
