package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	authapplication "github.com/seu-usuario/bank-api/internal/auth/application"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	"github.com/seu-usuario/bank-api/internal/security/application"
	"github.com/seu-usuario/bank-api/internal/security/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
)

func init() {
	authapplication.RegisterErrors()
	application.RegisterErrors()
}

type createTransactionPasswordUseCaseMock struct {
	input  application.CreateTransactionPasswordInput
	output *application.CreateTransactionPasswordOutput
	err    error
	calls  int
}

func (m *createTransactionPasswordUseCaseMock) Execute(
	ctx context.Context,
	input application.CreateTransactionPasswordInput,
) (*application.CreateTransactionPasswordOutput, error) {
	m.calls++
	m.input = input
	return m.output, m.err
}

func TestHandler_CreateTransactionPassword_Success(t *testing.T) {
	userID := uuid.New()
	createdAt := time.Date(2026, 5, 28, 10, 0, 0, 0, time.UTC)
	useCase := &createTransactionPasswordUseCaseMock{
		output: &application.CreateTransactionPasswordOutput{
			UserID:    userID.String(),
			Status:    string(domain.TransactionPasswordActive),
			CreatedAt: createdAt,
		},
	}
	handler := New(useCase)
	req := httptest.NewRequest(
		http.MethodPost,
		"/security/transaction-password",
		strings.NewReader(`{"transaction_password":"123456","transaction_password_confirmation":"123456"}`),
	)
	req = req.WithContext(sharedauthctx.WithAuthenticatedUser(req.Context(), sharedauthctx.AuthenticatedUser{
		UserID: userID,
		Role:   authdomain.RoleCustomer,
	}))
	rec := httptest.NewRecorder()

	handler.CreateTransactionPassword(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
	if useCase.calls != 1 {
		t.Fatalf("expected Execute to be called once, got %d", useCase.calls)
	}
	if useCase.input.User == nil || useCase.input.User.UserID != userID {
		t.Fatalf("expected authenticated user %q, got %+v", userID, useCase.input.User)
	}
	if useCase.input.TransactionPassword != "123456" {
		t.Fatalf("expected transaction password to be passed to use case")
	}
	if useCase.input.TransactionPasswordConfirmation != "123456" {
		t.Fatalf("expected transaction password confirmation to be passed to use case")
	}

	var got struct {
		Data struct {
			UserID    string `json:"user_id"`
			Status    string `json:"status"`
			CreatedAt string `json:"created_at"`
		} `json:"data"`
		Error any `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if got.Error != nil {
		t.Fatalf("expected nil error, got %+v", got.Error)
	}
	if got.Data.UserID != userID.String() {
		t.Fatalf("expected user_id %q, got %q", userID.String(), got.Data.UserID)
	}
	if got.Data.Status != string(domain.TransactionPasswordActive) {
		t.Fatalf("expected status %q, got %q", domain.TransactionPasswordActive, got.Data.Status)
	}
	if got.Data.CreatedAt == "" {
		t.Fatal("expected created_at to be present")
	}
}

func TestHandler_CreateTransactionPassword_Unauthorized(t *testing.T) {
	useCase := &createTransactionPasswordUseCaseMock{}
	handler := New(useCase)
	req := httptest.NewRequest(
		http.MethodPost,
		"/security/transaction-password",
		strings.NewReader(`{"transaction_password":"123456","transaction_password_confirmation":"123456"}`),
	)
	rec := httptest.NewRecorder()

	handler.CreateTransactionPassword(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
	if useCase.calls != 0 {
		t.Fatalf("expected Execute not to be called, got %d", useCase.calls)
	}
}

func TestHandler_CreateTransactionPassword_InvalidPayload(t *testing.T) {
	userID := uuid.New()
	useCase := &createTransactionPasswordUseCaseMock{}
	handler := New(useCase)
	req := httptest.NewRequest(
		http.MethodPost,
		"/security/transaction-password",
		strings.NewReader(`{"transaction_password":"123456","transaction_password_confirmation":"123456","extra":true}`),
	)
	req = req.WithContext(sharedauthctx.WithAuthenticatedUser(req.Context(), sharedauthctx.AuthenticatedUser{
		UserID: userID,
		Role:   authdomain.RoleCustomer,
	}))
	rec := httptest.NewRecorder()

	handler.CreateTransactionPassword(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if useCase.calls != 0 {
		t.Fatalf("expected Execute not to be called, got %d", useCase.calls)
	}
}

func TestHandler_CreateTransactionPassword_MapsDomainError(t *testing.T) {
	userID := uuid.New()
	useCase := &createTransactionPasswordUseCaseMock{err: domain.ErrTransactionPasswordAlreadySet}
	handler := New(useCase)
	req := httptest.NewRequest(
		http.MethodPost,
		"/security/transaction-password",
		strings.NewReader(`{"transaction_password":"123456","transaction_password_confirmation":"123456"}`),
	)
	req = req.WithContext(sharedauthctx.WithAuthenticatedUser(req.Context(), sharedauthctx.AuthenticatedUser{
		UserID: userID,
		Role:   authdomain.RoleCustomer,
	}))
	rec := httptest.NewRecorder()

	handler.CreateTransactionPassword(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}

	var got struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if got.Error.Code != "TRANSACTION_PASSWORD_ALREADY_SET" {
		t.Fatalf("expected error code %q, got %q", "TRANSACTION_PASSWORD_ALREADY_SET", got.Error.Code)
	}
}

func TestHandler_CreateTransactionPassword_InternalErrorOnNilUseCase(t *testing.T) {
	handler := New(nil)
	req := httptest.NewRequest(
		http.MethodPost,
		"/security/transaction-password",
		strings.NewReader(`{"transaction_password":"123456","transaction_password_confirmation":"123456"}`),
	)
	req = req.WithContext(sharedauthctx.WithAuthenticatedUser(req.Context(), sharedauthctx.AuthenticatedUser{
		UserID: uuid.New(),
		Role:   authdomain.RoleCustomer,
	}))
	rec := httptest.NewRecorder()

	handler.CreateTransactionPassword(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestHandler_CreateTransactionPassword_UnknownError(t *testing.T) {
	userID := uuid.New()
	useCase := &createTransactionPasswordUseCaseMock{err: errors.New("boom")}
	handler := New(useCase)
	req := httptest.NewRequest(
		http.MethodPost,
		"/security/transaction-password",
		strings.NewReader(`{"transaction_password":"123456","transaction_password_confirmation":"123456"}`),
	)
	req = req.WithContext(sharedauthctx.WithAuthenticatedUser(req.Context(), sharedauthctx.AuthenticatedUser{
		UserID: userID,
		Role:   authdomain.RoleCustomer,
	}))
	rec := httptest.NewRecorder()

	handler.CreateTransactionPassword(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}
