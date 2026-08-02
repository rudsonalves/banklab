package sharederrors

import (
	"errors"
	"net/http"
)

var ErrInvalidInstallationID = errors.New("invalid installation id")

func init() {
	Register(func(err error) bool {
		return errors.Is(err, ErrInvalidInstallationID)
	}, AppError{
		Code:    ErrCodeInvalidInstallationID,
		Message: "X-Installation-Id must be a canonical UUID v4.",
		Status:  http.StatusBadRequest,
	})
}
