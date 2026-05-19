package domain

import (
	"context"

	"github.com/google/uuid"
)

// CustomerRepository defines the canonical customer persistence contract.
type CustomerRepository interface {
	Create(ctx context.Context, c *Customer) error
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	GetByID(ctx context.Context, id uuid.UUID) (*CustomerProfile, error)
}

type CustomerDocumentRepository interface {
	CreateDocument(ctx context.Context, document *CustomerDocument) error
	ExistsCPF(ctx context.Context, cpf string) (bool, error)
	GetPrimaryDocumentByCustomerID(ctx context.Context, customerID uuid.UUID) (*CustomerDocument, error)
	GetCPFByCustomerID(ctx context.Context, customerID uuid.UUID) (*CustomerDocument, error)
}
