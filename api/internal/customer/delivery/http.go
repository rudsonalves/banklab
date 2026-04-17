package delivery

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	authdelivery "github.com/seu-usuario/bank-api/internal/auth/delivery"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	"github.com/seu-usuario/bank-api/internal/customer/application"
	customerdomain "github.com/seu-usuario/bank-api/internal/customer/domain"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
	sharedhttp "github.com/seu-usuario/bank-api/internal/shared/http"
)

type createCustomerUseCase interface {
	Execute(ctx context.Context, input application.Input) (*customerdomain.Customer, error)
}

type getCustomerMeUseCase interface {
	Execute(ctx context.Context, input application.GetCustomerMeInput) (*customerdomain.Customer, string, error)
}

type Handler struct {
	createUC createCustomerUseCase
	getMeUC  getCustomerMeUseCase
}

type createCustomerRequest struct {
	Name  string `json:"name"`
	CPF   string `json:"cpf"`
	Email string `json:"email"`
}

type createCustomerData struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CPF       string `json:"cpf"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func New(createUC createCustomerUseCase, getMeUC getCustomerMeUseCase) *Handler {
	return &Handler{createUC: createUC, getMeUC: getMeUC}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if h.createUC == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	var req createCustomerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	customer, err := h.createUC.Execute(r.Context(), application.Input{
		Name: req.Name,
		CPF:  req.CPF,
	})
	if err != nil {
		log.Println("create customer error:", err)

		if err == customerdomain.ErrNameRequired || err == customerdomain.ErrCPFRequired {
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
		Email:     req.Email,
		CreatedAt: customer.CreatedAt.Format(time.RFC3339),
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	if h.getMeUC == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	user, ok := authdelivery.GetAuthenticatedUser(r.Context())
	if !ok || user == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(authdomain.ErrUnauthorized))
		return
	}

	if user.CustomerID == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(authdomain.ErrInvalidUserState))
		return
	}

	customer, email, err := h.getMeUC.Execute(r.Context(), application.GetCustomerMeInput{CustomerID: *user.CustomerID})
	if err != nil {
		log.Println("get customer me error:", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusOK, createCustomerData{
		ID:        customer.ID.String(),
		Name:      customer.Name,
		CPF:       customer.CPF,
		Email:     email,
		CreatedAt: customer.CreatedAt.Format(time.RFC3339),
	})
}
