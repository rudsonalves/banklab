package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	accountDelivery "github.com/seu-usuario/bank-api/internal/account/bankaccount/delivery"
	statementDelivery "github.com/seu-usuario/bank-api/internal/account/statement/delivery"
	transactionDelivery "github.com/seu-usuario/bank-api/internal/account/transaction/delivery"
	adminDelivery "github.com/seu-usuario/bank-api/internal/admin/delivery"
	customerDelivery "github.com/seu-usuario/bank-api/internal/customer/delivery"
)

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
