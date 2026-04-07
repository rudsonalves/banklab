package application

import (
	"net/http"

	"github.com/seu-usuario/bank-api/internal/auth/domain"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
)

func RegisterErrors() {
	sharederrors.RegisterDomainError(
		domain.ErrEmailAlreadyExists,
		sharederrors.ErrCodeUserAlreadyExists,
		"User already exists",
		http.StatusConflict,
	)

	sharederrors.RegisterDomainError(
		domain.ErrInvalidCredentials,
		sharederrors.ErrCodeInvalidCredentials,
		"Invalid credentials",
		http.StatusUnauthorized,
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
		domain.ErrInvalidPassword,
		sharederrors.ErrCodeInvalidData,
		"Invalid data",
		http.StatusBadRequest,
	)

	sharederrors.RegisterDomainError(
		domain.ErrInvalidUserState,
		sharederrors.ErrCodeInvalidUserState,
		"Invalid user state",
		http.StatusInternalServerError,
	)
}
