package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/seu-usuario/bank-api/internal/customer/application"
	customerdomain "github.com/seu-usuario/bank-api/internal/customer/domain"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
	sharedhttp "github.com/seu-usuario/bank-api/internal/shared/http"
)

type Handler struct {
	uc *application.CreateCustomer
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
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	customer, err := h.uc.Execute(r.Context(), input)
	if err != nil {
		log.Println("create customer error:", err)

		if err == customerdomain.ErrNameRequired || err == customerdomain.ErrCPFRequired || err == customerdomain.ErrEmailRequired {
			sharedhttp.WriteError(w, sharederrors.MapError(customerdomain.ErrInvalidData))
			return
		}

		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusCreated, createCustomerData{
		ID:        customer.ID.String(),
		Name:      customer.Name,
		CPF:       customer.CPF,
		Email:     customer.Email,
		CreatedAt: customer.CreatedAt.Format(time.RFC3339),
	})
}
