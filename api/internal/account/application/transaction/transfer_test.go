package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/domain"
)

type transferAccountRepositoryMock struct {
	beginTxCalls int
	beginTxErr   error
	tx           domain.Tx
}

func (m *transferAccountRepositoryMock) Create(ctx context.Context, account *domain.Account) error {
	return nil
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
	getForUpdateErrs         map[uuid.UUID]error
	decreaseBalanceValue     int64
	decreaseBalanceErr       error
	updateBalanceValues      map[uuid.UUID]int64
	updateBalanceErr         error
	createTransactionErr     error
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

func (m *transferTxMock) Create(ctx context.Context, account *domain.Account) error {
	return nil
}

func (m *transferTxMock) CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
	m.createTransactionCalls++
	m.createdTransactions = append(m.createdTransactions, tx)
	return m.createTransactionErr
}

func (m *transferTxMock) GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*domain.Transaction, error) {
	m.getTransactionByKeyCalls++
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

func TestTransfer_Execute_InvalidSourceID(t *testing.T) {
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

func TestTransfer_Execute_InvalidDestinationID(t *testing.T) {
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

func TestTransfer_Execute_SameAccount(t *testing.T) {
	accountID := uuid.New()
	repo := &transferAccountRepositoryMock{}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{
		FromAccountID: accountID,
		ToAccountID:   accountID,
		Amount:        10,
	})

	if !errors.Is(err, domain.ErrSameAccountTransfer) {
		t.Fatalf("expected error %v, got %v", domain.ErrSameAccountTransfer, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
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

	result, err := useCase.Execute(context.Background(), TransferInput{FromAccountID: fromID, ToAccountID: toID, Amount: 10})

	if !errors.Is(err, domain.ErrAccountNotFound) {
		t.Fatalf("expected error %v, got %v", domain.ErrAccountNotFound, err)
	}

	if err != domain.ErrAccountNotFound {
		t.Fatalf("expected direct ErrAccountNotFound, got %v", err)
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

	result, err := useCase.Execute(context.Background(), TransferInput{FromAccountID: fromID, ToAccountID: toID, Amount: 10})

	if !errors.Is(err, domain.ErrAccountNotFound) {
		t.Fatalf("expected error %v, got %v", domain.ErrAccountNotFound, err)
	}

	if err != domain.ErrAccountNotFound {
		t.Fatalf("expected direct ErrAccountNotFound, got %v", err)
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
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 10},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{User: testCustomerUser(customerID), FromAccountID: fromID, ToAccountID: toID, Amount: 50})

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
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountInactive, Balance: 20},
		},
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{User: testCustomerUser(customerID), FromAccountID: fromID, ToAccountID: toID, Amount: 50})

	if !errors.Is(err, domain.ErrAccountInactive) {
		t.Fatalf("expected error %v, got %v", domain.ErrAccountInactive, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
	}
}

func TestTransfer_Execute_Success(t *testing.T) {
	fromID := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	toID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	customerID := uuid.New()
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		decreaseBalanceValue: 50,
		updateBalanceValues:  map[uuid.UUID]int64{toID: 70},
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{User: testCustomerUser(customerID), FromAccountID: fromID, ToAccountID: toID, Amount: 50})

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

	if tx.decreaseCalls != 1 {
		t.Fatalf("expected debit once, got %d", tx.decreaseCalls)
	}

	if tx.updateCalls != 1 {
		t.Fatalf("expected credit once, got %d", tx.updateCalls)
	}

	if tx.createTransactionCalls != 2 {
		t.Fatalf("expected two ledger writes, got %d", tx.createTransactionCalls)
	}

	outgoing := tx.createdTransactions[0]
	incoming := tx.createdTransactions[1]

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

func TestTransfer_Execute_IdempotencyKeyAlreadyProcessed(t *testing.T) {
	fromID := uuid.New()
	toID := uuid.New()
	customerID := uuid.New()
	key := "idem-key-1"
	referenceID := uuid.New()

	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 50},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 70},
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

	if tx.decreaseCalls != 0 || tx.updateCalls != 0 {
		t.Fatalf("expected no balance mutation calls, got decrease=%d increase=%d", tx.decreaseCalls, tx.updateCalls)
	}

	if tx.createTransactionCalls != 0 {
		t.Fatalf("expected no ledger writes, got %d", tx.createTransactionCalls)
	}

	if tx.commitCalls != 1 {
		t.Fatalf("expected commit once, got %d", tx.commitCalls)
	}
}

func TestTransfer_Execute_IdempotencyConflictRollsBackDuplicateMutation(t *testing.T) {
	fromID := uuid.New()
	toID := uuid.New()
	customerID := uuid.New()
	key := "idem-key-2"
	referenceID := uuid.New()

	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
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
	expectedErr := errors.New("debit failed")
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		decreaseBalanceErr: expectedErr,
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{User: testCustomerUser(customerID), FromAccountID: fromID, ToAccountID: toID, Amount: 50})

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
	expectedErr := errors.New("credit failed")
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		decreaseBalanceValue: 50,
		updateBalanceErr:     expectedErr,
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{User: testCustomerUser(customerID), FromAccountID: fromID, ToAccountID: toID, Amount: 50})

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
	expectedErr := errors.New("commit failed")
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		decreaseBalanceValue: 50,
		updateBalanceValues:  map[uuid.UUID]int64{toID: 70},
		commitErr:            expectedErr,
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{User: testCustomerUser(customerID), FromAccountID: fromID, ToAccountID: toID, Amount: 50})

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
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: uuid.New(), Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, CustomerID: uuid.New(), Status: domain.AccountActive, Balance: 20},
		},
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{User: testCustomerUser(uuid.New()), FromAccountID: fromID, ToAccountID: toID, Amount: 50})

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
	expectedErr := errors.New("ledger insert failed")
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, CustomerID: customerID, Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		decreaseBalanceValue: 50,
		updateBalanceValues:  map[uuid.UUID]int64{toID: 70},
		createTransactionErr: expectedErr,
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{User: testCustomerUser(customerID), FromAccountID: fromID, ToAccountID: toID, Amount: 50})

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
