package application

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

type installationLoginRepositoryMock struct {
	record             *InstallationLoginRecord
	recordErr          error
	knownCount         int
	knownCountErr      error
	hasAny             bool
	hasAnyErr          error
	findCalls          int
	countKnownCalls    int
	hasAnyCalls        int
	findUserID         uuid.UUID
	findInstallationID uuid.UUID
	hasAnyUserID       uuid.UUID
	countKnownUserID   uuid.UUID
}

func (m *installationLoginRepositoryMock) FindByUserIDAndInstallationID(
	ctx context.Context,
	userID uuid.UUID,
	installationID uuid.UUID,
) (*InstallationLoginRecord, error) {
	m.findCalls++
	m.findUserID = userID
	m.findInstallationID = installationID
	if m.recordErr != nil {
		return nil, m.recordErr
	}
	return m.record, nil
}

func (m *installationLoginRepositoryMock) CountKnownByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	m.countKnownCalls++
	m.countKnownUserID = userID
	if m.knownCountErr != nil {
		return 0, m.knownCountErr
	}
	return m.knownCount, nil
}

func (m *installationLoginRepositoryMock) HasAnyByUserID(ctx context.Context, userID uuid.UUID) (bool, error) {
	m.hasAnyCalls++
	m.hasAnyUserID = userID
	if m.hasAnyErr != nil {
		return false, m.hasAnyErr
	}
	return m.hasAny, nil
}

func TestDefaultInstallationLoginClassifier_Classify(t *testing.T) {
	userID := uuid.New()
	installationID := uuid.New()

	tests := []struct {
		name            string
		repo            *installationLoginRepositoryMock
		wantClass       InstallationLoginClassification
		wantKnownCount  int
		wantHasAny      bool
		wantCountCalls  int
		wantHasAnyCalls int
	}{
		{
			name: "known installation",
			repo: &installationLoginRepositoryMock{
				record: &InstallationLoginRecord{
					UserID:         userID,
					InstallationID: installationID,
					Status:         InstallationLoginRecordKnown,
				},
			},
			wantClass: InstallationLoginKnown,
		},
		{
			name: "revoked installation",
			repo: &installationLoginRepositoryMock{
				record: &InstallationLoginRecord{
					UserID:         userID,
					InstallationID: installationID,
					Status:         InstallationLoginRecordRevoked,
				},
			},
			wantClass: InstallationLoginRevoked,
		},
		{
			name: "first installation",
			repo: &installationLoginRepositoryMock{
				hasAny: false,
			},
			wantClass:       InstallationLoginFirst,
			wantHasAny:      false,
			wantHasAnyCalls: 1,
		},
		{
			name: "new installation with vacancy",
			repo: &installationLoginRepositoryMock{
				hasAny:     true,
				knownCount: 2,
			},
			wantClass:       InstallationLoginNew,
			wantKnownCount:  2,
			wantHasAny:      true,
			wantHasAnyCalls: 1,
			wantCountCalls:  1,
		},
		{
			name: "installation limit reached",
			repo: &installationLoginRepositoryMock{
				hasAny:     true,
				knownCount: 3,
			},
			wantClass:       InstallationLoginLimitReached,
			wantKnownCount:  3,
			wantHasAny:      true,
			wantHasAnyCalls: 1,
			wantCountCalls:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			classifier := NewDefaultInstallationLoginClassifier(tt.repo)

			got, err := classifier.Classify(context.Background(), userID, installationID)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if got.Classification != tt.wantClass {
				t.Fatalf("expected classification %q, got %q", tt.wantClass, got.Classification)
			}
			if got.KnownInstallationsCount != tt.wantKnownCount {
				t.Fatalf("expected known count %d, got %d", tt.wantKnownCount, got.KnownInstallationsCount)
			}
			if got.HasAssociatedInstallations != tt.wantHasAny {
				t.Fatalf("expected hasAny %v, got %v", tt.wantHasAny, got.HasAssociatedInstallations)
			}
			if tt.repo.findCalls != 1 {
				t.Fatalf("expected find calls 1, got %d", tt.repo.findCalls)
			}
			if tt.repo.hasAnyCalls != tt.wantHasAnyCalls {
				t.Fatalf("expected hasAny calls %d, got %d", tt.wantHasAnyCalls, tt.repo.hasAnyCalls)
			}
			if tt.repo.countKnownCalls != tt.wantCountCalls {
				t.Fatalf("expected countKnown calls %d, got %d", tt.wantCountCalls, tt.repo.countKnownCalls)
			}
		})
	}
}

func TestDefaultInstallationLoginClassifier_Classify_InvalidIDs(t *testing.T) {
	classifier := NewDefaultInstallationLoginClassifier(&installationLoginRepositoryMock{})

	_, err := classifier.Classify(context.Background(), uuid.Nil, uuid.New())
	if !errors.Is(err, authdomain.ErrInvalidData) {
		t.Fatalf("expected ErrInvalidData, got %v", err)
	}

	_, err = classifier.Classify(context.Background(), uuid.New(), uuid.Nil)
	if !errors.Is(err, authdomain.ErrInvalidData) {
		t.Fatalf("expected ErrInvalidData, got %v", err)
	}
}

func TestDefaultInstallationLoginClassifier_Classify_RepositoryErrorsAreWrapped(t *testing.T) {
	repoErr := errors.New("db down")
	classifier := NewDefaultInstallationLoginClassifier(&installationLoginRepositoryMock{
		recordErr: repoErr,
	})

	_, err := classifier.Classify(context.Background(), uuid.New(), uuid.New())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected wrapped error %v, got %v", repoErr, err)
	}
}
