package infrastructure

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestNewBcryptPasswordHasher_DefaultCostFallback(t *testing.T) {
	hasher := NewBcryptPasswordHasher(-1)

	hash, err := hasher.Hash("S3curePass!")
	if err != nil {
		t.Fatalf("expected no error hashing password, got %v", err)
	}

	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		t.Fatalf("expected valid bcrypt hash, got %v", err)
	}

	if cost != bcrypt.DefaultCost {
		t.Fatalf("expected bcrypt cost %d, got %d", bcrypt.DefaultCost, cost)
	}
}

func TestBcryptPasswordHasher_HashAndCompare_Success(t *testing.T) {
	hasher := NewBcryptPasswordHasher(bcrypt.MinCost)
	password := "S3curePass!"

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

	if err := hasher.Compare(hash, password); err != nil {
		t.Fatalf("expected password comparison to succeed, got %v", err)
	}
}

func TestBcryptPasswordHasher_Compare_WrongPassword(t *testing.T) {
	hasher := NewBcryptPasswordHasher(bcrypt.MinCost)

	hash, err := hasher.Hash("S3curePass!")
	if err != nil {
		t.Fatalf("expected no error hashing password, got %v", err)
	}

	err = hasher.Compare(hash, "wrong-password")
	if err == nil {
		t.Fatal("expected compare to fail with wrong password")
	}
}
