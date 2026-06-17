package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const RestrictedAccessTokenType = "restricted_access"

type RestrictedAccessTokenClaims struct {
	UserID         uuid.UUID
	InstallationID InstallationID
	JTI            string
	TokenType      string
	Scope          string
	IssuedAt       time.Time
	ExpiresAt      time.Time
}

func (c *RestrictedAccessTokenClaims) Validate() error {
	if c == nil ||
		c.UserID == uuid.Nil ||
		c.InstallationID.IsZero() ||
		c.JTI == "" ||
		c.TokenType != RestrictedAccessTokenType ||
		c.Scope != RestrictedAuthorizationScopeInstallationRegister ||
		c.IssuedAt.IsZero() ||
		c.ExpiresAt.IsZero() ||
		!c.ExpiresAt.After(c.IssuedAt) {
		return ErrRestrictedAuthorizationInvalid
	}

	return nil
}

type RestrictedAccessTokenSigner interface {
	SignRestrictedAccessToken(claims *RestrictedAccessTokenClaims) (string, error)
}

type RestrictedAccessTokenVerifier interface {
	VerifyRestrictedAccessToken(ctx context.Context, token string) (*RestrictedAccessTokenClaims, error)
}
