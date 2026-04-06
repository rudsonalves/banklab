package sharederrors

import "errors"

type AppError struct {
	Code    string
	Message string
	Status  int
}

func RegisterDomainError(err error, code, message string, status int) {
	Register(func(e error) bool {
		return errors.Is(e, err)
	}, AppError{
		Code:    code,
		Message: message,
		Status:  status,
	})
}
