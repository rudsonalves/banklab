package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	accountDelivery "github.com/seu-usuario/bank-api/internal/account/bankaccount/delivery"
	statementDelivery "github.com/seu-usuario/bank-api/internal/account/statement/delivery"
	transactionDelivery "github.com/seu-usuario/bank-api/internal/account/transaction/delivery"
	adminDelivery "github.com/seu-usuario/bank-api/internal/admin/delivery"
	authDelivery "github.com/seu-usuario/bank-api/internal/auth/delivery"
	customerDelivery "github.com/seu-usuario/bank-api/internal/customer/delivery"
	sharedhttpmiddleware "github.com/seu-usuario/bank-api/internal/shared/http/middleware"
)

func TestAuthRouter_OnboardingRoutesRequireAppToken(t *testing.T) {
	router := newAuthRouter(
		authDelivery.New(nil, nil, nil, nil, nil, nil),
		sharedhttpmiddleware.AppToken("expected-token"),
		func(next http.Handler) http.Handler { return next },
	)

	tests := []struct {
		name string
		path string
		body string
	}{
		{
			name: "request contact verification",
			path: "/auth/contact-verifications",
			body: `{"channel":"email","target":"user@example.com"}`,
		},
		{
			name: "confirm contact verification",
			path: "/auth/contact-verifications/confirm",
			body: `{"verification_id":"00000000-0000-0000-0000-000000000000","token":"123456"}`,
		},
		{
			name: "register",
			path: "/auth/register",
			body: `{"email":"user@example.com","phone":"+5511999999999","password":"password123","name":"Maria Silva","birth_date":"1990-01-15","cpf":"12345678901","email_verification_token":"email-token","phone_verification_token":"phone-token"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name+" missing app token", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tc.path, strings.NewReader(tc.body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assertInvalidAppToken(t, rec)
		})

		t.Run(tc.name+" invalid app token", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tc.path, strings.NewReader(tc.body))
			req.Header.Set("X-App-Token", "wrong-token")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assertInvalidAppToken(t, rec)
		})

		t.Run(tc.name+" valid app token reaches handler", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tc.path, strings.NewReader(tc.body))
			req.Header.Set("X-App-Token", "expected-token")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code == http.StatusUnauthorized {
				t.Fatal("expected valid app token to pass middleware")
			}
		})
	}
}

func TestAPIRouter_OperationalAccountRoutes(t *testing.T) {
	router := newAPIRouter(
		func(next http.Handler) http.Handler { return next },
		adminDelivery.New(nil),
		accountDelivery.New(nil, nil, nil, nil),
		customerDelivery.New(nil, nil),
		statementDelivery.New(nil),
		transactionDelivery.New(nil, nil, nil, nil),
	)

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{
			name:       "customer account creation is not registered",
			method:     http.MethodPost,
			path:       "/accounts",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "terminal deposit is not registered",
			method:     http.MethodPost,
			path:       "/terminal/accounts/00000000-0000-0000-0000-000000000000/deposit",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "terminal withdraw is not registered",
			method:     http.MethodPost,
			path:       "/terminal/accounts/00000000-0000-0000-0000-000000000000/withdraw",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("expected status %d, got %d", tc.wantStatus, rec.Code)
			}
		})
	}
}

func assertInvalidAppToken(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()

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

	if got.Error.Code != "INVALID_APP_TOKEN" {
		t.Fatalf("expected error code %q, got %q", "INVALID_APP_TOKEN", got.Error.Code)
	}
}
