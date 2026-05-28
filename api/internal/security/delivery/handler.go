package delivery

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	"github.com/seu-usuario/bank-api/internal/security/application"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
	sharedhttp "github.com/seu-usuario/bank-api/internal/shared/http"
)

type createTransactionPasswordUseCase interface {
	Execute(ctx context.Context, input application.CreateTransactionPasswordInput) (*application.CreateTransactionPasswordOutput, error)
}

type Handler struct {
	createTransactionPassword createTransactionPasswordUseCase
}

type createTransactionPasswordRequest struct {
	TransactionPassword             string `json:"transaction_password"`
	TransactionPasswordConfirmation string `json:"transaction_password_confirmation"`
}

type transactionPasswordData struct {
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func New(createTransactionPassword createTransactionPasswordUseCase) *Handler {
	return &Handler{createTransactionPassword: createTransactionPassword}
}

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
