package application

import (
	"net/http"

	"github.com/seu-usuario/bank-api/internal/customer/domain"
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
		domain.ErrCPFAlreadyExists,
		sharederrors.ErrCodeUserAlreadyExists,
		"User already exists",
		http.StatusConflict,
	)

	sharederrors.RegisterDomainError(
		domain.ErrEmailAlreadyExists,
		sharederrors.ErrCodeUserAlreadyExists,
		"User already exists",
		http.StatusConflict,
	)
}
