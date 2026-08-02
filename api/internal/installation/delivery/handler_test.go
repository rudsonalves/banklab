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
	"github.com/seu-usuario/bank-api/internal/installation/application"
	sharedheaders "github.com/seu-usuario/bank-api/internal/shared/http/headers"
)

type registerInstallationUseCaseMock struct {
	input  application.RegisterInstallationInput
	output *application.RegisterInstallationOutput
	err    error
	calls  int
}

func (m *registerInstallationUseCaseMock) Execute(
	ctx context.Context,
	input application.RegisterInstallationInput,
) (*application.RegisterInstallationOutput, error) {
	m.calls++
	m.input = input
	return m.output, m.err
}

type listInstallationsUseCaseMock struct {
	output *application.ListInstallationsOutput
	err    error
	calls  int
}

func (m *listInstallationsUseCaseMock) Execute(ctx context.Context) (*application.ListInstallationsOutput, error) {
	m.calls++
	return m.output, m.err
}

type revokeInstallationUseCaseMock struct {
	input  application.RevokeInstallationInput
	output *application.RevokeInstallationOutput
	err    error
	calls  int
}

func (m *revokeInstallationUseCaseMock) Execute(
	ctx context.Context,
	input application.RevokeInstallationInput,
) (*application.RevokeInstallationOutput, error) {
	m.calls++
	m.input = input
	return m.output, m.err
}

func TestHandler_Register_Success(t *testing.T) {
	installationID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	resourceID := uuid.New()
	registerUC := &registerInstallationUseCaseMock{output: &application.RegisterInstallationOutput{
		AccessToken:            "access-token",
		RefreshToken:           "refresh-token",
		InstallationResourceID: resourceID,
		InstallationStatus:     "known",
	}}
	handler := New(registerUC, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/security/installations", strings.NewReader(`{}`))
	req.Header.Set(sharedheaders.InstallationID, installationID.String())
	req.Header.Set(sharedheaders.StepUpToken, "step-up-token")
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
	if registerUC.calls != 1 {
		t.Fatalf("expected use case once, got %d", registerUC.calls)
	}
	if registerUC.input.PresentedInstallationID != installationID {
		t.Fatalf("expected installation id %q, got %q", installationID, registerUC.input.PresentedInstallationID)
	}
	if registerUC.input.StepUpToken != "step-up-token" {
		t.Fatalf("expected step-up token to be propagated")
	}

	var got struct {
		Data struct {
			AccessToken            string `json:"access_token"`
			RefreshToken           string `json:"refresh_token"`
			InstallationResourceID string `json:"installation_resource_id"`
			InstallationStatus     string `json:"installation_status"`
		} `json:"data"`
		Error any `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Data.AccessToken != "access-token" || got.Data.RefreshToken != "refresh-token" {
		t.Fatalf("unexpected token response: %#v", got.Data)
	}
	if got.Data.InstallationResourceID != resourceID.String() || got.Data.InstallationStatus != "known" {
		t.Fatalf("unexpected installation response: %#v", got.Data)
	}
}

func TestHandler_List_DoesNotExposeInstallationID(t *testing.T) {
	resourceID := uuid.New()
	now := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	listUC := &listInstallationsUseCaseMock{output: &application.ListInstallationsOutput{
		Installations: []application.InstallationSummary{{
			ResourceID:  resourceID,
			Status:      "known",
			FirstSeenAt: now,
			LastSeenAt:  now,
			CreatedAt:   now,
			UpdatedAt:   now,
		}},
	}}
	handler := New(nil, listUC, nil)
	req := httptest.NewRequest(http.MethodGet, "/security/installations", nil)
	rec := httptest.NewRecorder()

	handler.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if strings.Contains(rec.Body.String(), "installation_id") {
		t.Fatalf("expected response not to expose installation_id: %s", rec.Body.String())
	}
}

func TestHandler_Revoke_InvalidResourceID(t *testing.T) {
	revokeUC := &revokeInstallationUseCaseMock{}
	handler := New(nil, nil, revokeUC)
	req := httptest.NewRequest(http.MethodDelete, "/security/installations/not-a-uuid", nil)
	req.SetPathValue("installation_resource_id", "not-a-uuid")
	rec := httptest.NewRecorder()

	handler.Revoke(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if revokeUC.calls != 0 {
		t.Fatalf("expected use case not to be called, got %d", revokeUC.calls)
	}
}
