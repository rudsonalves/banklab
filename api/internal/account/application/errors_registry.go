package application

import (
	"net/http"

	"github.com/seu-usuario/bank-api/internal/account/domain"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
)

func RegisterErrors() {
	sharederrors.RegisterDomainError(
		domain.ErrInvalidData,
		sharederrors.ErrCodeInvalidData,
		"Invalid data",
		http.StatusBadRequest,
	)

	sharederrors.RegisterDomainError(
		domain.ErrInvalidAmount,
		sharederrors.ErrCodeInvalidAmount,
		"Invalid amount",
		http.StatusBadRequest,
	)

	sharederrors.RegisterDomainError(
		domain.ErrAccountNotFound,
		sharederrors.ErrCodeAccountNotFound,
		"Account not found",
		http.StatusNotFound,
	)

	sharederrors.RegisterDomainError(
		domain.ErrCustomerNotFound,
		sharederrors.ErrCodeCustomerNotFound,
		"Customer not found",
		http.StatusNotFound,
	)

	sharederrors.RegisterDomainError(
		domain.ErrInsufficientBalance,
		sharederrors.ErrCodeInsufficientFunds,
		"Insufficient balance",
		http.StatusUnprocessableEntity,
	)

	sharederrors.RegisterDomainError(
		domain.ErrAccountInactive,
		sharederrors.ErrCodeAccountInactive,
		"Account is not active",
		http.StatusUnprocessableEntity,
	)

	sharederrors.RegisterDomainError(
		domain.ErrSameAccountTransfer,
		sharederrors.ErrCodeSameAccount,
		"Source and destination accounts must be different",
		http.StatusBadRequest,
	)

	sharederrors.RegisterDomainError(
		domain.ErrForbidden,
		sharederrors.ErrCodeForbidden,
		"Access denied",
		http.StatusForbidden,
	)
}
