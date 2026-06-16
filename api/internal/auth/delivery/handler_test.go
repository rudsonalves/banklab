package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/application"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
	sharedheaders "github.com/seu-usuario/bank-api/internal/shared/http/headers"
)

type registerUserUseCaseMock struct {
	output *application.RegisterUserOutput
	err    error
	input  application.RegisterUserInput
	called bool
}

func (m *registerUserUseCaseMock) Execute(ctx context.Context, input application.RegisterUserInput) (*application.RegisterUserOutput, error) {
	m.called = true
	m.input = input
	return m.output, m.err
}

type loginUserUseCaseMock struct {
	output *application.LoginUserOutput
	err    error
	input  application.LoginUserInput
	called bool
}

func (m *loginUserUseCaseMock) Execute(ctx context.Context, input application.LoginUserInput) (*application.LoginUserOutput, error) {
	m.called = true
	m.input = input
	return m.output, m.err
}

type getCurrentUserUseCaseMock struct {
	output *application.GetCurrentUserOutput
	err    error
	called bool
}

func (m *getCurrentUserUseCaseMock) Execute(ctx context.Context) (*application.GetCurrentUserOutput, error) {
	m.called = true
	return m.output, m.err
}

type getSessionUseCaseMock struct {
	output *application.GetSessionOutput
	err    error
	called bool
}

func (m *getSessionUseCaseMock) Execute(ctx context.Context) (*application.GetSessionOutput, error) {
	m.called = true
	return m.output, m.err
}

type refreshAccessTokenUseCaseMock struct {
	output *application.RefreshAccessTokenOutput
	err    error
	input  application.RefreshAccessTokenInput
	called bool
}

type requestContactVerificationUseCaseMock struct {
	output *application.RequestContactVerificationOutput
	err    error
	input  application.RequestContactVerificationInput
	called bool
}

func (m *requestContactVerificationUseCaseMock) Execute(
	ctx context.Context,
	input application.RequestContactVerificationInput,
) (*application.RequestContactVerificationOutput, error) {
	m.called = true
	m.input = input
	return m.output, m.err
}

type confirmContactVerificationUseCaseMock struct {
	output *application.ConfirmContactVerificationOutput
	err    error
	input  application.ConfirmContactVerificationInput
	called bool
}

func (m *confirmContactVerificationUseCaseMock) Execute(
	ctx context.Context,
	input application.ConfirmContactVerificationInput,
) (*application.ConfirmContactVerificationOutput, error) {
	m.called = true
	m.input = input
	return m.output, m.err
}

func (m *refreshAccessTokenUseCaseMock) Execute(ctx context.Context, input application.RefreshAccessTokenInput) (*application.RefreshAccessTokenOutput, error) {
	m.called = true
	m.input = input
	return m.output, m.err
}

func TestHandler_Register_Success(t *testing.T) {
	userID := uuid.New()
	customerID := uuid.New()
	registerUC := &registerUserUseCaseMock{
		output: &application.RegisterUserOutput{
			ID:         userID,
			Email:      "user@example.com",
			Role:       "customer",
			CustomerID: &customerID,
		},
	}
	handler := New(registerUC, nil, nil, nil, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(`{"email":"user@example.com","phone":"+5511999999999","password":"password123","name":"Maria Silva","birth_date":"1990-01-15","cpf":"12345678901","email_verification_token":"email-token","phone_verification_token":"phone-token"}`))
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	if !registerUC.called {
		t.Fatal("expected use case to be called")
	}

	if registerUC.input.Email != "user@example.com" {
		t.Fatalf("expected email %q, got %q", "user@example.com", registerUC.input.Email)
	}

	if registerUC.input.Phone != "+5511999999999" {
		t.Fatalf("expected phone %q, got %q", "+5511999999999", registerUC.input.Phone)
	}

	if registerUC.input.Name != "Maria Silva" {
		t.Fatalf("expected name %q, got %q", "Maria Silva", registerUC.input.Name)
	}

	if registerUC.input.BirthDate.Format("2006-01-02") != "1990-01-15" {
		t.Fatalf("expected birth date %q, got %q", "1990-01-15", registerUC.input.BirthDate.Format("2006-01-02"))
	}

	if registerUC.input.CPF != "12345678901" {
		t.Fatalf("expected cpf %q, got %q", "12345678901", registerUC.input.CPF)
	}

	if registerUC.input.EmailVerificationToken != "email-token" {
		t.Fatalf("expected email verification token %q, got %q", "email-token", registerUC.input.EmailVerificationToken)
	}

	if registerUC.input.PhoneVerificationToken != "phone-token" {
		t.Fatalf("expected phone verification token %q, got %q", "phone-token", registerUC.input.PhoneVerificationToken)
	}

	var got struct {
		Data struct {
			ID         string `json:"id"`
			Email      string `json:"email"`
			Role       string `json:"role"`
			CustomerID string `json:"customer_id"`
		} `json:"data"`
		Error any `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Data.ID != userID.String() {
		t.Fatalf("expected id %q, got %q", userID.String(), got.Data.ID)
	}

	if got.Data.CustomerID != customerID.String() {
		t.Fatalf("expected customer_id %q, got %q", customerID.String(), got.Data.CustomerID)
	}

	if got.Error != nil {
		t.Fatalf("expected nil error, got %#v", got.Error)
	}
}

func TestHandler_RequestContactVerification_Success(t *testing.T) {
	verificationID := uuid.New()
	debugToken := "123456"
	expiresAt := time.Date(2026, 5, 18, 12, 10, 0, 0, time.UTC)
	requestUC := &requestContactVerificationUseCaseMock{
		output: &application.RequestContactVerificationOutput{
			VerificationID: verificationID,
			Channel:        "email",
			Target:         "user@example.com",
			ExpiresAt:      expiresAt,
			DebugToken:     &debugToken,
		},
	}
	handler := New(nil, nil, nil, nil, requestUC, nil)
	req := httptest.NewRequest(http.MethodPost, "/auth/contact-verifications", strings.NewReader(`{"channel":"email","target":"user@example.com"}`))
	rec := httptest.NewRecorder()

	handler.RequestContactVerification(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
	if !requestUC.called {
		t.Fatal("expected use case to be called")
	}
	if requestUC.input.Channel != "email" {
		t.Fatalf("expected channel email, got %q", requestUC.input.Channel)
	}
	if requestUC.input.Target != "user@example.com" {
		t.Fatalf("expected target user@example.com, got %q", requestUC.input.Target)
	}

	var got struct {
		Data struct {
			VerificationID string `json:"verification_id"`
			Channel        string `json:"channel"`
			Target         string `json:"target"`
			DebugToken     string `json:"debug_token"`
		} `json:"data"`
		Error any `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if got.Data.VerificationID != verificationID.String() {
		t.Fatalf("expected verification ID %q, got %q", verificationID.String(), got.Data.VerificationID)
	}
	if got.Data.DebugToken != "123456" {
		t.Fatalf("expected debug_token %q, got %q", "123456", got.Data.DebugToken)
	}
	if got.Error != nil {
		t.Fatalf("expected nil error, got %#v", got.Error)
	}
}

func TestHandler_ConfirmContactVerification_Success(t *testing.T) {
	verificationID := uuid.New()
	verifiedAt := time.Date(2026, 5, 18, 12, 5, 0, 0, time.UTC)
	confirmUC := &confirmContactVerificationUseCaseMock{
		output: &application.ConfirmContactVerificationOutput{
			VerificationToken: "verified-token",
			Channel:           "phone",
			Target:            "+5511999999999",
			VerifiedAt:        verifiedAt,
		},
	}
	handler := New(nil, nil, nil, nil, nil, confirmUC)
	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/contact-verifications/confirm",
		strings.NewReader(`{"verification_id":"`+verificationID.String()+`","token":"123456"}`),
	)
	rec := httptest.NewRecorder()

	handler.ConfirmContactVerification(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !confirmUC.called {
		t.Fatal("expected use case to be called")
	}
	if confirmUC.input.VerificationID != verificationID {
		t.Fatalf("expected verification ID %v, got %v", verificationID, confirmUC.input.VerificationID)
	}
	if confirmUC.input.Token != "123456" {
		t.Fatalf("expected token %q, got %q", "123456", confirmUC.input.Token)
	}

	var got struct {
		Data struct {
			VerificationToken string `json:"verification_token"`
			Channel           string `json:"channel"`
			Target            string `json:"target"`
		} `json:"data"`
		Error any `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if got.Data.VerificationToken != "verified-token" {
		t.Fatalf("expected verification token %q, got %q", "verified-token", got.Data.VerificationToken)
	}
	if got.Error != nil {
		t.Fatalf("expected nil error, got %#v", got.Error)
	}
}

func TestHandler_Register_UserAlreadyExists(t *testing.T) {
	registerUC := &registerUserUseCaseMock{err: domain.ErrEmailAlreadyExists}
	handler := New(registerUC, nil, nil, nil, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(`{"email":"user@example.com","phone":"+5511999999999","password":"password123","name":"Maria Silva","birth_date":"1990-01-15","cpf":"12345678901","email_verification_token":"email-token","phone_verification_token":"phone-token"}`))
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}

	var got struct {
		Data  any `json:"data"`
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Data != nil {
		t.Fatalf("expected nil data, got %#v", got.Data)
	}

	if got.Error.Code != "USER_ALREADY_EXISTS" {
		t.Fatalf("expected error code %q, got %q", "USER_ALREADY_EXISTS", got.Error.Code)
	}
}

func TestHandler_Register_InvalidInput(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "legacy payload is rejected",
			body: `{"email":"user@example.com","password":"password123"}`,
		},
		{
			name: "empty password",
			body: `{"email":"user@example.com","phone":"+5511999999999","password":"   ","name":"Maria Silva","birth_date":"1990-01-15","cpf":"12345678901","email_verification_token":"email-token","phone_verification_token":"phone-token"}`,
		},
		{
			name: "empty name",
			body: `{"email":"user@example.com","phone":"+5511999999999","password":"password123","name":"   ","birth_date":"1990-01-15","cpf":"12345678901","email_verification_token":"email-token","phone_verification_token":"phone-token"}`,
		},
		{
			name: "empty birth date",
			body: `{"email":"user@example.com","phone":"+5511999999999","password":"password123","name":"Maria Silva","birth_date":"   ","cpf":"12345678901","email_verification_token":"email-token","phone_verification_token":"phone-token"}`,
		},
		{
			name: "invalid birth date",
			body: `{"email":"user@example.com","phone":"+5511999999999","password":"password123","name":"Maria Silva","birth_date":"15/01/1990","cpf":"12345678901","email_verification_token":"email-token","phone_verification_token":"phone-token"}`,
		},
		{
			name: "empty phone",
			body: `{"email":"user@example.com","phone":"   ","password":"password123","name":"Maria Silva","birth_date":"1990-01-15","cpf":"12345678901","email_verification_token":"email-token","phone_verification_token":"phone-token"}`,
		},
		{
			name: "empty email verification token",
			body: `{"email":"user@example.com","phone":"+5511999999999","password":"password123","name":"Maria Silva","birth_date":"1990-01-15","cpf":"12345678901","email_verification_token":"   ","phone_verification_token":"phone-token"}`,
		},
		{
			name: "empty phone verification token",
			body: `{"email":"user@example.com","phone":"+5511999999999","password":"password123","name":"Maria Silva","birth_date":"1990-01-15","cpf":"12345678901","email_verification_token":"email-token","phone_verification_token":"   "}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			registerUC := &registerUserUseCaseMock{}
			handler := New(registerUC, nil, nil, nil, nil, nil)
			req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(tc.body))
			rec := httptest.NewRecorder()

			handler.Register(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
			}

			if registerUC.called {
				t.Fatal("expected use case not to be called")
			}

			var got struct {
				Error struct {
					Code string `json:"code"`
				} `json:"error"`
			}

			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if got.Error.Code != "INVALID_REQUEST" {
				t.Fatalf("expected error code %q, got %q", "INVALID_REQUEST", got.Error.Code)
			}
		})
	}
}

func TestHandler_Login_Success(t *testing.T) {
	userID := uuid.New()
	customerID := uuid.New()
	loginUC := &loginUserUseCaseMock{
		output: &application.LoginUserOutput{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			UserID:       userID,
			Email:        "user@example.com",
			Role:         string(domain.RoleCustomer),
			CustomerID:   &customerID,
		},
	}
	handler := New(nil, loginUC, nil, nil, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"user@example.com","password":"password123"}`))
	req.Header.Set(sharedheaders.InstallationID, "550e8400-e29b-41d4-a716-446655440000")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if !loginUC.called {
		t.Fatal("expected use case to be called")
	}

	var got struct {
		Data struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			UserID       string `json:"user_id"`
			Email        string `json:"email"`
			Role         string `json:"role"`
			CustomerID   string `json:"customer_id"`
		} `json:"data"`
		Error any `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Data.AccessToken != "access-token" {
		t.Fatalf("expected access token %q, got %q", "access-token", got.Data.AccessToken)
	}

	if got.Data.RefreshToken != "refresh-token" {
		t.Fatalf("expected refresh token %q, got %q", "refresh-token", got.Data.RefreshToken)
	}

	if got.Data.UserID != userID.String() {
		t.Fatalf("expected user id %q, got %q", userID.String(), got.Data.UserID)
	}

	if loginUC.input.InstallationID.String() != "550e8400-e29b-41d4-a716-446655440000" {
		t.Fatalf("expected installation id to be propagated, got %q", loginUC.input.InstallationID.String())
	}

	if got.Error != nil {
		t.Fatalf("expected nil error, got %#v", got.Error)
	}
}

func TestHandler_Login_InvalidCredentials(t *testing.T) {
	loginUC := &loginUserUseCaseMock{err: domain.ErrInvalidCredentials}
	handler := New(nil, loginUC, nil, nil, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"user@example.com","password":"wrong"}`))
	req.Header.Set(sharedheaders.InstallationID, "550e8400-e29b-41d4-a716-446655440000")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	var got struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Error.Code != "INVALID_CREDENTIALS" {
		t.Fatalf("expected error code %q, got %q", "INVALID_CREDENTIALS", got.Error.Code)
	}
}

func TestHandler_Login_AccountApprovalRequired(t *testing.T) {
	loginUC := &loginUserUseCaseMock{err: domain.ErrAccountApprovalRequired}
	handler := New(nil, loginUC, nil, nil, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"user@example.com","password":"password123"}`))
	req.Header.Set(sharedheaders.InstallationID, "550e8400-e29b-41d4-a716-446655440000")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}

	var got struct {
		Data  any `json:"data"`
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Data != nil {
		t.Fatalf("expected nil data, got %#v", got.Data)
	}

	if got.Error.Code != "ACCOUNT_APPROVAL_REQUIRED" {
		t.Fatalf("expected error code %q, got %q", "ACCOUNT_APPROVAL_REQUIRED", got.Error.Code)
	}

	if got.Error.Message != "Account approval required" {
		t.Fatalf("expected error message %q, got %q", "Account approval required", got.Error.Message)
	}
}

func TestHandler_Login_ContactNotVerified(t *testing.T) {
	loginUC := &loginUserUseCaseMock{err: domain.NewContactNotVerifiedError(false, true)}
	handler := New(nil, loginUC, nil, nil, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"user@example.com","password":"password123"}`))
	req.Header.Set(sharedheaders.InstallationID, "550e8400-e29b-41d4-a716-446655440000")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}

	var got struct {
		Data  any `json:"data"`
		Error struct {
			Code    string          `json:"code"`
			Message string          `json:"message"`
			Details map[string]bool `json:"details"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Data != nil {
		t.Fatalf("expected nil data, got %#v", got.Data)
	}

	if got.Error.Code != "CONTACT_NOT_VERIFIED" {
		t.Fatalf("expected error code %q, got %q", "CONTACT_NOT_VERIFIED", got.Error.Code)
	}

	if got.Error.Message != "Contact not verified" {
		t.Fatalf("expected error message %q, got %q", "Contact not verified", got.Error.Message)
	}

	if got.Error.Details["email_verified"] {
		t.Fatal("expected email_verified to be false")
	}

	if !got.Error.Details["phone_verified"] {
		t.Fatal("expected phone_verified to be true")
	}
}

func TestHandler_Login_MissingInstallationID(t *testing.T) {
	loginUC := &loginUserUseCaseMock{}
	handler := New(nil, loginUC, nil, nil, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"user@example.com","password":"password123"}`))
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	if loginUC.called {
		t.Fatal("expected use case not to be called")
	}

	var got struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Error.Code != "INVALID_INSTALLATION_ID" {
		t.Fatalf("expected error code %q, got %q", "INVALID_INSTALLATION_ID", got.Error.Code)
	}

	if got.Error.Message != "X-Installation-Id must be a canonical UUID v4." {
		t.Fatalf("expected error message %q, got %q", "X-Installation-Id must be a canonical UUID v4.", got.Error.Message)
	}
}

func TestHandler_Login_InvalidInstallationID(t *testing.T) {
	loginUC := &loginUserUseCaseMock{}
	handler := New(nil, loginUC, nil, nil, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"user@example.com","password":"password123"}`))
	req.Header.Set(sharedheaders.InstallationID, "not-a-uuid")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	if loginUC.called {
		t.Fatal("expected use case not to be called")
	}

	var got struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Error.Code != "INVALID_INSTALLATION_ID" {
		t.Fatalf("expected error code %q, got %q", "INVALID_INSTALLATION_ID", got.Error.Code)
	}
}

func TestHandler_Refresh_Success(t *testing.T) {
	refreshUC := &refreshAccessTokenUseCaseMock{output: &application.RefreshAccessTokenOutput{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
	}}
	handler := New(nil, nil, nil, refreshUC, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(`{"refresh_token":"valid-refresh-token"}`))
	rec := httptest.NewRecorder()

	handler.Refresh(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if !refreshUC.called {
		t.Fatal("expected use case to be called")
	}

	if refreshUC.input.RefreshToken != "valid-refresh-token" {
		t.Fatalf("expected refresh token %q, got %q", "valid-refresh-token", refreshUC.input.RefreshToken)
	}

	var got struct {
		Data struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"data"`
		Error any `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Data.AccessToken != "new-access-token" {
		t.Fatalf("expected access token %q, got %q", "new-access-token", got.Data.AccessToken)
	}

	if got.Data.RefreshToken != "new-refresh-token" {
		t.Fatalf("expected refresh token %q, got %q", "new-refresh-token", got.Data.RefreshToken)
	}

	if got.Error != nil {
		t.Fatalf("expected nil error, got %#v", got.Error)
	}
}

func TestHandler_Refresh_InvalidToken(t *testing.T) {
	refreshUC := &refreshAccessTokenUseCaseMock{err: domain.ErrInvalidToken}
	handler := New(nil, nil, nil, refreshUC, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(`{"refresh_token":"bad-token"}`))
	rec := httptest.NewRecorder()

	handler.Refresh(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	var got struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Error.Code != "INVALID_TOKEN" {
		t.Fatalf("expected error code %q, got %q", "INVALID_TOKEN", got.Error.Code)
	}
}

func TestHandler_Me_Unauthorized(t *testing.T) {
	currentUserUC := &getCurrentUserUseCaseMock{err: domain.ErrUnauthorized}
	handler := New(nil, nil, currentUserUC, nil, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	rec := httptest.NewRecorder()

	handler.Me(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	var got struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Error.Code != "UNAUTHORIZED" {
		t.Fatalf("expected error code %q, got %q", "UNAUTHORIZED", got.Error.Code)
	}
}

func TestHandler_Session_Success(t *testing.T) {
	userID := uuid.New()
	customerID := uuid.New()
	birthDate := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	createdAt := time.Date(2026, 6, 2, 10, 0, 0, 0, time.UTC)
	sessionUC := &getSessionUseCaseMock{
		output: &application.GetSessionOutput{
			User: application.GetSessionUserOutput{
				ID:    userID,
				Email: "user@example.com",
				Phone: "+5527999999999",
				Role:  "customer",
			},
			Customer: application.GetSessionCustomerOutput{
				ID:        customerID,
				Name:      "Maria Silva",
				CPF:       "12345678901",
				BirthDate: birthDate,
				CreatedAt: createdAt,
			},
			Readiness: application.GetSessionReadinessOutput{
				OnboardingCompleted:       true,
				Approved:                  true,
				HasOperationalAccount:     true,
				TransactionPasswordStatus: "active",
			},
		},
	}
	handler := New(nil, nil, nil, nil, nil, nil, sessionUC)
	req := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	rec := httptest.NewRecorder()

	handler.Session(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !sessionUC.called {
		t.Fatal("expected use case to be called")
	}

	var got struct {
		Data struct {
			User      map[string]any `json:"user"`
			Customer  map[string]any `json:"customer"`
			Readiness struct {
				OnboardingCompleted       bool   `json:"onboarding_completed"`
				Approved                  bool   `json:"approved"`
				HasOperationalAccount     bool   `json:"has_operational_account"`
				TransactionPasswordStatus string `json:"transaction_password_status"`
			} `json:"readiness"`
		} `json:"data"`
		Error any `json:"error"`
	}
	rawBody := rec.Body.String()
	if err := json.Unmarshal([]byte(rawBody), &got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Error != nil {
		t.Fatalf("expected nil error, got %#v", got.Error)
	}
	if got.Data.User["phone"] != "+5527999999999" {
		t.Fatalf("expected user phone %q, got %#v", "+5527999999999", got.Data.User["phone"])
	}
	if _, ok := got.Data.User["customer_id"]; ok {
		t.Fatal("expected user.customer_id to be absent")
	}
	if _, ok := got.Data.Customer["email"]; ok {
		t.Fatal("expected customer.email to be absent")
	}
	if got.Data.Customer["birth_date"] != "1990-01-15" {
		t.Fatalf("expected birth_date %q, got %#v", "1990-01-15", got.Data.Customer["birth_date"])
	}
	if got.Data.Readiness.TransactionPasswordStatus != "active" {
		t.Fatalf("expected transaction_password_status active, got %q", got.Data.Readiness.TransactionPasswordStatus)
	}
	if strings.Contains(rawBody, `"can_access_home"`) {
		t.Fatal("expected can_access_home to be absent")
	}
}

func TestHandler_Session_Unauthorized(t *testing.T) {
	sessionUC := &getSessionUseCaseMock{err: domain.ErrUnauthorized}
	handler := New(nil, nil, nil, nil, nil, nil, sessionUC)
	req := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	rec := httptest.NewRecorder()

	handler.Session(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	var got struct {
		Data  any `json:"data"`
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if got.Data != nil {
		t.Fatalf("expected nil data, got %#v", got.Data)
	}
	if got.Error.Code != "UNAUTHORIZED" {
		t.Fatalf("expected error code %q, got %q", "UNAUTHORIZED", got.Error.Code)
	}
}
