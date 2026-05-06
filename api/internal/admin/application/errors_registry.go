package application

import (
	"net/http"

	accountdomain "github.com/seu-usuario/bank-api/internal/account/bankaccount/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
)

// RegisterErrors registers the domain errors of the admin application with the
// shared error registry. This allows for consistent error handling and mapping to
// HTTP status codes across the application.
func RegisterErrors() {
	sharederrors.RegisterDomainError(
		authdomain.ErrForbidden,
		sharederrors.ErrCodeForbidden,
		"Access denied",
		http.StatusForbidden,
	)

	sharederrors.RegisterDomainError(
		authdomain.ErrUnauthorized,
		sharederrors.ErrCodeUnauthorized,
		"Authentication required",
		http.StatusUnauthorized,
	)

	sharederrors.RegisterDomainError(
		authdomain.ErrInvalidUserState,
		sharederrors.ErrCodeInvalidUserState,
		"Invalid user state",
		http.StatusConflict,
	)

	sharederrors.RegisterDomainError(
		authdomain.ErrUserNotFound,
		sharederrors.ErrCodeUserNotFound,
		"User not found",
		http.StatusNotFound,
	)

	sharederrors.RegisterDomainError(
		authdomain.ErrUserAlreadyActive,
		sharederrors.ErrCodeUserAlreadyActive,
		"User is already active",
		http.StatusConflict,
	)

	sharederrors.RegisterDomainError(
		accountdomain.ErrCustomerNotFound,
		sharederrors.ErrCodeCustomerNotFound,
		"Customer not found",
		http.StatusNotFound,
	)
}
