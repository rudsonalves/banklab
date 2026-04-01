package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/seu-usuario/bank-api/internal/customer/domain"
	"github.com/seu-usuario/bank-api/internal/customer/usecase"
)

type Handler struct {
	uc *usecase.CreateCustomer
}

func New(uc *usecase.CreateCustomer) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input usecase.Input

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.uc.Execute(r.Context(), input)
	if err != nil {
		log.Println("create customer error:", err)

		switch err {
		case domain.ErrInvalidData:
			http.Error(w, "invalid data", http.StatusBadRequest)
			return
		case domain.ErrCPFAlreadyExists:
			http.Error(w, "cpf already exists", http.StatusConflict)
			return
		case domain.ErrEmailAlreadyExists:
			http.Error(w, "email already exists", http.StatusConflict)
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
