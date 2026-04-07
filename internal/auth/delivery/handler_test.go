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

func TestHandler_Register_Success(t *testing.T) {
	userID := uuid.New()
	registerUC := &registerUserUseCaseMock{
		output: &application.RegisterUserOutput{
			ID:    userID,
			Email: "user@example.com",
			Role:  "customer",
		},
	}
	handler := New(registerUC, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(`{"email":"user@example.com","password":"password123"}`))
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

	var got struct {
		Data struct {
			ID    string `json:"id"`
			Email string `json:"email"`
			Role  string `json:"role"`
		} `json:"data"`
		Error any `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Data.ID != userID.String() {
		t.Fatalf("expected id %q, got %q", userID.String(), got.Data.ID)
	}

	if got.Error != nil {
		t.Fatalf("expected nil error, got %#v", got.Error)
	}
}

func TestHandler_Register_UserAlreadyExists(t *testing.T) {
	registerUC := &registerUserUseCaseMock{err: domain.ErrEmailAlreadyExists}
	handler := New(registerUC, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(`{"email":"user@example.com","password":"password123"}`))
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

func TestHandler_Login_InvalidCredentials(t *testing.T) {
	loginUC := &loginUserUseCaseMock{err: domain.ErrInvalidCredentials}
	handler := New(nil, loginUC, nil)
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

func TestHandler_Me_Unauthorized(t *testing.T) {
	currentUserUC := &getCurrentUserUseCaseMock{err: domain.ErrUnauthorized}
	handler := New(nil, nil, currentUserUC)
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
