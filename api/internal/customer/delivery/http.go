package delivery

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	"github.com/seu-usuario/bank-api/internal/customer/application"
	customerdomain "github.com/seu-usuario/bank-api/internal/customer/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
	sharedhttp "github.com/seu-usuario/bank-api/internal/shared/http"
)

type createCustomerUseCase interface {
	Execute(
		ctx context.Context,
		input application.Input,
	) (*customerdomain.Customer, error)
}

type getCustomerMeUseCase interface {
	Execute(
		ctx context.Context,
		input application.GetCustomerMeInput,
	) (*customerdomain.CustomerProfile, error)
}

type checkCPFUseCase interface {
	Execute(
		ctx context.Context,
		input application.CheckCPFInput,
	) (*application.CheckCPFOutput, error)
}

type Handler struct {
	createUC   createCustomerUseCase
	getMeUC    getCustomerMeUseCase
	checkCPFUC checkCPFUseCase
}

type createCustomerRequest struct {
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
}

type createCustomerData struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CPF       string `json:"cpf"`
	Email     string `json:"email"`
	BirthDate string `json:"birth_date,omitempty"`
	CreatedAt string `json:"created_at"`
}

type checkCPFRequest struct {
	CPF string `json:"cpf"`
}

// New creates a new instance of the customer Handler with the provided use cases.
func New(createUC createCustomerUseCase, getMeUC getCustomerMeUseCase, checkCPFUC ...checkCPFUseCase) *Handler {
	handler := &Handler{createUC: createUC, getMeUC: getMeUC}
	if len(checkCPFUC) > 0 {
		handler.checkCPFUC = checkCPFUC[0]
	}

	return handler
}

// Create handles the HTTP POST request for creating a new customer. It decodes the
// JSON request body containing name and birth_date, parses and validates the birth date,
// then calls the create customer use case. On success, it responds with HTTP 201 Created
// and the customer data (ID, name, CPF, email, birth_date, created_at). On validation
// or business logic errors, it logs the error and responds with an appropriate error
// message and status code.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if h.createUC == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	var req createCustomerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.WriteError(
			w,
			sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		sharedhttp.WriteError(
			w,
			sharederrors.MapError(customerdomain.ErrInvalidData),
		)
		return
	}

	customer, err := h.createUC.Execute(r.Context(), application.Input{
		Name:      req.Name,
		BirthDate: birthDate,
	})
	if err != nil {
		log.Println("create customer error:", err)

		if err == customerdomain.ErrNameRequired || err == customerdomain.ErrBirthDateRequired {
			sharedhttp.WriteError(
				w,
				sharederrors.MapError(customerdomain.ErrInvalidData),
			)
			return
		}

		sharedhttp.WriteError(
			w,
			sharederrors.MapError(err),
		)
		return
	}

	sharedhttp.WriteJSON(w, http.StatusCreated, createCustomerData{
		ID:        customer.ID.String(),
		Name:      customer.Name,
		BirthDate: formatBirthDate(customer.BirthDate),
		CreatedAt: customer.CreatedAt.Format(time.RFC3339),
	})
}

// Me retrieves the authenticated customer's profile information. It validates
// the authenticated user context and calls the getMe use case. Returns the customer
// profile with status 200, or an error response with appropriate status code.
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	if h.getMeUC == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	user, ok := sharedauthctx.GetAuthenticatedUser(r.Context())
	if !ok || user == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(authdomain.ErrUnauthorized))
		return
	}

	if user.CustomerID == nil {
		sharedhttp.WriteError(
			w, sharederrors.MapError(authdomain.ErrInvalidUserState))
		return
	}

	profile, err := h.getMeUC.Execute(
		r.Context(),
		application.GetCustomerMeInput{CustomerID: *user.CustomerID},
	)
	if err != nil {
		log.Println("get customer me error:", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusOK, createCustomerData{
		ID:        profile.Customer.ID.String(),
		Name:      profile.Customer.Name,
		CPF:       profile.CPF,
		Email:     profile.Email,
		BirthDate: formatBirthDate(profile.Customer.BirthDate),
		CreatedAt: profile.Customer.CreatedAt.Format(time.RFC3339),
	})
}

// CheckCPF verifies whether a CPF is already registered in customer documents.
func (h *Handler) CheckCPF(w http.ResponseWriter, r *http.Request) {
	if h.checkCPFUC == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	var req checkCPFRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	output, err := h.checkCPFUC.Execute(r.Context(), application.CheckCPFInput{CPF: req.CPF})
	if err != nil {
		log.Println("check cpf error:", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusOK, output)
}

// formatBirthDate formats the birth date as a string in "YYYY-MM-DD" format.
// If the birth date is zero (not set), it returns an empty string. This
// function is used to ensure that the birth date is consistently formatted
// in the API responses, and to handle cases where the birth date may not be provided.
func formatBirthDate(birthDate time.Time) string {
	if birthDate.IsZero() {
		return ""
	}

	return birthDate.Format("2006-01-02")
}
