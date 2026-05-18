package application

import (
	"net/http"

	"github.com/seu-usuario/bank-api/internal/auth/domain"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
)

// RegisterErrors registers the domain errors of the auth application with the
// shared error registry. This allows for consistent error handling and mapping to
// HTTP status codes across the application.
func RegisterErrors() {
	sharederrors.RegisterDomainError(
		domain.ErrEmailAlreadyExists,
		sharederrors.ErrCodeUserAlreadyExists,
		"User already exists",
		http.StatusConflict,
	)

	sharederrors.RegisterDomainError(
		domain.ErrForbidden,
		sharederrors.ErrCodeForbidden,
		"Access denied",
		http.StatusForbidden,
	)

	sharederrors.RegisterDomainError(
		domain.ErrInvalidCredentials,
		sharederrors.ErrCodeInvalidCredentials,
		"Invalid credentials",
		http.StatusUnauthorized,
	)

	sharederrors.RegisterDomainError(
		domain.ErrAccountApprovalRequired,
		sharederrors.ErrCodeAccountApprovalRequired,
		"Account approval required",
		http.StatusForbidden,
	)

	sharederrors.RegisterDomainError(
		domain.ErrUnauthorized,
		sharederrors.ErrCodeUnauthorized,
		"Authentication required",
		http.StatusUnauthorized,
	)

	sharederrors.RegisterDomainError(
		domain.ErrInvalidToken,
		sharederrors.ErrCodeInvalidToken,
		"Invalid token",
		http.StatusUnauthorized,
	)

	sharederrors.RegisterDomainError(
		domain.ErrInvalidEmail,
		sharederrors.ErrCodeInvalidData,
		"Invalid data",
		http.StatusBadRequest,
	)

	sharederrors.RegisterDomainError(
		domain.ErrInvalidData,
		sharederrors.ErrCodeInvalidData,
		"Invalid data",
		http.StatusBadRequest,
	)

	sharederrors.RegisterDomainError(
		domain.ErrInvalidPassword,
		sharederrors.ErrCodeInvalidData,
		"Invalid data",
		http.StatusBadRequest,
	)

	sharederrors.RegisterDomainError(
		domain.ErrInvalidUserState,
		sharederrors.ErrCodeInvalidUserState,
		"Invalid user state",
		http.StatusConflict,
	)

	sharederrors.RegisterDomainError(
		domain.ErrUserNotFound,
		sharederrors.ErrCodeUserNotFound,
		"User not found",
		http.StatusNotFound,
	)

	sharederrors.RegisterDomainError(
		domain.ErrUserAlreadyActive,
		sharederrors.ErrCodeUserAlreadyActive,
		"User is already active",
		http.StatusConflict,
	)

	sharederrors.RegisterDomainError(
		domain.ErrContactVerificationNotFound,
		sharederrors.ErrCodeInvalidData,
		"Invalid data",
		http.StatusBadRequest,
	)

	sharederrors.RegisterDomainError(
		domain.ErrInvalidVerificationToken,
		sharederrors.ErrCodeInvalidData,
		"Invalid data",
		http.StatusBadRequest,
	)

	sharederrors.RegisterDomainError(
		domain.ErrContactVerificationExpired,
		sharederrors.ErrCodeInvalidData,
		"Invalid data",
		http.StatusBadRequest,
	)
}
