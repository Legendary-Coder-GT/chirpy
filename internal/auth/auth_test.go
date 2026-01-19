package auth

import (
	"testing"
	"time"
	"net/http"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	if token == "" {
		t.Error("MakeJWT returned empty token")
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"

	// Create a valid token
	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	// Validate the token
	returnedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if returnedID != userID {
		t.Errorf("ValidateJWT returned wrong user ID: got %v, want %v", returnedID, userID)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"

	// Create a token that expires immediately (negative duration)
	token, err := MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	// Attempt to validate the expired token
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Error("ValidateJWT should reject expired tokens")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()
	correctSecret := "correct-secret"
	wrongSecret := "wrong-secret"

	// Create a token with the correct secret
	token, err := MakeJWT(userID, correctSecret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	// Attempt to validate with the wrong secret
	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("ValidateJWT should reject tokens signed with wrong secret")
	}
}

func TestGetBearerToken(t *testing.T) {
	// Valid bearer token
	headers := http.Header{}
	headers.Set("Authorization", "Bearer my-token-string")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("GetBearerToken failed: %v", err)
	}

	if token != "my-token-string" {
		t.Errorf("GetBearerToken returned wrong token: got %q, want %q", token, "my-token-string")
	}
}

func TestGetBearerToken_NoAuthHeader(t *testing.T) {
	// Missing Authorization header entirely
	headers := http.Header{}

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Error("GetBearerToken should return error when Authorization header is missing")
	}
}

func TestGetBearerToken_EmptyAuthHeader(t *testing.T) {
	// Empty Authorization header
	headers := http.Header{}
	headers.Set("Authorization", "")

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Error("GetBearerToken should return error when Authorization header is empty")
	}
}

func TestGetBearerToken_NoBearerPrefix(t *testing.T) {
	// Authorization header without "Bearer" prefix
	headers := http.Header{}
	headers.Set("Authorization", "my-token-string")

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Error("GetBearerToken should return error when Bearer prefix is missing")
	}
}

func TestGetBearerToken_BearerOnly(t *testing.T) {
	// Authorization header with just "Bearer" and no token
	headers := http.Header{}
	headers.Set("Authorization", "Bearer ")

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Error("GetBearerToken should return error when token is missing after Bearer")
	}
}