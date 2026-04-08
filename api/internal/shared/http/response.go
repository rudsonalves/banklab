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

func WriteJSON(w http.ResponseWriter, status int, data any) {
	writeResponse(w, status, Response{
		Data:  data,
		Error: nil,
	})
}

func WriteError(w http.ResponseWriter, appErr sharederrors.AppError) {
	writeResponse(w, appErr.Status, Response{
		Data: nil,
		Error: &ErrorBody{
			Code:    appErr.Code,
			Message: appErr.Message,
		},
	})
}

func writeResponse(w http.ResponseWriter, status int, payload Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Println("write response error:", err)
	}
}
