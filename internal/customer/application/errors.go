package application

import (
	"errors"

	"github.com/seu-usuario/bank-api/internal/customer/domain"
)

type ErrorCategory string

const (
	ErrorCategoryUnknown          ErrorCategory = "unknown"
	ErrorCategoryInvalidData      ErrorCategory = "invalid_data"
	ErrorCategoryAlreadyExistsCPF ErrorCategory = "already_exists_cpf"
	ErrorCategoryAlreadyExistsEML ErrorCategory = "already_exists_email"
)

func CategorizeError(err error) ErrorCategory {
	switch {
	case errors.Is(err, domain.ErrInvalidData):
		return ErrorCategoryInvalidData
	case errors.Is(err, domain.ErrCPFAlreadyExists):
		return ErrorCategoryAlreadyExistsCPF
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		return ErrorCategoryAlreadyExistsEML
	default:
		return ErrorCategoryUnknown
	}
}

func ValidationField(err error) (string, bool) {
	switch {
	case errors.Is(err, domain.ErrNameRequired):
		return "name", true
	case errors.Is(err, domain.ErrCPFRequired):
		return "cpf", true
	case errors.Is(err, domain.ErrEmailRequired):
		return "email", true
	default:
		return "", false
	}
}
