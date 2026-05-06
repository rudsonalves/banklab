package delivery

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	adminapplication "github.com/seu-usuario/bank-api/internal/admin/application"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
	sharedhttp "github.com/seu-usuario/bank-api/internal/shared/http"
)

type approveUserUseCase interface {
	Execute(ctx context.Context, input adminapplication.ApproveUserInput) (*adminapplication.ApproveUserOutput, error)
}

type Handler struct {
	approveUser approveUserUseCase
}

type approveUserData struct {
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	AccountID string `json:"account_id"`
}

// New creates a new instance of the account Handler with the provided use cases.
func New(approveUser approveUserUseCase) *Handler {
	return &Handler{approveUser: approveUser}
}

// ApproveUser handles the HTTP request for approving a pending user. It checks for
// authentication and authorization, validates the input, and calls the approveUser
// use case to perform the approval. The response is returned in JSON format, and
// appropriate error handling is performed for various failure scenarios.
func (h *Handler) ApproveUser(w http.ResponseWriter, r *http.Request) {
	if h.approveUser == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	user, err := sharedauthctx.RequireAuthenticatedUser(r.Context())
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(authdomain.ErrUnauthorized))
		return
	}

	if user.Role != authdomain.RoleAdmin {
		sharedhttp.WriteError(w, sharederrors.MapError(authdomain.ErrForbidden))
		return
	}

	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(authdomain.ErrInvalidData))
		return
	}

	output, err := h.approveUser.Execute(r.Context(), adminapplication.ApproveUserInput{UserID: userID})
	if err != nil {
		log.Printf("event=approve_user error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	if output == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusOK, approveUserData{
		UserID:    output.UserID.String(),
		Status:    output.Status,
		AccountID: output.AccountID.String(),
	})
}
