package sharedhttp

import (
	"encoding/json"
	"log"
	"net/http"

	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
)

type Response struct {
	Data  any        `json:"data"`
	Error *ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// WriteJSON writes a successful JSON response with the given status code and
// data. It sets the Content-Type header to application/json and encodes the
// response payload as JSON. If there is an error during encoding, it logs the
// error but does not return it to the caller, as this function is intended for
// writing responses rather than handling errors.
func WriteJSON(w http.ResponseWriter, status int, data any) {
	writeResponse(w, status, Response{
		Data:  data,
		Error: nil,
	})
}

// WriteError writes an error response in JSON format using the provided AppError.
// It sets the Content-Type header to application/json and encodes the error
// details in the response body. The HTTP status code is determined by the Status field of the AppError. If there is an error during encoding, it logs the
// error but does not return it to the caller, as this function is intended for
// writing responses rather than handling errors.
func WriteError(w http.ResponseWriter, appErr sharederrors.AppError) {
	writeResponse(w, appErr.Status, Response{
		Data: nil,
		Error: &ErrorBody{
			Code:    appErr.Code,
			Message: appErr.Message,
		},
	})
}

// writeResponse is a helper function that writes a JSON response with the specified
// HTTP status code and payload. It sets the Content-Type header to application/json
// and encodes the payload as JSON. If there is an error during encoding, it logs
// the error but does not return it to the caller, as this function is intended for
// writing responses rather than handling errors.
func writeResponse(w http.ResponseWriter, status int, payload Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Println("write response error:", err)
	}
}
