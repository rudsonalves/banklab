package application

import (
	"net/http"

	"github.com/seu-usuario/bank-api/internal/customer/domain"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
)

// RegisterErrors registers the domain errors of the customer application with
// the shared error registry. This allows for consistent error handling and
// mapping to HTTP status codes.
func RegisterErrors() {
	sharederrors.RegisterDomainError(
		domain.ErrInvalidData,
		sharederrors.ErrCodeInvalidData,
		"Invalid data",
		http.StatusBadRequest,
	)

	sharederrors.RegisterDomainError(
		domain.ErrNotFound,
		sharederrors.ErrCodeCustomerNotFound,
		"Customer not found",
		http.StatusNotFound,
	)
	sharederrors.RegisterDomainError(
		domain.ErrCPFInvalid,
		sharederrors.ErrCodeInvalidData,
		"Invalid CPF format",
		http.StatusBadRequest,
	)

	sharederrors.RegisterDomainError(
		domain.ErrCPFAlreadyExists,
		sharederrors.ErrCodeUserAlreadyExists,
		"User already exists",
		http.StatusConflict,
	)
}
