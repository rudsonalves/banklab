package delivery

import (
	"context"

	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
)

type AuthenticatedUser = authdomain.AuthenticatedUser

var ErrAuthenticatedUserNotFound = sharedauthctx.ErrAuthenticatedUserNotFound

// GetAuthenticatedUser retrieves the authenticated user information from the
// context. It returns the authenticated user and a boolean indicating whether
// the user was successfully retrieved from the context. If the user is not found
// in the context, it returns nil and false.
func GetAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, bool) {
	return sharedauthctx.GetAuthenticatedUser(ctx)
}

// WithAuthenticatedUser adds the authenticated user information to the context.
// It returns a new context with the authenticated user included.
func WithAuthenticatedUser(ctx context.Context, user AuthenticatedUser) context.Context {
	return sharedauthctx.WithAuthenticatedUser(ctx, user)
}

// RequireAuthenticatedUser retrieves the authenticated user information from the
// context. If the user is not found, it returns an error indicating that the
// authenticated user is required.
func RequireAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, error) {
	return sharedauthctx.RequireAuthenticatedUser(ctx)
}
