package domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestParseInstallationID(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "valid canonical v4", value: "550e8400-e29b-41d4-a716-446655440000"},
		{name: "blank", value: "", wantErr: true},
		{name: "not uuid", value: "not-a-uuid", wantErr: true},
		{name: "uuid v1", value: "550e8400-e29b-11d4-a716-446655440000", wantErr: true},
		{name: "without hyphens", value: "550e8400e29b41d4a716446655440000", wantErr: true},
		{name: "uppercase", value: "550E8400-E29B-41D4-A716-446655440000", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInstallationID(tt.value)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidInstallationID) {
					t.Fatalf("expected ErrInvalidInstallationID, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if got.String() != tt.value {
				t.Fatalf("expected %q, got %q", tt.value, got.String())
			}
		})
	}
}

func TestInstallationResourceIDIsSeparateFromInstallationID(t *testing.T) {
	raw := uuid.New()

	resourceID, err := NewInstallationResourceID(raw)
	if err != nil {
		t.Fatalf("expected resource id, got %v", err)
	}
	installationID, err := NewInstallationID(raw)
	if err != nil {
		t.Fatalf("expected installation id, got %v", err)
	}

	if resourceID.String() != installationID.String() {
		t.Fatalf("expected same raw UUID string, got resource=%q installation=%q", resourceID.String(), installationID.String())
	}
}

func TestNewInstallationResourceIDRejectsNil(t *testing.T) {
	_, err := NewInstallationResourceID(uuid.Nil)
	if !errors.Is(err, ErrInvalidInstallationResourceID) {
		t.Fatalf("expected ErrInvalidInstallationResourceID, got %v", err)
	}
}
