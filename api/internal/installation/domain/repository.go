package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type InstallationReader interface {
	FindByUserIDAndInstallationID(
		ctx context.Context,
		userID uuid.UUID,
		installationID InstallationID,
	) (*Installation, error)
	FindByResourceID(
		ctx context.Context,
		userID uuid.UUID,
		resourceID InstallationResourceID,
	) (*Installation, error)
	CountKnownByUserID(ctx context.Context, userID uuid.UUID) (int, error)
	HasAnyByUserID(ctx context.Context, userID uuid.UUID) (bool, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*Installation, error)
}

type InstallationWriter interface {
	BootstrapFirstInstallation(
		ctx context.Context,
		userID uuid.UUID,
		resourceID InstallationResourceID,
		installationID InstallationID,
		now time.Time,
	) (*Installation, error)
	ReserveKnownInstallation(
		ctx context.Context,
		userID uuid.UUID,
		resourceID InstallationResourceID,
		installationID InstallationID,
		maxKnownInstallations int,
		now time.Time,
	) (*Installation, error)
	RevokeByResourceID(
		ctx context.Context,
		userID uuid.UUID,
		resourceID InstallationResourceID,
		now time.Time,
	) (*Installation, error)
}

type InstallationRepository interface {
	InstallationReader
	InstallationWriter
}

type InstallationSessionInvalidator interface {
	InvalidateByInstallationID(ctx context.Context, userID uuid.UUID, installationID InstallationID, now time.Time) error
}

type RestrictedAuthorizationRepository interface {
	Create(ctx context.Context, authorization *RestrictedAuthorization) error
	FindByJTI(ctx context.Context, jti string) (*RestrictedAuthorization, error)
	ConsumeByJTI(ctx context.Context, jti string, now time.Time) (*RestrictedAuthorization, error)
	RevokeByJTI(ctx context.Context, jti string) error
	RevokeActiveByUserIDAndInstallationID(
		ctx context.Context,
		userID uuid.UUID,
		installationID InstallationID,
		scope string,
	) error
}
