package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	installationdomain "github.com/seu-usuario/bank-api/internal/installation/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
)

type RevokeInstallationUseCase struct {
	installations installationdomain.InstallationRepository
	invalidator   installationdomain.InstallationSessionInvalidator
	transactor    authdomain.Transactor
	now           func() time.Time
}

func NewRevokeInstallationUseCase(
	installations installationdomain.InstallationRepository,
	invalidator installationdomain.InstallationSessionInvalidator,
	transactor authdomain.Transactor,
) *RevokeInstallationUseCase {
	return &RevokeInstallationUseCase{
		installations: installations,
		invalidator:   invalidator,
		transactor:    transactor,
		now:           func() time.Time { return time.Now().UTC() },
	}
}

type RevokeInstallationInput struct {
	ResourceID uuid.UUID
	Now        time.Time
}

type RevokeInstallationOutput struct {
	ResourceID uuid.UUID
	Status     string
	RevokedAt  *time.Time
}

func (uc *RevokeInstallationUseCase) Execute(
	ctx context.Context,
	input RevokeInstallationInput,
) (*RevokeInstallationOutput, error) {
	session, err := sharedauthctx.RequireOperationalSession(ctx)
	if err != nil {
		return nil, err
	}
	if session.UserID == uuid.Nil || session.InstallationID == nil || *session.InstallationID == uuid.Nil {
		return nil, authdomain.ErrUnauthorized
	}
	if uc == nil || uc.installations == nil || uc.invalidator == nil || uc.transactor == nil {
		return nil, authdomain.ErrInvalidData
	}

	resourceID, err := installationdomain.NewInstallationResourceID(input.ResourceID)
	if err != nil {
		return nil, err
	}

	current, err := uc.installations.FindByResourceID(ctx, session.UserID, resourceID)
	if err != nil {
		return nil, err
	}
	if current.InstallationID.UUID() == *session.InstallationID {
		return nil, installationdomain.ErrInstallationMismatch
	}

	now := input.Now.UTC()
	if now.IsZero() {
		now = uc.now().UTC()
	}

	var revoked *installationdomain.Installation
	if err := uc.transactor.RunInTx(ctx, func(txCtx context.Context) error {
		installation, err := uc.installations.RevokeByResourceID(txCtx, session.UserID, resourceID, now)
		if err != nil {
			return fmt.Errorf("revoke installation: %w", err)
		}
		revoked = installation

		revokedAt := now
		if installation.RevokedAt != nil {
			revokedAt = *installation.RevokedAt
		}
		if err := uc.invalidator.InvalidateByInstallationID(
			txCtx,
			session.UserID,
			installation.InstallationID,
			revokedAt,
		); err != nil {
			return fmt.Errorf("invalidate installation sessions: %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}
	if revoked == nil {
		return nil, installationdomain.ErrInvalidInstallation
	}

	return &RevokeInstallationOutput{
		ResourceID: revoked.ResourceID.UUID(),
		Status:     string(revoked.Status),
		RevokedAt:  revoked.RevokedAt,
	}, nil
}
