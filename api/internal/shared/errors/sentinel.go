package sharederrors

import "net/http"

var ErrInvalidAppToken = AppError{
	Code:    "INVALID_APP_TOKEN",
	Message: "invalid application token",
	Status:  http.StatusUnauthorized,
}
