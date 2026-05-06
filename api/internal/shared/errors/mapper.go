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

// init registers the default error mapping for invalid request errors. It maps any
// error that matches ErrInvalidRequest to an AppError with a specific code, message,
// and HTTP status. This ensures that when an invalid request error occurs, it will
// be consistently handled and returned as a Bad Request response to the client.
func init() {
	Register(func(err error) bool {
		return errors.Is(err, ErrInvalidRequest)
	}, AppError{
		Code:    ErrCodeInvalidRequest,
		Message: "Invalid request body",
		Status:  http.StatusBadRequest,
	})
}

// Register allows you to register a custom error mapping by providing a matching
// function and an AppError. The matching function should return true for errors
// that you want to map to the specified AppError. This function appends the new
// error mapping to the registry, enabling consistent error handling across the
// application. When an error is mapped using MapError, it will be checked against
// all registered entries in the registry to find a matching AppError.
func Register(match func(error) bool, appErr AppError) {
	registry = append(registry, entry{
		match:  match,
		appErr: appErr,
	})
}

// MapError takes an error as input and checks it against the registered error
// mappings. If a matching error is found, it returns the corresponding AppError.
// If no matching error is found, it returns a generic internal server error AppError. This function is used to convert domain errors into standardized application errors
// that can be consistently handled and returned in HTTP responses throughout the
// application.
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

// internalError returns a generic AppError representing an internal server error.
// This is used as a fallback when no specific error mapping is found. It provides
// a standard error response for unexpected errors that occur within the
// application, ensuring that clients receive a consistent error message and
// status code for internal errors.
func internalError() AppError {
	return AppError{
		Code:    ErrCodeInternal,
		Message: "Internal server error",
		Status:  http.StatusInternalServerError,
	}
}
