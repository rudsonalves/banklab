package accounterrors

import (
	"net/http"

	bankaccountdomain "github.com/seu-usuario/bank-api/internal/account/bankaccount/domain"
	transactiondomain "github.com/seu-usuario/bank-api/internal/account/transaction/domain"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
)

func RegisterErrors() {
	sharederrors.RegisterDomainError(
		bankaccountdomain.ErrInvalidData,
		sharederrors.ErrCodeInvalidData,
		"Invalid data",
		http.StatusBadRequest,
	)
	sharederrors.RegisterDomainError(
		transactiondomain.ErrInvalidData,
		sharederrors.ErrCodeInvalidData,
		"Invalid data",
		http.StatusBadRequest,
	)

	sharederrors.RegisterDomainError(
		transactiondomain.ErrInvalidAmount,
		sharederrors.ErrCodeInvalidAmount,
		"Invalid amount",
		http.StatusBadRequest,
	)

	sharederrors.RegisterDomainError(
		bankaccountdomain.ErrAccountNotFound,
		sharederrors.ErrCodeAccountNotFound,
		"Account not found",
		http.StatusNotFound,
	)
	sharederrors.RegisterDomainError(
		transactiondomain.ErrAccountNotFound,
		sharederrors.ErrCodeAccountNotFound,
		"Account not found",
		http.StatusNotFound,
	)

	sharederrors.RegisterDomainError(
		bankaccountdomain.ErrCustomerNotFound,
		sharederrors.ErrCodeCustomerNotFound,
		"Customer not found",
		http.StatusNotFound,
	)

	sharederrors.RegisterDomainError(
		transactiondomain.ErrInsufficientBalance,
		sharederrors.ErrCodeInsufficientFunds,
		"Insufficient balance",
		http.StatusUnprocessableEntity,
	)

	sharederrors.RegisterDomainError(
		transactiondomain.ErrAccountInactive,
		sharederrors.ErrCodeAccountInactive,
		"Account is not active",
		http.StatusUnprocessableEntity,
	)

	sharederrors.RegisterDomainError(
		transactiondomain.ErrSameAccountTransfer,
		sharederrors.ErrCodeSameAccount,
		"Source and destination accounts must be different",
		http.StatusBadRequest,
	)

	sharederrors.RegisterDomainError(
		bankaccountdomain.ErrForbidden,
		sharederrors.ErrCodeForbidden,
		"Access denied",
		http.StatusForbidden,
	)
	sharederrors.RegisterDomainError(
		transactiondomain.ErrForbidden,
		sharederrors.ErrCodeForbidden,
		"Access denied",
		http.StatusForbidden,
	)
}
