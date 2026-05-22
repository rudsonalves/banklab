package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestValidateCPF(t *testing.T) {
	tests := []struct {
		name     string
		cpf      string
		expected bool
	}{
		{
			name:     "valid CPF with 11 digits",
			cpf:      "12345678909",
			expected: true,
		},
		{
			name:     "valid CPF with formatting",
			cpf:      "123.456.789-09",
			expected: true,
		},
		{
			name:     "invalid CPF with 10 digits",
			cpf:      "1234567890",
			expected: false,
		},
		{
			name:     "invalid CPF with 12 digits",
			cpf:      "123456789090",
			expected: false,
		},
		{
			name:     "invalid CPF with all same digits (11111111111)",
			cpf:      "11111111111",
			expected: false,
		},
		{
			name:     "invalid CPF with all same digits (00000000000)",
			cpf:      "00000000000",
			expected: false,
		},
		{
			name:     "invalid CPF with wrong first check digit",
			cpf:      "12345678900",
			expected: false,
		},
		{
			name:     "invalid CPF with wrong second check digit",
			cpf:      "12345678908",
			expected: false,
		},
		{
			name:     "empty CPF",
			cpf:      "",
			expected: false,
		},
		{
			name:     "CPF with spaces and formatting",
			cpf:      "123.456.789 - 09",
			expected: true,
		},
		{
			name:     "real valid CPF example",
			cpf:      "98765432100",
			expected: true,
		},
		{
			name:     "CPF with letters",
			cpf:      "1234567890a",
			expected: false,
		},
		{
			name:     "CPF with special characters",
			cpf:      "123@456#789$9",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateCPF(tt.cpf)
			if result != tt.expected {
				t.Errorf("ValidateCPF(%q) = %v, want %v", tt.cpf, result, tt.expected)
			}
		})
	}
}

func TestNewCPFDocument(t *testing.T) {
	customerID := uuid.New()

	document, err := NewCPFDocument(customerID, "123.456.789-09", true)
	if err != nil {
		t.Fatalf("NewCPFDocument() error = %v", err)
	}

	if document.CustomerID != customerID {
		t.Errorf("CustomerID = %v, want %v", document.CustomerID, customerID)
	}
	if document.Type != DocumentTypeCPF {
		t.Errorf("Type = %q, want %q", document.Type, DocumentTypeCPF)
	}
	if document.Value != "12345678909" {
		t.Errorf("Value = %q, want %q", document.Value, "12345678909")
	}
	if document.Country != DefaultCountry {
		t.Errorf("Country = %q, want %q", document.Country, DefaultCountry)
	}
	if !document.IsPrimary {
		t.Error("IsPrimary = false, want true")
	}
	if document.Issuer != nil {
		t.Errorf("Issuer = %v, want nil", document.Issuer)
	}
	if document.IssuerState != nil {
		t.Errorf("IssuerState = %v, want nil", document.IssuerState)
	}
}

func TestNewCPFDocument_InvalidInput(t *testing.T) {
	tests := []struct {
		name       string
		customerID uuid.UUID
		cpf        string
		wantErr    error
	}{
		{
			name:       "customer id vazio",
			customerID: uuid.Nil,
			cpf:        "12345678909",
			wantErr:    ErrInvalidData,
		},
		{
			name:       "cpf vazio",
			customerID: uuid.New(),
			cpf:        "",
			wantErr:    ErrCPFRequired,
		},
		{
			name:       "cpf inválido",
			customerID: uuid.New(),
			cpf:        "12345678901",
			wantErr:    ErrCPFInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCPFDocument(tt.customerID, tt.cpf, true)
			if err != tt.wantErr {
				t.Errorf("NewCPFDocument() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
