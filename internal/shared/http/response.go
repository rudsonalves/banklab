package sharedhttp

import (
	"encoding/json"
	"log"
	"net/http"

	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
)

type response struct {
	Data  interface{}            `json:"data"`
	Error *sharederrors.AppError `json:"error"`
}

func WriteSuccess(w http.ResponseWriter, status int, data interface{}) {
	writeJSON(w, status, response{
		Data:  data,
		Error: nil,
	})
}

func WriteError(w http.ResponseWriter, status int, err *sharederrors.AppError) {
	if err == nil {
		err = sharederrors.ErrInternal
	}

	writeJSON(w, status, response{
		Data:  nil,
		Error: err,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Println("write response error:", err)
	}
}
