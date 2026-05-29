package application

import (
	"net/http"

	"github.com/seu-usuario/bank-api/internal/security/domain"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
)

func RegisterErrors() {
	sharederrors.RegisterDomainError(
		domain.ErrTransactionPasswordAlreadySet,
		sharederrors.ErrCodeTransactionPasswordAlreadySet,
		"Transaction password already set",
		http.StatusConflict,
	)

	sharederrors.RegisterDomainError(
		domain.ErrTransactionPasswordNotSet,
		sharederrors.ErrCodeTransactionPasswordNotSet,
		"Transaction password not set",
		http.StatusConflict,
	)

	sharederrors.RegisterDomainError(
		domain.ErrTransactionPasswordInvalid,
		sharederrors.ErrCodeTransactionPasswordInvalid,
		"Invalid transaction password",
		http.StatusUnauthorized,
	)

	sharederrors.RegisterDomainError(
		domain.ErrInvalidTransactionPasswordPIN,
		sharederrors.ErrCodeInvalidData,
		"Invalid data",
		http.StatusBadRequest,
	)

	sharederrors.RegisterDomainError(
		domain.ErrInvalidTransactionPassword,
		sharederrors.ErrCodeInvalidData,
		"Invalid data",
		http.StatusBadRequest,
	)

	sharederrors.RegisterDomainError(
		domain.ErrTransactionPasswordLocked,
		sharederrors.ErrCodeTransactionPasswordLocked,
		"Transaction password locked",
		http.StatusForbidden,
	)

	sharederrors.RegisterDomainError(
		domain.ErrStepUpEndpointNotAllowed,
		sharederrors.ErrCodeStepUpEndpointNotAllowed,
		"Step-up endpoint not allowed",
		http.StatusForbidden,
	)

	sharederrors.RegisterDomainError(
		domain.ErrStepUpTokenRequired,
		sharederrors.ErrCodeStepUpTokenRequired,
		"Step-up token required",
		http.StatusUnauthorized,
	)
}
