package domain

import (
	"time"
	"unicode"

	"github.com/google/uuid"
)

type TransactionPasswordStatus string

const (
	TransactionPasswordActive  TransactionPasswordStatus = "active"
	TransactionPasswordBlocked TransactionPasswordStatus = "blocked"

	TransactionPasswordPINLength    = 6
	TransactionPasswordMaxFailures  = 3
	TransactionPasswordLockDuration = 30 * time.Minute
)

type TransactionPassword struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	PasswordHash   string
	Status         TransactionPasswordStatus
	FailedAttempts int
	LockedUntil    *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	ChangedAt *time.Time
}

func NewTransactionPassword(
	userID uuid.UUID,
	passwordHash string,
	now time.Time,
) (*TransactionPassword, error) {
	if userID == uuid.Nil || passwordHash == "" || now.IsZero() {
		return nil, ErrInvalidTransactionPassword
	}

	changedAt := now.UTC()

	return &TransactionPassword{
		UserID:       userID,
		PasswordHash: passwordHash,
		Status:       TransactionPasswordActive,
		CreatedAt:    changedAt,
		UpdatedAt:    changedAt,
		ChangedAt:    &changedAt,
	}, nil
}

func ValidateTransactionPasswordPIN(pin string) error {
	if len(pin) != TransactionPasswordPINLength {
		return ErrInvalidTransactionPasswordPIN
	}

	for _, char := range pin {
		if !unicode.IsDigit(char) {
			return ErrInvalidTransactionPasswordPIN
		}
	}

	return nil
}

// NormalizeLock checks if the transaction password is currently blocked and
// if the lock duration has expired.
func (p *TransactionPassword) NormalizeLock(now time.Time) {
	if p.Status != TransactionPasswordBlocked || p.LockedUntil == nil {
		return
	}

	now = now.UTC()

	if now.Before(*p.LockedUntil) {
		return
	}

	p.Status = TransactionPasswordActive
	p.FailedAttempts = 0
	p.LockedUntil = nil
	p.UpdatedAt = now
}

// CanValidate checks if the transaction password can be validated
// (i.e., not blocked).
func (p *TransactionPassword) CanValidate(now time.Time) error {
	p.NormalizeLock(now)

	if p.Status == TransactionPasswordBlocked {
		return ErrTransactionPasswordLocked
	}

	return nil
}

// RegisterFailure registers a failed validation attempt and updates the
// status if necessary.
func (p *TransactionPassword) RegisterFailure(now time.Time) error {
	now = now.UTC()

	p.NormalizeLock(now)

	p.FailedAttempts++
	p.UpdatedAt = now

	if p.FailedAttempts >= TransactionPasswordMaxFailures {
		lockedUntil := now.Add(TransactionPasswordLockDuration)

		p.Status = TransactionPasswordBlocked
		p.LockedUntil = &lockedUntil
		p.FailedAttempts = TransactionPasswordMaxFailures

		return ErrTransactionPasswordLocked
	}

	return ErrTransactionPasswordInvalid
}

// RegisterSuccess resets the failed attempts and unlocks the transaction
// password if it was blocked.
func (p *TransactionPassword) RegisterSuccess(now time.Time) {
	now = now.UTC()

	p.Status = TransactionPasswordActive
	p.FailedAttempts = 0
	p.LockedUntil = nil
	p.UpdatedAt = now
}
