package application

import (
	"context"
	"errors"
	"testing"

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

func (m *transferAccountRepositoryMock) UpdateBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return 0, nil
}

func (m *transferAccountRepositoryMock) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) error {
	return nil
}

func (m *transferAccountRepositoryMock) BeginTx(ctx context.Context) (domain.Tx, error) {
	m.beginTxCalls++
	if m.beginTxErr != nil {
		return nil, m.beginTxErr
	}
	return m.tx, nil
}

type transferTxMock struct {
	lockedOrder         []uuid.UUID
	accounts            map[uuid.UUID]*domain.Account
	getForUpdateErrs    map[uuid.UUID]error
	decreaseBalanceErr  error
	updateBalanceValues map[uuid.UUID]int64
	updateBalanceErr    error
	commitErr           error
	rollbackErr         error
	decreaseCalls       int
	updateCalls         int
	commitCalls         int
	rollbackCalls       int
}

func (m *transferTxMock) Create(ctx context.Context, account *domain.Account) error {
	return nil
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

func (m *transferTxMock) UpdateBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	m.updateCalls++
	if m.updateBalanceErr != nil {
		return 0, m.updateBalanceErr
	}
	if balance, ok := m.updateBalanceValues[id]; ok {
		return balance, nil
	}
	return 0, nil
}

func (m *transferTxMock) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) error {
	m.decreaseCalls++
	return m.decreaseBalanceErr
}

func (m *transferTxMock) BeginTx(ctx context.Context) (domain.Tx, error) {
	return nil, nil
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
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, Status: domain.AccountActive, Balance: 10},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{FromAccountID: fromID, ToAccountID: toID, Amount: 50})

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
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountInactive, Balance: 20},
		},
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{FromAccountID: fromID, ToAccountID: toID, Amount: 50})

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
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		updateBalanceValues: map[uuid.UUID]int64{toID: 70},
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{FromAccountID: fromID, ToAccountID: toID, Amount: 50})

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

	if tx.commitCalls != 1 {
		t.Fatalf("expected commit once, got %d", tx.commitCalls)
	}

	firstLocked, secondLocked := orderedUUIDs(fromID, toID)
	if len(tx.lockedOrder) != 2 || tx.lockedOrder[0] != firstLocked || tx.lockedOrder[1] != secondLocked {
		t.Fatalf("expected deterministic lock order [%s %s], got %+v", firstLocked, secondLocked, tx.lockedOrder)
	}
}

func TestTransfer_Execute_DebitFailure(t *testing.T) {
	fromID := uuid.New()
	toID := uuid.New()
	expectedErr := errors.New("debit failed")
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		decreaseBalanceErr: expectedErr,
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{FromAccountID: fromID, ToAccountID: toID, Amount: 50})

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
	expectedErr := errors.New("credit failed")
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		updateBalanceErr: expectedErr,
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{FromAccountID: fromID, ToAccountID: toID, Amount: 50})

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
}

func TestTransfer_Execute_CommitFailure(t *testing.T) {
	fromID := uuid.New()
	toID := uuid.New()
	expectedErr := errors.New("commit failed")
	tx := &transferTxMock{
		accounts: map[uuid.UUID]*domain.Account{
			fromID: {ID: fromID, Status: domain.AccountActive, Balance: 100},
			toID:   {ID: toID, Status: domain.AccountActive, Balance: 20},
		},
		updateBalanceValues: map[uuid.UUID]int64{toID: 70},
		commitErr:           expectedErr,
	}
	repo := &transferAccountRepositoryMock{tx: tx}
	useCase := NewTransfer(repo)

	result, err := useCase.Execute(context.Background(), TransferInput{FromAccountID: fromID, ToAccountID: toID, Amount: 50})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
	}

	if tx.commitCalls != 1 {
		t.Fatalf("expected commit once, got %d", tx.commitCalls)
	}

	if tx.rollbackCalls != 1 {
		t.Fatalf("expected rollback once after commit failure, got %d", tx.rollbackCalls)
	}
}
