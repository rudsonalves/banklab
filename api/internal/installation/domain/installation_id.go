package domain

import (
	"strings"

	"github.com/google/uuid"
)

type InstallationID struct {
	value uuid.UUID
}

func NewInstallationID(value uuid.UUID) (InstallationID, error) {
	if value == uuid.Nil || value.Version() != 4 {
		return InstallationID{}, ErrInvalidInstallationID
	}

	return InstallationID{value: value}, nil
}

func ParseInstallationID(raw string) (InstallationID, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return InstallationID{}, ErrInvalidInstallationID
	}

	parsed, err := uuid.Parse(value)
	if err != nil {
		return InstallationID{}, ErrInvalidInstallationID
	}
	if parsed.Version() != 4 || parsed.String() != value {
		return InstallationID{}, ErrInvalidInstallationID
	}

	return InstallationID{value: parsed}, nil
}

func (id InstallationID) UUID() uuid.UUID {
	return id.value
}

func (id InstallationID) String() string {
	if id.value == uuid.Nil {
		return ""
	}

	return id.value.String()
}

func (id InstallationID) IsZero() bool {
	return id.value == uuid.Nil
}

type InstallationResourceID struct {
	value uuid.UUID
}

func NewInstallationResourceID(value uuid.UUID) (InstallationResourceID, error) {
	if value == uuid.Nil {
		return InstallationResourceID{}, ErrInvalidInstallationResourceID
	}

	return InstallationResourceID{value: value}, nil
}

func ParseInstallationResourceID(raw string) (InstallationResourceID, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return InstallationResourceID{}, ErrInvalidInstallationResourceID
	}

	parsed, err := uuid.Parse(value)
	if err != nil || parsed == uuid.Nil || parsed.String() != value {
		return InstallationResourceID{}, ErrInvalidInstallationResourceID
	}

	return InstallationResourceID{value: parsed}, nil
}

func (id InstallationResourceID) UUID() uuid.UUID {
	return id.value
}

func (id InstallationResourceID) String() string {
	if id.value == uuid.Nil {
		return ""
	}

	return id.value.String()
}

func (id InstallationResourceID) IsZero() bool {
	return id.value == uuid.Nil
}
