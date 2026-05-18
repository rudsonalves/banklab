package sharederrors

import "errors"

type AppError struct {
	Code    string
	Message string
	Status  int
	Details any
}

// RegisterDomainError registers a domain error with the provided code, message,
// and HTTP status. It uses the Register function to associate the domain error
// with an AppError that contains the specified code, message, and status. This
// allows for consistent error handling and mapping of domain errors to appropriate
// HTTP responses throughout the application.
func RegisterDomainError(err error, code, message string, status int) {
	Register(func(e error) bool {
		return errors.Is(e, err)
	}, AppError{
		Code:    code,
		Message: message,
		Status:  status,
	})
}

// RegisterDomainErrorWithDetails registers a domain error whose AppError.Details
// is populated dynamically by calling detailsFn with the original error.
// Use this for errors that carry structured context (e.g. which contacts are missing).
func RegisterDomainErrorWithDetails(sentinel error, code, message string, status int, detailsFn func(error) any) {
	registry = append(registry, entry{
		match: func(e error) bool {
			return errors.Is(e, sentinel)
		},
		appErr: AppError{
			Code:    code,
			Message: message,
			Status:  status,
		},
		detailsFn: detailsFn,
	})
}
