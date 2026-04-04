package application

import (
	"errors"

	"github.com/seu-usuario/bank-api/internal/account/domain"
)

type ErrorCategory string

const (
	ErrorCategoryUnknown             ErrorCategory = "unknown"
	ErrorCategoryForbidden          ErrorCategory = "forbidden"
	ErrorCategoryInvalidData        ErrorCategory = "invalid_data"
	ErrorCategoryInvalidAmount      ErrorCategory = "invalid_amount"
	ErrorCategoryCustomerNotFound   ErrorCategory = "customer_not_found"
	ErrorCategoryAccountNotFound    ErrorCategory = "account_not_found"
	ErrorCategoryAccountInactive    ErrorCategory = "account_inactive"
	ErrorCategoryInsufficientAmount ErrorCategory = "insufficient_balance"
	ErrorCategorySameAccount        ErrorCategory = "same_account_transfer"
)

func CategorizeError(err error) ErrorCategory {
	switch {
	case errors.Is(err, domain.ErrForbidden):
		return ErrorCategoryForbidden
	case errors.Is(err, domain.ErrInvalidData):
		return ErrorCategoryInvalidData
	case errors.Is(err, domain.ErrInvalidAmount):
		return ErrorCategoryInvalidAmount
	case errors.Is(err, domain.ErrCustomerNotFound):
		return ErrorCategoryCustomerNotFound
	case errors.Is(err, domain.ErrAccountNotFound):
		return ErrorCategoryAccountNotFound
	case errors.Is(err, domain.ErrAccountInactive):
		return ErrorCategoryAccountInactive
	case errors.Is(err, domain.ErrInsufficientBalance):
		return ErrorCategoryInsufficientAmount
	case errors.Is(err, domain.ErrSameAccountTransfer):
		return ErrorCategorySameAccount
	default:
		return ErrorCategoryUnknown
	}
}
