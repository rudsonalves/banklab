package application

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	installationdomain "github.com/seu-usuario/bank-api/internal/installation/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
)

func TestListInstallationsUseCase_Execute_ReturnsSafeSummaries(t *testing.T) {
	userID := uuid.New()
	now := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	knownResourceID := mustResourceID(t, uuid.New())
	revokedResourceID := mustResourceID(t, uuid.New())
	known := mustKnownInstallation(t, userID, knownResourceID, mustInstallationID(t, uuid.New()), now)
	revoked := mustKnownInstallation(t, userID, revokedResourceID, mustInstallationID(t, uuid.New()), now)
	revokedAt := now.Add(time.Minute)
	if err := revoked.Revoke(revokedAt); err != nil {
		t.Fatalf("expected revoke, got %v", err)
	}
	repo := &installationRepositoryMock{listOut: []*installationdomain.Installation{known, revoked}}
	uc := NewListInstallationsUseCase(repo)

	ctx := sharedauthctx.WithOperationalSession(context.Background(), sharedauthctx.OperationalSession{
		UserID: userID,
	})
	output, err := uc.Execute(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.listCalls != 1 || repo.listUser != userID {
		t.Fatalf("expected list by user %q, got calls=%d user=%q", userID, repo.listCalls, repo.listUser)
	}
	if len(output.Installations) != 2 {
		t.Fatalf("expected two installations, got %d", len(output.Installations))
	}
	if output.Installations[0].ResourceID != knownResourceID.UUID() {
		t.Fatalf("expected known resource id %q, got %q", knownResourceID.UUID(), output.Installations[0].ResourceID)
	}
	if output.Installations[1].Status != string(installationdomain.InstallationStatusRevoked) {
		t.Fatalf("expected revoked status, got %q", output.Installations[1].Status)
	}
	if output.Installations[1].RevokedAt == nil || !output.Installations[1].RevokedAt.Equal(revokedAt) {
		t.Fatalf("expected revoked_at %s, got %#v", revokedAt, output.Installations[1].RevokedAt)
	}
}
