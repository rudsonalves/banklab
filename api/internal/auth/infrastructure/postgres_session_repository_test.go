package infrastructure

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

func TestPostgresSessionRepository_InstallationSessionIntegration(t *testing.T) {
	ctx := context.Background()
	pool := newAuthTestPool(t, ctx)
	defer pool.Close()

	ensureAuthRepoTestSchema(t, ctx, pool)

	user := testUser(time.Now().UTC())
	userRepo := NewPostgresUserRepository(pool)
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("expected no error creating user, got %v", err)
	}
	defer cleanupUserByID(t, ctx, pool, user.ID)

	repo := NewPostgresSessionRepository(pool)
	installationID := uuid.New()
	otherInstallationID := uuid.New()
	activeHash := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	otherHash := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	expiresAt := time.Now().UTC().Add(time.Hour)

	if err := repo.CreateWithInstallation(ctx, domain.CreateSessionInput{
		UserID:         user.ID,
		TokenHash:      activeHash,
		ExpiresAt:      expiresAt,
		InstallationID: &installationID,
	}); err != nil {
		t.Fatalf("expected create with installation to succeed, got %v", err)
	}
	if err := repo.CreateWithInstallation(ctx, domain.CreateSessionInput{
		UserID:         user.ID,
		TokenHash:      otherHash,
		ExpiresAt:      expiresAt,
		InstallationID: &otherInstallationID,
	}); err != nil {
		t.Fatalf("expected create with other installation to succeed, got %v", err)
	}

	record, err := repo.FindByTokenHashWithInstallation(ctx, activeHash)
	if err != nil {
		t.Fatalf("expected find with installation to succeed, got %v", err)
	}
	if record == nil {
		t.Fatal("expected session record, got nil")
	}
	if record.InstallationID == nil || *record.InstallationID != installationID {
		t.Fatalf("expected installation id %q, got %#v", installationID, record.InstallationID)
	}
	if record.Revoked {
		t.Fatal("expected session to be active before revocation")
	}

	revokedAt := time.Now().UTC()
	if err := repo.RevokeByUserIDAndInstallationID(ctx, user.ID, installationID, revokedAt); err != nil {
		t.Fatalf("expected revoke by installation to succeed, got %v", err)
	}

	revokedRecord, err := repo.FindByTokenHashWithInstallation(ctx, activeHash)
	if err != nil {
		t.Fatalf("expected find revoked session to succeed, got %v", err)
	}
	if revokedRecord == nil || !revokedRecord.Revoked {
		t.Fatalf("expected target session to be revoked, got %#v", revokedRecord)
	}

	otherRecord, err := repo.FindByTokenHashWithInstallation(ctx, otherHash)
	if err != nil {
		t.Fatalf("expected find other session to succeed, got %v", err)
	}
	if otherRecord == nil || otherRecord.Revoked {
		t.Fatalf("expected other installation session to remain active, got %#v", otherRecord)
	}
}
