package bcrypt

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("test123456")
	if err != nil {
		t.Errorf("HashPassword() error = %v", err)
	}
	if hash == "" {
		t.Error("HashPassword() returned empty hash")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "test123456"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if !CheckPasswordHash(password, hash) {
		t.Error("CheckPasswordHash() failed for correct password")
	}

	if CheckPasswordHash("wrong_password", hash) {
		t.Error("CheckPasswordHash() passed for wrong password")
	}
}
