package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

type contactVerificationRepositoryMock struct {
	createCalls         int
	createErr           error
	createdVerification *domain.ContactVerification
	findValue           *domain.ContactVerification
	findErr             error
	confirmCalls        int
	confirmErr          error
	confirmedID         uuid.UUID
	confirmedToken      string
	confirmedVerifiedAt time.Time
}

func (m *contactVerificationRepositoryMock) CreateContactVerification(
	ctx context.Context,
	verification *domain.ContactVerification,
) error {
	m.createCalls++
	m.createdVerification = verification
	return m.createErr
}

func (m *contactVerificationRepositoryMock) FindContactVerificationByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.ContactVerification, error) {
	return m.findValue, m.findErr
}

func (m *contactVerificationRepositoryMock) ConfirmContactVerification(
	ctx context.Context,
	id uuid.UUID,
	verificationToken string,
	verifiedAt time.Time,
) error {
	m.confirmCalls++
	m.confirmedID = id
	m.confirmedToken = verificationToken
	m.confirmedVerifiedAt = verifiedAt
	return m.confirmErr
}

func TestRequestContactVerificationUseCase_Execute_Success(t *testing.T) {
	repo := &contactVerificationRepositoryMock{}
	uc := NewRequestContactVerificationUseCase(repo)
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	uc.now = func() time.Time { return now }

	output, err := uc.Execute(context.Background(), RequestContactVerificationInput{
		Channel: " EMAIL ",
		Target:  " user@example.com ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output == nil {
		t.Fatal("expected output, got nil")
	}
	if repo.createCalls != 1 {
		t.Fatalf("expected CreateContactVerification to be called once, got %d", repo.createCalls)
	}
	if output.Channel != "email" {
		t.Fatalf("expected channel email, got %q", output.Channel)
	}
	if output.Target != "user@example.com" {
		t.Fatalf("expected trimmed target, got %q", output.Target)
	}
	if len(output.Token) != 6 {
		t.Fatalf("expected token length 6, got %q", output.Token)
	}
	if !output.ExpiresAt.Equal(now.Add(contactVerificationTTL)) {
		t.Fatalf("expected expires_at %v, got %v", now.Add(contactVerificationTTL), output.ExpiresAt)
	}
}

func TestRequestContactVerificationUseCase_Execute_InvalidInput(t *testing.T) {
	repo := &contactVerificationRepositoryMock{}
	uc := NewRequestContactVerificationUseCase(repo)

	output, err := uc.Execute(context.Background(), RequestContactVerificationInput{
		Channel: "fax",
		Target:  "user@example.com",
	})
	if !errors.Is(err, domain.ErrInvalidData) {
		t.Fatalf("expected ErrInvalidData, got %v", err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
	if repo.createCalls != 0 {
		t.Fatalf("expected repository not to be called, got %d calls", repo.createCalls)
	}
}

func TestConfirmContactVerificationUseCase_Execute_Success(t *testing.T) {
	verificationID := uuid.New()
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	repo := &contactVerificationRepositoryMock{
		findValue: &domain.ContactVerification{
			ID:        verificationID,
			Channel:   domain.ContactVerificationChannelPhone,
			Target:    "+5511999999999",
			Token:     "123456",
			ExpiresAt: now.Add(time.Minute),
			CreatedAt: now.Add(-time.Minute),
		},
	}
	uc := NewConfirmContactVerificationUseCase(repo)
	uc.now = func() time.Time { return now }

	output, err := uc.Execute(context.Background(), ConfirmContactVerificationInput{
		VerificationID: verificationID,
		Token:          "123456",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output == nil {
		t.Fatal("expected output, got nil")
	}
	if output.Channel != "phone" {
		t.Fatalf("expected channel phone, got %q", output.Channel)
	}
	if output.Target != "+5511999999999" {
		t.Fatalf("expected target +5511999999999, got %q", output.Target)
	}
	if output.VerificationToken == "" {
		t.Fatal("expected verification token to be set")
	}
	if repo.confirmCalls != 1 {
		t.Fatalf("expected confirm to be called once, got %d", repo.confirmCalls)
	}
	if repo.confirmedID != verificationID {
		t.Fatalf("expected confirmed ID %v, got %v", verificationID, repo.confirmedID)
	}
	if repo.confirmedToken == "" {
		t.Fatal("expected confirmed token to be set")
	}
	if !repo.confirmedVerifiedAt.Equal(now) {
		t.Fatalf("expected verified_at %v, got %v", now, repo.confirmedVerifiedAt)
	}
}

func TestConfirmContactVerificationUseCase_Execute_InvalidToken(t *testing.T) {
	verificationID := uuid.New()
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	repo := &contactVerificationRepositoryMock{
		findValue: &domain.ContactVerification{
			ID:        verificationID,
			Channel:   domain.ContactVerificationChannelEmail,
			Target:    "user@example.com",
			Token:     "123456",
			ExpiresAt: now.Add(time.Minute),
			CreatedAt: now.Add(-time.Minute),
		},
	}
	uc := NewConfirmContactVerificationUseCase(repo)
	uc.now = func() time.Time { return now }

	output, err := uc.Execute(context.Background(), ConfirmContactVerificationInput{
		VerificationID: verificationID,
		Token:          "000000",
	})
	if !errors.Is(err, domain.ErrInvalidVerificationToken) {
		t.Fatalf("expected ErrInvalidVerificationToken, got %v", err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
	if repo.confirmCalls != 0 {
		t.Fatalf("expected confirm not to be called, got %d calls", repo.confirmCalls)
	}
}
