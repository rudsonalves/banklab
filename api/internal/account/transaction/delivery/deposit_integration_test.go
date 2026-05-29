package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	transactionApplication "github.com/seu-usuario/bank-api/internal/account/transaction/application"
	"github.com/seu-usuario/bank-api/internal/account/transaction/domain"
	transactionInfrastructure "github.com/seu-usuario/bank-api/internal/account/transaction/infrastructure"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
)

func TestHandler_Deposit_Integration(t *testing.T) {
	ctx := context.Background()
	pool := newTestPool(t, ctx)
	defer pool.Close()

	ensureDepositTestSchema(t, ctx, pool)

	accountID := uuid.New()
	customerID := uuid.New()
	seedDepositTestData(t, ctx, pool, customerID, accountID, 100, domain.AccountActive)
	defer cleanupDepositTestData(t, ctx, pool, customerID, accountID)

	repo := transactionInfrastructure.New(pool)
	depositUC := transactionApplication.NewDeposit(repo)
	handler := New(depositUC, nil, nil, nil, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /terminal/accounts/{id}/deposit", func(w http.ResponseWriter, r *http.Request) {
		ctx := sharedauthctx.WithAuthenticatedUser(r.Context(), sharedauthctx.AuthenticatedUser{
			UserID:     uuid.New(),
			Role:       authdomain.RoleCustomer,
			CustomerID: &customerID,
		})
		handler.Deposit(w, r.WithContext(ctx))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	payload := bytes.NewBufferString(`{"amount": 50}`)
	resp, err := http.Post(server.URL+"/terminal/accounts/"+accountID.String()+"/deposit", "application/json", payload)
	if err != nil {
		t.Fatalf("failed to call deposit endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&body)
		t.Fatalf("expected status %d, got %d with body %+v", http.StatusOK, resp.StatusCode, body)
	}

	var got struct {
		Data struct {
			ID      string `json:"id"`
			Balance int64  `json:"balance"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Data.ID != accountID.String() {
		t.Fatalf("expected account id %q, got %q", accountID.String(), got.Data.ID)
	}

	if got.Data.Balance != 150 {
		t.Fatalf("expected balance %d in response, got %d", 150, got.Data.Balance)
	}

	balance := queryAccountBalance(t, ctx, pool, accountID)
	if balance != 150 {
		t.Fatalf("expected persisted balance %d, got %d", 150, balance)
	}
}

func TestHandler_Transfer_Integration(t *testing.T) {
	ctx := context.Background()
	pool := newTestPool(t, ctx)
	defer pool.Close()

	ensureDepositTestSchema(t, ctx, pool)

	sourceCustomerID := uuid.New()
	destinationCustomerID := uuid.New()
	sourceAccountID := uuid.New()
	destinationAccountID := uuid.New()
	sourceBranch := "0001"
	destinationBranch := "0001"
	sourceNumber := uniqueAccountNumber("11")
	destinationNumber := uniqueAccountNumber("22")

	seedTransferTestData(
		t,
		ctx,
		pool,
		sourceCustomerID,
		destinationCustomerID,
		sourceAccountID,
		destinationAccountID,
		sourceBranch,
		sourceNumber,
		destinationBranch,
		destinationNumber,
		10000,
		5000,
	)
	defer cleanupTransferTestData(t, ctx, pool, sourceCustomerID, destinationCustomerID, sourceAccountID, destinationAccountID)

	repo := transactionInfrastructure.New(pool)
	transferUC := transactionApplication.NewTransfer(repo)
	receiptUC := transactionApplication.NewGetTransferReceipt(repo)
	handler := New(nil, nil, transferUC, receiptUC, &enforceStepUpUseCaseMock{})

	mux := http.NewServeMux()
	mux.HandleFunc("POST /accounts/transfer", func(w http.ResponseWriter, r *http.Request) {
		authCtx := sharedauthctx.WithAuthenticatedUser(r.Context(), sharedauthctx.AuthenticatedUser{
			UserID:     uuid.New(),
			Role:       authdomain.RoleCustomer,
			CustomerID: &sourceCustomerID,
		})
		handler.Transfer(w, r.WithContext(authCtx))
	})
	mux.HandleFunc("GET /accounts/transfer/{transaction_reference}/receipt", func(w http.ResponseWriter, r *http.Request) {
		authCtx := sharedauthctx.WithAuthenticatedUser(r.Context(), sharedauthctx.AuthenticatedUser{
			UserID:     uuid.New(),
			Role:       authdomain.RoleCustomer,
			CustomerID: &sourceCustomerID,
		})
		handler.TransferReceipt(w, r.WithContext(authCtx))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	payload := fmt.Sprintf(`{
		"from_branch": %q,
		"from_account_number": %q,
		"to_branch": %q,
		"to_account_number": %q,
		"amount": 2500,
		"idempotency_key": "transfer-integration-key",
		"description": "Aluguel de maio"
	}`, sourceBranch, sourceNumber, destinationBranch, destinationNumber)

	first := postTransfer(t, server.URL, payload)
	if first.Data.FromBalance != 7500 || first.Data.ToBalance != 7500 {
		t.Fatalf("expected transfer balances 7500 and 7500, got %+v", first.Data)
	}
	if first.Data.TransactionReference == "" {
		t.Fatal("expected transaction_reference to be returned")
	}

	if sourceBalance := queryAccountBalance(t, ctx, pool, sourceAccountID); sourceBalance != 7500 {
		t.Fatalf("expected persisted source balance %d, got %d", 7500, sourceBalance)
	}
	if destinationBalance := queryAccountBalance(t, ctx, pool, destinationAccountID); destinationBalance != 7500 {
		t.Fatalf("expected persisted destination balance %d, got %d", 7500, destinationBalance)
	}

	retryPayload := fmt.Sprintf(`{
		"from_branch": %q,
		"from_account_number": %q,
		"to_branch": %q,
		"to_account_number": %q,
		"amount": 2500,
		"idempotency_key": "transfer-integration-key",
		"description": "Descricao alterada no retry"
	}`, sourceBranch, sourceNumber, destinationBranch, destinationNumber)

	second := postTransfer(t, server.URL, retryPayload)
	if second.Data.TransactionReference != first.Data.TransactionReference {
		t.Fatalf("expected idempotent replay reference %q, got %q", first.Data.TransactionReference, second.Data.TransactionReference)
	}
	if sourceBalance := queryAccountBalance(t, ctx, pool, sourceAccountID); sourceBalance != 7500 {
		t.Fatalf("expected replay to preserve source balance %d, got %d", 7500, sourceBalance)
	}
	if destinationBalance := queryAccountBalance(t, ctx, pool, destinationAccountID); destinationBalance != 7500 {
		t.Fatalf("expected replay to preserve destination balance %d, got %d", 7500, destinationBalance)
	}

	receiptResp, err := http.Get(server.URL + "/accounts/transfer/" + first.Data.TransactionReference + "/receipt")
	if err != nil {
		t.Fatalf("failed to call receipt endpoint: %v", err)
	}
	defer receiptResp.Body.Close()

	if receiptResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(receiptResp.Body)
		t.Fatalf("expected receipt status %d, got %d with body %s", http.StatusOK, receiptResp.StatusCode, string(body))
	}

	var receipt struct {
		Data TransferReceiptData `json:"data"`
	}
	if err := json.NewDecoder(receiptResp.Body).Decode(&receipt); err != nil {
		t.Fatalf("failed to decode receipt response: %v", err)
	}
	if receipt.Data.TransactionReference != first.Data.TransactionReference {
		t.Fatalf("expected receipt reference %q, got %q", first.Data.TransactionReference, receipt.Data.TransactionReference)
	}
	if receipt.Data.SourceAccountNumber != sourceNumber || receipt.Data.DestinationAccountNumber != destinationNumber {
		t.Fatalf("unexpected receipt account numbers: %+v", receipt.Data)
	}
	if receipt.Data.RecipientName != "Destination Transfer Test" {
		t.Fatalf("expected recipient name %q, got %q", "Destination Transfer Test", receipt.Data.RecipientName)
	}
	if receipt.Data.Description == nil || *receipt.Data.Description != "Aluguel de maio" {
		t.Fatalf("expected original receipt description %q, got %+v", "Aluguel de maio", receipt.Data.Description)
	}
}

func newTestPool(t *testing.T, ctx context.Context) *pgxpool.Pool {
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

func ensureDepositTestSchema(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	statements := []string{
		`CREATE EXTENSION IF NOT EXISTS pgcrypto`,
		`CREATE TABLE IF NOT EXISTS customers (
			id UUID PRIMARY KEY,
			name VARCHAR(120) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS customer_documents (
			id UUID PRIMARY KEY,
			customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
			type VARCHAR(30) NOT NULL,
			value VARCHAR(80) NOT NULL,
			issuer VARCHAR(80),
			issuer_state VARCHAR(30),
			country CHAR(2) NOT NULL DEFAULT 'BR',
			is_primary BOOLEAN NOT NULL DEFAULT false,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			CONSTRAINT customer_documents_unique_document UNIQUE (type, value, country)
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS customer_documents_one_primary_per_customer
			ON customer_documents(customer_id)
			WHERE is_primary = true`,
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
		`CREATE TABLE IF NOT EXISTS transactions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			account_id UUID NOT NULL REFERENCES accounts(id),
			type VARCHAR(20) NOT NULL,
			amount BIGINT NOT NULL,
			balance_after BIGINT NOT NULL,
			reference_id UUID,
			related_account_id UUID,
			idempotency_key VARCHAR(100),
			description TEXT,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)`,
		`CREATE SEQUENCE IF NOT EXISTS account_number_seq START WITH 10000000 INCREMENT BY 1`,
	}

	for _, statement := range statements {
		if _, err := pool.Exec(ctx, statement); err != nil {
			t.Fatalf("failed to ensure test schema: %v", err)
		}
	}

	if _, err := pool.Exec(ctx, `ALTER TABLE transactions ADD COLUMN IF NOT EXISTS related_account_id UUID`); err != nil {
		t.Fatalf("failed to ensure transactions.related_account_id column: %v", err)
	}

	if _, err := pool.Exec(ctx, `ALTER TABLE transactions ADD COLUMN IF NOT EXISTS idempotency_key VARCHAR(100)`); err != nil {
		t.Fatalf("failed to ensure transactions.idempotency_key column: %v", err)
	}

	if _, err := pool.Exec(ctx, `ALTER TABLE transactions ADD COLUMN IF NOT EXISTS description TEXT`); err != nil {
		t.Fatalf("failed to ensure transactions.description column: %v", err)
	}

	if _, err := pool.Exec(ctx, `ALTER TABLE transactions ALTER COLUMN id SET DEFAULT gen_random_uuid()`); err != nil {
		t.Fatalf("failed to ensure transactions.id default: %v", err)
	}

	if _, err := pool.Exec(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS ux_transactions_idempotency
		ON transactions(account_id, idempotency_key)
		WHERE idempotency_key IS NOT NULL`); err != nil {
		t.Fatalf("failed to ensure transactions idempotency index: %v", err)
	}
}

func seedDepositTestData(t *testing.T, ctx context.Context, pool *pgxpool.Pool, customerID, accountID uuid.UUID, balance int64, status domain.AccountStatus) {
	t.Helper()

	uniqueNumber := time.Now().UnixNano()
	cpfSuffix := fmt.Sprintf("%011d", uniqueNumber%100000000000)
	accountNumber := fmt.Sprintf("%08d", uniqueNumber%100000000)

	if _, err := pool.Exec(ctx, `
		INSERT INTO customers (id, name, created_at)
		VALUES ($1, $2, $3)
	`, customerID, "Deposit Test", time.Now().UTC()); err != nil {
		t.Fatalf("failed to insert customer: %v", err)
	}

	if _, err := pool.Exec(ctx, `
		INSERT INTO customer_documents (
			id,
			customer_id,
			type,
			value,
			country,
			is_primary,
			created_at,
			updated_at
		)
		VALUES ($1, $2, 'cpf', $3, 'BR', true, $4, $4)
	`, uuid.New(), customerID, cpfSuffix, time.Now().UTC()); err != nil {
		t.Fatalf("failed to insert customer document: %v", err)
	}

	if _, err := pool.Exec(ctx, `
		INSERT INTO accounts (id, customer_id, number, branch, balance, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, accountID, customerID, accountNumber, "0001", balance, status, time.Now().UTC()); err != nil {
		t.Fatalf("failed to insert account: %v", err)
	}
}

func seedTransferTestData(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	sourceCustomerID uuid.UUID,
	destinationCustomerID uuid.UUID,
	sourceAccountID uuid.UUID,
	destinationAccountID uuid.UUID,
	sourceBranch string,
	sourceNumber string,
	destinationBranch string,
	destinationNumber string,
	sourceBalance int64,
	destinationBalance int64,
) {
	t.Helper()

	insertCustomer(t, ctx, pool, sourceCustomerID, "Source Transfer Test", uniqueCPF("1"))
	insertCustomer(t, ctx, pool, destinationCustomerID, "Destination Transfer Test", uniqueCPF("2"))
	insertAccount(t, ctx, pool, sourceAccountID, sourceCustomerID, sourceNumber, sourceBranch, sourceBalance, domain.AccountActive)
	insertAccount(t, ctx, pool, destinationAccountID, destinationCustomerID, destinationNumber, destinationBranch, destinationBalance, domain.AccountActive)
}

func insertCustomer(t *testing.T, ctx context.Context, pool *pgxpool.Pool, customerID uuid.UUID, name, cpf string) {
	t.Helper()

	if _, err := pool.Exec(ctx, `
		INSERT INTO customers (id, name, created_at)
		VALUES ($1, $2, $3)
	`, customerID, name, time.Now().UTC()); err != nil {
		t.Fatalf("failed to insert customer %s: %v", name, err)
	}

	if _, err := pool.Exec(ctx, `
		INSERT INTO customer_documents (
			id,
			customer_id,
			type,
			value,
			country,
			is_primary,
			created_at,
			updated_at
		)
		VALUES ($1, $2, 'cpf', $3, 'BR', true, $4, $4)
	`, uuid.New(), customerID, cpf, time.Now().UTC()); err != nil {
		t.Fatalf("failed to insert customer document %s: %v", name, err)
	}
}

func insertAccount(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	accountID uuid.UUID,
	customerID uuid.UUID,
	number string,
	branch string,
	balance int64,
	status domain.AccountStatus,
) {
	t.Helper()

	if _, err := pool.Exec(ctx, `
		INSERT INTO accounts (id, customer_id, number, branch, balance, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, accountID, customerID, number, branch, balance, status, time.Now().UTC()); err != nil {
		t.Fatalf("failed to insert account %s/%s: %v", branch, number, err)
	}
}

func cleanupDepositTestData(t *testing.T, ctx context.Context, pool *pgxpool.Pool, customerID, accountID uuid.UUID) {
	t.Helper()

	if _, err := pool.Exec(ctx, `DELETE FROM transactions WHERE account_id = $1 OR reference_id = $1`, accountID); err != nil {
		t.Logf("cleanup warning: failed to delete transactions: %v", err)
	}

	if _, err := pool.Exec(ctx, `DELETE FROM accounts WHERE id = $1`, accountID); err != nil {
		t.Logf("cleanup warning: failed to delete account: %v", err)
	}

	if _, err := pool.Exec(ctx, `DELETE FROM customers WHERE id = $1`, customerID); err != nil {
		t.Logf("cleanup warning: failed to delete customer: %v", err)
	}
}

func cleanupTransferTestData(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	sourceCustomerID uuid.UUID,
	destinationCustomerID uuid.UUID,
	sourceAccountID uuid.UUID,
	destinationAccountID uuid.UUID,
) {
	t.Helper()

	if _, err := pool.Exec(ctx, `
		DELETE FROM transactions
		WHERE account_id IN ($1, $2)
		   OR related_account_id IN ($1, $2)
	`, sourceAccountID, destinationAccountID); err != nil {
		t.Logf("cleanup warning: failed to delete transfer transactions: %v", err)
	}

	if _, err := pool.Exec(ctx, `DELETE FROM accounts WHERE id IN ($1, $2)`, sourceAccountID, destinationAccountID); err != nil {
		t.Logf("cleanup warning: failed to delete transfer accounts: %v", err)
	}

	if _, err := pool.Exec(ctx, `DELETE FROM customers WHERE id IN ($1, $2)`, sourceCustomerID, destinationCustomerID); err != nil {
		t.Logf("cleanup warning: failed to delete transfer customers: %v", err)
	}
}

func queryAccountBalance(t *testing.T, ctx context.Context, pool *pgxpool.Pool, accountID uuid.UUID) int64 {
	t.Helper()

	var balance int64
	if err := pool.QueryRow(ctx, `SELECT balance FROM accounts WHERE id = $1`, accountID).Scan(&balance); err != nil {
		t.Fatalf("failed to query account balance: %v", err)
	}

	return balance
}

func postTransfer(t *testing.T, serverURL, payload string) struct {
	Data TransferData `json:"data"`
} {
	t.Helper()

	resp, err := http.Post(serverURL+"/accounts/transfer", "application/json", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatalf("failed to call transfer endpoint: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read transfer response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected transfer status %d, got %d with body %s", http.StatusOK, resp.StatusCode, string(body))
	}

	if strings.Contains(string(body), "from_account_id") || strings.Contains(string(body), "to_account_id") {
		t.Fatalf("transfer response must not expose internal account id fields: %s", string(body))
	}

	var got struct {
		Data TransferData `json:"data"`
	}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("failed to decode transfer response: %v", err)
	}

	return got
}

func uniqueCPF(prefix string) string {
	value := time.Now().UnixNano() % 10000000000
	return fmt.Sprintf("%s%010d", prefix, value)
}

func uniqueAccountNumber(prefix string) string {
	value := time.Now().UnixNano() % 1000000
	return fmt.Sprintf("%s%06d", prefix, value)
}
