package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	transactionapplication "github.com/seu-usuario/bank-api/internal/account/transaction/application"
	transactioninfrastructure "github.com/seu-usuario/bank-api/internal/account/transaction/infrastructure"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	authinfrastructure "github.com/seu-usuario/bank-api/internal/auth/infrastructure"
	securityapplication "github.com/seu-usuario/bank-api/internal/security/application"
	securitydelivery "github.com/seu-usuario/bank-api/internal/security/delivery"
	securitydomain "github.com/seu-usuario/bank-api/internal/security/domain"
	securityinfrastructure "github.com/seu-usuario/bank-api/internal/security/infrastructure"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
	"golang.org/x/crypto/bcrypt"
)

func TestIntegration_StepUpProtectedInternalTransfer(t *testing.T) {
	ctx := context.Background()
	pool := newTestPool(t, ctx)
	defer pool.Close()

	ensureStepUpTransferTestSchema(t, ctx, pool)

	userID := uuid.New()
	sourceCustomerID := uuid.New()
	destinationCustomerID := uuid.New()
	sourceAccountID := uuid.New()
	destinationAccountID := uuid.New()

	seedTransferTestData(
		t,
		ctx,
		pool,
		sourceCustomerID,
		destinationCustomerID,
		sourceAccountID,
		destinationAccountID,
		"0001",
		uniqueAccountNumber("31"),
		"0001",
		uniqueAccountNumber("32"),
		10000,
		5000,
	)
	seedStepUpTransferUser(t, ctx, pool, userID, sourceCustomerID)
	defer cleanupStepUpTransferTestData(
		t,
		ctx,
		pool,
		userID,
		sourceCustomerID,
		destinationCustomerID,
		sourceAccountID,
		destinationAccountID,
	)

	const (
		transactionPIN = "123456"
		idempotencyKey = "step-up-transfer-integration-key"
		jwtSecret      = "step-up-transfer-integration-secret"
		pepper         = "step-up-transfer-integration-pepper-32-bytes"
	)

	userRepo := authinfrastructure.NewPostgresUserRepository(pool)
	passwordRepo := securityinfrastructure.NewPostgresTransactionPasswordRepository(pool)
	tokenRepo := securityinfrastructure.NewPostgresStepUpTokenRepository(pool)
	passwordHasher := securityinfrastructure.NewBcryptTransactionPasswordHasher(
		bcrypt.MinCost,
		pepper,
	)
	tokenSigner := securityinfrastructure.NewJWTStepUpTokenSigner(jwtSecret)
	tokenVerifier := securityinfrastructure.NewJWTStepUpTokenVerifier(jwtSecret)

	createPasswordUC := securityapplication.NewCreateTransactionPasswordUseCase(
		passwordRepo,
		userRepo,
		passwordHasher,
	)
	authorizeStepUpUC := securityapplication.NewAuthorizeStepUpUseCase(
		passwordRepo,
		userRepo,
		passwordHasher,
		tokenRepo,
		tokenSigner,
		securitydomain.NewDefaultStepUpPublicOperationResolver(),
	)
	enforceStepUpUC := securityapplication.NewEnforceStepUpUseCase(
		tokenVerifier,
		tokenRepo,
	)

	transactionRepo := transactioninfrastructure.New(pool)
	transferUC := transactionapplication.NewTransfer(transactionRepo)
	transactionHandler := New(nil, nil, transferUC, nil, enforceStepUpUC)
	securityHandler := securitydelivery.New(createPasswordUC, authorizeStepUpUC)

	authenticatedUser := sharedauthctx.AuthenticatedUser{
		UserID:     userID,
		Role:       authdomain.RoleCustomer,
		CustomerID: &sourceCustomerID,
	}
	withAuthenticatedUser := func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authCtx := sharedauthctx.WithAuthenticatedUser(
				r.Context(),
				authenticatedUser,
			)
			handler(w, r.WithContext(authCtx))
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc(
		"POST /security/transaction-password",
		withAuthenticatedUser(securityHandler.CreateTransactionPassword),
	)
	mux.HandleFunc(
		"POST /security/step-up/authorize",
		withAuthenticatedUser(securityHandler.AuthorizeStepUp),
	)
	mux.HandleFunc(
		"POST /accounts/internal-transfers",
		withAuthenticatedUser(transactionHandler.Transfer),
	)
	server := httptest.NewServer(mux)
	defer server.Close()

	createPasswordResponse := performStepUpJSONRequest(
		t,
		http.MethodPost,
		server.URL+"/security/transaction-password",
		map[string]any{
			"transaction_password":              transactionPIN,
			"transaction_password_confirmation": transactionPIN,
		},
		"",
	)
	assertStepUpResponseStatus(
		t,
		createPasswordResponse,
		http.StatusCreated,
		"",
	)

	firstToken := authorizeInternalTransferStepUp(
		t,
		server.URL,
		transactionPIN,
	)
	transferPayload := map[string]any{
		"from_account_id": sourceAccountID.String(),
		"to_account_id":   destinationAccountID.String(),
		"amount":          2500,
		"idempotency_key": idempotencyKey,
		"description":     "Aluguel de maio",
	}

	firstTransferResponse := performStepUpJSONRequest(
		t,
		http.MethodPost,
		server.URL+"/accounts/internal-transfers",
		transferPayload,
		firstToken,
	)
	var firstTransfer struct {
		Data  TransferData    `json:"data"`
		Error *stepUpAPIError `json:"error"`
	}
	decodeStepUpResponse(
		t,
		firstTransferResponse,
		http.StatusOK,
		&firstTransfer,
	)
	if firstTransfer.Error != nil {
		t.Fatalf("expected first transfer error to be nil, got %+v", firstTransfer.Error)
	}
	if firstTransfer.Data.TransactionReference == "" {
		t.Fatal("expected first transfer to return transaction_reference")
	}

	consumedTokenResponse := performStepUpJSONRequest(
		t,
		http.MethodPost,
		server.URL+"/accounts/internal-transfers",
		transferPayload,
		firstToken,
	)
	assertStepUpResponseStatus(
		t,
		consumedTokenResponse,
		http.StatusUnauthorized,
		"STEP_UP_TOKEN_CONSUMED",
	)

	secondToken := authorizeInternalTransferStepUp(
		t,
		server.URL,
		transactionPIN,
	)
	if secondToken == firstToken {
		t.Fatal("expected a newly authorized step-up token")
	}

	replayResponse := performStepUpJSONRequest(
		t,
		http.MethodPost,
		server.URL+"/accounts/internal-transfers",
		transferPayload,
		secondToken,
	)
	var replay struct {
		Data  TransferData    `json:"data"`
		Error *stepUpAPIError `json:"error"`
	}
	decodeStepUpResponse(t, replayResponse, http.StatusOK, &replay)
	if replay.Error != nil {
		t.Fatalf("expected replay error to be nil, got %+v", replay.Error)
	}
	if replay.Data.TransactionReference != firstTransfer.Data.TransactionReference {
		t.Fatalf(
			"expected idempotent replay reference %q, got %q",
			firstTransfer.Data.TransactionReference,
			replay.Data.TransactionReference,
		)
	}

	if balance := queryAccountBalance(t, ctx, pool, sourceAccountID); balance != 7500 {
		t.Fatalf("expected source balance 7500 after replay, got %d", balance)
	}
	if balance := queryAccountBalance(t, ctx, pool, destinationAccountID); balance != 7500 {
		t.Fatalf("expected destination balance 7500 after replay, got %d", balance)
	}
}

type stepUpAPIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func authorizeInternalTransferStepUp(
	t *testing.T,
	serverURL string,
	transactionPIN string,
) string {
	t.Helper()

	response := performStepUpJSONRequest(
		t,
		http.MethodPost,
		serverURL+"/security/step-up/authorize",
		map[string]any{
			"method":               http.MethodPost,
			"path":                 "/accounts/internal-transfers",
			"transaction_password": transactionPIN,
		},
		"",
	)

	var body struct {
		Data struct {
			StepUpToken string `json:"step_up_token"`
			ExpiresIn   int    `json:"expires_in"`
		} `json:"data"`
		Error *stepUpAPIError `json:"error"`
	}
	decodeStepUpResponse(t, response, http.StatusOK, &body)

	if body.Error != nil {
		t.Fatalf("expected step-up authorization error to be nil, got %+v", body.Error)
	}
	if body.Data.StepUpToken == "" {
		t.Fatal("expected non-empty step_up_token")
	}
	if body.Data.ExpiresIn != 120 {
		t.Fatalf("expected expires_in 120, got %d", body.Data.ExpiresIn)
	}

	return body.Data.StepUpToken
}

func performStepUpJSONRequest(
	t *testing.T,
	method string,
	url string,
	payload any,
	stepUpToken string,
) *http.Response {
	t.Helper()

	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal request payload: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if stepUpToken != "" {
		req.Header.Set("X-Step-Up-Token", stepUpToken)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}

	return response
}

func assertStepUpResponseStatus(
	t *testing.T,
	response *http.Response,
	wantStatus int,
	wantErrorCode string,
) {
	t.Helper()

	var body struct {
		Error *stepUpAPIError `json:"error"`
	}
	decodeStepUpResponse(t, response, wantStatus, &body)

	if wantErrorCode == "" {
		if body.Error != nil {
			t.Fatalf("expected response error to be nil, got %+v", body.Error)
		}
		return
	}

	if body.Error == nil {
		t.Fatalf("expected error code %q, got nil error", wantErrorCode)
	}
	if body.Error.Code != wantErrorCode {
		t.Fatalf("expected error code %q, got %q", wantErrorCode, body.Error.Code)
	}
}

func decodeStepUpResponse(
	t *testing.T,
	response *http.Response,
	wantStatus int,
	target any,
) {
	t.Helper()
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	if response.StatusCode != wantStatus {
		t.Fatalf(
			"expected status %d, got %d with body %s",
			wantStatus,
			response.StatusCode,
			string(body),
		)
	}
	if err := json.Unmarshal(body, target); err != nil {
		t.Fatalf("failed to decode response body %s: %v", string(body), err)
	}
}

func ensureStepUpTransferTestSchema(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
) {
	t.Helper()

	ensureDepositTestSchema(t, ctx, pool)

	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(120) NOT NULL UNIQUE,
			phone VARCHAR(20) UNIQUE,
			password_hash TEXT NOT NULL,
			role VARCHAR(20) NOT NULL,
			customer_id UUID UNIQUE REFERENCES customers(id) ON DELETE SET NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			email_verified_at TIMESTAMP WITH TIME ZONE,
			phone_verified_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(20)`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS customer_id UUID REFERENCES customers(id) ON DELETE SET NULL`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'active'`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified_at TIMESTAMP WITH TIME ZONE`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS phone_verified_at TIMESTAMP WITH TIME ZONE`,
		`CREATE TABLE IF NOT EXISTS transaction_passwords (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			password_hash TEXT NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			failed_attempts INTEGER NOT NULL DEFAULT 0,
			locked_until TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			changed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			CONSTRAINT chk_transaction_passwords_status
				CHECK (status IN ('active', 'blocked')),
			CONSTRAINT chk_transaction_passwords_failed_attempts
				CHECK (failed_attempts >= 0),
			CONSTRAINT chk_transaction_passwords_blocked_locked_until
				CHECK (status <> 'blocked' OR locked_until IS NOT NULL)
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS ux_transaction_passwords_user_id
			ON transaction_passwords(user_id)`,
		`CREATE TABLE IF NOT EXISTS step_up_tokens (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			jti VARCHAR(120) NOT NULL,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			endpoint_key VARCHAR(120) NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			consumed_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			CONSTRAINT chk_step_up_tokens_status
				CHECK (status IN ('active', 'consumed')),
			CONSTRAINT chk_step_up_tokens_jti_not_blank
				CHECK (length(trim(jti)) > 0),
			CONSTRAINT chk_step_up_tokens_endpoint_key_not_blank
				CHECK (length(trim(endpoint_key)) > 0),
			CONSTRAINT chk_step_up_tokens_expires_after_created
				CHECK (expires_at > created_at),
			CONSTRAINT chk_step_up_tokens_consumed_at_consistency
				CHECK (
					(status = 'active' AND consumed_at IS NULL)
					OR
					(status = 'consumed' AND consumed_at IS NOT NULL)
				)
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS ux_step_up_tokens_jti
			ON step_up_tokens(jti)`,
	}

	for _, statement := range statements {
		if _, err := pool.Exec(ctx, statement); err != nil {
			t.Fatalf("failed to ensure step-up transfer test schema: %v", err)
		}
	}
}

func seedStepUpTransferUser(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	userID uuid.UUID,
	customerID uuid.UUID,
) {
	t.Helper()

	now := time.Now().UTC()
	email := fmt.Sprintf("step-up-transfer-%s@example.com", userID.String())
	if _, err := pool.Exec(ctx, `
		INSERT INTO users (
			id,
			email,
			password_hash,
			role,
			customer_id,
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
	`, userID, email, "unused-login-password-hash", "customer", customerID, "active", now); err != nil {
		t.Fatalf("failed to insert step-up transfer user: %v", err)
	}
}

func cleanupStepUpTransferTestData(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	userID uuid.UUID,
	sourceCustomerID uuid.UUID,
	destinationCustomerID uuid.UUID,
	sourceAccountID uuid.UUID,
	destinationAccountID uuid.UUID,
) {
	t.Helper()

	if _, err := pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID); err != nil {
		t.Logf("cleanup warning: failed to delete step-up transfer user: %v", err)
	}

	cleanupTransferTestData(
		t,
		ctx,
		pool,
		sourceCustomerID,
		destinationCustomerID,
		sourceAccountID,
		destinationAccountID,
	)
}
