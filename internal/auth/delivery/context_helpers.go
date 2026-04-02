package delivery

import (
	"context"

	"github.com/seu-usuario/bank-api/internal/auth/application"
)

type AuthenticatedUser = application.AuthenticatedUser

func GetAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, bool) {
	return application.GetAuthenticatedUser(ctx)
}

func MustGetAuthenticatedUser(ctx context.Context) *AuthenticatedUser {
	user, ok := GetAuthenticatedUser(ctx)
	if !ok {
		panic("authenticated user not found in context")
	}

	return user
}
