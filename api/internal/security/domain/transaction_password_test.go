package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateTransactionPasswordPIN(t *testing.T) {
	tests := []struct {
		name    string
		pin     string
		wantErr bool
	}{
		{name: "valid six digits", pin: "123456", wantErr: false},
		{name: "empty", pin: "", wantErr: true},
		{name: "too short", pin: "12345", wantErr: true},
		{name: "too long", pin: "1234567", wantErr: true},
		{name: "letters", pin: "12345a", wantErr: true},
		{name: "spaces", pin: "123 56", wantErr: true},
		{name: "symbols", pin: "12345-", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTransactionPasswordPIN(tt.pin)

			if tt.wantErr && !errors.Is(err, ErrInvalidTransactionPasswordPIN) {
				t.Fatalf("expected ErrInvalidTransactionPasswordPIN, got %v", err)
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestNewTransactionPassword(t *testing.T) {
	now := time.Date(2026, 5, 28, 10, 0, 0, 0, time.UTC)
	userID := uuid.New()

	password, err := NewTransactionPassword(userID, "hash-value", now)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if password.ID != uuid.Nil {
		t.Fatalf("expected empty ID before persistence, got %v", password.ID)
	}

	if password.UserID != userID {
		t.Fatalf("expected user id %v, got %v", userID, password.UserID)
	}

	if password.Status != TransactionPasswordActive {
		t.Fatalf("expected active status, got %s", password.Status)
	}

	if password.FailedAttempts != 0 {
		t.Fatalf("expected zero failed attempts, got %d", password.FailedAttempts)
	}

	if password.LockedUntil != nil {
		t.Fatalf("expected nil locked_until, got %v", password.LockedUntil)
	}
}

func TestNewTransactionPassword_InvalidData(t *testing.T) {
	now := time.Date(2026, 5, 28, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		userID       uuid.UUID
		passwordHash string
		now          time.Time
	}{
		{name: "empty user id", userID: uuid.Nil, passwordHash: "hash", now: now},
		{name: "empty hash", userID: uuid.New(), passwordHash: "", now: now},
		{name: "empty time", userID: uuid.New(), passwordHash: "hash", now: time.Time{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := NewTransactionPassword(tt.userID, tt.passwordHash, tt.now)

			if !errors.Is(err, ErrInvalidTransactionPassword) {
				t.Fatalf("expected ErrInvalidTransactionPassword, got %v", err)
			}

			if password != nil {
				t.Fatalf("expected nil password, got %+v", password)
			}
		})
	}
}

func TestTransactionPassword_RegisterFailureLocksAfterThreeAttempts(t *testing.T) {
	now := time.Date(2026, 5, 28, 10, 0, 0, 0, time.UTC)
	password, err := NewTransactionPassword(uuid.New(), "hash-value", now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = password.RegisterFailure(now.Add(time.Minute))
	if !errors.Is(err, ErrTransactionPasswordInvalid) {
		t.Fatalf("expected invalid error, got %v", err)
	}

	err = password.RegisterFailure(now.Add(2 * time.Minute))
	if !errors.Is(err, ErrTransactionPasswordInvalid) {
		t.Fatalf("expected invalid error, got %v", err)
	}

	lockTime := now.Add(3 * time.Minute)
	err = password.RegisterFailure(lockTime)
	if !errors.Is(err, ErrTransactionPasswordLocked) {
		t.Fatalf("expected locked error, got %v", err)
	}

	if password.Status != TransactionPasswordBlocked {
		t.Fatalf("expected blocked status, got %s", password.Status)
	}

	if password.FailedAttempts != TransactionPasswordMaxFailures {
		t.Fatalf("expected %d failed attempts, got %d", TransactionPasswordMaxFailures, password.FailedAttempts)
	}

	if password.LockedUntil == nil {
		t.Fatal("expected locked_until to be set")
	}

	expectedLockedUntil := lockTime.Add(TransactionPasswordLockDuration)
	if !password.LockedUntil.Equal(expectedLockedUntil) {
		t.Fatalf("expected locked_until %v, got %v", expectedLockedUntil, *password.LockedUntil)
	}
}

func TestTransactionPassword_CanValidateReturnsLockedBeforeExpiration(t *testing.T) {
	now := time.Date(2026, 5, 28, 10, 0, 0, 0, time.UTC)
	lockedUntil := now.Add(TransactionPasswordLockDuration)

	password := &TransactionPassword{
		UserID:         uuid.New(),
		PasswordHash:   "hash-value",
		Status:         TransactionPasswordBlocked,
		FailedAttempts: TransactionPasswordMaxFailures,
		LockedUntil:    &lockedUntil,
	}

	err := password.CanValidate(now.Add(10 * time.Minute))

	if !errors.Is(err, ErrTransactionPasswordLocked) {
		t.Fatalf("expected ErrTransactionPasswordLocked, got %v", err)
	}
}

func TestTransactionPassword_NormalizeLockUnlocksAfterExpiration(t *testing.T) {
	now := time.Date(2026, 5, 28, 10, 0, 0, 0, time.UTC)
	lockedUntil := now.Add(TransactionPasswordLockDuration)

	password := &TransactionPassword{
		UserID:         uuid.New(),
		PasswordHash:   "hash-value",
		Status:         TransactionPasswordBlocked,
		FailedAttempts: TransactionPasswordMaxFailures,
		LockedUntil:    &lockedUntil,
	}

	password.NormalizeLock(lockedUntil)

	if password.Status != TransactionPasswordActive {
		t.Fatalf("expected active status, got %s", password.Status)
	}

	if password.FailedAttempts != 0 {
		t.Fatalf("expected failed attempts reset, got %d", password.FailedAttempts)
	}

	if password.LockedUntil != nil {
		t.Fatalf("expected locked_until reset, got %v", password.LockedUntil)
	}
}

func TestTransactionPassword_RegisterSuccessResetsFailures(t *testing.T) {
	now := time.Date(2026, 5, 28, 10, 0, 0, 0, time.UTC)

	password := &TransactionPassword{
		UserID:         uuid.New(),
		PasswordHash:   "hash-value",
		Status:         TransactionPasswordActive,
		FailedAttempts: 2,
		UpdatedAt:      now,
	}

	successAt := now.Add(5 * time.Minute)
	password.RegisterSuccess(successAt)

	if password.Status != TransactionPasswordActive {
		t.Fatalf("expected active status, got %s", password.Status)
	}

	if password.FailedAttempts != 0 {
		t.Fatalf("expected failed attempts reset, got %d", password.FailedAttempts)
	}

	if password.LockedUntil != nil {
		t.Fatalf("expected locked_until nil, got %v", password.LockedUntil)
	}

	if !password.UpdatedAt.Equal(successAt) {
		t.Fatalf("expected updated_at %v, got %v", successAt, password.UpdatedAt)
	}
}
