package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	installationdomain "github.com/seu-usuario/bank-api/internal/installation/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
)

func TestRevokeInstallationUseCase_Execute_SuccessInvalidatesSessions(t *testing.T) {
	userID := uuid.New()
	currentInstallationUUID := uuid.New()
	targetInstallationUUID := uuid.New()
	resourceUUID := uuid.New()
	now := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	resourceID := mustResourceID(t, resourceUUID)
	targetInstallationID := mustInstallationID(t, targetInstallationUUID)
	target := mustKnownInstallation(t, userID, resourceID, targetInstallationID, now)
	revoked := mustKnownInstallation(t, userID, resourceID, targetInstallationID, now)
	revokedAt := now.Add(time.Minute)
	if err := revoked.Revoke(revokedAt); err != nil {
		t.Fatalf("expected revoke, got %v", err)
	}
	repo := &installationRepositoryMock{
		findByResourceOut: target,
		revokeOut:         revoked,
	}
	sessions := &sessionRepositoryMock{}
	tx := &transactorMock{}
	uc := NewRevokeInstallationUseCase(repo, sessions, tx)

	ctx := sharedauthctx.WithOperationalSession(context.Background(), sharedauthctx.OperationalSession{
		UserID:         userID,
		InstallationID: &currentInstallationUUID,
	})
	output, err := uc.Execute(ctx, RevokeInstallationInput{
		ResourceID: resourceUUID,
		Now:        now,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.ResourceID != resourceUUID || output.Status != string(installationdomain.InstallationStatusRevoked) {
		t.Fatalf("unexpected output: %#v", output)
	}
	if repo.findByResourceCalls != 1 || repo.findByResourceUser != userID {
		t.Fatalf("expected find by current user, got calls=%d user=%q", repo.findByResourceCalls, repo.findByResourceUser)
	}
	if tx.calls != 1 {
		t.Fatalf("expected one transaction, got %d", tx.calls)
	}
	if repo.revokeCalls != 1 || repo.revokeUser != userID {
		t.Fatalf("expected revoke by current user, got calls=%d user=%q", repo.revokeCalls, repo.revokeUser)
	}
	if sessions.invalidateCalls != 1 {
		t.Fatalf("expected invalidation once, got %d", sessions.invalidateCalls)
	}
	if sessions.invalidateUser != userID || sessions.invalidateInstallID != targetInstallationUUID {
		t.Fatalf("unexpected invalidation target user=%q installation=%q", sessions.invalidateUser, sessions.invalidateInstallID)
	}
	if !sessions.invalidateRevokedAt.Equal(revokedAt) {
		t.Fatalf("expected invalidation timestamp %s, got %s", revokedAt, sessions.invalidateRevokedAt)
	}
}

func TestRevokeInstallationUseCase_Execute_CannotRevokeCurrentInstallation(t *testing.T) {
	userID := uuid.New()
	currentInstallationUUID := uuid.New()
	resourceUUID := uuid.New()
	now := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	current := mustKnownInstallation(
		t,
		userID,
		mustResourceID(t, resourceUUID),
		mustInstallationID(t, currentInstallationUUID),
		now,
	)
	repo := &installationRepositoryMock{findByResourceOut: current}
	sessions := &sessionRepositoryMock{}
	uc := NewRevokeInstallationUseCase(repo, sessions, &transactorMock{})

	ctx := sharedauthctx.WithOperationalSession(context.Background(), sharedauthctx.OperationalSession{
		UserID:         userID,
		InstallationID: &currentInstallationUUID,
	})
	output, err := uc.Execute(ctx, RevokeInstallationInput{ResourceID: resourceUUID})
	if !errors.Is(err, installationdomain.ErrInstallationMismatch) {
		t.Fatalf("expected ErrInstallationMismatch, got %v", err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %#v", output)
	}
	if repo.revokeCalls != 0 {
		t.Fatalf("expected no revoke, got %d", repo.revokeCalls)
	}
	if sessions.invalidateCalls != 0 {
		t.Fatalf("expected no invalidation, got %d", sessions.invalidateCalls)
	}
}
