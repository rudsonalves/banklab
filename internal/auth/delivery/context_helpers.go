package delivery

import (
	"context"

	"github.com/seu-usuario/bank-api/internal/auth/application"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

type AuthenticatedUser = authdomain.AuthenticatedUser

func GetAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, bool) {
	return application.GetAuthenticatedUser(ctx)
}

func WithAuthenticatedUser(ctx context.Context, user AuthenticatedUser) context.Context {
	return application.WithAuthenticatedUser(ctx, user)
}

func MustGetAuthenticatedUser(ctx context.Context) *AuthenticatedUser {
	user, ok := GetAuthenticatedUser(ctx)
	if !ok {
		panic("authenticated user not found in context")
	}

	return user
}
