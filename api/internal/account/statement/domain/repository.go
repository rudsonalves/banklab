package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	bankaccountdomain "github.com/seu-usuario/bank-api/internal/account/bankaccount/domain"
	transactiondomain "github.com/seu-usuario/bank-api/internal/account/transaction/domain"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*bankaccountdomain.Account, error)
	GetTransactions(
		ctx context.Context,
		accountID uuid.UUID,
		limit int,
		cursorTime *time.Time,
		cursorID *uuid.UUID,
		from *time.Time,
		to *time.Time,
	) ([]transactiondomain.Transaction, error)
}
