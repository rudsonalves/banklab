package delivery

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/seu-usuario/bank-api/internal/customer/application"
	"github.com/seu-usuario/bank-api/internal/customer/domain"
)

type Handler struct {
	uc *application.CreateCustomer
}

type apiError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type response struct {
	Data  interface{} `json:"data"`
	Error *apiError   `json:"error"`
}

type createCustomerData struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CPF       string `json:"cpf"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func New(uc *application.CreateCustomer) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input application.Input

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	customer, err := h.uc.Execute(r.Context(), input)
	if err != nil {
		log.Println("create customer error:", err)

		if field, ok := validationField(err); ok {
			writeErrorWithDetails(w, http.StatusBadRequest, "INVALID_DATA", "invalid data", map[string]string{
				"field": field,
			})
			return
		}

		switch {
		case errors.Is(err, domain.ErrInvalidData):
			writeError(w, http.StatusBadRequest, "INVALID_DATA", "invalid data")
			return

		case errors.Is(err, domain.ErrCPFAlreadyExists):
			writeError(w, http.StatusConflict, "CUSTOMER_ALREADY_EXISTS", "cpf already exists")
			return

		case errors.Is(err, domain.ErrEmailAlreadyExists):
			writeError(w, http.StatusConflict, "CUSTOMER_ALREADY_EXISTS", "email already exists")
			return
		}

		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, response{
		Data: createCustomerData{
			ID:        customer.ID.String(),
			Name:      customer.Name,
			CPF:       customer.CPF,
			Email:     customer.Email,
			CreatedAt: customer.CreatedAt.Format(time.RFC3339),
		},
		Error: nil,
	})
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

func writeErrorWithDetails(w http.ResponseWriter, status int, code, message string, details interface{}) {
	writeJSON(w, status, response{
		Data: nil,
		Error: &apiError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

func validationField(err error) (string, bool) {
	switch {
	case errors.Is(err, domain.ErrNameRequired):
		return "name", true
	case errors.Is(err, domain.ErrCPFRequired):
		return "cpf", true
	case errors.Is(err, domain.ErrEmailRequired):
		return "email", true
	}

	return "", false
}

func writeJSON(w http.ResponseWriter, status int, payload response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Println("write response error:", err)
	}
}
