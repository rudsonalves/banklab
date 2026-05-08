package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/transaction/domain"
)

type transferAccountRepositoryMock struct {
	beginTxCalls int
	beginTxErr   error
	tx           domain.Tx
}

func (m *transferAccountRepositoryMock) Create(ctx context.Context, account *domain.Account) error {
	return nil
}

func (m *transferAccountRepositoryMock) ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]domain.Account, error) {
	return nil, nil
}

func (m *transferAccountRepositoryMock) CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
	return nil
}

func (m *transferAccountRepositoryMock) GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*domain.Transaction, error) {
	return nil, nil
}

func (m *transferAccountRepositoryMock) GetTransactionByReference(ctx context.Context, accountID uuid.UUID, referenceID uuid.UUID, typeName domain.TransactionType) (*domain.Transaction, error) {
	return nil, nil
}

func (m *transferAccountRepositoryMock) GetTransferReceiptByReference(ctx context.Context, referenceID uuid.UUID) (*domain.TransferReceipt, error) {
	return nil, nil
}

func (m *transferAccountRepositoryMock) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
	return false, nil
}

func (m *transferAccountRepositoryMock) NextAccountNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (m *transferAccountRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return nil, nil
}

func (m *transferAccountRepositoryMock) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return nil, nil
}

func (m *transferAccountRepositoryMock) GetByBranchAndNumber(ctx context.Context, branch, number string) (*domain.Account, error) {
	return nil, nil
}

func (m *transferAccountRepositoryMock) GetTransactions(
	ctx context.Context,
	accountID uuid.UUID,
	limit int,
	cursorTime *time.Time,
	cursorID *uuid.UUID,
	from *time.Time,
	to *time.Time,
) ([]domain.Transaction, error) {
	return nil, nil
}

func (m *transferAccountRepositoryMock) IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return 0, nil
}

func (m *transferAccountRepositoryMock) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return 0, nil
}

func (m *transferAccountRepositoryMock) BeginTx(ctx context.Context) (domain.Tx, error) {
	m.beginTxCalls++
	if m.beginTxErr != nil {
		return nil, m.beginTxErr
	}
	return m.tx, nil
}

func (m *transferAccountRepositoryMock) WithTransaction(ctx context.Context, fn func(tx domain.Tx) error) error {
	m.beginTxCalls++
	if m.beginTxErr != nil {
		return m.beginTxErr
	}

	if err := fn(m.tx); err != nil {
		_ = m.tx.Rollback(ctx)
		return err
	}

	return m.tx.Commit(ctx)
}

type transferTxMock struct {
	lockedOrder              []uuid.UUID
	accounts                 map[uuid.UUID]*domain.Account
	accountsByBranchNumber   map[string]*domain.Account
	getForUpdateErrs         map[uuid.UUID]error
	decreaseBalanceValue     int64
	decreaseBalanceErr       error
	updateBalanceValues      map[uuid.UUID]int64
	updateBalanceErr         error
	createTransactionErr     error
	getTransactionByKeyFn    func(accountID uuid.UUID, key string) (*domain.Transaction, error)
	getTransactionResult     *domain.Transaction
	getTransactionResults    []*domain.Transaction
	getTransactionErr        error
	getTransactionByRef      *domain.Transaction
	getTransactionByKeyCalls int
	commitErr                error
	rollbackErr              error
	decreaseCalls            int
	updateCalls              int
	createTransactionCalls   int
	commitCalls              int
	rollbackCalls            int
	createdTransactions      []*domain.Transaction
}

func branchNumberKey(branch, number string) string {
	return branch + ":" + number
}

func transferInput(userID uuid.UUID, fromID, toID uuid.UUID, amount int64) TransferInput {
	return TransferInput{
		User:           testCustomerUser(userID),
		FromAccountID:  fromID,
		ToAccountID:    toID,
		Amount:         amount,
		IdempotencyKey: "test-idempotency-key",
	}
}

func (m *transferTxMock) Create(ctx context.Context, account *domain.Account) error {
	return nil
}

func (m *transferTxMock) ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]domain.Account, error) {
	return nil, nil
}

func (m *transferTxMock) CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
	m.createTransactionCalls++
	m.createdTransactions = append(m.createdTransactions, tx)
	return m.createTransactionErr
}

func (m *transferTxMock) GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*domain.Transaction, error) {
	m.getTransactionByKeyCalls++
	if m.getTransactionByKeyFn != nil {
		return m.getTransactionByKeyFn(accountID, key)
	}

	if m.getTransactionErr != nil {
		return nil, m.getTransactionErr
	}

	if len(m.getTransactionResults) > 0 {
		result := m.getTransactionResults[0]
		m.getTransactionResults = m.getTransactionResults[1:]
		return result, nil
	}

	return m.getTransactionResult, nil
}

func (m *transferTxMock) GetTransactionByReference(ctx context.Context, accountID uuid.UUID, referenceID uuid.UUID, typeName domain.TransactionType) (*domain.Transaction, error) {
	if m.getTransactionByRef != nil && m.getTransactionByRef.AccountID == accountID && m.getTransactionByRef.ReferenceID != nil && *m.getTransactionByRef.ReferenceID == referenceID && m.getTransactionByRef.Type == typeName {
		return m.getTransactionByRef, nil
	}

	for _, t := range m.createdTransactions {
		if t.AccountID == accountID && t.ReferenceID != nil && *t.ReferenceID == referenceID && t.Type == typeName {
			return t, nil
		}
	}

	if m.getTransactionResult != nil && m.getTransactionResult.AccountID == accountID && m.getTransactionResult.ReferenceID != nil && *m.getTransactionResult.ReferenceID == referenceID && m.getTransactionResult.Type == typeName {
		return m.getTransactionResult, nil
	}

	for _, t := range m.getTransactionResults {
		if t != nil && t.AccountID == accountID && t.ReferenceID != nil && *t.ReferenceID == referenceID && t.Type == typeName {
			return t, nil
		}
	}

	return nil, nil
}

func (m *transferTxMock) GetTransferReceiptByReference(ctx context.Context, referenceID uuid.UUID) (*domain.TransferReceipt, error) {
	return nil, nil
}

func (m *transferTxMock) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
	return false, nil
}

func (m *transferTxMock) NextAccountNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (m *transferTxMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	if err := m.getForUpdateErrs[id]; err != nil {
		return nil, err
	}
	account, ok := m.accounts[id]
	if !ok {
		return nil, domain.ErrAccountNotFound
	}
	return account, nil
}

func (m *transferTxMock) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	m.lockedOrder = append(m.lockedOrder, id)
	if err := m.getForUpdateErrs[id]; err != nil {
		return nil, err
	}
	account, ok := m.accounts[id]
	if !ok {
		return nil, domain.ErrAccountNotFound
	}
	return account, nil
}

func (m *transferTxMock) GetByBranchAndNumber(ctx context.Context, branch, number string) (*domain.Account, error) {
	if m.accountsByBranchNumber == nil {
		return nil, domain.ErrAccountNotFound
	}

	account, ok := m.accountsByBranchNumber[branchNumberKey(branch, number)]
	if !ok {
		return nil, domain.ErrAccountNotFound
	}

	return account, nil
}

func (m *transferTxMock) GetTransactions(
	ctx context.Context,
	accountID uuid.UUID,
	limit int,
	cursorTime *time.Time,
	cursorID *uuid.UUID,
	from *time.Time,
	to *time.Time,
) ([]domain.Transaction, error) {
	return nil, nil
}

func (m *transferTxMock) IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	m.updateCalls++
	if m.updateBalanceErr != nil {
		return 0, m.updateBalanceErr
	}
	if balance, ok := m.updateBalanceValues[id]; ok {
		return balance, nil
	}
	return 0, nil
}

func (m *transferTxMock) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	m.decreaseCalls++
	return m.decreaseBalanceValue, m.decreaseBalanceErr
}

func (m *transferTxMock) BeginTx(ctx context.Context) (domain.Tx, error) {
	return nil, nil
}

func (m *transferTxMock) WithTransaction(ctx context.Context, fn func(tx domain.Tx) error) error {
	return errors.New("nested transactions are not supported")
}

func (m *transferTxMock) Commit(ctx context.Context) error {
	m.commitCalls++
	return m.commitErr
}

func (m *transferTxMock) Rollback(ctx context.Context) error {
	m.rollbackCalls++
	return m.rollbackErr
}

func TestTransfer_Execute_InvalidSourceAccountData(t *testing.T) {
	repo := &transferAccountRepositoryMock{}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{
		FromAccountID: uuid.Nil,
		ToAccountID:   uuid.New(),
		Amount:        10,
	})

	if !errors.Is(err, domain.ErrInvalidData) {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidData, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
	}
}

func TestTransfer_Execute_InvalidDestinationAccountData(t *testing.T) {
	repo := &transferAccountRepositoryMock{}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{
		FromAccountID: uuid.New(),
		ToAccountID:   uuid.Nil,
		Amount:        10,
	})

	if !errors.Is(err, domain.ErrInvalidData) {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidData, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
	}
}

func TestTransfer_Execute_MissingIdempotencyKey(t *testing.T) {
	repo := &transferAccountRepositoryMock{}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{
		FromAccountID:  uuid.New(),
		ToAccountID:    uuid.New(),
		Amount:         10,
		IdempotencyKey: "   ",
	})

	if !errors.Is(err, domain.ErrInvalidData) {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidData, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
	}
}

func TestTransfer_Execute_SameAccount(t *testing.T) {
	accountID := uuid.New()
	customerID := uuid.New()
	repo := &transferAccountRepositoryMock{}
	useCase := NewTransfer(repo)
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			accountID: {ID: accountID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
		},
	}
	repo.tx = tx

	result, err := useCase.Execute(context.Background(), TransferInput{
		User:           testCustomerUser(customerID),
		FromAccountID:  accountID,
		ToAccountID:    accountID,
		Amount:         10,
		IdempotencyKey: "same-account-key",
	})

	if !errors.Is(err, domain.ErrSameAccountTransfer) {
		t.Fatalf("expected error %v, got %v", domain.ErrSameAccountTransfer, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
	}

	if tx.getTransactionByKeyCalls != 0 {
		t.Fatalf("expected idempotency lookup not to be called, got %d calls", tx.getTransactionByKeyCalls)
	}
}

func TestTransfer_Execute_SameAccountBeforeIdempotencyReplay(t *testing.T) {
	accountID := uuid.New()
	customerID := uuid.New()
	key := "same-account-key"
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			accountID: {ID: accountID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
		},
		getTransactionResult: &domain.Transaction{
			ID:             uuid.New(),
			AccountID:      accountID,
			Type:           domain.TransactionTransferOut,
			Amount:         10,
			BalanceAfter:   90,
			IdempotencyKey: &key,
		},
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{
		User:           testCustomerUser(customerID),
		FromAccountID:  accountID,
		ToAccountID:    accountID,
		Amount:         10,
		IdempotencyKey: key,
	})

	if !errors.Is(err, domain.ErrSameAccountTransfer) {
		t.Fatalf("expected error %v, got %v", domain.ErrSameAccountTransfer, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
	}

	if tx.getTransactionByKeyCalls != 0 {
		t.Fatalf("expected idempotency lookup not to be called, got %d calls", tx.getTransactionByKeyCalls)
	}
}

func TestTransfer_Execute_SourceAccountNotFound(t *testing.T) {
	fromID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	toID := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			toID: {ID: toID, Status: domain.AccountActive, Balance: 100},
		},
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{
		FromAccountID:  fromID,
		ToAccountID:    toID,
		Amount:         10,
		IdempotencyKey: "source-not-found-key",
	})

	if !errors.Is(err, domain.ErrAccountNotFound) {
		t.Fatalf("expected error %v, got %v", domain.ErrAccountNotFound, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
	}

	if tx.rollbackCalls != 1 {
		t.Fatalf("expected rollback once, got %d", tx.rollbackCalls)
	}
}

func TestTransfer_Execute_DestinationAccountNotFound(t *testing.T) {
	fromID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	toID := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, Status: domain.AccountActive, Balance: 100},
		},
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{
		FromAccountID:  fromID,
		ToAccountID:    toID,
		Amount:         10,
		IdempotencyKey: "destination-not-found-key",
	})

	if !errors.Is(err, domain.ErrAccountNotFound) {
		t.Fatalf("expected error %v, got %v", domain.ErrAccountNotFound, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
	}

	if tx.rollbackCalls != 1 {
		t.Fatalf("expected rollback once, got %d", tx.rollbackCalls)
	}
}

func TestTransfer_Execute_InsufficientBalance(t *testing.T) {
	fromID := uuid.New()
	toID := uuid.New()
	customerID := uuid.New()
	fromBranch := "0001"
	fromNumber := "111111"
	toBranch := "0001"
	toNumber := "222222"
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 10},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		accountsByBranchNumber: map[string]*domain.Account{
			branchNumberKey(fromBranch, fromNumber): {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 10},
			branchNumberKey(toBranch, toNumber):     {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), transferInput(customerID, fromID, toID, 50))

	if !errors.Is(err, domain.ErrInsufficientBalance) {
		t.Fatalf("expected error %v, got %v", domain.ErrInsufficientBalance, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
	}

	if tx.decreaseCalls != 0 {
		t.Fatalf("expected no debit call, got %d", tx.decreaseCalls)
	}

	if tx.rollbackCalls != 1 {
		t.Fatalf("expected rollback once, got %d", tx.rollbackCalls)
	}
}

func TestTransfer_Execute_DestinationInactive(t *testing.T) {
	fromID := uuid.New()
	toID := uuid.New()
	customerID := uuid.New()
	fromBranch := "0001"
	fromNumber := "111111"
	toBranch := "0001"
	toNumber := "222222"
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountInactive, Balance: 20},
		},
		accountsByBranchNumber: map[string]*domain.Account{
			branchNumberKey(fromBranch, fromNumber): {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			branchNumberKey(toBranch, toNumber):     {ID: toID, Status: domain.AccountInactive, Balance: 20},
		},
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), transferInput(customerID, fromID, toID, 50))

	if !errors.Is(err, domain.ErrAccountInactive) {
		t.Fatalf("expected error %v, got %v", domain.ErrAccountInactive, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
	}
}

func TestTransfer_Execute_SourceInactive(t *testing.T) {
	fromID := uuid.New()
	toID := uuid.New()
	customerID := uuid.New()
	fromBranch := "0001"
	fromNumber := "111111"
	toBranch := "0001"
	toNumber := "222222"
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: customerID, Status: domain.AccountInactive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		accountsByBranchNumber: map[string]*domain.Account{
			branchNumberKey(fromBranch, fromNumber): {ID: fromID, CustomerID: customerID, Status: domain.AccountInactive, Balance: 100},
			branchNumberKey(toBranch, toNumber):     {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), transferInput(customerID, fromID, toID, 50))

	if !errors.Is(err, domain.ErrAccountInactive) {
		t.Fatalf("expected error %v, got %v", domain.ErrAccountInactive, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
	}

	if tx.decreaseCalls != 0 {
		t.Fatalf("expected no debit call, got %d", tx.decreaseCalls)
	}
}

func TestTransfer_Execute_Success(t *testing.T) {
	fromID := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	toID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	customerID := uuid.New()
	fromBranch := "0001"
	fromNumber := "111111"
	toBranch := "0001"
	toNumber := "222222"
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		accountsByBranchNumber: map[string]*domain.Account{
			branchNumberKey(fromBranch, fromNumber): {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			branchNumberKey(toBranch, toNumber):     {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		decreaseBalanceValue: 50,
		updateBalanceValues:  map[uuid.UUID]int64{toID: 70},
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), transferInput(customerID, fromID, toID, 50))

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result to be non-nil")
	}

	if result.FromBalance != 50 {
		t.Fatalf("expected source balance %d, got %d", 50, result.FromBalance)
	}

	if result.ToBalance != 70 {
		t.Fatalf("expected destination balance %d, got %d", 70, result.ToBalance)
	}

	outgoing := tx.createdTransactions[0]
	incoming := tx.createdTransactions[1]

	if result.TransactionReference == uuid.Nil {
		t.Fatal("expected transaction reference to be populated")
	}

	if outgoing.ReferenceID == nil || result.TransactionReference != *outgoing.ReferenceID {
		t.Fatalf("expected result transaction reference to match outgoing reference id")
	}

	if tx.decreaseCalls != 1 {
		t.Fatalf("expected debit once, got %d", tx.decreaseCalls)
	}

	if tx.updateCalls != 1 {
		t.Fatalf("expected credit once, got %d", tx.updateCalls)
	}

	if tx.createTransactionCalls != 2 {
		t.Fatalf("expected two ledger writes, got %d", tx.createTransactionCalls)
	}

	if outgoing.Type != domain.TransactionTransferOut {
		t.Fatalf("expected first ledger type %s, got %s", domain.TransactionTransferOut, outgoing.Type)
	}

	if incoming.Type != domain.TransactionTransferIn {
		t.Fatalf("expected second ledger type %s, got %s", domain.TransactionTransferIn, incoming.Type)
	}

	if outgoing.BalanceAfter != 50 {
		t.Fatalf("expected outgoing balance_after %d, got %d", 50, outgoing.BalanceAfter)
	}

	if incoming.BalanceAfter != 70 {
		t.Fatalf("expected incoming balance_after %d, got %d", 70, incoming.BalanceAfter)
	}

	if outgoing.ReferenceID == nil || incoming.ReferenceID == nil {
		t.Fatalf("expected both ledger entries to have reference id")
	}

	if *outgoing.ReferenceID != *incoming.ReferenceID {
		t.Fatalf("expected same reference id on both ledger entries")
	}

	if outgoing.RelatedAccountID == nil || *outgoing.RelatedAccountID != toID {
		t.Fatalf("expected outgoing related_account_id=%s, got %+v", toID, outgoing.RelatedAccountID)
	}

	if incoming.RelatedAccountID == nil || *incoming.RelatedAccountID != fromID {
		t.Fatalf("expected incoming related_account_id=%s, got %+v", fromID, incoming.RelatedAccountID)
	}

	if tx.commitCalls != 1 {
		t.Fatalf("expected commit once, got %d", tx.commitCalls)
	}

	firstLocked, secondLocked := orderedUUIDs(fromID, toID)
	if len(tx.lockedOrder) != 2 || tx.lockedOrder[0] != firstLocked || tx.lockedOrder[1] != secondLocked {
		t.Fatalf("expected deterministic lock order [%s %s], got %+v", firstLocked, secondLocked, tx.lockedOrder)
	}
}

func TestTransfer_Execute_PersistsDescriptionOnTransferOut(t *testing.T) {
	fromID := uuid.New()
	toID := uuid.New()
	customerID := uuid.New()
	fromBranch := "0001"
	fromNumber := "111111"
	toBranch := "0001"
	toNumber := "222222"
	description := "  Aluguel de maio  "
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		accountsByBranchNumber: map[string]*domain.Account{
			branchNumberKey(fromBranch, fromNumber): {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			branchNumberKey(toBranch, toNumber):     {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		decreaseBalanceValue: 50,
		updateBalanceValues:  map[uuid.UUID]int64{toID: 70},
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	_, err := useCase.Execute(context.Background(), TransferInput{
		User:           testCustomerUser(customerID),
		FromAccountID:  fromID,
		ToAccountID:    toID,
		Amount:         50,
		IdempotencyKey: "description-key",
		Description:    &description,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(tx.createdTransactions) != 2 {
		t.Fatalf("expected two ledger writes, got %d", len(tx.createdTransactions))
	}

	outgoing := tx.createdTransactions[0]
	incoming := tx.createdTransactions[1]
	if outgoing.Description == nil || *outgoing.Description != "Aluguel de maio" {
		t.Fatalf("expected normalized outgoing description, got %+v", outgoing.Description)
	}
	if incoming.Description != nil {
		t.Fatalf("expected incoming description to be nil, got %+v", incoming.Description)
	}
}

func TestTransfer_Execute_IdempotencyKeyAlreadyProcessed(t *testing.T) {
	fromID := uuid.New()
	toID := uuid.New()
	customerID := uuid.New()
	fromBranch := "0001"
	fromNumber := "111111"
	toBranch := "0001"
	toNumber := "222222"
	key := "idem-key-1"
	originalDescription := "Aluguel de maio"
	retryDescription := "Descricao alterada no retry"
	referenceID := uuid.New()

	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 50},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 70},
		},
		accountsByBranchNumber: map[string]*domain.Account{
			branchNumberKey(fromBranch, fromNumber): {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 50},
			branchNumberKey(toBranch, toNumber):     {ID: toID, Status: domain.AccountActive, Balance: 70},
		},
		getTransactionResult: &domain.Transaction{
			ID:               uuid.New(),
			AccountID:        fromID,
			Type:             domain.TransactionTransferOut,
			Amount:           50,
			BalanceAfter:     50,
			ReferenceID:      &referenceID,
			RelatedAccountID: &toID,
			IdempotencyKey:   &key,
			Description:      &originalDescription,
		},
		getTransactionByRef: &domain.Transaction{
			ID:               uuid.New(),
			AccountID:        toID,
			Type:             domain.TransactionTransferIn,
			Amount:           50,
			BalanceAfter:     70,
			ReferenceID:      &referenceID,
			RelatedAccountID: &fromID,
		},
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{
		User:           testCustomerUser(customerID),
		FromAccountID:  fromID,
		ToAccountID:    toID,
		Amount:         50,
		IdempotencyKey: key,
		Description:    &retryDescription,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result to be non-nil")
	}

	if result.FromBalance != 50 || result.ToBalance != 70 {
		t.Fatalf("expected balances from replay to be 50 and 70, got %d and %d", result.FromBalance, result.ToBalance)
	}

	if result.TransactionReference != referenceID {
		t.Fatalf("expected replay transaction reference %s, got %s", referenceID, result.TransactionReference)
	}

	if tx.decreaseCalls != 0 || tx.updateCalls != 0 {
		t.Fatalf("expected no balance mutation calls, got decrease=%d increase=%d", tx.decreaseCalls, tx.updateCalls)
	}

	if tx.createTransactionCalls != 0 {
		t.Fatalf("expected no ledger writes, got %d", tx.createTransactionCalls)
	}

	if *tx.getTransactionResult.Description != originalDescription {
		t.Fatalf("expected idempotent replay to preserve original description %q, got %q", originalDescription, *tx.getTransactionResult.Description)
	}

	if tx.commitCalls != 1 {
		t.Fatalf("expected commit once, got %d", tx.commitCalls)
	}
}

func TestTransfer_Execute_IdempotencyKeyScopedByResolvedSourceAccount(t *testing.T) {
	sourceAID := uuid.New()
	sourceBID := uuid.New()
	toID := uuid.New()
	customerID := uuid.New()
	fromBranchA := "0001"
	fromNumberA := "111111"
	fromBranchB := "0001"
	fromNumberB := "333333"
	toBranch := "0001"
	toNumber := "222222"
	key := "shared-idempotency-key"
	historicalReferenceID := uuid.New()
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			sourceBID: {ID: sourceBID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			toID:      {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		accountsByBranchNumber: map[string]*domain.Account{
			branchNumberKey(fromBranchA, fromNumberA): {ID: sourceAID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			branchNumberKey(fromBranchB, fromNumberB): {ID: sourceBID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			branchNumberKey(toBranch, toNumber):       {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		getTransactionByKeyFn: func(accountID uuid.UUID, gotKey string) (*domain.Transaction, error) {
			if gotKey != key {
				t.Fatalf("expected idempotency key %q, got %q", key, gotKey)
			}

			if accountID == sourceAID {
				return &domain.Transaction{
					ID:             uuid.New(),
					AccountID:      sourceAID,
					Type:           domain.TransactionTransferOut,
					Amount:         50,
					BalanceAfter:   50,
					ReferenceID:    &historicalReferenceID,
					IdempotencyKey: &key,
				}, nil
			}

			return nil, nil
		},
		decreaseBalanceValue: 50,
		updateBalanceValues:  map[uuid.UUID]int64{toID: 70},
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{
		User:           testCustomerUser(customerID),
		FromAccountID:  sourceBID,
		ToAccountID:    toID,
		Amount:         50,
		IdempotencyKey: key,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result to be non-nil")
	}

	if result.FromAccountID != sourceBID {
		t.Fatalf("expected transfer to execute from source B %s, got %s", sourceBID, result.FromAccountID)
	}

	if result.TransactionReference == historicalReferenceID {
		t.Fatalf("expected a new transaction reference for different source account, got historical reference %s", historicalReferenceID)
	}

	if tx.decreaseCalls != 1 || tx.updateCalls != 1 {
		t.Fatalf("expected financial effects once, got decrease=%d increase=%d", tx.decreaseCalls, tx.updateCalls)
	}

	if tx.createTransactionCalls != 2 {
		t.Fatalf("expected two ledger writes, got %d", tx.createTransactionCalls)
	}
}

func TestTransfer_Execute_IdempotencyConflictRollsBackDuplicateMutation(t *testing.T) {
	fromID := uuid.New()
	toID := uuid.New()
	customerID := uuid.New()
	fromBranch := "0001"
	fromNumber := "111111"
	toBranch := "0001"
	toNumber := "222222"
	key := "idem-key-2"
	referenceID := uuid.New()

	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		accountsByBranchNumber: map[string]*domain.Account{
			branchNumberKey(fromBranch, fromNumber): {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			branchNumberKey(toBranch, toNumber):     {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		decreaseBalanceValue: 50,
		updateBalanceValues:  map[uuid.UUID]int64{toID: 70},
		// First call returns nil (race: both requests passed the initial check).
		// Second call (after ErrTransferDuplicate) returns the committed entry.
		getTransactionResults: []*domain.Transaction{
			nil,
			{
				ID:               uuid.New(),
				AccountID:        fromID,
				Type:             domain.TransactionTransferOut,
				Amount:           50,
				BalanceAfter:     50,
				ReferenceID:      &referenceID,
				RelatedAccountID: &toID,
				IdempotencyKey:   &key,
			},
		},
		getTransactionByRef: &domain.Transaction{
			ID:               uuid.New(),
			AccountID:        toID,
			Type:             domain.TransactionTransferIn,
			Amount:           50,
			BalanceAfter:     70,
			ReferenceID:      &referenceID,
			RelatedAccountID: &fromID,
		},
		createTransactionErr: domain.ErrTransferDuplicate,
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{
		User:           testCustomerUser(customerID),
		FromAccountID:  fromID,
		ToAccountID:    toID,
		Amount:         50,
		IdempotencyKey: key,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result to be non-nil")
	}

	if result.TransactionReference != referenceID {
		t.Fatalf("expected conflict replay transaction reference %s, got %s", referenceID, result.TransactionReference)
	}

	if tx.decreaseCalls != 1 || tx.updateCalls != 1 {
		t.Fatalf("expected attempted mutation once before conflict, got decrease=%d increase=%d", tx.decreaseCalls, tx.updateCalls)
	}

	if tx.createTransactionCalls != 1 {
		t.Fatalf("expected one attempted ledger write before conflict rollback, got %d", tx.createTransactionCalls)
	}

	if tx.rollbackCalls != 1 {
		t.Fatalf("expected rollback once for duplicate execution, got %d", tx.rollbackCalls)
	}

	if tx.commitCalls != 0 {
		t.Fatalf("expected no commit on duplicate execution, got %d", tx.commitCalls)
	}
}

func TestTransfer_Execute_DebitFailure(t *testing.T) {
	fromID := uuid.New()
	toID := uuid.New()
	customerID := uuid.New()
	fromBranch := "0001"
	fromNumber := "111111"
	toBranch := "0001"
	toNumber := "222222"
	expectedErr := errors.New("debit failed")
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		accountsByBranchNumber: map[string]*domain.Account{
			branchNumberKey(fromBranch, fromNumber): {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			branchNumberKey(toBranch, toNumber):     {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		decreaseBalanceErr: expectedErr,
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), transferInput(customerID, fromID, toID, 50))

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
	}

	if tx.rollbackCalls != 1 {
		t.Fatalf("expected rollback once, got %d", tx.rollbackCalls)
	}
}

func TestTransfer_Execute_CreditFailure(t *testing.T) {
	fromID := uuid.New()
	toID := uuid.New()
	customerID := uuid.New()
	fromBranch := "0001"
	fromNumber := "111111"
	toBranch := "0001"
	toNumber := "222222"
	expectedErr := errors.New("credit failed")
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		accountsByBranchNumber: map[string]*domain.Account{
			branchNumberKey(fromBranch, fromNumber): {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			branchNumberKey(toBranch, toNumber):     {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		decreaseBalanceValue: 50,
		updateBalanceErr:     expectedErr,
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), transferInput(customerID, fromID, toID, 50))

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
	}

	if tx.decreaseCalls != 1 {
		t.Fatalf("expected debit once, got %d", tx.decreaseCalls)
	}

	if tx.rollbackCalls != 1 {
		t.Fatalf("expected rollback once, got %d", tx.rollbackCalls)
	}

	if tx.commitCalls != 0 {
		t.Fatalf("expected no commit, got %d", tx.commitCalls)
	}

	if tx.createTransactionCalls != 0 {
		t.Fatalf("expected no ledger writes, got %d", tx.createTransactionCalls)
	}
}

func TestTransfer_Execute_CommitFailure(t *testing.T) {
	fromID := uuid.New()
	toID := uuid.New()
	customerID := uuid.New()
	fromBranch := "0001"
	fromNumber := "111111"
	toBranch := "0001"
	toNumber := "222222"
	expectedErr := errors.New("commit failed")
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		accountsByBranchNumber: map[string]*domain.Account{
			branchNumberKey(fromBranch, fromNumber): {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			branchNumberKey(toBranch, toNumber):     {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		decreaseBalanceValue: 50,
		updateBalanceValues:  map[uuid.UUID]int64{toID: 70},
		commitErr:            expectedErr,
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), transferInput(customerID, fromID, toID, 50))

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
	}

	if tx.commitCalls != 1 {
		t.Fatalf("expected commit once, got %d", tx.commitCalls)
	}
}

func TestTransfer_Execute_ForbiddenForDifferentCustomer(t *testing.T) {
	fromID := uuid.New()
	toID := uuid.New()
	fromBranch := "0001"
	fromNumber := "111111"
	toBranch := "0001"
	toNumber := "222222"
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: uuid.New(), Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, CustomerID: uuid.New(), Status: domain.AccountActive, Balance: 20},
		},
		accountsByBranchNumber: map[string]*domain.Account{
			branchNumberKey(fromBranch, fromNumber): {ID: fromID, CustomerID: uuid.New(), Status: domain.AccountActive, Balance: 100},
			branchNumberKey(toBranch, toNumber):     {ID: toID, CustomerID: uuid.New(), Status: domain.AccountActive, Balance: 20},
		},
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{
		User:           testCustomerUser(uuid.New()),
		FromAccountID:  fromID,
		ToAccountID:    toID,
		Amount:         50,
		IdempotencyKey: "forbidden-key",
	})

	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected error %v, got %v", domain.ErrForbidden, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
	}

	if tx.decreaseCalls != 0 {
		t.Fatalf("expected no debit call, got %d", tx.decreaseCalls)
	}

	if tx.rollbackCalls != 1 {
		t.Fatalf("expected rollback once, got %d", tx.rollbackCalls)
	}
}

func TestTransfer_Execute_LedgerInsertFailure(t *testing.T) {
	fromID := uuid.New()
	toID := uuid.New()
	customerID := uuid.New()
	fromBranch := "0001"
	fromNumber := "111111"
	toBranch := "0001"
	toNumber := "222222"
	expectedErr := errors.New("ledger insert failed")
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		accountsByBranchNumber: map[string]*domain.Account{
			branchNumberKey(fromBranch, fromNumber): {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			branchNumberKey(toBranch, toNumber):     {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		decreaseBalanceValue: 50,
		updateBalanceValues:  map[uuid.UUID]int64{toID: 70},
		createTransactionErr: expectedErr,
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), transferInput(customerID, fromID, toID, 50))

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
	}

	if tx.createTransactionCalls != 1 {
		t.Fatalf("expected ledger write to fail on first insert, got %d calls", tx.createTransactionCalls)
	}

	if tx.rollbackCalls != 1 {
		t.Fatalf("expected rollback once, got %d", tx.rollbackCalls)
	}

	if tx.commitCalls != 0 {
		t.Fatalf("expected no commit, got %d", tx.commitCalls)
	}
}
