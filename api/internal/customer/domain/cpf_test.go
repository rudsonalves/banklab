package domain

import "testing"

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

func TestNewCustomer_InvalidCPF(t *testing.T) {
	tests := []struct {
		name         string
		customerName string
		cpf          string
		expectedErr  error
	}{
		{
			name:         "empty name",
			customerName: "",
			cpf:          "12345678909",
			expectedErr:  ErrNameRequired,
		},
		{
			name:         "empty CPF",
			customerName: "John Doe",
			cpf:          "",
			expectedErr:  ErrCPFRequired,
		},
		{
			name:         "invalid CPF format",
			customerName: "John Doe",
			cpf:          "12345678901",
			expectedErr:  ErrCPFInvalid,
		},
		{
			name:         "valid formatted CPF",
			customerName: "John Doe",
			cpf:          "123.456.789-09",
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customer, err := NewCustomer(tt.customerName, tt.cpf)
			if err != tt.expectedErr {
				t.Errorf("NewCustomer() error = %v, want %v", err, tt.expectedErr)
			}
			if err == nil && customer == nil {
				t.Error("NewCustomer() returned nil customer with nil error")
			}
			if err == nil && customer.CPF != "12345678909" {
				t.Errorf("NewCustomer() CPF = %q, want %q", customer.CPF, "12345678909")
			}
		})
	}
}
