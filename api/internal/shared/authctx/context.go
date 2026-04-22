package authctx

import (
	"context"
	"errors"

	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

type AuthenticatedUser = authdomain.AuthenticatedUser

type contextKey string

const authenticatedUserKey contextKey = "authenticatedUser"

var ErrAuthenticatedUserNotFound = errors.New("authenticated user not found in context")

func WithAuthenticatedUser(ctx context.Context, user AuthenticatedUser) context.Context {
	return context.WithValue(ctx, authenticatedUserKey, user)
}

func GetAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, bool) {
	user, ok := ctx.Value(authenticatedUserKey).(AuthenticatedUser)
	if ok {
		return &user, true
	}

	userPtr, ok := ctx.Value(authenticatedUserKey).(*AuthenticatedUser)
	if !ok || userPtr == nil {
		return nil, false
	}

	return userPtr, true
}

func RequireAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, error) {
	user, ok := GetAuthenticatedUser(ctx)
	if !ok {
		return nil, ErrAuthenticatedUserNotFound
	}

	return user, nil
}
