package sharederrors

import (
	"errors"
	"net/http"
)

var ErrInvalidRequest = errors.New("invalid request")

type entry struct {
	match  func(error) bool
	appErr AppError
}

var registry []entry

func init() {
	Register(func(err error) bool {
		return errors.Is(err, ErrInvalidRequest)
	}, AppError{
		Code:    ErrCodeInvalidRequest,
		Message: "Invalid request body",
		Status:  http.StatusBadRequest,
	})
}

func Register(match func(error) bool, appErr AppError) {
	registry = append(registry, entry{
		match:  match,
		appErr: appErr,
	})
}

func MapError(err error) AppError {
	if err == nil {
		return internalError()
	}

	for _, e := range registry {
		if e.match(err) {
			return e.appErr
		}
	}

	return internalError()
}

func internalError() AppError {
	return AppError{
		Code:    ErrCodeInternal,
		Message: "Internal server error",
		Status:  http.StatusInternalServerError,
	}
}
