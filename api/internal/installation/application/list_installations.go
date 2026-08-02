package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	installationdomain "github.com/seu-usuario/bank-api/internal/installation/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
)

type ListInstallationsUseCase struct {
	installations installationdomain.InstallationReader
}

func NewListInstallationsUseCase(
	installations installationdomain.InstallationReader,
) *ListInstallationsUseCase {
	return &ListInstallationsUseCase{installations: installations}
}

type ListInstallationsOutput struct {
	Installations []InstallationSummary
}

type InstallationSummary struct {
	ResourceID  uuid.UUID
	Status      string
	FirstSeenAt time.Time
	LastSeenAt  time.Time
	RevokedAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (uc *ListInstallationsUseCase) Execute(ctx context.Context) (*ListInstallationsOutput, error) {
	session, err := sharedauthctx.RequireOperationalSession(ctx)
	if err != nil {
		return nil, err
	}
	if session.UserID == uuid.Nil {
		return nil, authdomain.ErrUnauthorized
	}
	if uc == nil || uc.installations == nil {
		return nil, authdomain.ErrInvalidData
	}

	installations, err := uc.installations.ListByUserID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	output := &ListInstallationsOutput{
		Installations: make([]InstallationSummary, 0, len(installations)),
	}
	for _, installation := range installations {
		if installation == nil {
			continue
		}
		output.Installations = append(output.Installations, InstallationSummary{
			ResourceID:  installation.ResourceID.UUID(),
			Status:      string(installation.Status),
			FirstSeenAt: installation.FirstSeenAt,
			LastSeenAt:  installation.LastSeenAt,
			RevokedAt:   installation.RevokedAt,
			CreatedAt:   installation.CreatedAt,
			UpdatedAt:   installation.UpdatedAt,
		})
	}

	return output, nil
}
