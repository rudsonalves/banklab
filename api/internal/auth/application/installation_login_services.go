package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	installationdomain "github.com/seu-usuario/bank-api/internal/installation/domain"
)

type InstallationLoginRepositoryAdapter struct {
	reader installationdomain.InstallationReader
}

func NewInstallationLoginRepositoryAdapter(
	reader installationdomain.InstallationReader,
) *InstallationLoginRepositoryAdapter {
	return &InstallationLoginRepositoryAdapter{reader: reader}
}

func (a *InstallationLoginRepositoryAdapter) FindByUserIDAndInstallationID(
	ctx context.Context,
	userID uuid.UUID,
	installationUUID uuid.UUID,
) (*InstallationLoginRecord, error) {
	if a == nil || a.reader == nil {
		return nil, fmt.Errorf("installation reader not configured")
	}

	installationID, err := installationdomain.NewInstallationID(installationUUID)
	if err != nil {
		return nil, err
	}

	installation, err := a.reader.FindByUserIDAndInstallationID(ctx, userID, installationID)
	if err != nil {
		if errors.Is(err, installationdomain.ErrInstallationNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if installation == nil {
		return nil, nil
	}

	return &InstallationLoginRecord{
		UserID:         installation.UserID,
		InstallationID: installation.InstallationID.UUID(),
		Status:         InstallationLoginRecordStatus(installation.Status),
	}, nil
}

func (a *InstallationLoginRepositoryAdapter) CountKnownByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	if a == nil || a.reader == nil {
		return 0, fmt.Errorf("installation reader not configured")
	}

	return a.reader.CountKnownByUserID(ctx, userID)
}

func (a *InstallationLoginRepositoryAdapter) HasAnyByUserID(ctx context.Context, userID uuid.UUID) (bool, error) {
	if a == nil || a.reader == nil {
		return false, fmt.Errorf("installation reader not configured")
	}

	return a.reader.HasAnyByUserID(ctx, userID)
}

type FirstInstallationBootstrapperAdapter struct {
	writer installationdomain.InstallationWriter
}

func NewFirstInstallationBootstrapperAdapter(
	writer installationdomain.InstallationWriter,
) *FirstInstallationBootstrapperAdapter {
	return &FirstInstallationBootstrapperAdapter{writer: writer}
}

func (a *FirstInstallationBootstrapperAdapter) BootstrapFirstInstallation(
	ctx context.Context,
	userID uuid.UUID,
	installationUUID uuid.UUID,
	now time.Time,
) error {
	if a == nil || a.writer == nil {
		return fmt.Errorf("installation writer not configured")
	}

	resourceID, err := installationdomain.NewInstallationResourceID(uuid.New())
	if err != nil {
		return err
	}
	installationID, err := installationdomain.NewInstallationID(installationUUID)
	if err != nil {
		return err
	}

	_, err = a.writer.BootstrapFirstInstallation(ctx, userID, resourceID, installationID, now)
	return err
}
