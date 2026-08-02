package services

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestPasswordHashingAndVerification(t *testing.T) {
	plainPassword := "ClashMaster2026!"

	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("expected successful password hash generation, got error: %v", err)
	}

	// Verify exact correct match
	err = bcrypt.CompareHashAndPassword(hashed, []byte(plainPassword))
	if err != nil {
		t.Errorf("expected generated bcrypt hash to match valid plain password, got error: %v", err)
	}

	// Verify wrong password fails match
	wrongPassword := "WrongPassword123!"
	err = bcrypt.CompareHashAndPassword(hashed, []byte(wrongPassword))
	if err == nil {
		t.Errorf("expected authentication error when comparing against incorrect password, but validation succeeded")
	}
}
