package delivery

import (
	"context"
	"errors"

	"github.com/seu-usuario/bank-api/internal/auth/application"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

type AuthenticatedUser = authdomain.AuthenticatedUser

var ErrAuthenticatedUserNotFound = errors.New("authenticated user not found in context")

func GetAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, bool) {
	return application.GetAuthenticatedUser(ctx)
}

func WithAuthenticatedUser(ctx context.Context, user AuthenticatedUser) context.Context {
	return application.WithAuthenticatedUser(ctx, user)
}

func RequireAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, error) {
	user, ok := GetAuthenticatedUser(ctx)
	if !ok {
		return nil, ErrAuthenticatedUserNotFound
	}

	return user, nil
}
