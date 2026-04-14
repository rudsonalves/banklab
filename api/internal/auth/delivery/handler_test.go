package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/application"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
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

type refreshAccessTokenUseCaseMock struct {
	output *application.RefreshAccessTokenOutput
	err    error
	input  application.RefreshAccessTokenInput
	called bool
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
	handler := New(registerUC, nil, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(`{"email":"user@example.com","password":"password123","name":"Maria Silva","cpf":"12345678901"}`))
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

	if registerUC.input.Name != "Maria Silva" {
		t.Fatalf("expected name %q, got %q", "Maria Silva", registerUC.input.Name)
	}

	if registerUC.input.CPF != "12345678901" {
		t.Fatalf("expected cpf %q, got %q", "12345678901", registerUC.input.CPF)
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

func TestHandler_Register_UserAlreadyExists(t *testing.T) {
	registerUC := &registerUserUseCaseMock{err: domain.ErrEmailAlreadyExists}
	handler := New(registerUC, nil, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(`{"email":"user@example.com","password":"password123","name":"Maria Silva","cpf":"12345678901"}`))
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
			body: `{"email":"user@example.com","password":"   ","name":"Maria Silva","cpf":"12345678901"}`,
		},
		{
			name: "empty name",
			body: `{"email":"user@example.com","password":"password123","name":"   ","cpf":"12345678901"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			registerUC := &registerUserUseCaseMock{}
			handler := New(registerUC, nil, nil, nil)
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
	handler := New(nil, loginUC, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"user@example.com","password":"password123"}`))
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

	if got.Error != nil {
		t.Fatalf("expected nil error, got %#v", got.Error)
	}
}

func TestHandler_Login_InvalidCredentials(t *testing.T) {
	loginUC := &loginUserUseCaseMock{err: domain.ErrInvalidCredentials}
	handler := New(nil, loginUC, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"user@example.com","password":"wrong"}`))
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

func TestHandler_Refresh_Success(t *testing.T) {
	refreshUC := &refreshAccessTokenUseCaseMock{output: &application.RefreshAccessTokenOutput{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
	}}
	handler := New(nil, nil, nil, refreshUC)
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
	handler := New(nil, nil, nil, refreshUC)
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
	handler := New(nil, nil, currentUserUC, nil)
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
