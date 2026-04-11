package delivery

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/application"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
	sharedhttp "github.com/seu-usuario/bank-api/internal/shared/http"
)

type registerUserUseCase interface {
	Execute(ctx context.Context, input application.RegisterUserInput) (*application.RegisterUserOutput, error)
}

type loginUserUseCase interface {
	Execute(ctx context.Context, input application.LoginUserInput) (*application.LoginUserOutput, error)
}

type getCurrentUserUseCase interface {
	Execute(ctx context.Context) (*application.GetCurrentUserOutput, error)
}

type refreshAccessTokenUseCase interface {
	Execute(ctx context.Context, input application.RefreshAccessTokenInput) (*application.RefreshAccessTokenOutput, error)
}

type Handler struct {
	registerUser       registerUserUseCase
	loginUser          loginUserUseCase
	getCurrentUser     getCurrentUserUseCase
	refreshAccessToken refreshAccessTokenUseCase
}

type registerUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	CPF      string `json:"cpf"`
}

type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type userData struct {
	ID         uuid.UUID  `json:"id"`
	Email      string     `json:"email"`
	Role       string     `json:"role"`
	CustomerID *uuid.UUID `json:"customer_id,omitempty"`
}

type loginData struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	UserID       uuid.UUID  `json:"user_id"`
	Email        string     `json:"email"`
	Role         string     `json:"role"`
	CustomerID   *uuid.UUID `json:"customer_id,omitempty"`
}

type refreshAccessTokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func New(
	registerUser registerUserUseCase,
	loginUser loginUserUseCase,
	getCurrentUser getCurrentUserUseCase,
) *Handler {
	return &Handler{
		registerUser:       registerUser,
		loginUser:          loginUser,
		getCurrentUser:     getCurrentUser,
		refreshAccessToken: nil,
	}
}

func (h *Handler) SetRefreshAccessTokenUseCase(useCase refreshAccessTokenUseCase) {
	h.refreshAccessToken = useCase
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if h.registerUser == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	var req registerUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	if !isValidRegisterRequest(req) {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	output, err := h.registerUser.Execute(r.Context(), application.RegisterUserInput{
		Email:    strings.TrimSpace(req.Email),
		Password: strings.TrimSpace(req.Password),
		Name:     strings.TrimSpace(req.Name),
		CPF:      strings.TrimSpace(req.CPF),
	})
	if err != nil {
		log.Printf("event=register_user error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	if output == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusCreated, userData{
		ID:         output.ID,
		Email:      output.Email,
		Role:       output.Role,
		CustomerID: output.CustomerID,
	})
}

func isValidRegisterRequest(req registerUserRequest) bool {
	email := strings.TrimSpace(req.Email)
	password := strings.TrimSpace(req.Password)
	name := strings.TrimSpace(req.Name)
	cpf := strings.TrimSpace(req.CPF)

	if email == "" || password == "" || name == "" || cpf == "" {
		return false
	}

	return true
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if h.loginUser == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	var req loginUserRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	output, err := h.loginUser.Execute(r.Context(), application.LoginUserInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.Printf("event=login_user error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusOK, loginData{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		UserID:       output.UserID,
		Email:        output.Email,
		Role:         output.Role,
		CustomerID:   output.CustomerID,
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	if h.getCurrentUser == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	output, err := h.getCurrentUser.Execute(r.Context())
	if err != nil {
		log.Printf("event=get_current_user error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusOK, userData{
		ID:         output.ID,
		Email:      output.Email,
		Role:       output.Role,
		CustomerID: output.CustomerID,
	})
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	if h.refreshAccessToken == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	var req refreshAccessTokenRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	output, err := h.refreshAccessToken.Execute(r.Context(), application.RefreshAccessTokenInput{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		log.Printf("event=refresh_access_token error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	if output == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusOK, refreshAccessTokenData{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	})
}
