package infrastructure

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestNewBcryptTransactionPasswordHasher_PanicsWhenPepperEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when pepper is empty")
		}
	}()

	_ = NewBcryptTransactionPasswordHasher(bcrypt.MinCost, "  ")
}

func TestBcryptTransactionPasswordHasher_HashAndCompare_Success(t *testing.T) {
	hasher := NewBcryptTransactionPasswordHasher(bcrypt.MinCost, "test-pepper-1234567890")
	password := "123456"

	hash, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("expected no error hashing password, got %v", err)
	}

	if hash == "" {
		t.Fatal("expected non-empty hash")
	}

	if hash == password {
		t.Fatal("expected hash to differ from password")
	}

	if !hasher.Compare(hash, password) {
		t.Fatal("expected compare to succeed for correct password and pepper")
	}
}

func TestBcryptTransactionPasswordHasher_Compare_WrongPassword(t *testing.T) {
	hasher := NewBcryptTransactionPasswordHasher(bcrypt.MinCost, "test-pepper-1234567890")

	hash, err := hasher.Hash("123456")
	if err != nil {
		t.Fatalf("expected no error hashing password, got %v", err)
	}

	if hasher.Compare(hash, "654321") {
		t.Fatal("expected compare to fail for wrong password")
	}
}

func TestBcryptTransactionPasswordHasher_Compare_FailsWithDifferentPepper(t *testing.T) {
	creator := NewBcryptTransactionPasswordHasher(bcrypt.MinCost, "pepper-a-123456789012345678901234")
	validator := NewBcryptTransactionPasswordHasher(bcrypt.MinCost, "pepper-b-123456789012345678901234")

	hash, err := creator.Hash("123456")
	if err != nil {
		t.Fatalf("expected no error hashing password, got %v", err)
	}

	if validator.Compare(hash, "123456") {
		t.Fatal("expected compare to fail when pepper differs")
	}
}

func TestBcryptTransactionPasswordHasher_HashProducesValidBcrypt(t *testing.T) {
	hasher := NewBcryptTransactionPasswordHasher(bcrypt.MinCost, "test-pepper-1234567890")

	hash, err := hasher.Hash("123456")
	if err != nil {
		t.Fatalf("expected no error hashing password, got %v", err)
	}

	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		t.Fatalf("expected valid bcrypt hash, got %v", err)
	}

	if cost != bcrypt.MinCost {
		t.Fatalf("expected bcrypt cost %d, got %d", bcrypt.MinCost, cost)
	}
}
