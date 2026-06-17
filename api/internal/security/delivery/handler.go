package delivery

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	"github.com/seu-usuario/bank-api/internal/security/application"
	"github.com/seu-usuario/bank-api/internal/security/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
	sharedhttp "github.com/seu-usuario/bank-api/internal/shared/http"
)

type createTransactionPasswordUseCase interface {
	Execute(ctx context.Context, input application.CreateTransactionPasswordInput) (*application.CreateTransactionPasswordOutput, error)
}

type authorizeStepUpUseCase interface {
	Execute(ctx context.Context, input application.AuthorizeStepUpInput) (*application.AuthorizeStepUpOutput, error)
}

// Handler exposes the HTTP endpoints for the security module.
type Handler struct {
	// createTransactionPassword handles the creation of transaction passwords.
	createTransactionPassword createTransactionPasswordUseCase
	authorizeStepUp           authorizeStepUpUseCase
}

type createTransactionPasswordRequest struct {
	TransactionPassword             string `json:"transaction_password"`
	TransactionPasswordConfirmation string `json:"transaction_password_confirmation"`
}

type authorizeStepUpRequest struct {
	Method              string `json:"method"`
	Path                string `json:"path"`
	TransactionPassword string `json:"transaction_password"`
}

type transactionPasswordData struct {
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type stepUpAuthorizationData struct {
	StepUpToken string `json:"step_up_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// New creates a Handler with the use cases required by security endpoints.
//
// The step-up authorization use case is optional to support gradual module
// composition in scenarios where only transaction password creation is enabled.
func New(
	createTransactionPassword createTransactionPasswordUseCase,
	authorizeStepUp ...authorizeStepUpUseCase,
) *Handler {
	handler := &Handler{createTransactionPassword: createTransactionPassword}
	if len(authorizeStepUp) > 0 {
		handler.authorizeStepUp = authorizeStepUp[0]
	}

	return handler
}

// CreateTransactionPassword creates or updates the authenticated user's
// transaction password.
//
// The handler validates authentication, decodes a strict JSON payload, and
// delegates business rules to the corresponding use case. On success, it
// returns HTTP 201 with transaction password metadata.
func (h *Handler) CreateTransactionPassword(w http.ResponseWriter, r *http.Request) {
	if h.createTransactionPassword == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	user, ok := sharedauthctx.GetAuthenticatedUser(r.Context())
	if !ok || user == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(authdomain.ErrUnauthorized))
		return
	}

	var req createTransactionPasswordRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	output, err := h.createTransactionPassword.Execute(r.Context(), application.CreateTransactionPasswordInput{
		User:                            user,
		TransactionPassword:             strings.TrimSpace(req.TransactionPassword),
		TransactionPasswordConfirmation: strings.TrimSpace(req.TransactionPasswordConfirmation),
	})
	if err != nil {
		log.Printf("event=create_transaction_password error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	if output == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusCreated, transactionPasswordData{
		UserID:    output.UserID,
		Status:    output.Status,
		CreatedAt: output.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// AuthorizeStepUp authorizes a sensitive operation through step-up
// authentication.
//
// The handler requires an authenticated user, validates the request body,
// and delegates transaction password validation to the use case. On success,
// it returns HTTP 200 with a step-up token and expiration time in seconds.
func (h *Handler) AuthorizeStepUp(w http.ResponseWriter, r *http.Request) {
	if h.authorizeStepUp == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	var req authorizeStepUpRequest
	if err := sharedhttp.DecodeJSON(r.Body, &req); err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	user, ok := sharedauthctx.GetAuthenticatedUser(r.Context())
	if !ok || user == nil {
		restricted, restrictedOK := sharedauthctx.GetRestrictedSession(r.Context())
		if !restrictedOK || restricted == nil {
			sharedhttp.WriteError(w, sharederrors.MapError(authdomain.ErrUnauthorized))
			return
		}
		if strings.TrimSpace(req.Method) != domain.StepUpPublicMethodInstallationRegisterCreate ||
			strings.TrimSpace(req.Path) != domain.StepUpPublicPathInstallationRegisterCreate {
			sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrStepUpEndpointNotAllowed))
			return
		}
		user = &authdomain.AuthenticatedUser{UserID: restricted.UserID}
	}

	output, err := h.authorizeStepUp.Execute(r.Context(), application.AuthorizeStepUpInput{
		User:                user,
		Method:              strings.TrimSpace(req.Method),
		Path:                strings.TrimSpace(req.Path),
		TransactionPassword: strings.TrimSpace(req.TransactionPassword),
	})
	if err != nil {
		log.Printf("event=authorize_step_up error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	if output == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusOK, stepUpAuthorizationData{
		StepUpToken: output.StepUpToken,
		ExpiresIn:   output.ExpiresIn,
	})
}
