package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "my-secret-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error during hashing, got %v", err)
	}

	if hash == password {
		t.Errorf("Hash should not be identical to plain password")
	}

	// Test correct password
	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("Expected no error during check, got %v", err)
	}
	if !match {
		t.Errorf("Expected password to match it hash")
	}

	// Test incorrect password
	match, err = CheckPasswordHash("wrong-password", hash)
	if err != nil {
		t.Fatalf("Expected no error during check, got %v", err)
	}
	if match {
		t.Errorf("Expected wrong password NOT to match hash")
	}
}

func TestJWT(t *testing.T) {
	secret := "this-is-a-very-secure-secret-key"
	userID := uuid.New()
	expiry := time.Hour

	// Test generating token
	token, err := MakeJWT(userID, secret, expiry)
	if err != nil {
		t.Fatalf("Faild to create JWT, %v", err)
	}

	// Test validating valid token
	returnUUID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Errorf("Validation failed on valid token: %v", err)
	}
	if returnUUID != userID {
		t.Errorf("Expected uuid %v, but got %v", userID, returnUUID)
	}

	// Test vaildation with wrong token
	_, err = ValidateJWT(token, "this-is-the-wrong-secret")
	if err == nil {
		t.Errorf("Expected error when when validating with wrong secret, but got nil")
	}

	// Test expired token
	expiredToken, _ := MakeJWT(userID, secret, -time.Minute)
	_, err = ValidateJWT(expiredToken, secret)
	if err == nil {
		t.Errorf("Expected error for expired token, but got nil")
	}
}
