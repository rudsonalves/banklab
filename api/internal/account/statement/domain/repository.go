package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	accountdomain "github.com/seu-usuario/bank-api/internal/account/domain"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*accountdomain.Account, error)
	GetTransactions(
		ctx context.Context,
		accountID uuid.UUID,
		limit int,
		cursorTime *time.Time,
		cursorID *uuid.UUID,
		from *time.Time,
		to *time.Time,
	) ([]accountdomain.Transaction, error)
}
