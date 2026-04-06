package domain

import (
	"time"

	"github.com/google/uuid"
)

// Operation represents a logical operation persisted in the transactions table
// for idempotency and request replay handling.
type Operation struct {
	ID               uuid.UUID
	AccountID        uuid.UUID
	Type             TransactionType
	Amount           int64
	Description      *string
	RelatedAccountID *uuid.UUID
	ReferenceID      *uuid.UUID
	IdempotencyKey   *string
	CreatedAt        time.Time
}

func NewOperation(
	accountID uuid.UUID,
	typeName TransactionType,
	amount int64,
	relatedAccountID *uuid.UUID,
	referenceID *uuid.UUID,
	idempotencyKey string,
) *Operation {
	key := idempotencyKey
	return &Operation{
		ID:               uuid.New(),
		AccountID:        accountID,
		Type:             typeName,
		Amount:           amount,
		RelatedAccountID: relatedAccountID,
		ReferenceID:      referenceID,
		IdempotencyKey:   &key,
		CreatedAt:        time.Now().UTC(),
	}
}
