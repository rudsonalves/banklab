package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewCustomer(t *testing.T) {
	birthDate := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)

	customer, err := NewCustomer(" Maria Silva ", birthDate)
	if err != nil {
		t.Fatalf("NewCustomer() error = %v", err)
	}

	if customer.Name != "Maria Silva" {
		t.Errorf("Name = %q, want %q", customer.Name, "Maria Silva")
	}
	if !customer.BirthDate.Equal(birthDate) {
		t.Errorf("BirthDate = %v, want %v", customer.BirthDate, birthDate)
	}
	if customer.ID == uuid.Nil {
		t.Error("ID is nil")
	}
	if customer.CreatedAt.IsZero() {
		t.Error("CreatedAt is zero")
	}
}

func TestNewCustomer_InvalidInput(t *testing.T) {
	validBirthDate := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		inputName string
		birthDate time.Time
		wantErr   error
	}{
		{
			name:      "nome vazio",
			inputName: "",
			birthDate: validBirthDate,
			wantErr:   ErrNameRequired,
		},
		{
			name:      "data de nascimento vazia",
			inputName: "Maria Silva",
			birthDate: time.Time{},
			wantErr:   ErrBirthDateRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCustomer(tt.inputName, tt.birthDate)
			if err != tt.wantErr {
				t.Errorf("NewCustomer() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
