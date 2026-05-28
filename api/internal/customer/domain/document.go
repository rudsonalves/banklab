package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	DocumentTypeCPF = "cpf"
	DefaultCountry  = "BR"
)

type CustomerDocument struct {
	ID          uuid.UUID
	CustomerID  uuid.UUID
	Type        string
	Value       string
	Issuer      *string
	IssuerState *string
	Country     string
	IsPrimary   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewCPFDocument creates a new CustomerDocument entity of type CPF with the
// provided customer ID, CPF value, and primary status.
// It validates the input parameters and returns an error if any validation
// fails, such as an empty customer ID,
// an empty CPF value, or an invalid CPF format. If the input is valid, it
// returns a pointer to the newly created
// CustomerDocument entity with the appropriate fields populated. The document ID
// is assigned by the persistence layer.
func NewCPFDocument(
	customerID uuid.UUID,
	cpf string,
	isPrimary bool,
) (*CustomerDocument, error) {
	if customerID == uuid.Nil {
		return nil, ErrInvalidData
	}

	normalizedCPF := normalizeCPF(cpf)
	if normalizedCPF == "" {
		return nil, ErrCPFRequired
	}

	if !ValidateCPF(normalizedCPF) {
		return nil, ErrCPFInvalid
	}

	now := time.Now().UTC()

	return &CustomerDocument{
		CustomerID: customerID,
		Type:       DocumentTypeCPF,
		Value:      normalizedCPF,
		Country:    DefaultCountry,
		IsPrimary:  isPrimary,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// NormalizeDocumentType normalizes the document type string by trimming
// whitespace and converting it to lowercase.
// This function ensures that document type comparisons are case-insensitive
// and not affected by leading or trailing whitespace.
func NormalizeDocumentType(documentType string) string {
	return strings.ToLower(strings.TrimSpace(documentType))
}
