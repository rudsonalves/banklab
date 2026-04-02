package delivery

import (
	"encoding/json"
	"log"
	"net/http"
)

type apiError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type response struct {
	Data  interface{} `json:"data"`
	Error *apiError   `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, payload response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Println("write response error:", err)
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, response{
		Data: nil,
		Error: &apiError{
			Code:    code,
			Message: message,
		},
	})
}
