// Package authctx provides helpers to propagate authenticated user
// information through context.Context across application layers.
//
// This package is typically used at the delivery layer (HTTP middleware)
// to attach the authenticated user to the request context, and later
// consumed by application use cases.
//
// The authenticated user is treated as part of the request execution context,
// not as a global or implicit dependency.
package authctx

import (
	"context"
	"errors"

	"github.com/google/uuid"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

// AuthenticatedUser represents the authenticated identity extracted from JWT
// and used across application use cases.
//
// It is defined as an alias to the domain type to avoid duplication
// while keeping dependency direction explicit.
type AuthenticatedUser = authdomain.AuthenticatedUser

type OperationalSession struct {
	UserID         uuid.UUID
	Role           authdomain.Role
	CustomerID     *uuid.UUID
	InstallationID *uuid.UUID
}

type RestrictedSession struct {
	UserID         uuid.UUID
	InstallationID uuid.UUID
	JTI            string
	Scope          string
}

// contextKey is a private type used to avoid key collisions
// when storing values in context.
type contextKey string

// authenticatedUserKey is the key used to store the authenticated user in context.
const (
	authenticatedUserKey  contextKey = "authenticatedUser"
	operationalSessionKey contextKey = "operationalSession"
	restrictedSessionKey  contextKey = "restrictedSession"
)

// ErrAuthenticatedUserNotFound is returned when the context does not contain
// a valid authenticated user.
var (
	ErrAuthenticatedUserNotFound  = errors.New("authenticated user not found in context")
	ErrOperationalSessionNotFound = errors.New("operational session not found in context")
	ErrRestrictedSessionNotFound  = errors.New("restricted session not found in context")
)

// WithAuthenticatedUser returns a new context derived from ctx that carries
// the given authenticated user.
//
// This function is typically used by authentication middleware after
// validating a JWT token.
//
// Example:
//
//	ctx = authctx.WithAuthenticatedUser(ctx, user)
func WithAuthenticatedUser(ctx context.Context, user AuthenticatedUser) context.Context {
	return context.WithValue(ctx, authenticatedUserKey, user)
}

// GetAuthenticatedUser retrieves the authenticated user from the context.
//
// It supports both value and pointer storage for flexibility.
// Returns the user and true if found, otherwise nil and false.
//
// This function should be used when the caller can handle the absence
// of authentication context.
//
// Example:
//
//	user, ok := authctx.GetAuthenticatedUser(ctx)
//	if !ok {
//	    // handle unauthenticated request
//	}
func GetAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, bool) {
	if user, ok := ctx.Value(authenticatedUserKey).(AuthenticatedUser); ok {
		return &user, true
	}

	if userPtr, ok := ctx.Value(authenticatedUserKey).(*AuthenticatedUser); ok && userPtr != nil {
		return userPtr, true
	}

	return nil, false
}

// RequireAuthenticatedUser retrieves the authenticated user from the context.
//
// It returns an error if the user is not found, making it suitable for use
// in application use cases that require authentication.
//
// Example:
//
//	user, err := authctx.RequireAuthenticatedUser(ctx)
//	if err != nil {
//	    // handle unauthenticated request
//	}
func RequireAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, error) {
	if user, ok := GetAuthenticatedUser(ctx); ok {
		return user, nil
	}

	return nil, ErrAuthenticatedUserNotFound
}

func WithOperationalSession(ctx context.Context, session OperationalSession) context.Context {
	return context.WithValue(ctx, operationalSessionKey, session)
}

func GetOperationalSession(ctx context.Context) (*OperationalSession, bool) {
	if session, ok := ctx.Value(operationalSessionKey).(OperationalSession); ok {
		return &session, true
	}
	if sessionPtr, ok := ctx.Value(operationalSessionKey).(*OperationalSession); ok && sessionPtr != nil {
		return sessionPtr, true
	}

	return nil, false
}

func RequireOperationalSession(ctx context.Context) (*OperationalSession, error) {
	if session, ok := GetOperationalSession(ctx); ok {
		return session, nil
	}

	return nil, ErrOperationalSessionNotFound
}

func WithRestrictedSession(ctx context.Context, session RestrictedSession) context.Context {
	return context.WithValue(ctx, restrictedSessionKey, session)
}

func GetRestrictedSession(ctx context.Context) (*RestrictedSession, bool) {
	if session, ok := ctx.Value(restrictedSessionKey).(RestrictedSession); ok {
		return &session, true
	}
	if sessionPtr, ok := ctx.Value(restrictedSessionKey).(*RestrictedSession); ok && sessionPtr != nil {
		return sessionPtr, true
	}

	return nil, false
}

func RequireRestrictedSession(ctx context.Context) (*RestrictedSession, error) {
	if session, ok := GetRestrictedSession(ctx); ok {
		return session, nil
	}

	return nil, ErrRestrictedSessionNotFound
}
