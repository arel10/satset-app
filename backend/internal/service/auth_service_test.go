package service_test

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"satset-app/internal/service"
)

func TestHashPassword(t *testing.T) {
	password := "SecurePassword123!"

	hash, err := service.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if len(hash) == 0 {
		t.Fatal("HashPassword returned empty string")
	}

	// Verify using bcrypt directly
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		t.Errorf("bcrypt verification failed: %v", err)
	}

	// Verify with wrong password fails
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("WrongPassword")); err == nil {
		t.Error("bcrypt verification should have failed with wrong password")
	}
}
