package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	adminapplication "github.com/seu-usuario/bank-api/internal/admin/application"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
	"github.com/seu-usuario/bank-api/internal/bootstrap"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
)

func TestMain(m *testing.M) {
	bootstrap.RegisterErrors()
	os.Exit(m.Run())
}

type approveUserUseCaseMock struct {
	output *adminapplication.ApproveUserOutput
	err    error
	input  adminapplication.ApproveUserInput
	called bool
}

func (m *approveUserUseCaseMock) Execute(ctx context.Context, input adminapplication.ApproveUserInput) (*adminapplication.ApproveUserOutput, error) {
	m.called = true
	m.input = input
	return m.output, m.err
}

func TestHandler_ApproveUser_Success(t *testing.T) {
	userID := uuid.New()
	accountID := uuid.New()
	approveUC := &approveUserUseCaseMock{
		output: &adminapplication.ApproveUserOutput{
			UserID:    userID,
			Status:    string(domain.UserStatusActive),
			AccountID: accountID,
		},
	}
	handler := New(approveUC)
	req := httptest.NewRequest(http.MethodPost, "/admin/users/"+userID.String()+"/approve", nil)
	req.SetPathValue("id", userID.String())
	req = req.WithContext(sharedauthctx.WithAuthenticatedUser(req.Context(), sharedauthctx.AuthenticatedUser{UserID: uuid.New(), Role: domain.RoleAdmin}))
	rec := httptest.NewRecorder()

	handler.ApproveUser(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if !approveUC.called {
		t.Fatal("expected use case to be called")
	}

	if approveUC.input.UserID != userID {
		t.Fatalf("expected user id %v, got %v", userID, approveUC.input.UserID)
	}

	var got struct {
		Data struct {
			UserID    string `json:"user_id"`
			Status    string `json:"status"`
			AccountID string `json:"account_id"`
		} `json:"data"`
		Error any `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Data.UserID != userID.String() {
		t.Fatalf("expected user id %q, got %q", userID.String(), got.Data.UserID)
	}

	if got.Data.Status != string(domain.UserStatusActive) {
		t.Fatalf("expected status %q, got %q", domain.UserStatusActive, got.Data.Status)
	}

	if got.Data.AccountID != accountID.String() {
		t.Fatalf("expected account id %q, got %q", accountID.String(), got.Data.AccountID)
	}

	if got.Error != nil {
		t.Fatalf("expected nil error, got %#v", got.Error)
	}
}

func TestHandler_ApproveUser_Unauthorized(t *testing.T) {
	approveUC := &approveUserUseCaseMock{}
	handler := New(approveUC)
	req := httptest.NewRequest(http.MethodPost, "/admin/users/"+uuid.NewString()+"/approve", nil)
	req.SetPathValue("id", uuid.NewString())
	rec := httptest.NewRecorder()

	handler.ApproveUser(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	if approveUC.called {
		t.Fatal("expected use case not to be called")
	}
}

func TestHandler_ApproveUser_RejectsNonAdmin(t *testing.T) {
	approveUC := &approveUserUseCaseMock{}
	userID := uuid.New()
	handler := New(approveUC)
	req := httptest.NewRequest(http.MethodPost, "/admin/users/"+userID.String()+"/approve", nil)
	req.SetPathValue("id", userID.String())
	req = req.WithContext(sharedauthctx.WithAuthenticatedUser(req.Context(), sharedauthctx.AuthenticatedUser{UserID: uuid.New(), Role: domain.RoleCustomer}))
	rec := httptest.NewRecorder()

	handler.ApproveUser(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}

	if approveUC.called {
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

	if got.Error.Code != "FORBIDDEN" {
		t.Fatalf("expected error code %q, got %q", "FORBIDDEN", got.Error.Code)
	}
}

func TestHandler_ApproveUser_InvalidUserID(t *testing.T) {
	approveUC := &approveUserUseCaseMock{}
	handler := New(approveUC)
	req := httptest.NewRequest(http.MethodPost, "/admin/users/invalid/approve", nil)
	req.SetPathValue("id", "invalid")
	req = req.WithContext(sharedauthctx.WithAuthenticatedUser(req.Context(), sharedauthctx.AuthenticatedUser{UserID: uuid.New(), Role: domain.RoleAdmin}))
	rec := httptest.NewRecorder()

	handler.ApproveUser(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	if approveUC.called {
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

	if got.Error.Code != "INVALID_DATA" {
		t.Fatalf("expected error code %q, got %q", "INVALID_DATA", got.Error.Code)
	}
}

func TestHandler_ApproveUser_MapsUseCaseErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{name: "user not found", err: domain.ErrUserNotFound, wantStatus: http.StatusNotFound, wantCode: "USER_NOT_FOUND"},
		{name: "user already active", err: domain.ErrUserAlreadyActive, wantStatus: http.StatusConflict, wantCode: "USER_ALREADY_ACTIVE"},
		{name: "forbidden", err: domain.ErrForbidden, wantStatus: http.StatusForbidden, wantCode: "FORBIDDEN"},
		{name: "internal", err: context.DeadlineExceeded, wantStatus: http.StatusInternalServerError, wantCode: "INTERNAL_ERROR"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			approveUC := &approveUserUseCaseMock{err: tc.err}
			targetUserID := uuid.New()
			handler := New(approveUC)
			req := httptest.NewRequest(http.MethodPost, "/admin/users/"+targetUserID.String()+"/approve", nil)
			req.SetPathValue("id", targetUserID.String())
			req = req.WithContext(sharedauthctx.WithAuthenticatedUser(req.Context(), sharedauthctx.AuthenticatedUser{UserID: uuid.New(), Role: domain.RoleAdmin}))
			rec := httptest.NewRecorder()

			handler.ApproveUser(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("expected status %d, got %d", tc.wantStatus, rec.Code)
			}

			var got struct {
				Error struct {
					Code string `json:"code"`
				} `json:"error"`
			}

			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if got.Error.Code != tc.wantCode {
				t.Fatalf("expected error code %q, got %q", tc.wantCode, got.Error.Code)
			}
		})
	}
}
