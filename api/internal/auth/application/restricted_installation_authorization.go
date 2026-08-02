package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	installationdomain "github.com/seu-usuario/bank-api/internal/installation/domain"
)

type DefaultRestrictedInstallationAuthorizationIssuer struct {
	repo   installationdomain.RestrictedAuthorizationRepository
	signer installationdomain.RestrictedAccessTokenSigner
}

func NewDefaultRestrictedInstallationAuthorizationIssuer(
	repo installationdomain.RestrictedAuthorizationRepository,
	signer installationdomain.RestrictedAccessTokenSigner,
) *DefaultRestrictedInstallationAuthorizationIssuer {
	return &DefaultRestrictedInstallationAuthorizationIssuer{
		repo:   repo,
		signer: signer,
	}
}

func (i *DefaultRestrictedInstallationAuthorizationIssuer) Issue(
	ctx context.Context,
	userID uuid.UUID,
	installationUUID uuid.UUID,
	now time.Time,
) (*RestrictedInstallationAuthorization, error) {
	if i == nil || i.repo == nil || i.signer == nil {
		return nil, fmt.Errorf("restricted installation authorization issuer not configured")
	}

	installationID, err := installationdomain.NewInstallationID(installationUUID)
	if err != nil {
		return nil, err
	}

	if err := i.repo.RevokeActiveByUserIDAndInstallationID(
		ctx,
		userID,
		installationID,
		installationdomain.RestrictedAuthorizationScopeInstallationRegister,
	); err != nil {
		return nil, fmt.Errorf("revoke previous restricted authorizations: %w", err)
	}

	authorization, err := installationdomain.NewRestrictedAuthorization(uuid.NewString(), userID, installationID, now)
	if err != nil {
		return nil, err
	}
	if err := i.repo.Create(ctx, authorization); err != nil {
		return nil, err
	}

	claims := &installationdomain.RestrictedAccessTokenClaims{
		UserID:         authorization.UserID,
		InstallationID: authorization.InstallationID,
		JTI:            authorization.JTI,
		TokenType:      installationdomain.RestrictedAccessTokenType,
		Scope:          authorization.Scope,
		IssuedAt:       authorization.CreatedAt,
		ExpiresAt:      authorization.ExpiresAt,
	}
	token, err := i.signer.SignRestrictedAccessToken(claims)
	if err != nil {
		return nil, err
	}

	return &RestrictedInstallationAuthorization{
		Token:     token,
		TokenType: claims.TokenType,
		Scope:     claims.Scope,
		ExpiresAt: claims.ExpiresAt,
	}, nil
}
