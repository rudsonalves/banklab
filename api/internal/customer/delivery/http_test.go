package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	"github.com/seu-usuario/bank-api/internal/customer/application"
	"github.com/seu-usuario/bank-api/internal/customer/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
)

type getCustomerMeUseCaseMock struct {
	output *domain.CustomerProfile
	err    error
	called bool
	input  application.GetCustomerMeInput
}

func (m *getCustomerMeUseCaseMock) Execute(ctx context.Context, input application.GetCustomerMeInput) (*domain.CustomerProfile, error) {
	m.called = true
	m.input = input
	return m.output, m.err
}

var registerErrorsOnce sync.Once

func ensureErrorsRegistered() {
	registerErrorsOnce.Do(func() {
		application.RegisterErrors()
		sharederrors.RegisterDomainError(
			authdomain.ErrInvalidUserState,
			sharederrors.ErrCodeInvalidUserState,
			"Invalid user state",
			http.StatusConflict,
		)
		sharederrors.Register(func(err error) bool {
			return errors.Is(err, authdomain.ErrUnauthorized)
		}, sharederrors.AppError{
			Code:    sharederrors.ErrCodeUnauthorized,
			Message: "Authentication required",
			Status:  http.StatusUnauthorized,
		})
	})
}

func TestHandler_Me_Success(t *testing.T) {
	ensureErrorsRegistered()

	customerID := uuid.New()
	createdAt := time.Now().UTC().Truncate(time.Second)
	birthDate := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	uc := &getCustomerMeUseCaseMock{output: &domain.CustomerProfile{
		Customer: domain.Customer{
			ID:        customerID,
			Name:      "Maria Silva",
			BirthDate: birthDate,
			CreatedAt: createdAt,
		},
		Email: "maria@example.com",
		CPF:   "12345678901",
	}}
	h := &Handler{getMeUC: uc}

	req := httptest.NewRequest(http.MethodGet, "/customers/me", nil)
	req = req.WithContext(sharedauthctx.WithAuthenticatedUser(req.Context(), sharedauthctx.AuthenticatedUser{
		UserID:     uuid.New(),
		Role:       authdomain.RoleCustomer,
		CustomerID: &customerID,
	}))
	rec := httptest.NewRecorder()

	h.Me(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !uc.called {
		t.Fatal("expected use case to be called")
	}
	if uc.input.CustomerID != customerID {
		t.Fatalf("expected customer ID %v, got %v", customerID, uc.input.CustomerID)
	}

	var got struct {
		Data struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			CPF       string `json:"cpf"`
			Email     string `json:"email"`
			BirthDate string `json:"birth_date"`
			CreatedAt string `json:"created_at"`
		} `json:"data"`
		Error any `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Data.ID != customerID.String() {
		t.Fatalf("expected id %q, got %q", customerID.String(), got.Data.ID)
	}
	if got.Data.Name != "Maria Silva" {
		t.Fatalf("expected name Maria Silva, got %q", got.Data.Name)
	}
	if got.Data.Email != "maria@example.com" {
		t.Fatalf("expected email %q, got %q", "maria@example.com", got.Data.Email)
	}
	if got.Data.CPF != "12345678901" {
		t.Fatalf("expected cpf %q, got %q", "12345678901", got.Data.CPF)
	}
	if got.Data.BirthDate != "1990-01-15" {
		t.Fatalf("expected birth date %q, got %q", "1990-01-15", got.Data.BirthDate)
	}
	if got.Error != nil {
		t.Fatalf("expected nil error, got %#v", got.Error)
	}
}

func TestHandler_Me_InvalidStateWhenCustomerIDMissing(t *testing.T) {
	ensureErrorsRegistered()

	h := &Handler{getMeUC: &getCustomerMeUseCaseMock{}}
	req := httptest.NewRequest(http.MethodGet, "/customers/me", nil)
	req = req.WithContext(sharedauthctx.WithAuthenticatedUser(req.Context(), sharedauthctx.AuthenticatedUser{
		UserID: uuid.New(),
		Role:   authdomain.RoleCustomer,
	}))
	rec := httptest.NewRecorder()

	h.Me(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}

	var got struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Error.Code != "INVALID_USER_STATE" {
		t.Fatalf("expected error code %q, got %q", "INVALID_USER_STATE", got.Error.Code)
	}
}

func TestHandler_Me_NotFound(t *testing.T) {
	ensureErrorsRegistered()

	customerID := uuid.New()
	uc := &getCustomerMeUseCaseMock{err: domain.ErrNotFound}
	h := &Handler{getMeUC: uc}
	req := httptest.NewRequest(http.MethodGet, "/customers/me", nil)
	req = req.WithContext(sharedauthctx.WithAuthenticatedUser(req.Context(), sharedauthctx.AuthenticatedUser{
		UserID:     uuid.New(),
		Role:       authdomain.RoleCustomer,
		CustomerID: &customerID,
	}))
	rec := httptest.NewRecorder()

	h.Me(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
