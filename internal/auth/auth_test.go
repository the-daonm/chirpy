package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

// HashPassword and CheckPasswordHash tests
func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if hash == "password123" {
		t.Fatal("hash should not equal the original password")
	}
}

func TestCheckPasswordHash_Correct(t *testing.T) {
	hash, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("expected no error hashing, got: %v", err)
	}
	match, err := CheckPasswordHash("password123", hash)
	if err != nil {
		t.Fatalf("expected no error checking, got: %v", err)
	}
	if !match {
		t.Fatal("expected password to match hash")
	}
}

func TestCheckPasswordHash_Wrong(t *testing.T) {
	hash, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("expected no error hashing, got: %v", err)
	}
	match, err := CheckPasswordHash("wrongpassword", hash)
	if err != nil {
		t.Fatalf("expected no error checking, got: %v", err)
	}
	if match {
		t.Fatal("expected password to not match hash")
	}
}

func TestHashPassword_Unique(t *testing.T) {
	hash1, _ := HashPassword("password123")
	hash2, _ := HashPassword("password123")
	if hash1 == hash2 {
		t.Fatal("expected two hashes of same password to be different (salted)")
	}
}

// MakeJWT and ValidateJWT tests
func TestMakeJWT_ValidToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestValidateJWT_Valid(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("error creating token: %v", err)
	}
	parsedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("expected no error validating, got: %v", err)
	}
	if parsedID != userID {
		t.Fatalf("expected userID %v, got %v", userID, parsedID)
	}
}

func TestValidateJWT_Expired(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	token, err := MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("error creating token: %v", err)
	}
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()
	token, err := MakeJWT(userID, "correct-secret", time.Hour)
	if err != nil {
		t.Fatalf("error creating token: %v", err)
	}
	_, err = ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Fatal("expected error for wrong secret, got nil")
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	_, err := ValidateJWT("not.a.valid.token", "secret")
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}

// GetBearerToken tests
func TestGetBearerToken_Valid(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer mytoken123")
	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if token != "mytoken123" {
		t.Fatalf("expected mytoken123, got: %s", token)
	}
}

func TestGetBearerToken_Missing(t *testing.T) {
	headers := http.Header{}
	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("expected error for missing header, got nil")
	}
}

func TestGetBearerToken_InvalidFormat(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Basic mytoken123")
	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
}

func TestGetBearerToken_EmptyToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer ")
	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "" {
		t.Fatalf("expected empty token, got: %s", token)
	}
}

// MakeRefreshToken tests
func TestMakeRefreshToken_Length(t *testing.T) {
	token, err := MakeRefreshToken()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	// 32 bytes hex encoded = 64 characters
	if len(token) != 64 {
		t.Fatalf("expected token length 64, got: %d", len(token))
	}
}

func TestMakeRefreshToken_Unique(t *testing.T) {
	token1, _ := MakeRefreshToken()
	token2, _ := MakeRefreshToken()
	if token1 == token2 {
		t.Fatal("expected two refresh tokens to be different")
	}
}
