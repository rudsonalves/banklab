package delivery_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	accountapplication "github.com/seu-usuario/bank-api/internal/account/application"
	accountdelivery "github.com/seu-usuario/bank-api/internal/account/delivery"
	accountdomain "github.com/seu-usuario/bank-api/internal/account/domain"
	accountinfrastructure "github.com/seu-usuario/bank-api/internal/account/infrastructure"
	authapplication "github.com/seu-usuario/bank-api/internal/auth/application"
	authdelivery "github.com/seu-usuario/bank-api/internal/auth/delivery"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	authinfrastructure "github.com/seu-usuario/bank-api/internal/auth/infrastructure"
	customerinfrastructure "github.com/seu-usuario/bank-api/internal/customer/infrastructure"
	"golang.org/x/crypto/bcrypt"
)

func TestIntegration_AuthAndAuthorizationFlows(t *testing.T) {
	ctx := context.Background()
	pool := newIntegrationPool(t, ctx)
	defer pool.Close()

	ensureIntegrationSchema(t, ctx, pool)

	server, cleanup := newIntegrationServer(t, pool)
	defer cleanup()

	t.Run("register then login then me", func(t *testing.T) {
		email := fmt.Sprintf("register-%s@example.com", uuid.NewString())
		password := "P@ssword123"

		registerResp := performJSONRequest(t, server.URL+"/auth/register", http.MethodPost, map[string]any{
			"email":    email,
			"password": password,
			"name":     "Integration User",
			"cpf":      fmt.Sprintf("%011d", time.Now().UnixNano()%100000000000),
		}, "")
		if registerResp.StatusCode != http.StatusCreated {
			t.Fatalf("expected register status %d, got %d", http.StatusCreated, registerResp.StatusCode)
		}

		var registerBody envelope
		decodeResponseBody(t, registerResp, &registerBody)
		if registerBody.Error != nil {
			t.Fatalf("expected register error to be nil, got %+v", registerBody.Error)
		}

		loginResp := performJSONRequest(t, server.URL+"/auth/login", http.MethodPost, map[string]any{
			"email":    email,
			"password": password,
		}, "")
		if loginResp.StatusCode != http.StatusOK {
			t.Fatalf("expected login status %d, got %d", http.StatusOK, loginResp.StatusCode)
		}

		var loginBody struct {
			Data struct {
				AccessToken  string `json:"access_token"`
				RefreshToken string `json:"refresh_token"`
				UserID       string `json:"user_id"`
				Email        string `json:"email"`
				Role         string `json:"role"`
			} `json:"data"`
			Error *apiError `json:"error"`
		}
		decodeResponseBody(t, loginResp, &loginBody)
		if loginBody.Error != nil {
			t.Fatalf("expected login error to be nil, got %+v", loginBody.Error)
		}
		if loginBody.Data.AccessToken == "" {
			t.Fatal("expected non-empty access token")
		}
		if loginBody.Data.RefreshToken == "" {
			t.Fatal("expected non-empty refresh token")
		}

		meResp := performJSONRequest(t, server.URL+"/auth/me", http.MethodGet, nil, loginBody.Data.AccessToken)
		if meResp.StatusCode != http.StatusOK {
			t.Fatalf("expected /auth/me status %d, got %d", http.StatusOK, meResp.StatusCode)
		}

		var meBody struct {
			Data struct {
				ID    string `json:"id"`
				Email string `json:"email"`
				Role  string `json:"role"`
			} `json:"data"`
			Error *apiError `json:"error"`
		}
		decodeResponseBody(t, meResp, &meBody)

		if meBody.Error != nil {
			t.Fatalf("expected /auth/me error to be nil, got %+v", meBody.Error)
		}
		if meBody.Data.Email != email {
			t.Fatalf("expected /auth/me email %q, got %q", email, meBody.Data.Email)
		}
		if meBody.Data.Role != string(authdomain.RoleCustomer) {
			t.Fatalf("expected /auth/me role %q, got %q", authdomain.RoleCustomer, meBody.Data.Role)
		}
	})

	t.Run("own account access succeeds", func(t *testing.T) {
		customerID := seedCustomer(t, ctx, pool)
		accountID := seedAccount(t, ctx, pool, customerID, 100)
		password := "OwnPass123!"
		email := seedUser(t, ctx, pool, password, authdomain.RoleCustomer, &customerID)

		token := loginAndGetToken(t, server.URL, email, password)
		resp := performJSONRequest(t, server.URL+"/accounts/"+accountID.String()+"/deposit", http.MethodPost, map[string]any{
			"amount": 25,
		}, token)

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected own-account deposit status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var body struct {
			Data struct {
				Balance int64 `json:"balance"`
			} `json:"data"`
			Error *apiError `json:"error"`
		}
		decodeResponseBody(t, resp, &body)
		if body.Error != nil {
			t.Fatalf("expected own-account deposit error to be nil, got %+v", body.Error)
		}
		if body.Data.Balance != 125 {
			t.Fatalf("expected balance %d, got %d", 125, body.Data.Balance)
		}
	})

	t.Run("another account access is forbidden", func(t *testing.T) {
		ownerCustomerID := seedCustomer(t, ctx, pool)
		otherCustomerID := seedCustomer(t, ctx, pool)
		otherAccountID := seedAccount(t, ctx, pool, otherCustomerID, 100)

		password := "NoAccess123!"
		email := seedUser(t, ctx, pool, password, authdomain.RoleCustomer, &ownerCustomerID)
		token := loginAndGetToken(t, server.URL, email, password)

		resp := performJSONRequest(t, server.URL+"/accounts/"+otherAccountID.String()+"/deposit", http.MethodPost, map[string]any{
			"amount": 20,
		}, token)

		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("expected forbidden status %d, got %d", http.StatusForbidden, resp.StatusCode)
		}

		var body envelope
		decodeResponseBody(t, resp, &body)
		if body.Error == nil {
			t.Fatal("expected error body for forbidden response")
		}
		if body.Error.Code != "FORBIDDEN" {
			t.Fatalf("expected error code %q, got %q", "FORBIDDEN", body.Error.Code)
		}
		if body.Error.Message != "Access denied" {
			t.Fatalf("expected error message %q, got %q", "Access denied", body.Error.Message)
		}
	})

	t.Run("admin access is allowed", func(t *testing.T) {
		customerID := seedCustomer(t, ctx, pool)
		accountID := seedAccount(t, ctx, pool, customerID, 200)

		password := "AdminPass123!"
		email := seedUser(t, ctx, pool, password, authdomain.RoleAdmin, nil)
		token := loginAndGetToken(t, server.URL, email, password)

		resp := performJSONRequest(t, server.URL+"/accounts/"+accountID.String()+"/deposit", http.MethodPost, map[string]any{
			"amount": 30,
		}, token)

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected admin deposit status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var body struct {
			Data struct {
				Balance int64 `json:"balance"`
			} `json:"data"`
			Error *apiError `json:"error"`
		}
		decodeResponseBody(t, resp, &body)
		if body.Error != nil {
			t.Fatalf("expected admin deposit error to be nil, got %+v", body.Error)
		}
		if body.Data.Balance != 230 {
			t.Fatalf("expected admin-updated balance %d, got %d", 230, body.Data.Balance)
		}
	})
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type envelope struct {
	Data  json.RawMessage `json:"data"`
	Error *apiError       `json:"error"`
}

func newIntegrationServer(t *testing.T, pool *pgxpool.Pool) (*httptest.Server, func()) {
	t.Helper()

	userRepo := authinfrastructure.NewPostgresUserRepository(pool)
	sessionRepo := authinfrastructure.NewPostgresSessionRepository(pool)
	transactor := authinfrastructure.NewPostgresTransactor(pool)
	customerRepo := customerinfrastructure.New(pool)
	hasher := authinfrastructure.NewBcryptPasswordHasher(bcrypt.MinCost)
	tokenService := authinfrastructure.NewJWTTokenService("integration-secret", 20*time.Minute)

	registerUserUC := authapplication.NewRegisterUserUseCase(userRepo, customerRepo, hasher, transactor)
	loginUserUC := authapplication.NewLoginUserUseCase(userRepo, hasher, tokenService, sessionRepo)
	refreshAccessTokenUC := authapplication.NewRefreshAccessTokenUseCase(userRepo, tokenService, sessionRepo, transactor)
	getCurrentUserUC := authapplication.NewGetCurrentUserUseCase(userRepo)
	accountRepo := accountinfrastructure.New(pool)
	approveUserUC := authapplication.NewApproveUserUseCase(userRepo, accountRepo, customerRepo, transactor)
	authHandler := authdelivery.New(registerUserUC, loginUserUC, getCurrentUserUC, refreshAccessTokenUC, approveUserUC)
	authMiddleware := authdelivery.NewJWTMiddleware(tokenService)

	depositUC := accountapplication.NewDeposit(accountRepo)
	accountHandler := accountdelivery.New(nil, depositUC, nil, nil, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/refresh", authHandler.Refresh)
	mux.Handle("GET /auth/me", authMiddleware.RequireAuth(http.HandlerFunc(authHandler.Me)))
	mux.Handle("POST /accounts/{id}/deposit", authMiddleware.RequireAuth(http.HandlerFunc(accountHandler.Deposit)))

	server := httptest.NewServer(mux)

	return server, func() {
		server.Close()
	}
}

func newIntegrationPool(t *testing.T, ctx context.Context) *pgxpool.Pool {
	t.Helper()

	connString := os.Getenv("BANK_TEST_DATABASE_URL")
	if connString == "" {
		connString = "postgres://postgres:postgres@localhost:5432/bank_test?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		t.Skipf("skipping integration test: cannot create pool: %v", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		t.Skipf("skipping integration test: database unavailable: %v", err)
	}

	return pool
}

func ensureIntegrationSchema(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	statements := []string{
		`CREATE TABLE IF NOT EXISTS customers (
			id UUID PRIMARY KEY,
			name VARCHAR(120) NOT NULL,
			cpf VARCHAR(11) NOT NULL UNIQUE,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			CONSTRAINT chk_cpf_format CHECK (cpf ~ '^\d{11}$')
		)`,
		// Repair constraint if a previous test run created it with a broken regex.
		`DO $$ BEGIN
			ALTER TABLE customers DROP CONSTRAINT IF EXISTS chk_cpf_format;
			ALTER TABLE customers ADD CONSTRAINT chk_cpf_format CHECK (cpf ~ '^\d{11}$');
		EXCEPTION WHEN duplicate_object THEN NULL;
		END $$`,
		`CREATE TABLE IF NOT EXISTS accounts (
			id UUID PRIMARY KEY,
			customer_id UUID NOT NULL REFERENCES customers(id),
			number VARCHAR(20) NOT NULL UNIQUE,
			branch VARCHAR(10) NOT NULL,
			balance BIGINT NOT NULL DEFAULT 0,
			status VARCHAR(20) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			CONSTRAINT chk_account_status CHECK (status IN ('active', 'inactive', 'blocked'))
		)`,
		`CREATE SEQUENCE IF NOT EXISTS account_number_seq START WITH 10000000 INCREMENT BY 1`,
		`CREATE TABLE IF NOT EXISTS account_transactions (
			id UUID PRIMARY KEY,
			account_id UUID NOT NULL REFERENCES accounts(id),
			type VARCHAR(20) NOT NULL,
			amount BIGINT NOT NULL,
			balance_after BIGINT NOT NULL,
			reference_id UUID,
			related_account_id UUID,
			idempotency_key VARCHAR(100),
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			email VARCHAR(120) NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			role VARCHAR(20) NOT NULL,
			customer_id UUID UNIQUE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			CONSTRAINT fk_users_customer_id FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL
		)`,
		`CREATE TABLE IF NOT EXISTS user_sessions (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			token_hash CHAR(64) NOT NULL UNIQUE,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			revoked_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)`,
	}

	for _, statement := range statements {
		if _, err := pool.Exec(ctx, statement); err != nil {
			t.Fatalf("failed to ensure integration schema: %v", err)
		}
	}

	// Keep compatibility with pre-existing test databases created before ledger consolidation.
	if _, err := pool.Exec(ctx, `ALTER TABLE account_transactions ADD COLUMN IF NOT EXISTS related_account_id UUID`); err != nil {
		t.Fatalf("failed to ensure account_transactions.related_account_id column: %v", err)
	}

	if _, err := pool.Exec(ctx, `ALTER TABLE account_transactions ADD COLUMN IF NOT EXISTS idempotency_key VARCHAR(100)`); err != nil {
		t.Fatalf("failed to ensure account_transactions.idempotency_key column: %v", err)
	}

	if _, err := pool.Exec(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS ux_account_transactions_idempotency
		ON account_transactions(account_id, idempotency_key)
		WHERE idempotency_key IS NOT NULL`); err != nil {
		t.Fatalf("failed to ensure account_transactions idempotency index: %v", err)
	}

	if _, err := pool.Exec(ctx, `ALTER TABLE users ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'pending'`); err != nil {
		t.Fatalf("failed to ensure users.status column: %v", err)
	}
}

func performJSONRequest(t *testing.T, url, method string, payload any, token string) *http.Response {
	t.Helper()

	var bodyReader *bytes.Reader
	if payload == nil {
		bodyReader = bytes.NewReader(nil)
	} else {
		raw, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}
		bodyReader = bytes.NewReader(raw)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}

	t.Cleanup(func() {
		_ = resp.Body.Close()
	})

	return resp
}

func decodeResponseBody(t *testing.T, resp *http.Response, target any) {
	t.Helper()

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
}

func loginAndGetToken(t *testing.T, baseURL, email, password string) string {
	t.Helper()

	resp := performJSONRequest(t, baseURL+"/auth/login", http.MethodPost, map[string]any{
		"email":    email,
		"password": password,
	}, "")

	if resp.StatusCode != http.StatusOK {
		var body envelope
		decodeResponseBody(t, resp, &body)
		t.Fatalf("expected login status %d, got %d with error %+v", http.StatusOK, resp.StatusCode, body.Error)
	}

	var body struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
		Error *apiError `json:"error"`
	}
	decodeResponseBody(t, resp, &body)

	if body.Error != nil {
		t.Fatalf("expected login error nil, got %+v", body.Error)
	}
	if body.Data.AccessToken == "" {
		t.Fatal("expected non-empty access token")
	}

	return body.Data.AccessToken
}

func seedCustomer(t *testing.T, ctx context.Context, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()

	id := uuid.New()
	unique := time.Now().UnixNano()
	cpf := fmt.Sprintf("%011d", unique%100000000000)

	if _, err := pool.Exec(ctx, `
		INSERT INTO customers (id, name, cpf, created_at)
		VALUES ($1, $2, $3, $4)
	`, id, "Integration Customer", cpf, time.Now().UTC()); err != nil {
		t.Fatalf("failed to seed customer: %v", err)
	}

	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM customers WHERE id = $1`, id)
	})

	return id
}

func seedAccount(t *testing.T, ctx context.Context, pool *pgxpool.Pool, customerID uuid.UUID, balance int64) uuid.UUID {
	t.Helper()

	id := uuid.New()
	number := fmt.Sprintf("%08d", time.Now().UnixNano()%100000000)

	if _, err := pool.Exec(ctx, `
		INSERT INTO accounts (id, customer_id, number, branch, balance, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, id, customerID, number, "0001", balance, accountdomain.AccountActive, time.Now().UTC()); err != nil {
		t.Fatalf("failed to seed account: %v", err)
	}

	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM account_transactions WHERE account_id = $1`, id)
		_, _ = pool.Exec(ctx, `DELETE FROM accounts WHERE id = $1`, id)
	})

	return id
}

func seedUser(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	password string,
	role authdomain.Role,
	customerID *uuid.UUID,
) string {
	t.Helper()

	id := uuid.NewString()
	email := fmt.Sprintf("user-%s@example.com", id)
	hasher := authinfrastructure.NewBcryptPasswordHasher(bcrypt.MinCost)
	hash, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("failed to hash password for seeded user: %v", err)
	}

	var nullableCustomerID any
	if customerID != nil {
		nullableCustomerID = customerID.String()
	}

	now := time.Now().UTC()
	if _, err := pool.Exec(ctx, `
		INSERT INTO users (id, email, password_hash, role, customer_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, id, email, hash, string(role), nullableCustomerID, now, now); err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	})

	return email
}
