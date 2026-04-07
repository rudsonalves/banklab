package application

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/domain"
)

func TestWithdraw_Execute_InvalidAmount(t *testing.T) {
	repo := &depositAccountRepositoryMock{}
	useCase := NewWithdraw(repo)

	account, err := useCase.Execute(context.Background(), WithdrawInput{
		AccountID: uuid.New(),
		Amount:    0,
	})

	if !errors.Is(err, domain.ErrInvalidAmount) {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidAmount, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if repo.beginTxCalls != 0 {
		t.Fatalf("expected BeginTx not to be called, got %d calls", repo.beginTxCalls)
	}
}

func TestWithdraw_Execute_InvalidAccountID(t *testing.T) {
	repo := &depositAccountRepositoryMock{}
	useCase := NewWithdraw(repo)

	account, err := useCase.Execute(context.Background(), WithdrawInput{
		AccountID: uuid.Nil,
		Amount:    10,
	})

	if !errors.Is(err, domain.ErrInvalidData) {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidData, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if repo.beginTxCalls != 0 {
		t.Fatalf("expected BeginTx not to be called, got %d calls", repo.beginTxCalls)
	}
}

func TestWithdraw_Execute_AccountNotFound(t *testing.T) {
	tx := &txMock{getByIDErr: domain.ErrAccountNotFound}
	repo := &depositAccountRepositoryMock{tx: tx}
	useCase := NewWithdraw(repo)

	account, err := useCase.Execute(context.Background(), WithdrawInput{
		AccountID: uuid.New(),
		Amount:    100,
	})

	if !errors.Is(err, domain.ErrAccountNotFound) {
		t.Fatalf("expected error %v, got %v", domain.ErrAccountNotFound, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if tx.rollbackCalls != 1 {
		t.Fatalf("expected Rollback to be called once, got %d calls", tx.rollbackCalls)
	}

	if tx.getByIDForUpdateCalls != 1 {
		t.Fatalf("expected GetByIDForUpdate to be called once, got %d calls", tx.getByIDForUpdateCalls)
	}

	if tx.decreaseBalanceCalls != 0 {
		t.Fatalf("expected DecreaseBalance not to be called, got %d calls", tx.decreaseBalanceCalls)
	}

	if tx.createTransactionCalls != 0 {
		t.Fatalf("expected CreateTransaction not to be called, got %d calls", tx.createTransactionCalls)
	}
}

func TestWithdraw_Execute_InsufficientBalance(t *testing.T) {
	customerID := uuid.New()
	tx := &txMock{
		account: &domain.Account{
			ID:         uuid.New(),
			CustomerID: customerID,
			Balance:    100,
			Status:     domain.AccountActive,
		},
	}
	repo := &depositAccountRepositoryMock{tx: tx}
	useCase := NewWithdraw(repo)

	account, err := useCase.Execute(context.Background(), WithdrawInput{
		User:      testCustomerUser(customerID),
		AccountID: uuid.New(),
		Amount:    150,
	})

	if !errors.Is(err, domain.ErrInsufficientBalance) {
		t.Fatalf("expected error %v, got %v", domain.ErrInsufficientBalance, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if tx.decreaseBalanceCalls != 0 {
		t.Fatalf("expected DecreaseBalance not to be called, got %d calls", tx.decreaseBalanceCalls)
	}

	if tx.rollbackCalls != 1 {
		t.Fatalf("expected Rollback to be called once, got %d calls", tx.rollbackCalls)
	}

	if tx.createTransactionCalls != 0 {
		t.Fatalf("expected CreateTransaction not to be called, got %d calls", tx.createTransactionCalls)
	}
}

func TestWithdraw_Execute_RepositoryFailure(t *testing.T) {
	expectedErr := errors.New("decrease failed")
	customerID := uuid.New()
	tx := &txMock{
		account: &domain.Account{
			ID:         uuid.New(),
			CustomerID: customerID,
			Balance:    200,
			Status:     domain.AccountActive,
		},
		decreaseBalanceErr: expectedErr,
	}
	repo := &depositAccountRepositoryMock{tx: tx}
	useCase := NewWithdraw(repo)

	account, err := useCase.Execute(context.Background(), WithdrawInput{
		User:      testCustomerUser(customerID),
		AccountID: uuid.New(),
		Amount:    10,
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if tx.rollbackCalls != 1 {
		t.Fatalf("expected Rollback to be called once, got %d calls", tx.rollbackCalls)
	}
}

func TestWithdraw_Execute_Success(t *testing.T) {
	initialBalance := int64(100)
	withdrawAmount := int64(50)
	accountID := uuid.New()
	customerID := uuid.New()

	tx := &txMock{
		account: &domain.Account{
			ID:         accountID,
			CustomerID: customerID,
			Balance:    initialBalance,
			Status:     domain.AccountActive,
		},
		decreaseBalanceValue: initialBalance - withdrawAmount,
	}
	repo := &depositAccountRepositoryMock{tx: tx}
	useCase := NewWithdraw(repo)

	account, err := useCase.Execute(context.Background(), WithdrawInput{
		User:      testCustomerUser(customerID),
		AccountID: accountID,
		Amount:    withdrawAmount,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if account == nil {
		t.Fatal("expected account to be non-nil")
	}

	if account.Balance != initialBalance-withdrawAmount {
		t.Fatalf("expected balance %d, got %d", initialBalance-withdrawAmount, account.Balance)
	}

	if tx.decreaseBalanceCalls != 1 {
		t.Fatalf("expected DecreaseBalance to be called once, got %d calls", tx.decreaseBalanceCalls)
	}

	if tx.getByIDForUpdateCalls != 1 {
		t.Fatalf("expected GetByIDForUpdate to be called once, got %d calls", tx.getByIDForUpdateCalls)
	}

	if tx.createTransactionCalls != 1 {
		t.Fatalf("expected CreateTransaction to be called once, got %d calls", tx.createTransactionCalls)
	}

	created := tx.createdTransactions[0]
	if created.Type != domain.TransactionWithdraw {
		t.Fatalf("expected ledger type %s, got %s", domain.TransactionWithdraw, created.Type)
	}

	if created.BalanceAfter != initialBalance-withdrawAmount {
		t.Fatalf("expected ledger balance_after %d, got %d", initialBalance-withdrawAmount, created.BalanceAfter)
	}

	if tx.commitCalls != 1 {
		t.Fatalf("expected Commit to be called once, got %d calls", tx.commitCalls)
	}

	if tx.rollbackCalls != 0 {
		t.Fatalf("expected Rollback not to be called, got %d calls", tx.rollbackCalls)
	}
}

func TestWithdraw_Execute_LedgerInsertFailure(t *testing.T) {
	expectedErr := errors.New("ledger insert failed")
	customerID := uuid.New()
	tx := &txMock{
		account: &domain.Account{
			ID:         uuid.New(),
			CustomerID: customerID,
			Balance:    200,
			Status:     domain.AccountActive,
		},
		decreaseBalanceValue: 190,
		createTransactionErr: expectedErr,
	}
	repo := &depositAccountRepositoryMock{tx: tx}
	useCase := NewWithdraw(repo)

	account, err := useCase.Execute(context.Background(), WithdrawInput{
		User:      testCustomerUser(customerID),
		AccountID: uuid.New(),
		Amount:    10,
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}

	if tx.createTransactionCalls != 1 {
		t.Fatalf("expected CreateTransaction to be called once, got %d calls", tx.createTransactionCalls)
	}

	if tx.rollbackCalls != 1 {
		t.Fatalf("expected Rollback to be called once, got %d calls", tx.rollbackCalls)
	}

	if tx.commitCalls != 0 {
		t.Fatalf("expected Commit not to be called, got %d calls", tx.commitCalls)
	}
}

func TestWithdraw_Execute_AdminAllowed(t *testing.T) {
	accountID := uuid.New()
	tx := &txMock{
		account: &domain.Account{
			ID:         accountID,
			CustomerID: uuid.New(),
			Balance:    100,
			Status:     domain.AccountActive,
		},
		decreaseBalanceValue: 90,
	}
	repo := &depositAccountRepositoryMock{tx: tx}
	useCase := NewWithdraw(repo)

	account, err := useCase.Execute(context.Background(), WithdrawInput{
		User:      testAdminUser(),
		AccountID: accountID,
		Amount:    10,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if account == nil {
		t.Fatal("expected account to be non-nil")
	}

	if tx.decreaseBalanceCalls != 1 {
		t.Fatalf("expected DecreaseBalance to be called once, got %d calls", tx.decreaseBalanceCalls)
	}
}

func TestWithdraw_Execute_ForbiddenForDifferentCustomer(t *testing.T) {
	accountID := uuid.New()
	accountCustomerID := uuid.New() // different customer owns this account
	tx := &txMock{
		account: &domain.Account{
			ID:         accountID,
			CustomerID: accountCustomerID,
			Balance:    100,
			Status:     domain.AccountActive,
		},
	}
	repo := &depositAccountRepositoryMock{tx: tx}
	useCase := NewWithdraw(repo)

	account, err := useCase.Execute(context.Background(), WithdrawInput{
		User:      testCustomerUser(uuid.New()), // different customer
		AccountID: accountID,
		Amount:    10,
	})

	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected error %v, got %v", domain.ErrForbidden, err)
	}

	if account != nil {
		t.Fatalf("expected account to be nil, got %+v", account)
	}
}
