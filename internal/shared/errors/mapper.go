package sharederrors

import (
	"errors"
	"net/http"

	accountdomain "github.com/seu-usuario/bank-api/internal/account/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	customerdomain "github.com/seu-usuario/bank-api/internal/customer/domain"
)

var ErrInvalidRequest = errors.New("invalid request")

func MapError(err error) AppError {
	switch {
	case err == nil:
		return AppError{Code: ErrCodeInternal, Message: "Internal server error", Status: http.StatusInternalServerError}
	case errors.Is(err, ErrInvalidRequest):
		return AppError{Code: ErrCodeInvalidRequest, Message: "Invalid request body", Status: http.StatusBadRequest}
	case errors.Is(err, accountdomain.ErrInvalidData),
		errors.Is(err, customerdomain.ErrInvalidData),
		errors.Is(err, authdomain.ErrInvalidEmail),
		errors.Is(err, authdomain.ErrInvalidPassword):
		return AppError{Code: ErrCodeInvalidData, Message: "Invalid data", Status: http.StatusBadRequest}
	case errors.Is(err, accountdomain.ErrInvalidAmount):
		return AppError{Code: ErrCodeInvalidAmount, Message: "Invalid amount", Status: http.StatusBadRequest}
	case errors.Is(err, accountdomain.ErrAccountNotFound):
		return AppError{Code: ErrCodeAccountNotFound, Message: "Account not found", Status: http.StatusNotFound}
	case errors.Is(err, accountdomain.ErrCustomerNotFound):
		return AppError{Code: ErrCodeCustomerNotFound, Message: "Customer not found", Status: http.StatusNotFound}
	case errors.Is(err, accountdomain.ErrInsufficientBalance):
		return AppError{Code: ErrCodeInsufficientFunds, Message: "Insufficient balance", Status: http.StatusUnprocessableEntity}
	case errors.Is(err, accountdomain.ErrAccountInactive):
		return AppError{Code: ErrCodeAccountInactive, Message: "Account is not active", Status: http.StatusUnprocessableEntity}
	case errors.Is(err, accountdomain.ErrSameAccountTransfer):
		return AppError{Code: ErrCodeSameAccount, Message: "Source and destination accounts must be different", Status: http.StatusBadRequest}
	case errors.Is(err, accountdomain.ErrForbidden):
		return AppError{Code: ErrCodeForbidden, Message: "Access denied to account", Status: http.StatusForbidden}
	case errors.Is(err, customerdomain.ErrCPFAlreadyExists), errors.Is(err, customerdomain.ErrEmailAlreadyExists), errors.Is(err, authdomain.ErrEmailAlreadyExists):
		return AppError{Code: ErrCodeUserAlreadyExists, Message: "User already exists", Status: http.StatusConflict}
	case errors.Is(err, authdomain.ErrInvalidCredentials):
		return AppError{Code: ErrCodeInvalidCredentials, Message: "Invalid credentials", Status: http.StatusUnauthorized}
	case errors.Is(err, authdomain.ErrUnauthorized):
		return AppError{Code: ErrCodeUnauthorized, Message: "Authentication required", Status: http.StatusUnauthorized}
	case errors.Is(err, authdomain.ErrInvalidToken):
		return AppError{Code: ErrCodeInvalidToken, Message: "Invalid token", Status: http.StatusUnauthorized}
	default:
		return AppError{Code: ErrCodeInternal, Message: "Internal server error", Status: http.StatusInternalServerError}
	}
}
