package domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestNewAccount(t *testing.T) {
	tests := []struct {
		name       string
		customerID uuid.UUID
		number     string
		branch     string
		wantErr    error
	}{
		{
			name:       "returns ErrInvalidData when customerID is nil UUID",
			customerID: uuid.Nil,
			number:     "123456",
			branch:     "0001",
			wantErr:    ErrInvalidData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := NewAccount(tt.customerID, tt.number, tt.branch)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if account != nil {
				t.Fatalf("expected account to be nil, got %+v", account)
			}
		})
	}
}

func TestNewAccount_ValidInput(t *testing.T) {
	customerID := uuid.New()
	number := "123456"
	branch := "0001"

	account, err := NewAccount(customerID, number, branch)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if account == nil {
		t.Fatal("expected account to be non-nil")
	}

	if account.ID == uuid.Nil {
		t.Fatal("expected ID to be non-zero UUID")
	}

	if account.CustomerID != customerID {
		t.Fatalf("expected CustomerID %v, got %v", customerID, account.CustomerID)
	}

	if account.Number != number {
		t.Fatalf("expected Number %q, got %q", number, account.Number)
	}

	if account.Branch != branch {
		t.Fatalf("expected Branch %q, got %q", branch, account.Branch)
	}

	if account.Status != "active" {
		t.Fatalf("expected Status %q, got %q", "active", account.Status)
	}

	if account.Balance != 0 {
		t.Fatalf("expected Balance %d, got %d", 0, account.Balance)
	}

	if account.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be non-zero")
	}
}

func TestAccount_CanDeposit(t *testing.T) {
	tests := []struct {
		name    string
		account Account
		amount  int64
		wantErr error
	}{
		{
			name:    "invalid amount",
			account: Account{Status: AccountActive},
			amount:  0,
			wantErr: ErrInvalidAmount,
		},
		{
			name:    "inactive account",
			account: Account{Status: AccountInactive},
			amount:  10,
			wantErr: ErrAccountInactive,
		},
		{
			name:    "success",
			account: Account{Status: AccountActive},
			amount:  10,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.account.CanDeposit(tt.amount)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestAccount_CanWithdraw(t *testing.T) {
	tests := []struct {
		name    string
		account Account
		amount  int64
		wantErr error
	}{
		{
			name:    "invalid amount",
			account: Account{Status: AccountActive, Balance: 100},
			amount:  0,
			wantErr: ErrInvalidAmount,
		},
		{
			name:    "inactive account",
			account: Account{Status: AccountInactive, Balance: 100},
			amount:  10,
			wantErr: ErrAccountInactive,
		},
		{
			name:    "insufficient balance",
			account: Account{Status: AccountActive, Balance: 50},
			amount:  100,
			wantErr: ErrInsufficientBalance,
		},
		{
			name:    "success",
			account: Account{Status: AccountActive, Balance: 100},
			amount:  10,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.account.CanWithdraw(tt.amount)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}
