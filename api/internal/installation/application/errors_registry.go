package application

import (
	"net/http"

	"github.com/seu-usuario/bank-api/internal/installation/domain"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
)

func RegisterErrors() {
	sharederrors.RegisterDomainError(
		domain.ErrInstallationMismatch,
		sharederrors.ErrCodeInstallationMismatch,
		"Installation mismatch",
		http.StatusForbidden,
	)

	sharederrors.RegisterDomainError(
		domain.ErrInvalidInstallationResourceID,
		sharederrors.ErrCodeInvalidData,
		"Invalid data",
		http.StatusBadRequest,
	)

	sharederrors.RegisterDomainError(
		domain.ErrInstallationNotFound,
		sharederrors.ErrCodeInvalidData,
		"Invalid data",
		http.StatusNotFound,
	)

	sharederrors.RegisterDomainError(
		domain.ErrRestrictedAuthorizationNotFound,
		sharederrors.ErrCodeInvalidToken,
		"Invalid token",
		http.StatusUnauthorized,
	)

	sharederrors.RegisterDomainError(
		domain.ErrRestrictedAuthorizationConsumed,
		sharederrors.ErrCodeInvalidToken,
		"Invalid token",
		http.StatusUnauthorized,
	)

	sharederrors.RegisterDomainError(
		domain.ErrRestrictedAuthorizationRevoked,
		sharederrors.ErrCodeInvalidToken,
		"Invalid token",
		http.StatusUnauthorized,
	)
}
