package delivery

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
	accountApplication "github.com/seu-usuario/bank-api/internal/account/application"
	"github.com/seu-usuario/bank-api/internal/account/domain"
	accountInfrastructure "github.com/seu-usuario/bank-api/internal/account/infrastructure"
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

	repo := accountInfrastructure.New(pool)
	depositUC := accountApplication.NewDeposit(repo)
	handler := New(nil, depositUC)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /accounts/{id}/deposit", handler.Deposit)
	server := httptest.NewServer(mux)
	defer server.Close()

	payload := bytes.NewBufferString(`{"amount": 50}`)
	resp, err := http.Post(server.URL+"/accounts/"+accountID.String()+"/deposit", "application/json", payload)
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

func newTestPool(t *testing.T, ctx context.Context) *pgxpool.Pool {
	t.Helper()

	connString := os.Getenv("BANK_TEST_DATABASE_URL")
	if connString == "" {
		connString = "postgres://postgres:postgres@localhost:5432/bank?sslmode=disable"
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
		`CREATE TABLE IF NOT EXISTS customers (
			id UUID PRIMARY KEY,
			name VARCHAR(120) NOT NULL,
			cpf VARCHAR(11) NOT NULL UNIQUE,
			email VARCHAR(120) NOT NULL UNIQUE,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			CONSTRAINT chk_cpf_format CHECK (cpf ~ '^\d{11}$')
		)`,
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
	}

	for _, statement := range statements {
		if _, err := pool.Exec(ctx, statement); err != nil {
			t.Fatalf("failed to ensure test schema: %v", err)
		}
	}
}

func seedDepositTestData(t *testing.T, ctx context.Context, pool *pgxpool.Pool, customerID, accountID uuid.UUID, balance int64, status domain.AccountStatus) {
	t.Helper()

	uniqueNumber := time.Now().UnixNano()
	cpfSuffix := fmt.Sprintf("%011d", uniqueNumber%100000000000)
	accountNumber := fmt.Sprintf("%08d", uniqueNumber%100000000)
	email := fmt.Sprintf("deposit-%s@example.com", customerID.String())

	if _, err := pool.Exec(ctx, `
		INSERT INTO customers (id, name, cpf, email, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, customerID, "Deposit Test", cpfSuffix, email, time.Now().UTC()); err != nil {
		t.Fatalf("failed to insert customer: %v", err)
	}

	if _, err := pool.Exec(ctx, `
		INSERT INTO accounts (id, customer_id, number, branch, balance, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, accountID, customerID, accountNumber, "0001", balance, status, time.Now().UTC()); err != nil {
		t.Fatalf("failed to insert account: %v", err)
	}
}

func cleanupDepositTestData(t *testing.T, ctx context.Context, pool *pgxpool.Pool, customerID, accountID uuid.UUID) {
	t.Helper()

	if _, err := pool.Exec(ctx, `DELETE FROM accounts WHERE id = $1`, accountID); err != nil {
		t.Fatalf("failed to delete account: %v", err)
	}

	if _, err := pool.Exec(ctx, `DELETE FROM customers WHERE id = $1`, customerID); err != nil {
		t.Fatalf("failed to delete customer: %v", err)
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
