package delivery

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/application"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
	sharedhttp "github.com/seu-usuario/bank-api/internal/shared/http"
	sharedheaders "github.com/seu-usuario/bank-api/internal/shared/http/headers"
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

type getSessionUseCase interface {
	Execute(ctx context.Context) (*application.GetSessionOutput, error)
}

type refreshAccessTokenUseCase interface {
	Execute(ctx context.Context, input application.RefreshAccessTokenInput) (*application.RefreshAccessTokenOutput, error)
}

type requestContactVerificationUseCase interface {
	Execute(ctx context.Context, input application.RequestContactVerificationInput) (*application.RequestContactVerificationOutput, error)
}

type confirmContactVerificationUseCase interface {
	Execute(ctx context.Context, input application.ConfirmContactVerificationInput) (*application.ConfirmContactVerificationOutput, error)
}

type Handler struct {
	registerUser               registerUserUseCase
	loginUser                  loginUserUseCase
	getCurrentUser             getCurrentUserUseCase
	getSession                 getSessionUseCase
	refreshAccessToken         refreshAccessTokenUseCase
	requestContactVerification requestContactVerificationUseCase
	confirmContactVerification confirmContactVerificationUseCase
}

type registerUserRequest struct {
	Email                  string `json:"email"`
	Phone                  string `json:"phone"`
	Password               string `json:"password"`
	Name                   string `json:"name"`
	BirthDate              string `json:"birth_date"`
	CPF                    string `json:"cpf"`
	EmailVerificationToken string `json:"email_verification_token"`
	PhoneVerificationToken string `json:"phone_verification_token"`
}

type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type requestContactVerificationRequest struct {
	Channel string `json:"channel"`
	Target  string `json:"target"`
}

type confirmContactVerificationRequest struct {
	VerificationID string `json:"verification_id"`
	Token          string `json:"token"`
}

type userData struct {
	ID         uuid.UUID  `json:"id"`
	Email      string     `json:"email"`
	Role       string     `json:"role"`
	CustomerID *uuid.UUID `json:"customer_id,omitempty"`
}

type loginData struct {
	AccessToken           string     `json:"access_token,omitempty"`
	RefreshToken          string     `json:"refresh_token,omitempty"`
	RestrictedAccessToken string     `json:"restricted_access_token,omitempty"`
	RestrictedTokenType   string     `json:"restricted_token_type,omitempty"`
	RestrictedScope       string     `json:"restricted_scope,omitempty"`
	RestrictedExpiresAt   *time.Time `json:"restricted_expires_at,omitempty"`
	UserID                uuid.UUID  `json:"user_id"`
	Email                 string     `json:"email"`
	Role                  string     `json:"role"`
	CustomerID            *uuid.UUID `json:"customer_id,omitempty"`
}

type sessionData struct {
	User      sessionUserData      `json:"user"`
	Customer  sessionCustomerData  `json:"customer"`
	Readiness sessionReadinessData `json:"readiness"`
}

type sessionUserData struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Phone string    `json:"phone"`
	Role  string    `json:"role"`
}

type sessionCustomerData struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	BirthDate string    `json:"birth_date"`
	CreatedAt string    `json:"created_at"`
}

type sessionReadinessData struct {
	OnboardingCompleted       bool   `json:"onboarding_completed"`
	Approved                  bool   `json:"approved"`
	HasOperationalAccount     bool   `json:"has_operational_account"`
	TransactionPasswordStatus string `json:"transaction_password_status"`
}

type refreshAccessTokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// New creates a new instance of the Handler with the provided use cases.
// It requires use cases for registering a user, logging in a user, getting the
// current user, refreshing an access token, requesting contact verification,
// and confirming contact verification.
func New(
	registerUser registerUserUseCase,
	loginUser loginUserUseCase,
	getCurrentUser getCurrentUserUseCase,
	refreshAccessToken refreshAccessTokenUseCase,
	requestContactVerification requestContactVerificationUseCase,
	confirmContactVerification confirmContactVerificationUseCase,
	getSession ...getSessionUseCase,
) *Handler {
	handler := &Handler{
		registerUser:               registerUser,
		loginUser:                  loginUser,
		getCurrentUser:             getCurrentUser,
		refreshAccessToken:         refreshAccessToken,
		requestContactVerification: requestContactVerification,
		confirmContactVerification: confirmContactVerification,
	}
	if len(getSession) > 0 {
		handler.getSession = getSession[0]
	}

	return handler
}

// Register handles HTTP requests for user registration.
// It validates the request, parses the birth date, and executes the register use case.
// Returns appropriate HTTP responses with user data on success or errors on failure.
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

	birthDate, err := time.Parse("2006-01-02", strings.TrimSpace(req.BirthDate))
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	output, err := h.registerUser.Execute(r.Context(), application.RegisterUserInput{
		Email:                  strings.TrimSpace(req.Email),
		Phone:                  strings.TrimSpace(req.Phone),
		Password:               strings.TrimSpace(req.Password),
		Name:                   strings.TrimSpace(req.Name),
		BirthDate:              birthDate,
		CPF:                    strings.TrimSpace(req.CPF),
		EmailVerificationToken: strings.TrimSpace(req.EmailVerificationToken),
		PhoneVerificationToken: strings.TrimSpace(req.PhoneVerificationToken),
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

// RequestContactVerification handles the HTTP request for creating a contact
// verification attempt for the provided channel and target.
func (h *Handler) RequestContactVerification(w http.ResponseWriter, r *http.Request) {
	if h.requestContactVerification == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	var req requestContactVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	output, err := h.requestContactVerification.Execute(r.Context(), application.RequestContactVerificationInput{
		Channel: req.Channel,
		Target:  req.Target,
	})
	if err != nil {
		log.Printf("event=request_contact_verification error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusCreated, output)
}

// ConfirmContactVerification handles the HTTP request for confirming a contact
// verification attempt using its identifier and token.
func (h *Handler) ConfirmContactVerification(w http.ResponseWriter, r *http.Request) {
	if h.confirmContactVerification == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	var req confirmContactVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	verificationID, err := uuid.Parse(strings.TrimSpace(req.VerificationID))
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	output, err := h.confirmContactVerification.Execute(r.Context(), application.ConfirmContactVerificationInput{
		VerificationID: verificationID,
		Token:          req.Token,
	})
	if err != nil {
		log.Printf("event=confirm_contact_verification error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusOK, output)
}

// isValidRegisterRequest validates the registration request by checking if all
// required fields (email, phone, password, name, birth date, CPF, and contact
// verification tokens) are provided and not empty after whitespace trimming.
func isValidRegisterRequest(req registerUserRequest) bool {
	email := strings.TrimSpace(req.Email)
	phone := strings.TrimSpace(req.Phone)
	password := strings.TrimSpace(req.Password)
	name := strings.TrimSpace(req.Name)
	birthDate := strings.TrimSpace(req.BirthDate)
	cpf := strings.TrimSpace(req.CPF)
	emailVerificationToken := strings.TrimSpace(req.EmailVerificationToken)
	phoneVerificationToken := strings.TrimSpace(req.PhoneVerificationToken)

	if email == "" ||
		phone == "" ||
		password == "" ||
		name == "" ||
		birthDate == "" ||
		cpf == "" ||
		emailVerificationToken == "" ||
		phoneVerificationToken == "" {
		return false
	}

	return true
}

func formatDate(value time.Time) string {
	if value.IsZero() {
		return ""
	}

	return value.Format("2006-01-02")
}

// Login handles the HTTP request for user login. It validates the input, calls the
// loginUser use case to perform the login operation, and returns the access token,
// refresh token, and user information in the response. If any step in the process
// fails, it returns an appropriate error response.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if h.loginUser == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	installationID, err := parseCanonicalInstallationID(r.Header.Get(sharedheaders.InstallationID))
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(err))
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
		Email:          req.Email,
		Password:       req.Password,
		InstallationID: installationID,
	})
	if err != nil {
		log.Printf("event=login_user error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	if output == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusOK, loginData{
		AccessToken:           output.AccessToken,
		RefreshToken:          output.RefreshToken,
		RestrictedAccessToken: output.RestrictedAccessToken,
		RestrictedTokenType:   output.RestrictedTokenType,
		RestrictedScope:       output.RestrictedScope,
		RestrictedExpiresAt:   output.RestrictedExpiresAt,
		UserID:                output.UserID,
		Email:                 output.Email,
		Role:                  output.Role,
		CustomerID:            output.CustomerID,
	})
}

// parseCanonicalInstallationID validates and parses the installation ID
// from the request header.
// It ensures that the installation ID is a valid UUID version 4 and is
// in its canonical string representation.
// If the installation ID is invalid, it returns an error indicating an
// invalid installation ID.
func parseCanonicalInstallationID(raw string) (uuid.UUID, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return uuid.Nil, sharederrors.ErrInvalidInstallationID
	}

	installationID, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, sharederrors.ErrInvalidInstallationID
	}

	if installationID.Version() != 4 {
		return uuid.Nil, sharederrors.ErrInvalidInstallationID
	}

	if installationID.String() != value {
		return uuid.Nil, sharederrors.ErrInvalidInstallationID
	}

	return installationID, nil
}

// Me handles the HTTP request for retrieving the current authenticated user's
// information. It calls the getCurrentUser use case to fetch the user's details
// and returns them in the response. If any step in the process fails, it returns
// an appropriate error response.
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

func (h *Handler) Session(w http.ResponseWriter, r *http.Request) {
	if h.getSession == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	output, err := h.getSession.Execute(r.Context())
	if err != nil {
		log.Printf("event=get_session error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}
	if output == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusOK, sessionData{
		User: sessionUserData{
			ID:    output.User.ID,
			Email: output.User.Email,
			Phone: output.User.Phone,
			Role:  output.User.Role,
		},
		Customer: sessionCustomerData{
			ID:        output.Customer.ID,
			Name:      output.Customer.Name,
			CPF:       output.Customer.CPF,
			BirthDate: formatDate(output.Customer.BirthDate),
			CreatedAt: output.Customer.CreatedAt.Format(time.RFC3339),
		},
		Readiness: sessionReadinessData{
			OnboardingCompleted:       output.Readiness.OnboardingCompleted,
			Approved:                  output.Readiness.Approved,
			HasOperationalAccount:     output.Readiness.HasOperationalAccount,
			TransactionPasswordStatus: output.Readiness.TransactionPasswordStatus,
		},
	})
}

// Refresh handles the HTTP request for refreshing the access token. It validates the input, calls the
// refreshAccessToken use case to perform the token refresh operation, and returns the new access token
// and refresh token in the response. If any step in the process fails, it returns an appropriate error response.
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	if h.refreshAccessToken == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	installationID, err := parseCanonicalInstallationID(r.Header.Get(sharedheaders.InstallationID))
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(err))
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
		RefreshToken:   req.RefreshToken,
		InstallationID: installationID,
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
